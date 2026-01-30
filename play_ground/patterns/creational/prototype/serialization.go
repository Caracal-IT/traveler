package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

type Address2 struct {
	StreetAddress, City, Country string
}

type Person2 struct {
	Name    string
	Address *Address2
	Friends []string
}

func (p *Person2) DeepCopy() *Person2 {
	// note: no error handling below
	b := bytes.Buffer{}
	e := gob.NewEncoder(&b)
	_ = e.Encode(p)

	// peek into structure
	fmt.Println(string(b.Bytes()))

	d := gob.NewDecoder(&b)
	result := Person2{}
	_ = d.Decode(&result)
	return &result
}

func main() {
	john := Person2{"John",
		&Address2{"123 London Rd", "London", "UK"},
		[]string{"Chris", "Matt", "Sam"}}

	jane := john.DeepCopy()
	jane.Name = "Jane"
	jane.Address.StreetAddress = "321 Baker St"
	jane.Friends = append(jane.Friends, "Jill")

	fmt.Println(john, john.Address)
	fmt.Println(jane, jane.Address)

}
