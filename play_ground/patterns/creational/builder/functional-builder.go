package main

import "fmt"

type Person2 struct {
	name, position string
}

type personMod2 func(*Person2)

type PersonBuilder2 struct {
	actions []personMod2
}

func (b *PersonBuilder2) Called(name string) *PersonBuilder2 {
	b.actions = append(b.actions, func(p *Person2) {
		p.name = name
	})
	return b
}

func (b *PersonBuilder2) Build() *Person2 {
	p := Person2{}
	for _, a := range b.actions {
		a(&p)
	}
	return &p
}

// WorksAsA extend PersonBuilder
func (b *PersonBuilder2) WorksAsA(position string) *PersonBuilder2 {
	b.actions = append(b.actions, func(p *Person2) {
		p.position = position
	})
	return b
}

func main() {
	b := PersonBuilder2{}
	p := b.Called("Dmitri").WorksAsA("dev").Build()
	fmt.Println(*p)
}
