package hello

import (
	"fmt"
	"testing"
)

func TestRecover(t *testing.T) {
	defer func() { fmt.Println(recover()) }()
	defer func() { fmt.Println(recover()) }()
	defer panic(1)
	panic(2)
}
