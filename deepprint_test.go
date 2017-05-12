package deepprint

import (
	"fmt"
	"reflect"
	"testing"
)

type Info struct {
	updated bool
	eca     map[int]string
	extra   []string
}

type person struct {
	name    string
	age     int
	add     *string
	scores  map[string]float32
	info    *Info
	parents [2]string
	c       chan int
	f       F
}

type F func()

func dummy() {
}

func TestNestedStruct(t *testing.T) {
	add := "singapore"
	p := person{
		name:   "john",
		age:    21,
		add:    &add,
		scores: make(map[string]float32),
		info: &Info{
			eca: map[int]string{
				1000: "basketball",
				2000: "sockker",
			},
			extra: []string{"ABC", "DEF"},
		},
		parents: [2]string{"bill", "mary"},
		c:       make(chan int),
		f:       dummy,
	}
	p.scores["math"] = 98.3
	p.scores["history"] = 89.1
	p.scores["arts"] = 82.9
	s, err := DeepPrint(p)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(s)

	var i interface{} = p.c
	fmt.Printf("%v\n", reflect.ValueOf(i).Kind())
}
