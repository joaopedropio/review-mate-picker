package domain

import (
	"fmt"
	"math/rand"
)

type statefulPicker struct {
	personsNotPickedYet Persons
}

func NewStatefulPicker() Picker {
	return &statefulPicker{}
}

func (p *statefulPicker) Pick(persons Persons) (Person, error) {
	if len(persons) == 0 {
		return nil, fmt.Errorf("cannot pick an empty list")
	}
	p.refreshAvailablePersons(persons)
	return p.pick(), nil
}

func (p *statefulPicker) refreshAvailablePersons(personsAvailableToBePicked Persons) {
	p.removePersonsThatAreNotAvailableAnymore(personsAvailableToBePicked)
	if len(p.personsNotPickedYet) == 0 {
		p.personsNotPickedYet = personsAvailableToBePicked
	}
}

func (p *statefulPicker) removePersonsThatAreNotAvailableAnymore(personsAvailableToBePicked Persons) {
	for _, personNotPickedYet := range p.personsNotPickedYet {
		if personsAvailableToBePicked.FindByName(personNotPickedYet.Name()) == nil {
			p.personsNotPickedYet = p.personsNotPickedYet.Remove(personNotPickedYet)
		}
	}
}

func (p *statefulPicker) pick() Person {
	index := rand.Int() % len(p.personsNotPickedYet)
	person := p.personsNotPickedYet[index]
	p.personsNotPickedYet = p.personsNotPickedYet.Remove(person)
	return person
}
