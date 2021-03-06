// Code generated by thriftrw

package broker

import (
	"errors"
	"fmt"
	"github.com/thriftrw/thriftrw-go/wire"
	"strings"
)

type SubscribeArgs struct {
	Topic *string `json:"topic,omitempty"`
}

func (v *SubscribeArgs) ToWire() (wire.Value, error) {
	var (
		fields [1]wire.Field
		i      int = 0
		w      wire.Value
		err    error
	)
	if v.Topic != nil {
		w, err = wire.NewValueString(*(v.Topic)), error(nil)
		if err != nil {
			return w, err
		}
		fields[i] = wire.Field{ID: 1, Value: w}
		i++
	}
	return wire.NewValueStruct(wire.Struct{Fields: fields[:i]}), nil
}

func (v *SubscribeArgs) FromWire(w wire.Value) error {
	var err error
	for _, field := range w.GetStruct().Fields {
		switch field.ID {
		case 1:
			if field.Value.Type() == wire.TBinary {
				var x string
				x, err = field.Value.GetString(), error(nil)
				v.Topic = &x
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (v *SubscribeArgs) String() string {
	var fields [1]string
	i := 0
	if v.Topic != nil {
		fields[i] = fmt.Sprintf("Topic: %v", *(v.Topic))
		i++
	}
	return fmt.Sprintf("SubscribeArgs{%v}", strings.Join(fields[:i], ", "))
}

func (v *SubscribeArgs) MethodName() string {
	return "subscribe"
}

func (v *SubscribeArgs) EnvelopeType() wire.EnvelopeType {
	return wire.Call
}

type SubscribeResult struct {
	Success *string `json:"success,omitempty"`
}

func (v *SubscribeResult) ToWire() (wire.Value, error) {
	var (
		fields [1]wire.Field
		i      int = 0
		w      wire.Value
		err    error
	)
	if v.Success != nil {
		w, err = wire.NewValueString(*(v.Success)), error(nil)
		if err != nil {
			return w, err
		}
		fields[i] = wire.Field{ID: 0, Value: w}
		i++
	}
	if i != 1 {
		return wire.Value{}, fmt.Errorf("SubscribeResult should have exactly one field: got %v fields", i)
	}
	return wire.NewValueStruct(wire.Struct{Fields: fields[:i]}), nil
}

func (v *SubscribeResult) FromWire(w wire.Value) error {
	var err error
	for _, field := range w.GetStruct().Fields {
		switch field.ID {
		case 0:
			if field.Value.Type() == wire.TBinary {
				var x string
				x, err = field.Value.GetString(), error(nil)
				v.Success = &x
				if err != nil {
					return err
				}
			}
		}
	}
	count := 0
	if v.Success != nil {
		count++
	}
	if count != 1 {
		return fmt.Errorf("SubscribeResult should have exactly one field: got %v fields", count)
	}
	return nil
}

func (v *SubscribeResult) String() string {
	var fields [1]string
	i := 0
	if v.Success != nil {
		fields[i] = fmt.Sprintf("Success: %v", *(v.Success))
		i++
	}
	return fmt.Sprintf("SubscribeResult{%v}", strings.Join(fields[:i], ", "))
}

func (v *SubscribeResult) MethodName() string {
	return "subscribe"
}

func (v *SubscribeResult) EnvelopeType() wire.EnvelopeType {
	return wire.Reply
}

var SubscribeHelper = struct {
	IsException    func(error) bool
	Args           func(topic *string) *SubscribeArgs
	WrapResponse   func(string, error) (*SubscribeResult, error)
	UnwrapResponse func(*SubscribeResult) (string, error)
}{}

func init() {
	SubscribeHelper.IsException = func(err error) bool {
		switch err.(type) {
		default:
			return false
		}
	}
	SubscribeHelper.Args = func(topic *string) *SubscribeArgs {
		return &SubscribeArgs{Topic: topic}
	}
	SubscribeHelper.WrapResponse = func(success string, err error) (*SubscribeResult, error) {
		if err == nil {
			return &SubscribeResult{Success: &success}, nil
		}
		return nil, err
	}
	SubscribeHelper.UnwrapResponse = func(result *SubscribeResult) (success string, err error) {
		if result.Success != nil {
			success = *result.Success
			return
		}
		err = errors.New("expected a non-void result")
		return
	}
}
