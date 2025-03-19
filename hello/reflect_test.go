package hello

import (
	"fmt"
	"reflect"
	"testing"
)

type User struct {
	Id   int
	Name string
	Age  int
}

func (u *User) SayHello(str string) string {
	return "Hello " + str
}

func TestReflect(t *testing.T) {
	user := User{
		Name: "zs",
		Age:  38,
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
