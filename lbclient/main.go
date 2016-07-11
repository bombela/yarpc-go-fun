package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"golang.org/x/net/context"

	"github.com/bombela/yarpc-go-fun/broker/yarpc/loadbalancerclient"
	"github.com/uber/tchannel-go"
	"github.com/yarpc/yarpc-go"
	"github.com/yarpc/yarpc-go/transport"

	ytch "github.com/yarpc/yarpc-go/transport/tchannel"
)

func main() {
	lb := flag.String("lb", "localhost:28924", "load balancer endpoint")
	action := flag.String("action", "ls", "add/del/ls")
	endpoint := flag.String("endpoint", "", "endpoint to add/del")
	flag.Parse()

	channel, err := tchannel.NewChannel("lbclient", nil)
	if err != nil {
		log.Fatalln(err)
	}
	outbound := ytch.NewOutbound(channel, ytch.HostPort(*lb))

	rpc := yarpc.New(yarpc.Config{
		Name:      "lbclient",
		Outbounds: transport.Outbounds{"loadbalancer": outbound},
	})

	if err := rpc.Start(); err != nil {
		log.Fatalf("failed to start RPC: %v", err)
	}
	defer rpc.Stop()

	client := loadbalancerclient.New(rpc.Channel("loadbalancer"))
	rootCtx := context.Background()

	ctx, _ := context.WithTimeout(rootCtx, 250*time.Millisecond)

	switch *action {
	case "ls":
	case "add":
	case "del":
	default:
		fmt.Printf("Unknown action: %s\n", *action)
		return

	}

	if *action != "ls" && *endpoint == "" {
		fmt.Printf("Give me an endpoint! (--endpoint)")
		return
	}

	switch *action {
	case "add":
		fmt.Printf("Adding backend %s\n", *endpoint)
		_, err := client.AddBackend(yarpc.NewReqMeta(ctx), endpoint)
		if err != nil {
			fmt.Printf("rpc err: %s\n", err)
			return
		}
	case "del":
		fmt.Printf("Removing backend %s\n", *endpoint)
		_, err := client.DelBackend(yarpc.NewReqMeta(ctx), endpoint)
		if err != nil {
			fmt.Printf("rpc err: %s\n", err)
			return
		}
	}

	backends, _, err := client.Backends(yarpc.NewReqMeta(ctx))
	if err != nil {
		fmt.Printf("rpc err: %s\n", err)
		return
	}
	fmt.Printf("%v backend(s):\n", len(backends))
	for _, b := range backends {
		fmt.Printf("- %v (%v)\n", b.Endpoint, *b.Status)
	}
}
