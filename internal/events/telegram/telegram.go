package telegram

import (
	"context"
	"github.com/SerjLeo/storage_bot/internal/clients/telegram"
	"github.com/SerjLeo/storage_bot/internal/events"
	"github.com/SerjLeo/storage_bot/internal/storage"
	"github.com/pkg/errors"
)

var UnknownEventError = errors.New("unknown event occured")

type Meta struct {
	ChatId          int
	Username        string
	Caption         string
	Entities        []telegram.MessageEntity
	CaptionEntities []telegram.MessageEntity
}

type EventProcessor struct {
	client  *telegram.Client
	storage storage.Storage
	offset  int
}

func New(client *telegram.Client, storage storage.Storage) *EventProcessor {
	return &EventProcessor{
		client:  client,
		storage: storage,
		offset:  0,
	}
}

func (p *EventProcessor) Fetch(limit int) ([]events.Event, error) {
	updates, err := p.client.Updates(p.offset, limit)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	if len(updates) == 0 {
		return nil, nil
	}
	p.offset = updates[len(updates)-1].Id + 1
	res := make([]events.Event, 0, len(updates))
	for _, upd := range updates {
		res = append(res, buildEvent(upd))
	}
	return res, nil
}

func (p *EventProcessor) Process(ctx context.Context, e events.Event) error {

	switch e.Type {
	case events.Message:
		return p.processMessage(ctx, e)
	default:
		return UnknownEventError
	}
}

func (p *EventProcessor) processMessage(ctx context.Context, e events.Event) error {
	meta, ok := e.Meta.(Meta)
	if !ok {
		return errors.New("can't type cast event metadata")
	}

	return p.doCommand(ctx, e.Text, meta)
}

func buildEvent(update telegram.Update) events.Event {
	t := buildType(update)
	event := events.Event{
		Type: t,
		Text: buildText(update),
	}
	if t == events.Message {
		event.Meta = buildMeta(update)
	}

	return event
}

func buildType(update telegram.Update) events.Type {
	if update.Message == nil {
		return events.Unknown
	}
	return events.Message
}

func buildText(update telegram.Update) string {
	if update.Message == nil {
		return ""
	}
	return update.Message.Text
}

func buildMeta(update telegram.Update) Meta {
	return Meta{
		ChatId:          update.Message.Chat.Id,
		Username:        update.Message.From.Name,
		Caption:         update.Message.Caption,
		Entities:        update.Message.Entities,
		CaptionEntities: update.Message.CaptionEntities,
	}
}
