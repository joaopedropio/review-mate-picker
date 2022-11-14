package infrastructe

import (
	"fmt"
	"github.com/joaopedropio/review-mate-picker/domain"
	"os"
	"strconv"
	"strings"
)

type Environment interface {
	GetAllowedChannels() []string
	GetBannedUsers() []string
	GetSlackAuthToken() string
	GetSlackSigningSecret() string
	GetHttpPort() uint64
	GetPickingType() domain.PickingType
	IsProduction() bool
}

type environment struct {
	env             string
	token           string
	signingSecret   string
	httpPort        uint64
	pickingType     domain.PickingType
	allowedChannels []string
	bannedUsers     []string
}

func NewEnvironment() (Environment, error) {
	bu, ok := os.LookupEnv("BANNED_USERS")
	if !ok {
		return nil, fmt.Errorf("unable to get BANNED_USERS environment variable")
	}
	bannedUsers := parseStringList(bu)
	ac, ok := os.LookupEnv("ALLOWED_CHANNELS")
	if !ok {
		return nil, fmt.Errorf("unable to get ALLOWED_CHANNELS environment variable")
	}
	allowedChannels := parseStringList(ac)
	env, ok := os.LookupEnv("ENVIRONMENT")
	if !ok {
		return nil, fmt.Errorf("unable to get ENVIRONMENT environment variable")
	}
	pt, ok := os.LookupEnv("PICKING_TYPE")
	if !ok {
		return nil, fmt.Errorf("unable to get PICKING_TYPE environment variable")
	}
	pickingType, err := parsePickingType(pt)
	if err != nil {
		return nil, fmt.Errorf("unable to parse picking type: %w", err)
	}
	httpPort, ok := os.LookupEnv("PORT")
	if !ok {
		return nil, fmt.Errorf("unable to get PORT environment variable")
	}
	port, err := strconv.ParseUint(httpPort, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("unablt to parse http port: %w", err)
	}
	signingSecret, ok := os.LookupEnv("SLACK_SIGNING_SECRET")
	if !ok {
		return nil, fmt.Errorf("unable to get SLACK_SIGNING_SECRET environment variable")
	}
	token, ok := os.LookupEnv("SLACK_AUTH_TOKEN")
	if !ok {
		return nil, fmt.Errorf("unable to get SLACK_AUTH_TOKEN environment variable")
	}
	return &environment{
		token:           token,
		signingSecret:   signingSecret,
		httpPort:        port,
		pickingType:     pickingType,
		env:             env,
		allowedChannels: allowedChannels,
		bannedUsers:     bannedUsers,
	}, nil
}

func (e *environment) GetBannedUsers() []string {
	return e.bannedUsers
}

func (e *environment) GetAllowedChannels() []string {
	return e.allowedChannels
}

func (e *environment) GetSlackAuthToken() string {
	return e.token
}

func (e *environment) GetSlackSigningSecret() string {
	return e.signingSecret
}

func (e *environment) GetHttpPort() uint64 {
	return e.httpPort
}

func (e *environment) GetPickingType() domain.PickingType {
	return e.pickingType
}

func (e *environment) IsProduction() bool {
	return e.env == "production"
}

func parsePickingType(pickingType string) (domain.PickingType, error) {
	if pickingType == "stateless" {
		return domain.PickingTypeStateless, nil
	}
	if pickingType == "stateful" {
		return domain.PickingTypeStateful, nil
	}
	return "", fmt.Errorf("%s is no a valid picking type", pickingType)
}

func parseStringList(list string) []string {
	return strings.Split(list, ",")
}
