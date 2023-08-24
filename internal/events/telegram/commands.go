package telegram

import (
	"context"
	"github.com/SerjLeo/storage_bot/internal/models"
	"github.com/pkg/errors"
	"log"
	"net/url"
	"strings"
)

const (
	RndCmd          = "/rnd"
	HelpCmd         = "/help"
	StartCmd        = "/start"
	MsgCmdNotFound  = "Sorry, unknown command. To see the list of commands use /help command."
	MsgAlreadyExist = "Link is already stored"
	MsgIsNotUrl     = "Provided link is not url"
	MsgEmptyUrl     = "Please provide url as command argument"
	MsgSuccess      = "Link successfully saved"
	MsgNotFound     = "You dont have any stored links"
	MsgHelp         = `
	/rnd - pick random saved url
	/add url - add url to storage
	/help - list of commands
`
	MsgHello       = "Hi there! \n\n" + MsgHelp
	EntityTypeLink = "text_link"
)

func (p *EventProcessor) doCommand(ctx context.Context, text string, meta Meta) error {
	username, chatId := meta.Username, meta.ChatId
	log.Printf("got %s command from %s", text, username)
	values := strings.Split(strings.TrimSpace(text), " ")
	cmd := values[0]
	links := mustParseLinks(text, meta)

	if len(links) > 0 {
		return p.addCommand(ctx, links, chatId, username)
	}
	var err error
	switch cmd {
	case HelpCmd:
		err = p.client.SendMessage(chatId, MsgHelp)
	case StartCmd:
		err = p.client.SendMessage(chatId, MsgHello)
	case RndCmd:
		err = p.rndCommand(ctx, chatId, username)
	default:
		err = p.client.SendMessage(chatId, MsgCmdNotFound)
	}
	return err
}

func (p *EventProcessor) addCommand(ctx context.Context, links []string, chatId int, username string) error {
	for _, link := range links {
		page := &models.Page{
			URL:      link,
			UserName: username,
		}
		isExist, err := p.storage.IsExist(ctx, page)
		if err != nil {
			return errors.Wrap(err, "add command")
		}
		if isExist {
			continue
		}
		err = p.storage.Save(ctx, page)
		if err != nil {
			return errors.Wrap(err, "add command")
		}
	}
	return p.client.SendMessage(chatId, MsgSuccess)
}

func (p *EventProcessor) rndCommand(ctx context.Context, chatId int, username string) error {
	page, err := p.storage.Pick(ctx, username)
	if err != nil {
		return errors.Wrap(err, "rnd command")
	}
	if page == nil {
		return p.client.SendMessage(chatId, MsgNotFound)
	}
	err = p.storage.Remove(ctx, page)
	if err != nil {
		return errors.Wrap(err, "rnd command")
	}
	return p.client.SendMessage(chatId, page.URL)
}

func isUrl(text string) bool {
	u, err := url.Parse(text)
	return !(err != nil || u.Host == "")
}

func mustParseLinks(text string, meta Meta) []string {
	links := make([]string, 0, 10)
	words := strings.Split(text+" "+meta.Caption, " ")
	for _, word := range words {
		if isUrl(word) {
			links = append(links, word)
		}
	}
	for _, entity := range meta.Entities {
		if entity.Type == EntityTypeLink {
			links = append(links, entity.Url)
		}
	}
	for _, entity := range meta.CaptionEntities {
		if entity.Type == EntityTypeLink {
			links = append(links, entity.Url)
		}
	}
	return links
}
