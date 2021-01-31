package main

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/seehuhn/mt19937"
	"github.com/slack-go/slack"
)

// Lot decide member randomly and make slackAttachment.
type Lot struct {
	client          *slack.Client
	messageTemplate MessageTemplate
}

// DrawLots decide member randomly and send message to slack.
func (l *Lot) DrawLots(channelID string, members []Member, userGroupID string) error {
	if len(members) == 0 {
		return errors.New("no member found")
	}
	rng := rand.New(mt19937.New())
	rng.Seed(time.Now().UnixNano())
	winner := draw(members, rng.Intn)
	messages := buildLotMessage(winner, userGroupID, l.messageTemplate)
	if _, _, err := l.client.PostMessage(channelID, messages...); err != nil {
		return fmt.Errorf("failed to post message: %s", err)
	}
	return nil
}

func draw(members []Member, rngfn func(int) int) Member {
	return members[rngfn(len(members))]
}

func buildLotMessage(winner Member, userGroupID string, template MessageTemplate) []slack.MsgOption {
	var messages []slack.MsgOption
	blockID := "lot_result"
	messages = append(messages, slack.MsgOptionText(fmt.Sprintf("<@%s>", winner.ID), false))

	headerBlock := slack.NewHeaderBlock(slack.NewTextBlockObject("plain_text", template.LotTitle, false, false))
	chooseTextBlock := slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn", fmt.Sprintf(template.Choose, winner.ID), false, false), nil, nil)

	okBtnText := slack.NewTextBlockObject("plain_text", "OK!", false, false)
	okBtn := slack.NewButtonBlockElement(actionAccept, winner.ID, okBtnText)
	okBtn.WithStyle(slack.StylePrimary)
	ngBtnText := slack.NewTextBlockObject("plain_text", "NG:cry:", true, false)
	ngBtn := slack.NewButtonBlockElement(actionRepeat, userGroupID, ngBtnText)
	ngBtn.WithStyle(slack.StyleDanger)
	actionBlock := slack.NewActionBlock(blockID, okBtn, ngBtn)

	optBlocks := slack.MsgOptionBlocks(headerBlock, chooseTextBlock, actionBlock)
	messages = append(messages, optBlocks)

	return messages
}
