package telegram

import (
	"context"
	"github.com/SerjLeo/storage_bot/internal/events"
	"github.com/pkg/errors"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

type Consumer struct {
	fetcher   events.Fetcher
	processor events.Processor
	batchSize int
}

func New(fetcher events.Fetcher, processor events.Processor, batchSize int) Consumer {
	return Consumer{
		fetcher:   fetcher,
		processor: processor,
		batchSize: batchSize,
	}
}

func (c *Consumer) Start() error {
	for {
		eventList, err := c.fetcher.Fetch(c.batchSize)
		if err != nil {
			log.Print(errors.Wrap(err, "[consumer] error while processing event"))
			continue
		}
		if len(eventList) == 0 {
			time.Sleep(1 * time.Second)
			continue
		}
		c.processEvents(context.Background(), eventList)
	}
}

func (c *Consumer) processEvents(ctx context.Context, eventList []events.Event) {
	wg := sync.WaitGroup{}
	var errCount int64
	start := time.Now()
	for _, e := range eventList {
		wg.Add(1)
		e := e
		go func() {
			defer wg.Done()
			err := c.processor.Process(ctx, e)
			if err != nil {
				atomic.AddInt64(&errCount, 1)
				log.Print(errors.Wrap(err, "error while processing event"))
			}
		}()
	}
	wg.Wait()
	log.Printf("Processed batch of %d events in %s, got %d errors", len(eventList), time.Since(start), errCount)
}
