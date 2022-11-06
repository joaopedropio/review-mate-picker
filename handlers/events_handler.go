package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/joaopedropio/review-mate-picker/repositories"
	"github.com/joaopedropio/review-mate-picker/services"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"io/ioutil"
	"net/http"
)

type eventsHandler struct {
	slackClient   *slack.Client
	signingSecret string
	mateService   services.MateService
	tsCache       *repositories.MessageTimestampCache
}

func NewEventsHandler(slackClient *slack.Client, signingSecret string, mateService services.MateService, tsCache *repositories.MessageTimestampCache) Handler {
	return &eventsHandler{
		slackClient:   slackClient,
		signingSecret: signingSecret,
		mateService:   mateService,
		tsCache:       tsCache,
	}
}

func (h *eventsHandler) Handle(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		h.ErrorResponse(w, http.StatusBadRequest, fmt.Errorf("unable to read body: %w", err))
		return
	}
	if err = h.validateRequest(body, r); err != nil {
		h.ErrorResponse(w, http.StatusUnauthorized, fmt.Errorf("invalid request: %w", err))
		return
	}
	challenge, err := h.handleEvent(body)
	if err != nil {
		h.ErrorResponse(w, http.StatusInternalServerError, fmt.Errorf("unable to handle event: %w", err))
		return
	}
	if challenge != nil {
		w.Header().Set("Content-Type", "text")
		w.Write(challenge)
		return
	}
}

func (h *eventsHandler) ErrorResponse(w http.ResponseWriter, httpStatus int, err error) {
	w.WriteHeader(httpStatus)
	_, _ = w.Write([]byte(err.Error()))
}

func (h *eventsHandler) handleEvent(body []byte) ([]byte, error) {
	event, err := slackevents.ParseEvent(body, slackevents.OptionNoVerifyToken())
	if err != nil {
		return nil, fmt.Errorf("unable to parse event: %w", err)
	}
	if event.Type == slackevents.URLVerification {
		challenge, err := h.verifyURL(body)
		if err != nil {
			return nil, fmt.Errorf("unable to verify url: %w", err)
		}
		return challenge, nil
	}
	if event.Type == slackevents.CallbackEvent {
		innerEvent := event.InnerEvent
		switch ev := innerEvent.Data.(type) {
		case *slackevents.LinkSharedEvent:
			// verify if the message has already been replied
			if h.tsCache.IsSet(ev.MessageTimeStamp) {
				return nil, nil
			}
			h.tsCache.Set(ev.MessageTimeStamp)

			if err := h.mateService.PickMateToReview(ev.Channel, ev.User, ev.MessageTimeStamp); err != nil {
				return nil, fmt.Errorf("unable to pick mate to review: %w", err)
			}
		}
	}
	return nil, nil
}

func (h *eventsHandler) verifyURL(body []byte) (challenge []byte, err error) {
	var r *slackevents.ChallengeResponse
	err = json.Unmarshal(body, &r)
	if err != nil {
		return nil, fmt.Errorf("unable to parse json of type ChallengeResponse: body: %s", string(body))
	}
	return []byte(r.Challenge), nil
}

func (h *eventsHandler) validateRequest(body []byte, r *http.Request) error {
	sv, err := slack.NewSecretsVerifier(r.Header, h.signingSecret)
	if err != nil {
		return fmt.Errorf("unable to create secrets verifier: %w", err)
	}
	if _, err := sv.Write(body); err != nil {
		return fmt.Errorf("unable to write body to secrets verifier: %w", err)
	}
	if err := sv.Ensure(); err != nil {
		return fmt.Errorf("unable to ensure this request is valid: %w", err)
	}
	return nil
}
