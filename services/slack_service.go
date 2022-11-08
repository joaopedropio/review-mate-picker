package services

import (
	"fmt"
	"github.com/slack-go/slack"
	"sync"
)

type SlackService interface {
	GetUsers(userIDs []string) ([]SlackUser, error)
	GetChannelName(channelID string) (string, error)
	GetUser(userID string) (SlackUser, error)
	GetAllUsersFromChannel(channelID string, excludeBannedUsers bool) ([]SlackUser, error)
	ReplyMessage(channel string, text string, messageTimestamp string) error
}

func NewSlackService(client *slack.Client, allowedChannels []string, bannedUsers []string) SlackService {
	return &slackService{
		slackClient:     client,
		allowedChannels: allowedChannels,
		bannedUsers:     bannedUsers,
	}
}

type slackService struct {
	cache           sync.Map
	slackClient     *slack.Client
	allowedChannels []string
	bannedUsers     []string
}

func (s *slackService) GetChannelName(channelID string) (string, error) {
	channel, err := s.slackClient.GetConversationInfo(channelID, false)
	if err != nil {
		return "", fmt.Errorf("unable to get conversation info: %w", err)
	}
	return channel.Name, nil
}

func (s *slackService) GetUser(userID string) (SlackUser, error) {
	usr, ok := s.cache.Load(userID)
	if ok {
		return usr.(SlackUser), nil
	}
	user, err := s.getUser(userID)
	if err != nil {
		return SlackUser{}, err
	}
	s.cache.Store(userID, user)
	return user, nil
}

func (s *slackService) GetUsers(userIDs []string) ([]SlackUser, error) {
	var users []SlackUser
	for _, user := range userIDs {
		usr, err := s.GetUser(user)
		if err != nil {
			return nil, fmt.Errorf("unable to get user info: %w", err)
		}
		users = append(users, usr)
	}
	return users, nil
}

func (s *slackService) GetAllUsersFromChannel(channelID string, excludeBannedUsers bool) ([]SlackUser, error) {
	if err := s.validateChannel(channelID); err != nil {
		return nil, fmt.Errorf("unable to validate channel: channel id: %s: %w", channelID, err)
	}
	userIDs, _, err := s.slackClient.GetUsersInConversation(&slack.GetUsersInConversationParameters{
		ChannelID: channelID,
	})
	if err != nil {
		return nil, fmt.Errorf("unable to get user ids from slack: %w", err)
	}
	users, err := s.GetUsers(userIDs)
	if err != nil {
		return nil, fmt.Errorf("unable to get users from slack: %w", err)
	}
	if excludeBannedUsers {
		users = s.filterBannedUsers(users)
	}
	return users, nil
}

func (s *slackService) ReplyMessage(channel string, text string, messageTimestamp string) error {
	_, _, err := s.slackClient.PostMessage(channel, slack.MsgOptionText(text, false), slack.MsgOptionTS(messageTimestamp))
	return fmt.Errorf("unable to post message: %w", err)
}

func (s *slackService) getUser(userID string) (SlackUser, error) {
	user, err := s.slackClient.GetUserInfo(userID)
	if err != nil {
		return SlackUser{}, fmt.Errorf("unable to get users from slack: %w", err)
	}
	return SlackUser{
		ID:   user.ID,
		Name: user.Name,
	}, nil
}

func (s *slackService) validateChannel(channelID string) error {
	channelName, err := s.GetChannelName(channelID)
	if err != nil {
		return fmt.Errorf("unable to get channel name: channel id: %s: %w", channelID, err)
	}
	if !s.isChannelAllowed(channelName) {
		return fmt.Errorf("channel is not allowed: channel name: %s", channelName)
	}
	return nil
}

func (s *slackService) isChannelAllowed(channelName string) bool {
	for _, cName := range s.allowedChannels {
		if cName == channelName {
			return true
		}
	}
	return false
}

func (s *slackService) isBanned(user SlackUser) bool {
	return user.Name == "slackbot" || s.contains(s.bannedUsers, user.Name)
	//user.RealName == "Carlos Cabral" -> user.RealName is the name displayed on Slack
}

func (s *slackService) contains(list []string, value string) bool {
	for _, v := range list {
		if v == value {
			return true
		}
	}
	return false
}

func (s *slackService) filterBannedUsers(users []SlackUser) []SlackUser {
	var usrs []SlackUser
	for _, user := range users {
		if !s.isBanned(user) {
			usrs = append(usrs, user)
		}
	}
	return usrs
}
