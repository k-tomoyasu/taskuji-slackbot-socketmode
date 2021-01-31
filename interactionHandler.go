package main

import (
	"fmt"
	"log"

	"github.com/slack-go/slack"
)

// interactionHandler handles interactive message response.
type interactionHandler struct {
	slackClient     *slack.Client
	lot             *Lot
	memberCollector *MemberCollector
	env             EnvConfig
	messageTemplate MessageTemplate
}

func (h interactionHandler) HandleInteraction(interactionRequest slack.InteractionCallback) {
	action := interactionRequest.ActionCallback.BlockActions[0]
	h.reply(action, interactionRequest)
}

func (h interactionHandler) reply(action *slack.BlockAction, interactionRequest slack.InteractionCallback) {
	switch action.ActionID {
	case actionAccept:
		winnerResponsed := fmt.Sprintf("<@%s>", interactionRequest.User.ID) == interactionRequest.OriginalMessage.Text
		var value string
		if winnerResponsed {
			value = h.messageTemplate.WinnerResponded
		} else {
			value = fmt.Sprintf(h.messageTemplate.OtherResponded, interactionRequest.User.ID)
		}
		message := slack.MsgOptionText(value, false)
		if _, _, err := h.slackClient.PostMessage(interactionRequest.Channel.ID, message); err != nil {
			log.Printf("failed to post message: %s", err)
		}
		return
	case actionRepeat:
		userGroupID := action.Value
		var members []Member
		if len(userGroupID) > 0 {
			members, _ = h.memberCollector.CollectByUserGroup(userGroupID, interactionRequest.Channel.ID)
		} else {
			members, _ = h.memberCollector.Collect(interactionRequest.Channel.ID)
		}
		targetMembers := make([]Member, 0)
		// exclude member NG button pushed.
		for _, member := range members {
			if member.ID != interactionRequest.User.ID {
				targetMembers = append(targetMembers, member)
			}
		}
		h.lot.DrawLots(interactionRequest.Channel.ID, targetMembers, userGroupID)
		return
	default:
		log.Printf("Invalid action was submitted: %s", action.ActionID)
		return
	}
}
