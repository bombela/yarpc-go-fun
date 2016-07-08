package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"golang.org/x/net/context"

	"github.com/bombela/yarpc-go-fun/broker/yarpc/brokerclient"
	"github.com/uber/tchannel-go"
	"github.com/yarpc/yarpc-go"
	"github.com/yarpc/yarpc-go/transport"

	ytch "github.com/yarpc/yarpc-go/transport/tchannel"
)

func main() {
	topic := flag.String("topic", "", "topic to use")
	flag.Parse()

	if *topic == "" {
		log.Fatal("Give me a -topic <topic>")
	}

	channel, err := tchannel.NewChannel("broker-client", nil)
	if err != nil {
		log.Fatalln(err)
	}
	outbound := ytch.NewOutbound(channel, ytch.HostPort("localhost:28923"))

	rpc := yarpc.New(yarpc.Config{
		Name:      "broker-client",
		Outbounds: transport.Outbounds{"broker": outbound},
	})

	if err := rpc.Start(); err != nil {
		log.Fatalf("failed to start RPC: %v", err)
	}
	defer rpc.Stop()

	client := brokerclient.New(rpc.Channel("broker"))
	rootCtx := context.Background()

	fmt.Printf("# subscribing to topic %v\n", *topic)
	ctx, _ := context.WithTimeout(rootCtx, 250*time.Millisecond)
	key, _, err := client.Subscribe(yarpc.NewReqMeta(ctx), topic)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("# subscribed to topic %v with key %v\n", *topic, key)
	go func() {
		for {
			ctx, _ := context.WithTimeout(rootCtx, 15*time.Second)
			var maxMsgs int32 = 1
			msgs, _, err := client.Poll(yarpc.NewReqMeta(ctx), &key, &maxMsgs)
			if err != nil {
				if transport.IsTimeoutError(err) {
					log.Printf("Ignoring timeout %v", err)
				} else {
					log.Fatal(err)
				}
			}
			for _, msg := range msgs {
				fmt.Printf("<- %v:%v\n", *topic, msg)
			}
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		msg := scanner.Text()
		fmt.Printf("-> %v:%v\n", *topic, msg)
		ctx, _ := context.WithTimeout(rootCtx, 250*time.Millisecond)
		_, err := client.Publish(yarpc.NewReqMeta(ctx), topic, &msg)
		if err != nil {
			log.Println(err)
		}
	}
}
