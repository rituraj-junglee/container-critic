package slacker

import "context"

type Service interface {
	ReadMessages()
	CloseReader()
	SendMessage(ctx context.Context, message interface{}) (err error)
	SendFile(ctx context.Context, filename string, content string) (err error)
}
