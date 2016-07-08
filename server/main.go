package main

import (
	"crypto/rand"
	"crypto/sha1"
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"github.com/bombela/yarpc-go-fun/broker"
	"github.com/bombela/yarpc-go-fun/broker/yarpc/brokerserver"
	"github.com/yarpc/yarpc-go"
	"github.com/yarpc/yarpc-go/encoding/thrift"
	"github.com/yarpc/yarpc-go/transport"
	ytch "github.com/yarpc/yarpc-go/transport/tchannel"

	"github.com/uber/tchannel-go"
)

func genSubscriberKey() string {
	c := 32
	buf := make([]byte, c)
	_, err := rand.Read(buf)
	if err != nil {
		log.Fatal(err)
	}
	key := fmt.Sprintf("%x", sha1.Sum(buf))
	return key
}

type subscription struct {
	// a buffer for incoming message to the subscriber.
	buffer chan string

	// >0 if currently being polled.
	// must use atomic to read/write.
	active uint32

	// last time this subscription was active, unix timestamp in seconds.
	// must use atomic to read/write.
	lastActiveTs int64
}

func newActiveSubscription() *subscription {
	s := subscription{
		buffer: make(chan string, 16),
	}
	return &s
}

func (s *subscription) isActive(at int64) bool {
	if atomic.LoadUint32(&s.active) > 0 {
		log.Printf("s? currently active")
		return true
	}
	if at-atomic.LoadInt64(&s.lastActiveTs) < 10 {
		log.Printf("s? currently still fresh %v", at-atomic.LoadInt64(&s.lastActiveTs))
		return true
	}
	return false
}

type activeTopic struct {
	sync.Mutex

	topic string

	// incoming messages to the topic.
	input chan<- string

	// currently active subscriptions on the topic.
	activeSubscriptions map[string]*subscription

	// a way for the topic to remove itself and its subscriptions from the
	// broker.
	broker *brokerHandler
}

func newTopic(b *brokerHandler, topic string) *activeTopic {
	log.Printf("new topic %v", topic)
	input := make(chan string, 16)
	t := activeTopic{
		topic:               topic,
		input:               input,
		activeSubscriptions: make(map[string]*subscription),
		broker:              b,
	}
	go func() {
	loop:
		for {
			select {
			case msg := <-input:
				t.Lock()
				t.dispatchAndGc(msg)
				t.Unlock()
			case <-time.After(5 * time.Second):
				t.Lock()
				anyActive := t.isAnyActiveSubcribers()
				t.Unlock()
				if !anyActive {
					break loop
				}
			}
		}
		// cleanup all the things!
		t.Lock()
		t.cleanup()
		t.Unlock()
	}()
	return &t
}

func (t *activeTopic) isAnyActiveSubcribers() bool {
	now := time.Now().Unix()
	for skey, s := range t.activeSubscriptions {
		if s.isActive(now) {
			log.Printf("s %v active", skey)
			return true
		}
		log.Printf("s %v kapout", skey)
	}
	return false
}

func (t *activeTopic) dispatchAndGc(msg string) {
	for _, s := range t.activeSubscriptions {
		select {
		case s.buffer <- msg:
		default: // ignore if buffer is full
		}
	}
}

func (t *activeTopic) cleanup() {
	log.Printf("cleaning up topic %v...", t.topic)
	t.broker.Lock()
	defer t.broker.Unlock()

	close(t.input)
	delete(t.broker.activeTopics, t.topic)
	for skey, s := range t.activeSubscriptions {
		close(s.buffer)
		delete(t.broker.activeSubscriptions, skey)
	}
	log.Printf("cleaned up topic %v", t.topic)
}

type brokerHandler struct {
	sync.RWMutex

	// currently active topics on the broker.
	activeTopics map[string]*activeTopic

	// currently active subscriptions, all of them across topics.
	activeSubscriptions map[string]*subscription
}

func newBrokerHandler() *brokerHandler {
	b := brokerHandler{
		activeTopics:        make(map[string]*activeTopic),
		activeSubscriptions: make(map[string]*subscription),
	}
	return &b
}

func (b *brokerHandler) getOrCreateTopic(topic string) *activeTopic {
	if t, found := b.activeTopics[topic]; found {
		return t
	}
	t := newTopic(b, topic)
	b.activeTopics[topic] = t
	return t
}

func (b *brokerHandler) Publish(reqMeta yarpc.ReqMeta,
	topic *string, message *string) (yarpc.ResMeta, error) {

	b.Lock()
	t := b.getOrCreateTopic(*topic)
	b.Unlock()

	t.input <- *message
	return nil, nil
}

func (b *brokerHandler) Subscribe(reqMeta yarpc.ReqMeta,
	topic *string) (string, yarpc.ResMeta, error) {

	skey := genSubscriberKey()
	s := newActiveSubscription()
	log.Printf("subscribing on %v (%v)", *topic, skey)

	b.Lock()
	defer b.Unlock()

	t := b.getOrCreateTopic(*topic)

	t.Lock()
	defer t.Unlock()

	t.activeSubscriptions[skey] = s
	b.activeSubscriptions[skey] = s

	log.Printf("subscribed on %v (%v)", *topic, skey)
	return skey, nil, nil
}

func (b *brokerHandler) Poll(reqMeta yarpc.ReqMeta, key *string,
	maxMsgs *int32) ([]string, yarpc.ResMeta, error) {

	b.RLock()
	s, found := b.activeSubscriptions[*key]
	b.RUnlock()

	if !found {
		return nil, nil, fmt.Errorf("Invalid subscription key %v", *key)
	}

	atomic.StoreUint32(&s.active, 1)
	ctx := reqMeta.Context()

	var msgs []string
	select {
	case msg := <-s.buffer:
		msgs = append(msgs, msg)
	case <-ctx.Done():
	}

loop:
	for i := int32(1); i < *maxMsgs; i++ {
		select {
		case msg := <-s.buffer:
			msgs = append(msgs, msg)
		case <-time.After(2 * time.Second):
			break loop
		case <-ctx.Done():
			break loop
		}
	}

	now := time.Now().Unix()
	atomic.StoreInt64(&s.lastActiveTs, now)
	atomic.StoreUint32(&s.active, 0)

	t, ok := ctx.Deadline()
	log.Printf("context deadline: %v, %v", t, ok)
	return msgs, nil, ctx.Err()
}

func (b *brokerHandler) ActiveSubscription(reqMeta yarpc.ReqMeta,
) ([]*broker.SubscribedTopic, yarpc.ResMeta, error) {
	return nil, nil, fmt.Errorf("boo")
}

func main() {
	channel, err := tchannel.NewChannel("broker", nil)
	if err != nil {
		log.Fatalln(err)
	}

	const port = ":28923"
	rpc := yarpc.New(yarpc.Config{
		Name: "broker",
		Inbounds: []transport.Inbound{
			ytch.NewInbound(channel, ytch.ListenAddr(port)),
		},
	})

	handler := newBrokerHandler()
	thrift.Register(rpc, brokerserver.New(handler))

	if err := rpc.Start(); err != nil {
		fmt.Println("error:", err.Error())
	}

	fmt.Println("Server listening on", port)

	select {} // block forever
}
