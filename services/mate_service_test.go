package services_test

import (
	"github.com/golang/mock/gomock"
	"github.com/joaopedropio/review-mate-picker/domain"
	"github.com/joaopedropio/review-mate-picker/services"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMateService_PickMateToReview_ShouldSucceed(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	slackService := services.NewMockSlackService(ctrl)
	users := []string{"joao", "pedro", "marcos"}
	slackService.
		EXPECT().
		GetAllUsers().
		Return(users, nil)
	slackService.
		EXPECT().
		ReplyMessage("channel1", "@marcos can you review my pull request?", "121212121").
		Return(nil)
	picker := domain.NewMockPicker(ctrl)
	picker.EXPECT().Pick(gomock.Any()).Return(domain.NewPerson("", "marcos"), nil)
	mateService := services.NewMateService(slackService, picker)

	// Act
	err := mateService.PickMateToReview("channel1", "fulano", "121212121")

	// Assert
	require.NoError(t, err)
}

func TestMateService_PickMateToReview_ShouldDoNothingWhenMessageDoesNotHaveAValidGithubPullRequestURL(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	slackService := services.NewMockSlackService(ctrl)
	slackService.
		EXPECT().
		GetAllUsers().
		MaxTimes(0)
	picker := domain.NewStatelessPicker()
	mateService := services.NewMateService(slackService, picker)

	// Act
	err := mateService.PickMateToReview("channel1", "fulano", "121212121")

	// Assert
	require.NoError(t, err)
}
