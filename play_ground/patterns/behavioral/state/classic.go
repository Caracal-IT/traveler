package state

import "fmt"

type Switch1 struct {
	State State1
}

func NewSwitch1() *Switch1 {
	return &Switch1{NewOffState1()}
}

func (s *Switch1) On() {
	s.State.On(s)
}

func (s *Switch1) Off() {
	s.State.Off(s)
}

type State1 interface {
	On(sw *Switch1)
	Off(sw *Switch1)
}

type BaseState1 struct{}

func (s *BaseState1) On(sw *Switch1) {
	fmt.Println("Light is already on")
}

func (s *BaseState1) Off(sw *Switch1) {
	fmt.Println("Light is already off")
}

type OnState1 struct {
	BaseState1
}

func NewOnState1() *OnState1 {
	fmt.Println("Light turned on")
	return &OnState1{BaseState1{}}
}

func (o *OnState1) Off(sw *Switch1) {
	fmt.Println("Turning light off...")
	sw.State = NewOffState1()
}

type OffState1 struct {
	BaseState1
}

func NewOffState1() *OffState1 {
	fmt.Println("Light turned off")
	return &OffState1{BaseState1{}}
}

func (o *OffState1) On(sw *Switch1) {
	fmt.Println("Turning light on...")
	sw.State = NewOnState1()
}

func main() {
	sw := NewSwitch1()
	sw.On()
	sw.Off()
	sw.Off()
}
