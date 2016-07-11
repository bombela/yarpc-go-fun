package main

import (
	"flag"
	"fmt"
	"hash/fnv"
	"log"
	"sync"

	"github.com/bombela/yarpc-go-fun/broker"
	"github.com/bombela/yarpc-go-fun/broker/yarpc/brokerclient"
	"github.com/bombela/yarpc-go-fun/broker/yarpc/brokerserver"
	"github.com/bombela/yarpc-go-fun/broker/yarpc/loadbalancerserver"
	"github.com/yarpc/yarpc-go"
	"github.com/yarpc/yarpc-go/encoding/thrift"
	"github.com/yarpc/yarpc-go/transport"
	ytch "github.com/yarpc/yarpc-go/transport/tchannel"

	"github.com/uber/tchannel-go"
)

type backend struct {
	endpoint string
	con      *tchannel.Connection
}

type loadBalancerHandler struct {
	sync.RWMutex

	tch      *tchannel.Channel
	broker   brokerclient.Interface
	backends []*backend
}

func newLoadBalancerHandler(tch *tchannel.Channel, rpc yarpc.RPC) *loadBalancerHandler {
	b := loadBalancerHandler{
		tch:    tch,
		broker: brokerclient.New(rpc.Channel("broker")),
	}
	return &b
}

func (lb *loadBalancerHandler) selectBackend(topic string) *backend {
	hash := fnv.New32a()
	hash.Write([]byte(topic))
	topicHash := hash.Sum32()

	backendIdx := topicHash % uint32(len(lb.backends))
	backend := lb.backends[backendIdx]
	log.Printf("topic %v -> backend %v", topic, backend.endpoint)
	return backend
}

func (lb *loadBalancerHandler) Publish(reqMeta yarpc.ReqMeta,
	topic *string, message *string) (yarpc.ResMeta, error) {

	lb.RLock()
	defer lb.RUnlock()
	//backend := lb.selectBackend(*topic)
	//if backend == nil {
	//return nil, fmt.Errorf("No backends to shard onto!")
	//}
	ctx := reqMeta.Context()
	_, err := lb.broker.Publish(yarpc.NewReqMeta(ctx), topic, message)
	return nil, err
}

func (lb *loadBalancerHandler) Subscribe(reqMeta yarpc.ReqMeta,
	topic *string) (string, yarpc.ResMeta, error) {

	lb.RLock()
	defer lb.RUnlock()
	//backend := lb.selectBackend(*topic)
	//if backend == nil {
	//return "", nil, fmt.Errorf("No backends to shard onto!")
	//}
	ctx := reqMeta.Context()
	skey, _, err := lb.broker.Subscribe(yarpc.NewReqMeta(ctx), topic)
	return skey, nil, err
}

func (lb *loadBalancerHandler) Poll(reqMeta yarpc.ReqMeta, key *string,
	maxMsgs *int32) ([]string, yarpc.ResMeta, error) {

	ctx := reqMeta.Context()
	msgsChan := make(chan string)
	errChan := make(chan error)
	lb.RLock()
	for range lb.backends {
		go func() {
			msgs, _, err := lb.broker.Poll(yarpc.NewReqMeta(ctx), key, maxMsgs)
			if err != nil {
				errChan <- err
				return
			}
			for _, msg := range msgs {
				msgsChan <- msg
			}
		}()
	}
	lb.RUnlock()

	var msgs []string
	select {
	case msg := <-msgsChan:
		msgs = append(msgs, msg)
	case err := <-errChan:
		return nil, nil, err
	case <-ctx.Done():
	}
loop:
	for i := int32(1); i < *maxMsgs; i++ {
		select {
		case msg := <-msgsChan:
			msgs = append(msgs, msg)
		case err := <-errChan:
			return nil, nil, err
		case <-ctx.Done():
			break loop
		}
	}
	return msgs, nil, nil
}

func (lb *loadBalancerHandler) AddBackend(reqMeta yarpc.ReqMeta, endpoint *string) (yarpc.ResMeta, error) {
	log.Printf("Adding backend %v", *endpoint)
	lb.Lock()
	defer lb.Unlock()

	for _, backend := range lb.backends {
		if backend.endpoint == *endpoint {
			return nil, nil
		}
	}
	con, err := lb.tch.Connect(reqMeta.Context(), *endpoint)
	if err != nil {
		return nil, err
	}
	backend := backend{
		endpoint: *endpoint,
		con:      con,
	}
	lb.backends = append(lb.backends, &backend)
	return nil, nil
}

func (lb *loadBalancerHandler) DelBackend(reqMeta yarpc.ReqMeta,
	endpoint *string) (yarpc.ResMeta, error) {
	log.Printf("Removing backend %v", *endpoint)
	lb.Lock()
	defer lb.Unlock()
	for idx, backend := range lb.backends {
		if backend.endpoint == *endpoint {
			lb.backends = append(lb.backends[idx:], lb.backends[idx+1:]...)
		}
	}
	return nil, nil
}

func (lb *loadBalancerHandler) Backends(reqMeta yarpc.ReqMeta,
) ([]*broker.BackendStatus, yarpc.ResMeta, error) {
	log.Printf("Listing backends")
	lb.RLock()
	defer lb.RUnlock()
	r := make([]*broker.BackendStatus, 0, len(lb.backends))
	for _, backend := range lb.backends {
		status := "new"
		bs := broker.BackendStatus{
			Endpoint: backend.endpoint,
			Status:   &status,
		}
		r = append(r, &bs)
	}
	return r, nil, nil
}

func main() {
	bind := flag.String("bind", ":28924", "bind to this endpoint")
	flag.Parse()

	channel, err := tchannel.NewChannel("loadbalancer", nil)
	if err != nil {
		log.Fatalln(err)
	}

	backendsChannel, err := tchannel.NewChannel("loadbalancer", nil)
	if err != nil {
		log.Fatalln(err)
	}
	backendsOutbound := ytch.NewOutbound(backendsChannel)

	rpc := yarpc.New(yarpc.Config{
		Name: "loadbalancer",
		Inbounds: []transport.Inbound{
			ytch.NewInbound(channel, ytch.ListenAddr(*bind)),
		},
		Outbounds: transport.Outbounds{"broker": backendsOutbound},
	})

	loadBalancerHandler := newLoadBalancerHandler(backendsChannel, rpc)

	thrift.Register(rpc, loadbalancerserver.New(loadBalancerHandler))
	thrift.Register(rpc, brokerserver.New(loadBalancerHandler))

	if err := rpc.Start(); err != nil {
		fmt.Println("error:", err.Error())
	}

	fmt.Println("Server listening on", *bind)

	select {} // block forever
}
