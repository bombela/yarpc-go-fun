// Code generated by thriftrw

package brokerserver

import (
	"github.com/bombela/yarpc-go-fun/broker"
	broker2 "github.com/bombela/yarpc-go-fun/broker/service/broker"
	"github.com/thriftrw/thriftrw-go/protocol"
	"github.com/thriftrw/thriftrw-go/wire"
	yarpc "github.com/yarpc/yarpc-go"
	"github.com/yarpc/yarpc-go/encoding/thrift"
)

type Interface interface {
	ActiveSubscription(reqMeta yarpc.ReqMeta) ([]*broker.SubscribedTopic, yarpc.ResMeta, error)
	Poll(reqMeta yarpc.ReqMeta, key *string, maxMsgs *int32) ([]string, yarpc.ResMeta, error)
	Publish(reqMeta yarpc.ReqMeta, topic *string, message *string) (yarpc.ResMeta, error)
	Subscribe(reqMeta yarpc.ReqMeta, topic *string) (string, yarpc.ResMeta, error)
}

func New(impl Interface) thrift.Service {
	return service{handler{impl}}
}

type service struct{ h handler }

func (service) Name() string {
	return "Broker"
}

func (service) Protocol() protocol.Protocol {
	return protocol.Binary
}

func (s service) Handlers() map[string]thrift.Handler {
	return map[string]thrift.Handler{"activeSubscription": thrift.HandlerFunc(s.h.ActiveSubscription), "poll": thrift.HandlerFunc(s.h.Poll), "publish": thrift.HandlerFunc(s.h.Publish), "subscribe": thrift.HandlerFunc(s.h.Subscribe)}
}

type handler struct{ impl Interface }

func (h handler) ActiveSubscription(reqMeta yarpc.ReqMeta, body wire.Value) (thrift.Response, error) {
	var args broker2.ActiveSubscriptionArgs
	if err := args.FromWire(body); err != nil {
		return thrift.Response{}, err
	}
	success, resMeta, err := h.impl.ActiveSubscription(reqMeta)
	hadError := err != nil
	result, err := broker2.ActiveSubscriptionHelper.WrapResponse(success, err)
	var response thrift.Response
	if err == nil {
		response.IsApplicationError = hadError
		response.Meta = resMeta
		response.Body = result
	}
	return response, err
}

func (h handler) Poll(reqMeta yarpc.ReqMeta, body wire.Value) (thrift.Response, error) {
	var args broker2.PollArgs
	if err := args.FromWire(body); err != nil {
		return thrift.Response{}, err
	}
	success, resMeta, err := h.impl.Poll(reqMeta, args.Key, args.MaxMsgs)
	hadError := err != nil
	result, err := broker2.PollHelper.WrapResponse(success, err)
	var response thrift.Response
	if err == nil {
		response.IsApplicationError = hadError
		response.Meta = resMeta
		response.Body = result
	}
	return response, err
}

func (h handler) Publish(reqMeta yarpc.ReqMeta, body wire.Value) (thrift.Response, error) {
	var args broker2.PublishArgs
	if err := args.FromWire(body); err != nil {
		return thrift.Response{}, err
	}
	resMeta, err := h.impl.Publish(reqMeta, args.Topic, args.Message)
	hadError := err != nil
	result, err := broker2.PublishHelper.WrapResponse(err)
	var response thrift.Response
	if err == nil {
		response.IsApplicationError = hadError
		response.Meta = resMeta
		response.Body = result
	}
	return response, err
}

func (h handler) Subscribe(reqMeta yarpc.ReqMeta, body wire.Value) (thrift.Response, error) {
	var args broker2.SubscribeArgs
	if err := args.FromWire(body); err != nil {
		return thrift.Response{}, err
	}
	success, resMeta, err := h.impl.Subscribe(reqMeta, args.Topic)
	hadError := err != nil
	result, err := broker2.SubscribeHelper.WrapResponse(success, err)
	var response thrift.Response
	if err == nil {
		response.IsApplicationError = hadError
		response.Meta = resMeta
		response.Body = result
	}
	return response, err
}
