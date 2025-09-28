package main

import "fmt"

type Target interface {
	Operation()
}

// адаптируемая структура
type Adaptee struct {
}

func (adaptee *Adaptee) AdaptedOperation() {
	fmt.Println("I am AdaptedOperation()")
}

type ConcreteAdapter struct {
	*Adaptee
}

func (adapter *ConcreteAdapter) Operation() {
	adapter.AdaptedOperation()
}

func NewAdapter(adaptee *Adaptee) Target {
	return &ConcreteAdapter{adaptee}
}

func main() {
	fmt.Println("Adapter:")
	adapter := NewAdapter(&Adaptee{})
	adapter.Operation()
}
