package main

import (
	"log"
	"os"

	"github.com/kelseyhightower/envconfig"
	"github.com/slack-go/slack"
)

// EnvConfig https://api.slack.com/internal-integrations
type EnvConfig struct {
	// AppLevelToken is app-level-token to run socketmode.
	AppLevelToken string `envconfig:"APP_LEVEL_TOKEN" required:"true"`
	// BotToken is bot user token to access to slack API.
	BotToken string `envconfig:"BOT_TOKEN" required:"true"`

	WinnerResponded string `envconfig:"WINNER_RESPONDED" default:"Thank you:muscle:"`
	OtherResponded  string `envconfig:"OTHER_RESPONDED" default:"Oh,Thank you! <@%s>:muscle:"`
	Choose          string `envconfig:"CHOOSE" default:"I choose you <@%s>!"`
	LotTitle        string `envconfig:"LOT_TITLE" default:"Turn the Gacha!"`
}

// MessageTemplate Template messages bot speak
type MessageTemplate struct {
	WinnerResponded string
	OtherResponded  string
	Choose          string
	LotTitle        string
}

// Member to assign task
type Member struct {
	ID   string
	Name string
}

// MemberList to assign task
type MemberList struct {
	members []Member
}

func main() {
	os.Exit(_main(os.Args[1:]))
}

func _main(args []string) int {
	var env EnvConfig
	if err := envconfig.Process("", &env); err != nil {
		log.Printf("[ERROR] Failed to process env var: %s", err)
		return 1
	}
	messageTemplate := MessageTemplate{
		WinnerResponded: env.WinnerResponded,
		OtherResponded:  env.OtherResponded,
		Choose:          env.Choose,
		LotTitle:        env.LotTitle,
	}

	// Listening slack event and response
	log.Printf("[INFO] Start slack event listening")
	client := slack.New(
		env.BotToken,
		slack.OptionAppLevelToken(env.AppLevelToken),
		// slack.OptionDebug(true),
		// slack.OptionLog(log.New(os.Stdout, "api: ", log.Lshortfile|log.LstdFlags)),
	)

	lot := &Lot{client: client, messageTemplate: messageTemplate}
	memberCollector := &MemberCollector{client: client}

	interactionHandler := interactionHandler{
		slackClient:     client,
		lot:             lot,
		memberCollector: memberCollector,
		env:             env,
		messageTemplate: messageTemplate,
	}

	eventHandler := eventHandler{
		lot:             lot,
		memberCollector: memberCollector,
	}

	if err := runSocketMode(client, interactionHandler, eventHandler, env); err != nil {
		log.Printf("[ERROR failed run socketmode] %s", err)
		return 1
	}

	return 0
}
