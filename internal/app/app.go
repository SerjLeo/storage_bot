package app

import (
	"flag"
	"github.com/SerjLeo/storage_bot/internal/clients/telegram"
	tgconsumer "github.com/SerjLeo/storage_bot/internal/consumer/telegram"
	tgevents "github.com/SerjLeo/storage_bot/internal/events/telegram"
	"github.com/SerjLeo/storage_bot/internal/storage/sqlite"
	"github.com/pkg/errors"
	"log"
)

const maxBatchSize = 100

type Config struct {
	Host  string
	Token string
	Mode  string
}

func Run() error {
	config := mustConfig()
	client := telegram.New(config.Host, config.Token)
	storage, err := sqlite.New("./data/sqlite/storage.db")
	if err != nil {
		return errors.Wrap(err, "initializing storage")
	}
	eventProcessor := tgevents.New(&client, storage, config.Mode)
	consumer := tgconsumer.New(eventProcessor, eventProcessor, maxBatchSize)
	log.Printf("The bot is running!")
	return consumer.Start()
}

func mustConfig() Config {
	t := flag.String("token", "", "telegram token for runtime")
	h := flag.String("host", "api.telegram.org", "host for telegram api")
	m := flag.String("mode", "keeper", "mode for bot: scavenger (delete picked links) and keeper (default)")
	flag.Parse()
	return Config{
		Host:  *h,
		Token: *t,
		Mode:  *m,
	}
}
