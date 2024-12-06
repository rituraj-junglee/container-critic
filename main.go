package main

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rituraj-junglee/container-critic/repo/reportmeta"
	"github.com/rituraj-junglee/container-critic/services/report"
	"github.com/rituraj-junglee/container-critic/services/slacker"
	"github.com/rituraj-junglee/container-critic/services/trigger"
	triggersvc "github.com/rituraj-junglee/container-critic/services/trigger/service"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
)

func init() {
	godotenv.Load(".env")
}

const (
	SLACK_CHANNEL_NAME = "test-container-critic"
	SLACK_CHANNEL_ID   = "C083PH81V38"
)

func main() {
	appToken := os.Getenv("SOCKET_APP_TOKEN")
	botToken := os.Getenv("BOT_AUTH_TOKEN")

	opts := []slack.Option{
		slack.OptionDebug(true),
		slack.OptionAppLevelToken(appToken),
	}
	slackClient := slack.New(botToken, opts...)

	app := socketmode.New(slackClient, socketmode.OptionDebug(true))

	httpRouter := mux.NewRouter().StrictSlash(false)

	var slackservice slacker.Service
	{
		slackservice = slacker.NewService(app, SLACK_CHANNEL_ID)
	}

	var reportmetarepo reportmeta.Repository
	{
		reportmetarepo = reportmeta.NewRepository()
	}
	var reportservice report.Service
	{
		reportservice = report.NewService(reportmetarepo)
	}

	// triggerService := trigger.NewService(app, )
	var triggerservice trigger.Service
	{
		triggerservice = trigger.NewService(app, reportservice, slackservice)
		handler := triggersvc.MakeHTTPHandler(triggerservice, nil, nil)
		httpRouter.PathPrefix("/slack/trigger").Handler(handler)
	}

	go slackservice.ReadMessages()
	defer slackservice.CloseReader()

	// Start APP
	go app.Run()

	// Start HTTP Server
	http.ListenAndServe(":8000", httpRouter)

}
