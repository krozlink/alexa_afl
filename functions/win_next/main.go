package main

import (
	"context"
	"errors"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	alexa "github.com/ericdaugherty/alexa-skills-kit-golang"
)

var a = &alexa.Alexa{ApplicationID: "amzn1.ask.skill.f6128da6-7813-4025-a762-b78be8fc6863", RequestHandler: &Handler{}, IgnoreApplicationID: true, IgnoreTimestamp: true}

const cardTitle = "AFL Odds"

type Handler struct{}

// Handle processes calls from Lambda
func Handle(ctx context.Context, requestEnv *alexa.RequestEnvelope) (interface{}, error) {
	return a.ProcessRequest(ctx, requestEnv)
}

// OnSessionStarted called when a new session is created.
func (h *Handler) OnSessionStarted(ctx context.Context, request *alexa.Request, session *alexa.Session, ctxPtr *alexa.Context, response *alexa.Response) error {
	log.Printf("OnSessionStarted requestId=%s, sessionId=%s", request.RequestID, session.SessionID)
	return nil
}

// OnLaunch called with a reqeust is received of type LaunchRequest
func (h *Handler) OnLaunch(ctx context.Context, request *alexa.Request, session *alexa.Session, ctxPtr *alexa.Context, response *alexa.Response) error {
	// speechText := "Welcome to AFL odds"

	log.Printf("OnLaunch requestId=%s, sessionId=%s", request.RequestID, session.SessionID)

	// response.SetSimpleCard(cardTitle, speechText)
	// response.SetOutputText(speechText)
	// response.SetRepromptText(speechText)

	// response.ShouldSessionEnd = true

	return nil
}

// OnIntent called with a reqeust is received of type IntentRequest
func (h *Handler) OnIntent(ctx context.Context, request *alexa.Request, session *alexa.Session, ctxPtr *alexa.Context, response *alexa.Response) error {
	log.Printf("OnIntent requestId=%s, sessionId=%s, intent=%s", request.RequestID, session.SessionID, request.Intent.Name)

	switch request.Intent.Name {
	case "NEXT_WIN_CHANCE":
		log.Println("Next win chance intent triggered")
		speechText := "Hello World"

		teamSlot := request.Intent.Slots["TEAM"]
		log.Printf(teamSlot.ID)
		log.Printf(teamSlot.Name)
		log.Printf(teamSlot.Value)

		response.SetSimpleCard(cardTitle, speechText)
		response.SetOutputText(speechText)

		log.Printf("Set Output speech, value now: %s", response.OutputSpeech.Text)
	default:
		return errors.New("Invalid Intent")
	}

	return nil
}

// OnSessionEnded called with a reqeust is received of type SessionEndedRequest
func (h *Handler) OnSessionEnded(ctx context.Context, request *alexa.Request, session *alexa.Session, ctxPtr *alexa.Context, response *alexa.Response) error {
	log.Printf("OnSessionEnded requestId=%s, sessionId=%s", request.RequestID, session.SessionID)
	return nil
}

func main() {
	lambda.Start(Handle)
}
