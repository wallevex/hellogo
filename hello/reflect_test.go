package hello

import (
	"fmt"
	"reflect"
	"testing"
)

type User struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func (u *User) SayHello(str string) string {
	return "Hello " + str
}

func TestReflect(t *testing.T) {
	user := User{
		Id:   1,
		Name: "zs",
	}
	v := reflect.ValueOf(&user)

	method := v.MethodByName("SayHello")
	ret := method.Call([]reflect.Value{reflect.ValueOf("hello")})
	fmt.Println(ret)
}

func TestSetValue(t *testing.T) {
	var user User
	v := reflect.ValueOf(&user).Elem()

	v.FieldByName("Name").Set(reflect.ValueOf("张三"))

	fmt.Println(user)
}

func TestType(t *testing.T) {
	user := &User{}
	typ := reflect.TypeOf(user).Elem()
	tag := typ.Field(0).Tag.Get("json")
	fmt.Println(tag)
}
