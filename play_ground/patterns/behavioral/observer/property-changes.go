package main

import (
	"container/list"
	"fmt"
)

type Observable2 struct {
	subs *list.List
}

func (o *Observable2) Subscribe(x Observer2) {
	o.subs.PushBack(x)
}

func (o *Observable2) Unsubscribe(x Observer2) {
	for z := o.subs.Front(); z != nil; z = z.Next() {
		if z.Value.(Observer2) == x {
			o.subs.Remove(z)
		}
	}
}

func (o *Observable2) Fire(data interface{}) {
	for z := o.subs.Front(); z != nil; z = z.Next() {
		z.Value.(Observer2).Notify(data)
	}
}

type Observer2 interface {
	Notify(data interface{})
}

type Person2 struct {
	Observable2
	age int
}

func NewPerson2(age int) *Person2 {
	return &Person2{Observable2{new(list.List)}, age}
}

type PropertyChanged2 struct {
	Name  string
	Value interface{}
}

func (p *Person2) Age() int { return p.age }
func (p *Person2) SetAge(age int) {
	if age == p.age {
		return
	} // no change
	p.age = age
	p.Fire(PropertyChanged2{"Age", p.age})
}

type TrafficManagement struct {
	o Observable2
}

func (t *TrafficManagement) Notify(data interface{}) {
	if pc, ok := data.(PropertyChanged2); ok {
		if pc.Value.(int) >= 16 {
			fmt.Println("Congrats, you can drive now!")
			// we no longer care
			t.o.Unsubscribe(t)
		}
	}
}

func main() {
	p := NewPerson2(15)
	t := &TrafficManagement{p.Observable2}
	p.Subscribe(t)

	for i := 16; i <= 20; i++ {
		fmt.Println("Setting age to", i)
		p.SetAge(i)
	}
}
