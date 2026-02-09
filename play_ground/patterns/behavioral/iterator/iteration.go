package main

import "fmt"

type Person struct {
	FirstName, MiddleName, LastName string
}

func (p *Person) Names() []string {
	return []string{p.FirstName, p.MiddleName, p.LastName}
}

func (p *Person) NamesGenerator() <-chan string {
	ch := make(chan string)
	go func() {
		defer close(ch)
		for _, name := range []string{p.FirstName, p.MiddleName, p.LastName} {
			ch <- name
		}
	}()
	return ch
}

func main() {
	p := Person{"Alexander", "Graham", "Bell"}
	for _, name := range p.Names() {
		fmt.Println(name)
	}

	for name := range p.NamesGenerator() {
		fmt.Println(name)
	}
}
