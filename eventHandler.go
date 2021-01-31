package main

import (
	"log"

	"github.com/slack-go/slack/slackevents"
)

// eventHandler handles interactive message response.
type eventHandler struct {
	lot             *Lot
	memberCollector *MemberCollector
}

func (h eventHandler) HandleEvent(eventPayload slackevents.EventsAPIEvent, selfUserID string) {
	switch eventPayload.Type {
	case slackevents.CallbackEvent:
		switch event := eventPayload.InnerEvent.Data.(type) {
		case *slackevents.AppMentionEvent:
			if event.User == selfUserID {
				return
			}
			members, err := h.memberCollector.Collect(event.Channel)
			if err != nil {
				log.Printf("[ERROR] Failed to collect member request: %s", err)
			}
			h.lot.DrawLots(event.Channel, members, "")
		}
	}
}
