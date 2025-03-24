package hello

import (
	"fmt"
	"testing"
)

func TestFullSliceAppend(t *testing.T) {
	s := []int{100, 200, 300}
	fmt.Printf("Before append: len=%d cap=%d ptr=%p\n", len(s), cap(s), s)
	s = append(s, 4)
	fmt.Printf("After append: len=%d cap=%d ptr=%p\n", len(s), cap(s), s)
}

func TestRangeModifySlice(t *testing.T) {
	s := []int{100, 200, 300}
	var a [3]int
	for i, x := range s { // 切片拷贝仍然会引用相同的底层数组
		if i == 0 {
			s[1] = 201
			s[2] = 301
		}
		a[i] = x
	}
	fmt.Println(a) // [100, 201, 301]
}

func TestRangeModifyArray(t *testing.T) {
	s := [3]int{100, 200, 300}
	var a [3]int
	for i, x := range s { // 数组拷贝是整个值拷贝
		if i == 0 {
			s[1] = 201
			s[2] = 301
		}
		a[i] = x
	}
	fmt.Println(a) // [100, 200, 300]
}

func TestRangeAppend(t *testing.T) {
	s := []int{100, 200, 300}
	for i := range s {
		s = append(s, i)
	}
	fmt.Println(s)
}
