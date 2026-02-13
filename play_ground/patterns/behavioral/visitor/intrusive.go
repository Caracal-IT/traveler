package main

import (
	"fmt"
	"strings"
)

type Expression2 interface {
	Print(sb *strings.Builder)
}

type DoubleExpression2 struct {
	value float64
}

func (d *DoubleExpression2) Print(sb *strings.Builder) {
	sb.WriteString(fmt.Sprintf("%g", d.value))
}

type AdditionExpression2 struct {
	left, right Expression2
}

func (a *AdditionExpression2) Print(sb *strings.Builder) {
	sb.WriteString("(")
	a.left.Print(sb)
	sb.WriteString("+")
	a.right.Print(sb)
	sb.WriteString(")")
}

func main() {
	// 1+(2+3)
	e := AdditionExpression2{
		&DoubleExpression2{1},
		&AdditionExpression2{
			left:  &DoubleExpression2{2},
			right: &DoubleExpression2{3},
		},
	}
	sb := strings.Builder{}
	e.Print(&sb)
	fmt.Println(sb.String())
}
