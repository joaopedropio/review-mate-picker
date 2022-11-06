package services

import (
	"fmt"
	"github.com/slack-go/slack"
)

type SlackService interface {
	GetAllUsers() ([]SlackUser, error)
	ReplyMessage(channel string, text string, messageTimestamp string) error
}

func NewSlackService(client *slack.Client) SlackService {
	return &slackService{slackClient: client}
}

type slackService struct {
	slackClient *slack.Client
}

func (s *slackService) GetAllUsers() ([]SlackUser, error) {
	users, err := s.slackClient.GetUsers()
	if err != nil {
		return nil, fmt.Errorf("unable to get users from slack: %w", err)
	}
	return s.parseUsers(users), nil
}

func (s *slackService) parseUsers(users []slack.User) []SlackUser {
	var slackUsers []SlackUser
	for _, user := range users {
		if s.isBanned(user) {
			continue
		}
		slackUsers = append(slackUsers, NewSlackUser(user.ID, user.Name))
	}
	return slackUsers
}

func (s *slackService) isBanned(user slack.User) bool {
	return user.IsBot || user.Name == "slackbot"
}

func (s *slackService) ReplyMessage(channel string, text string, messageTimestamp string) error {
	_, _, err := s.slackClient.PostMessage(channel, slack.MsgOptionText(text, false), slack.MsgOptionTS(messageTimestamp))
	return fmt.Errorf("unable to post message: %w", err)
}

type SlackUser struct {
	ID   string
	Name string
}

func NewSlackUser(id, name string) SlackUser {
	return SlackUser{
		ID:   id,
		Name: name,
	}
}
