package domain

type Picker interface {
	Pick(persons Persons) (Person, error)
}
