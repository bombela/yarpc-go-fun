// Code generated by thriftrw

package broker

import (
	"errors"
	"fmt"
	"github.com/bombela/yarpc-go-fun/broker"
	"github.com/thriftrw/thriftrw-go/wire"
	"strings"
)

type ActiveSubscriptionArgs struct{}

func (v *ActiveSubscriptionArgs) ToWire() (wire.Value, error) {
	var (
		fields [0]wire.Field
		i      int = 0
	)
	return wire.NewValueStruct(wire.Struct{Fields: fields[:i]}), nil
}

func (v *ActiveSubscriptionArgs) FromWire(w wire.Value) error {
	for _, field := range w.GetStruct().Fields {
		switch field.ID {
		}
	}
	return nil
}

func (v *ActiveSubscriptionArgs) String() string {
	var fields [0]string
	i := 0
	return fmt.Sprintf("ActiveSubscriptionArgs{%v}", strings.Join(fields[:i], ", "))
}

func (v *ActiveSubscriptionArgs) MethodName() string {
	return "activeSubscription"
}

func (v *ActiveSubscriptionArgs) EnvelopeType() wire.EnvelopeType {
	return wire.Call
}

type ActiveSubscriptionResult struct {
	Success []*broker.SubscribedTopic `json:"success"`
}

type _List_SubscribedTopic_ValueList []*broker.SubscribedTopic

func (v _List_SubscribedTopic_ValueList) ForEach(f func(wire.Value) error) error {
	for _, x := range v {
		w, err := x.ToWire()
		if err != nil {
			return err
		}
		err = f(w)
		if err != nil {
			return err
		}
	}
	return nil
}

func (v _List_SubscribedTopic_ValueList) Size() int {
	return len(v)
}

func (_List_SubscribedTopic_ValueList) ValueType() wire.Type {
	return wire.TStruct
}

func (_List_SubscribedTopic_ValueList) Close() {
}

func (v *ActiveSubscriptionResult) ToWire() (wire.Value, error) {
	var (
		fields [1]wire.Field
		i      int = 0
		w      wire.Value
		err    error
	)
	if v.Success != nil {
		w, err = wire.NewValueList(_List_SubscribedTopic_ValueList(v.Success)), error(nil)
		if err != nil {
			return w, err
		}
		fields[i] = wire.Field{ID: 0, Value: w}
		i++
	}
	if i != 1 {
		return wire.Value{}, fmt.Errorf("ActiveSubscriptionResult should have exactly one field: got %v fields", i)
	}
	return wire.NewValueStruct(wire.Struct{Fields: fields[:i]}), nil
}

func _SubscribedTopic_Read(w wire.Value) (*broker.SubscribedTopic, error) {
	var v broker.SubscribedTopic
	err := v.FromWire(w)
	return &v, err
}

func _List_SubscribedTopic_Read(l wire.ValueList) ([]*broker.SubscribedTopic, error) {
	if l.ValueType() != wire.TStruct {
		return nil, nil
	}
	o := make([]*broker.SubscribedTopic, 0, l.Size())
	err := l.ForEach(func(x wire.Value) error {
		i, err := _SubscribedTopic_Read(x)
		if err != nil {
			return err
		}
		o = append(o, i)
		return nil
	})
	l.Close()
	return o, err
}

func (v *ActiveSubscriptionResult) FromWire(w wire.Value) error {
	var err error
	for _, field := range w.GetStruct().Fields {
		switch field.ID {
		case 0:
			if field.Value.Type() == wire.TList {
				v.Success, err = _List_SubscribedTopic_Read(field.Value.GetList())
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
		return fmt.Errorf("ActiveSubscriptionResult should have exactly one field: got %v fields", count)
	}
	return nil
}

func (v *ActiveSubscriptionResult) String() string {
	var fields [1]string
	i := 0
	if v.Success != nil {
		fields[i] = fmt.Sprintf("Success: %v", v.Success)
		i++
	}
	return fmt.Sprintf("ActiveSubscriptionResult{%v}", strings.Join(fields[:i], ", "))
}

func (v *ActiveSubscriptionResult) MethodName() string {
	return "activeSubscription"
}

func (v *ActiveSubscriptionResult) EnvelopeType() wire.EnvelopeType {
	return wire.Reply
}

var ActiveSubscriptionHelper = struct {
	IsException    func(error) bool
	Args           func() *ActiveSubscriptionArgs
	WrapResponse   func([]*broker.SubscribedTopic, error) (*ActiveSubscriptionResult, error)
	UnwrapResponse func(*ActiveSubscriptionResult) ([]*broker.SubscribedTopic, error)
}{}

func init() {
	ActiveSubscriptionHelper.IsException = func(err error) bool {
		switch err.(type) {
		default:
			return false
		}
	}
	ActiveSubscriptionHelper.Args = func() *ActiveSubscriptionArgs {
		return &ActiveSubscriptionArgs{}
	}
	ActiveSubscriptionHelper.WrapResponse = func(success []*broker.SubscribedTopic, err error) (*ActiveSubscriptionResult, error) {
		if err == nil {
			return &ActiveSubscriptionResult{Success: success}, nil
		}
		return nil, err
	}
	ActiveSubscriptionHelper.UnwrapResponse = func(result *ActiveSubscriptionResult) (success []*broker.SubscribedTopic, err error) {
		if result.Success != nil {
			success = result.Success
			return
		}
		err = errors.New("expected a non-void result")
		return
	}
}
