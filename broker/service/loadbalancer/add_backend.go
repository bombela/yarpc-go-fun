// Code generated by thriftrw

package loadbalancer

import (
	"fmt"
	"github.com/thriftrw/thriftrw-go/wire"
	"strings"
)

type AddBackendArgs struct {
	Endpoint *string `json:"endpoint,omitempty"`
}

func (v *AddBackendArgs) ToWire() (wire.Value, error) {
	var (
		fields [1]wire.Field
		i      int = 0
		w      wire.Value
		err    error
	)
	if v.Endpoint != nil {
		w, err = wire.NewValueString(*(v.Endpoint)), error(nil)
		if err != nil {
			return w, err
		}
		fields[i] = wire.Field{ID: 1, Value: w}
		i++
	}
	return wire.NewValueStruct(wire.Struct{Fields: fields[:i]}), nil
}

func (v *AddBackendArgs) FromWire(w wire.Value) error {
	var err error
	for _, field := range w.GetStruct().Fields {
		switch field.ID {
		case 1:
			if field.Value.Type() == wire.TBinary {
				var x string
				x, err = field.Value.GetString(), error(nil)
				v.Endpoint = &x
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (v *AddBackendArgs) String() string {
	var fields [1]string
	i := 0
	if v.Endpoint != nil {
		fields[i] = fmt.Sprintf("Endpoint: %v", *(v.Endpoint))
		i++
	}
	return fmt.Sprintf("AddBackendArgs{%v}", strings.Join(fields[:i], ", "))
}

func (v *AddBackendArgs) MethodName() string {
	return "add_backend"
}

func (v *AddBackendArgs) EnvelopeType() wire.EnvelopeType {
	return wire.Call
}

type AddBackendResult struct{}

func (v *AddBackendResult) ToWire() (wire.Value, error) {
	var (
		fields [0]wire.Field
		i      int = 0
	)
	return wire.NewValueStruct(wire.Struct{Fields: fields[:i]}), nil
}

func (v *AddBackendResult) FromWire(w wire.Value) error {
	for _, field := range w.GetStruct().Fields {
		switch field.ID {
		}
	}
	return nil
}

func (v *AddBackendResult) String() string {
	var fields [0]string
	i := 0
	return fmt.Sprintf("AddBackendResult{%v}", strings.Join(fields[:i], ", "))
}

func (v *AddBackendResult) MethodName() string {
	return "add_backend"
}

func (v *AddBackendResult) EnvelopeType() wire.EnvelopeType {
	return wire.Reply
}

var AddBackendHelper = struct {
	IsException    func(error) bool
	Args           func(endpoint *string) *AddBackendArgs
	WrapResponse   func(error) (*AddBackendResult, error)
	UnwrapResponse func(*AddBackendResult) error
}{}

func init() {
	AddBackendHelper.IsException = func(err error) bool {
		switch err.(type) {
		default:
			return false
		}
	}
	AddBackendHelper.Args = func(endpoint *string) *AddBackendArgs {
		return &AddBackendArgs{Endpoint: endpoint}
	}
	AddBackendHelper.WrapResponse = func(err error) (*AddBackendResult, error) {
		if err == nil {
			return &AddBackendResult{}, nil
		}
		return nil, err
	}
	AddBackendHelper.UnwrapResponse = func(result *AddBackendResult) (err error) {
		return
	}
}
