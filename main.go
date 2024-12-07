package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rituraj-junglee/container-critic/repo/reportconfig"
	reportconfigmongo "github.com/rituraj-junglee/container-critic/repo/reportconfig/mongo"
	"github.com/rituraj-junglee/container-critic/repo/reportmeta"
	"github.com/rituraj-junglee/container-critic/services/report"
	"github.com/rituraj-junglee/container-critic/services/slacker"
	"github.com/rituraj-junglee/container-critic/services/trigger"
	triggersvc "github.com/rituraj-junglee/container-critic/services/trigger/service"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	godotenv.Load(".env")
}

const (
	SLACK_CHANNEL_NAME = "test-container-critic"
	SLACK_CHANNEL_ID   = "C083PH81V38"
	TIMELAPSE_ENABLED  = true
	MONGO_DB           = "container-critic"
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

	var mongoClient *mongo.Client
	{
		opts := options.Client().ApplyURI(os.Getenv("MONGO_URI"))

		mongoCtx, mongoCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer mongoCancel()

		var err error
		mongoClient, err = mongo.Connect(mongoCtx, opts)
		if err != nil {
			panic(err)
		}
	}

	var slackservice slacker.Service
	{
		slackservice = slacker.NewService(app, SLACK_CHANNEL_ID)
	}

	var reportmetarepo reportmeta.Repository
	{
		reportmetarepo = reportmeta.NewRepository()
	}
	var reportconfigrepo reportconfig.Repository
	{
		reportconfigrepo = reportconfigmongo.NewRepository(mongoClient, MONGO_DB)
	}
	var reportservice report.Service
	{
		reportservice = report.NewService(TIMELAPSE_ENABLED, reportmetarepo, reportconfigrepo)
	}

	// triggerService := trigger.NewService(app, )
	var triggerservice trigger.Service
	{
		triggerservice = trigger.NewService(app, reportservice, slackservice, reportconfigrepo)
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
