package main

import (
	"container/list"
	"fmt"
)

type Observable3 struct {
	subs *list.List
}

func (o *Observable3) Subscribe(x Observer3) {
	o.subs.PushBack(x)
}

func (o *Observable3) Unsubscribe(x Observer3) {
	for z := o.subs.Front(); z != nil; z = z.Next() {
		if z.Value.(Observer3) == x {
			o.subs.Remove(z)
		}
	}
}

func (o *Observable3) Fire(data interface{}) {
	for z := o.subs.Front(); z != nil; z = z.Next() {
		z.Value.(Observer3).Notify(data)
	}
}

type Observer3 interface {
	Notify(data interface{})
}

type Person3 struct {
	Observable3
	age int
}

func NewPerson3(age int) *Person3 {
	return &Person3{Observable3{new(list.List)}, age}
}

type PropertyChanged3 struct {
	Name  string
	Value interface{}
}

func (p *Person3) Age() int { return p.age }
func (p *Person3) SetAge(age int) {
	if age == p.age {
		return
	} // no change

	oldCanVote := p.CanVote()

	p.age = age
	p.Fire(PropertyChanged3{"Age", p.age})

	if oldCanVote != p.CanVote() {
		p.Fire(PropertyChanged3{"CanVote", p.CanVote()})
	}
}

func (p *Person3) CanVote() bool {
	return p.age >= 18
}

type ElectrocalRoll3 struct {
}

func (e *ElectrocalRoll3) Notify(data interface{}) {
	if pc, ok := data.(PropertyChanged3); ok {
		if pc.Name == "CanVote" && pc.Value.(bool) {
			fmt.Println("Congratulations, you can vote!")
		}
	}
}

func main() {
	p := NewPerson3(0)
	er := &ElectrocalRoll3{}
	p.Subscribe(er)

	for i := 10; i < 20; i++ {
		fmt.Println("Setting age to", i)
		p.SetAge(i)
	}
}
