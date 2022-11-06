package services

import (
	"fmt"
	"github.com/joaopedropio/review-mate-picker/domain"
)

type MateService interface {
	PickMateToReview(channelID string, userID string, messageTimestamp string) error
}

type mateService struct {
	slackService SlackService
	picker       domain.Picker
}

func NewMateService(slackService SlackService, picker domain.Picker) MateService {
	return &mateService{
		slackService: slackService,
		picker:       picker,
	}
}

func (s *mateService) PickMateToReview(channelID string, userID string, messageTimestamp string) error {
	users, err := s.slackService.GetAllUsers()
	if err != nil {
		return fmt.Errorf("unable to get all users by channel: %w", err)
	}
	persons := s.parseUsersToPersons(users)
	persons = persons.RemoveByID(userID)
	mate, err := s.picker.Pick(persons)
	if err != nil {
		return fmt.Errorf("unable to pick mate: %w", err)
	}
	mateMention := domain.NewMateMention(mate.Name()).Build()
	if err := s.slackService.ReplyMessage(channelID, mateMention, messageTimestamp); err != nil {
		return fmt.Errorf("unable to reply message: %w", err)
	}
	return nil
}

func (s *mateService) parseUsersToPersons(users []SlackUser) domain.Persons {
	var persons domain.Persons
	for _, user := range users {
		persons = append(persons, domain.NewPerson(user.ID, user.Name))
	}
	return persons
}
