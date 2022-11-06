package domain

import (
	"math/rand"
)

type statelessPicker struct{}

func NewStatelessPicker() Picker {
	return &statelessPicker{}
}

func (p *statelessPicker) Pick(persons Persons) (Person, error) {
	index := rand.Int() % len(persons)
	return persons[index], nil
}
