package kvs

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"github.com/mailru/easyjson"
	"math"
	"reflect"
)

var (
	Raw     Codec = rawCodec{}
	String  Codec = stringCodec{}
	Int64   Codec = int64Codec{}
	Float64 Codec = float64Codec{}
	JSON          = func(v interface{}, opts ...jsonCodecOpt) Codec { return newJsonCodec(v, opts...) }
)

type Codec interface {
	Marshal(v interface{}) ([]byte, error)
	Unmarshal(data []byte) (interface{}, error)
	String() string
}

type rawCodec struct{}

func (rawCodec) Marshal(v interface{}) ([]byte, error) {
	b, ok := v.([]byte)
	if !ok {
		return nil, errors.New(reflect.TypeOf(v).String() + " is not []byte")
	}
	return b, nil
}

func (rawCodec) Unmarshal(data []byte) (interface{}, error) {
	return data, nil
}

func (rawCodec) String() string {
	return "[]byte"
}

type stringCodec struct{}

func (stringCodec) Marshal(v interface{}) ([]byte, error) {
	b, ok := v.(string)
	if !ok {
		return nil, errors.New(reflect.TypeOf(v).String() + " is not string")
	}
	return []byte(b), nil
}

func (stringCodec) Unmarshal(data []byte) (interface{}, error) {
	return string(data), nil
}

func (stringCodec) String() string {
	return "string"
}

type int64Codec struct{}

func (int64Codec) Marshal(v interface{}) ([]byte, error) {
	i, ok := v.(int64)
	if !ok {
		return nil, errors.New(reflect.TypeOf(v).String() + " is not int64")
	}

	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(i))

	return b, nil
}

func (int64Codec) Unmarshal(data []byte) (interface{}, error) {
	bits := binary.LittleEndian.Uint64(data)

	return int64(bits), nil
}

func (int64Codec) String() string {
	return "int64"
}

type float64Codec struct{}

func (float64Codec) Marshal(v interface{}) ([]byte, error) {
	f, ok := v.(float64)
	if !ok {
		return nil, errors.New(reflect.TypeOf(v).String() + " is not float64")
	}

	bits := math.Float64bits(f)
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, bits)

	return b, nil
}

func (float64Codec) Unmarshal(data []byte) (interface{}, error) {
	bits := binary.LittleEndian.Uint64(data)

	return math.Float64frombits(bits), nil
}

func (float64Codec) String() string {
	return "float64"
}

type jsonCodec struct {
	prefix, indent string

	typ     reflect.Type
	ptr     bool
	marshal func(v interface{}) ([]byte, error)
}

type jsonCodecOpt func(*jsonCodec)

func WithIndent(prefix, indent string) jsonCodecOpt {
	return func(codec *jsonCodec) {
		codec.prefix = prefix
		codec.indent = indent
	}
}

func newJsonCodec(v interface{}, opts ...jsonCodecOpt) jsonCodec {
	jc := jsonCodec{
		typ:     reflect.TypeOf(v),
		marshal: json.Marshal,
	}

	if jc.typ.Kind() == reflect.Ptr {
		jc.typ = jc.typ.Elem()
		jc.ptr = true
	}

	for _, opt := range opts {
		opt(&jc)
	}
	if jc.prefix != "" || jc.indent != "" {
		jc.marshal = func(v interface{}) ([]byte, error) {
			return json.MarshalIndent(v, jc.prefix, jc.indent)
		}
	}

	return jc
}

func (c jsonCodec) Marshal(v interface{}) ([]byte, error) {
	if m, ok := v.(easyjson.Marshaler); ok {
		return easyjson.Marshal(m)
	}
	return c.marshal(v)
}

func (c jsonCodec) Unmarshal(data []byte) (interface{}, error) {
	v := reflect.New(c.typ)

	if err := c.unmarshal(data, v); err != nil {
		return nil, err
	}

	if c.ptr {
		return v.Interface(), nil
	}
	return v.Elem().Interface(), nil
}

func (c jsonCodec) unmarshal(data []byte, ptr reflect.Value) error {
	if um, ok := ptr.Interface().(easyjson.Unmarshaler); ok {
		return easyjson.Unmarshal(data, um)
	}

	return json.Unmarshal(data, ptr.Interface())
}

func (jsonCodec) String() string {
	return "JSON"
}
