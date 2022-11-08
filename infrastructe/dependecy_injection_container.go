package infrastructe

import (
	"github.com/joaopedropio/review-mate-picker/domain"
	"github.com/joaopedropio/review-mate-picker/repositories"
	"github.com/joaopedropio/review-mate-picker/services"
	"github.com/slack-go/slack"
)

type Container struct {
	SlackClient    *slack.Client
	SlackService   services.SlackService
	Picker         domain.Picker
	MateService    services.MateService
	TimestampCache *repositories.MessageTimestampCache
}

func NewDependencyInjectionContainer(env Environment) *Container {
	client := slack.New(env.GetSlackAuthToken(), slack.OptionDebug(!env.IsProduction()))
	pickerFactory := domain.NewPickerFactory(env.GetPickingType())
	picker := pickerFactory.Build()
	slackService := services.NewSlackService(client, env.GetAllowedChannels(), env.GetBannedUsers())
	mateMentionBuilder := domain.NewMateMention()
	mateService := services.NewMateService(slackService, picker, mateMentionBuilder)
	tsCache := repositories.NewMessageTimestampCache()
	return &Container{
		SlackClient:    client,
		SlackService:   slackService,
		Picker:         picker,
		MateService:    mateService,
		TimestampCache: tsCache,
	}
}
