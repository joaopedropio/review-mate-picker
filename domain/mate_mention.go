package domain

import (
	"fmt"
	"math/rand"
)

type MateMentionBuilder interface {
	Build(name string) string
}

type mateMentionBuilder struct {
	mentionTemplates []string
}

func NewMateMention() MateMentionBuilder {
	mentionTemplates := []string{
		"<@%s> can you review my pull request?",
		"<@%s> help this friend by reviewing his/her pull request?",
	}
	return &mateMentionBuilder{
		mentionTemplates: mentionTemplates,
	}
}

func (mm mateMentionBuilder) Build(name string) string {
	index := rand.Int() % len(mm.mentionTemplates)
	return fmt.Sprintf(mm.mentionTemplates[index], name)
}
