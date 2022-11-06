package domain

func NewPerson(id string, name string) Person {
	return &person{
		id:   id,
		name: name,
	}
}

type person struct {
	id   string
	name string
}

type Person interface {
	Name() string
	ID() string
}

type Persons []Person

func (ps *Persons) RemoveByID(id string) Persons {
	var persons Persons
	for _, p := range *ps {
		if p.ID() != id {
			persons = append(persons, p)
		}
	}
	return persons
}

func (ps *Persons) Remove(person Person) Persons {
	var persons Persons
	for _, p := range *ps {
		if p != person {
			persons = append(persons, p)
		}
	}
	return persons
}

func (ps *Persons) FindByName(name string) Person {
	for _, p := range *ps {
		if p.Name() == name {
			return p
		}
	}
	return nil
}

func (p *person) Name() string {
	return p.name
}

func (p *person) ID() string {
	return p.id
}
