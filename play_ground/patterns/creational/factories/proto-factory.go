package main

import "fmt"

type Employee2 struct {
	Name, Position string
	AnnualIncome   int
}

type Role int

const (
	Developer Role = iota
	Manager
)

// functional
func NewEmployee2(role Role) *Employee2 {
	switch role {
	case Developer:
		return &Employee2{"", "Developer", 60000}
	case Manager:
		return &Employee2{"", "Manager", 80000}
	default:
		panic("unsupported role")
	}
}

func main() {
	mm := NewEmployee2(1)
	m := NewEmployee2(Manager)
	m.Name = "Sam"
	fmt.Println(m)

	fmt.Println(mm)
}
