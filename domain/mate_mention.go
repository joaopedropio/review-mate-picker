package domain

import (
	"fmt"
	"math/rand"
)

type MateMention struct {
	name             string
	mentionTemplates []string
}

func NewMateMention(name string) MateMention {
	mentionTemplates := []string{
		"<@%s> can you review my pull request?",
		"<@%s> help this friend by reviewing his/her pull request?",
	}
	return MateMention{
		name:             name,
		mentionTemplates: mentionTemplates,
	}
}

func (mm MateMention) Build() string {
	index := rand.Int() % len(mm.mentionTemplates)
	return fmt.Sprintf(mm.mentionTemplates[index], mm.name)
}
