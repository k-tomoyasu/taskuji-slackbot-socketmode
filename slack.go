package main

import (
	"fmt"
	"os"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

const (
	// action is used for slack attament action.
	actionAccept = "appcept"
	actionRepeat = "repeat"
)

func runSocketMode(client *slack.Client, interactionHandler interactionHandler, eventHandler eventHandler, env EnvConfig) error {
	//https://github.com/slack-go/slack/blob/master/examples/socketmode/socketmode.go
	socketMode := socketmode.New(
		client,
		// socketmode.OptionDebug(true),
		// socketmode.OptionLog(log.New(os.Stdout, "sm: ", log.Lshortfile|log.LstdFlags)),
	)
	authTest, authTestErr := client.AuthTest()
	if authTestErr != nil {
		fmt.Fprintf(os.Stderr, "SLACK_BOT_TOKEN is invalid: %v\n", authTestErr)
		os.Exit(1)
	}
	selfUserID := authTest.UserID

	go func() {
		for envelope := range socketMode.Events {
			switch envelope.Type {
			case socketmode.EventTypeEventsAPI:
				socketMode.Ack(*envelope.Request)
				eventPayload, _ := envelope.Data.(slackevents.EventsAPIEvent)
				eventHandler.HandleEvent(eventPayload, selfUserID)
			case socketmode.EventTypeInteractive:
				// See https://api.slack.com/apis/connections/socket-implement#button
				socketMode.Ack(*envelope.Request)
				payload, _ := envelope.Data.(slack.InteractionCallback)
				interactionHandler.HandleInteraction(payload)
			}
		}
	}()

	return socketMode.Run()
}
