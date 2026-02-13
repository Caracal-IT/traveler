package main

import (
	"fmt"
	"strings"
)

type ExpressionVisitor1 interface {
	VisitDoubleExpression1(de *DoubleExpression1)
	VisitAdditionExpression1(ae *AdditionExpression1)
}

type Expression1 interface {
	Accept(ev ExpressionVisitor1)
}

type DoubleExpression1 struct {
	value float64
}

func (d *DoubleExpression1) Accept(ev ExpressionVisitor1) {
	ev.VisitDoubleExpression1(d)
}

type AdditionExpression1 struct {
	left, right Expression1
}

func (a *AdditionExpression1) Accept(ev ExpressionVisitor1) {
	ev.VisitAdditionExpression1(a)
}

type ExpressionPrinter1 struct {
	sb strings.Builder
}

func (e *ExpressionPrinter1) VisitDoubleExpression1(de *DoubleExpression1) {
	e.sb.WriteString(fmt.Sprintf("%g", de.value))
}

func (e *ExpressionPrinter1) VisitAdditionExpression1(ae *AdditionExpression1) {
	e.sb.WriteString("(")
	ae.left.Accept(e)
	e.sb.WriteString("+")
	ae.right.Accept(e)
	e.sb.WriteString(")")
}

func NewExpressionPrinter1() *ExpressionPrinter1 {
	return &ExpressionPrinter1{strings.Builder{}}
}

func (e *ExpressionPrinter1) String() string {
	return e.sb.String()
}

func main() {
	// 1+(2+3)
	e := &AdditionExpression1{
		&DoubleExpression1{1},
		&AdditionExpression1{
			left:  &DoubleExpression1{2},
			right: &DoubleExpression1{3},
		},
	}
	ep := NewExpressionPrinter1()
	ep.VisitAdditionExpression1(e)
	fmt.Println(ep.String())
}
