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

type I interface {
	do()
}

type iImpl struct {
}

func (i *iImpl) do() {

}

type person struct {
	name      string
	age       int
	add       *string
	scores    map[string]float32
	info      *Info
	parents   [2]string
	c         chan int
	f         F
	p         *string
	intf      I
	intf2     I
	emptyIntf interface{}
}

type F func()

func dummy() {
}

func TestNestedStruct(t *testing.T) {
	add := "singapore"
	//var itf *iImpl
	var itf2 I
	var emptyInf interface{}
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
		parents:   [2]string{"bill", "mary"},
		c:         make(chan int),
		f:         dummy,
		intf:      &iImpl{},
		intf2:     itf2,
		emptyIntf: emptyInf,
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

type A struct {
	b *B
}

type B struct {
	a *A
}

func TestCircularReference(t *testing.T) {
	a := &A{}
	b := &B{}
	a.b = b
	b.a = a
	s, err := DeepPrint(a)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(s)
}
