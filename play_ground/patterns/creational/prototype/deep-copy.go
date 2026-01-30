package main

import "fmt"

type Address4 struct {
	StreetAddress, City, Country string
}

type Person4 struct {
	Name    string
	Address *Address4
}

func main() {
	john := Person4{"John",
		&Address4{"123 London Rd", "London", "UK"}}

	//jane := john

	// shallow copy
	//jane.Name = "Jane" // ok

	//jane.Address.StreetAddress = "321 Baker St"

	//fmt.Println(john.Name, john.Address)
	//fmt.Println(jane.Name, jane. Address)

	// what you really want
	jane := john
	jane.Address = &Address4{
		john.Address.StreetAddress,
		john.Address.City,
		john.Address.Country}

	jane.Name = "Jane" // ok

	jane.Address.StreetAddress = "321 Baker St"

	fmt.Println(john.Name, john.Address)
	fmt.Println(jane.Name, jane.Address)
}
