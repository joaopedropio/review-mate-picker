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
	channelID := "channel1"
	messageTimestamp := "121212121"
	senderID := "123"
	senderName := "marcos"
	users := []services.SlackUser{
		{ID: "456", Name: "joao"},
		{ID: "789", Name: "pedro"},
		{ID: senderID, Name: senderName},
	}
	expectedMateMention := "<@marcos> can you review my pull request?"
	slackService.
		EXPECT().
		GetAllUsers().
		Return(users, nil)
	slackService.
		EXPECT().
		ReplyMessage(channelID, expectedMateMention, messageTimestamp).
		Return(nil)
	picker := domain.NewMockPicker(ctrl)
	picker.EXPECT().Pick(gomock.Any()).Return(domain.NewPerson(senderID, senderName), nil)
	mateMention := domain.NewMockMateMentionBuilder(ctrl)
	mateMention.EXPECT().Build(senderName).Return(expectedMateMention)
	mateService := services.NewMateService(slackService, picker, mateMention)

	// Act
	err := mateService.PickMateToReview(channelID, senderID, messageTimestamp)

	// Assert
	require.NoError(t, err)
}
