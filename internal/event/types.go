package event

import "context"

type Producer interface {
	ProduceReadEvent(c context.Context, event ReadEvent) error
}
