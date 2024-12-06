package slacker

import (
	"context"
	"log"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
)

type service struct {
	socketmodeClient *socketmode.Client
	channelID        string
}

func NewService(socketmodeClient *socketmode.Client, channelID string) Service {
	return &service{
		socketmodeClient: socketmodeClient,
		channelID:        channelID,
	}
}

func (s *service) ReadMessages() {
	for evt := range s.socketmodeClient.Events {
		log.Println("Event Received", evt.Type)
		switch evt.Type {
		case socketmode.EventTypeInteractive:
			callback, ok := evt.Data.(slack.InteractionCallback)
			if !ok {
				return
			}
			s.socketmodeClient.Ack(*evt.Request)
			s.socketmodeClient.SendMessage(callback.Channel.ID, slack.MsgOptionText("Received", false))
			// case socketmode.EventTypeEventsAPI:
			// 	s.socketmodeClient.SendMessage("test-container-critic", slack.MsgOptionText("Received", false))

		}
	}
}

func (s *service) SendMessage(ctx context.Context, message interface{}) (err error) {
	opts := []slack.MsgOption{
		slack.MsgOptionText(message.(string), false),
	}
	_, _, _, err = s.socketmodeClient.SendMessage(s.channelID, opts...)
	return
}

func (s *service) SendFile(ctx context.Context, filename string, content string) (err error) {

	file := slack.UploadFileV2Parameters{
		File:     filename,
		Channel:  s.channelID,
		Filename: filename,
		Title:    filename,
		FileSize: len(content),
		Content:  content,
	}

	summary, err := s.socketmodeClient.UploadFileV2(file)
	if err != nil {
		log.Println("Error uploading file: ", err)
		log.Println("Summary: ", summary)
	}

	return
}

func (s *service) CloseReader() {
	s.socketmodeClient.CloseConversation(s.channelID)
}
