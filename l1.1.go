package main

import "fmt"

type Human struct {
	Name string
	Age  int
	Sex  string
}

func NewHuman(name string, age int, sex string) *Human {
	return &Human{
		Name: name,
		Age:  age,
		Sex:  sex,
	}

}

func (h *Human) setName(name string) {
	h.Name = name
}

func (h Human) getName() string {
	return h.Name
}

type Action struct {
	*Human
	Sleep bool
}

func NewAction(h *Human) *Action {
	if h == nil {
		h = &Human{}
	}
	return &Action{Human: h}
}

func main() {
	human1 := NewHuman("Bob", 21, "male")
	action := NewAction(human1)

	fmt.Println(action.getName())
	action.setName("Dave")
	fmt.Println(action.getName())
	fmt.Println(human1.getName())

}
