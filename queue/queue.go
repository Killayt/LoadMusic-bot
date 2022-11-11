package queue

import (
	"context"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type HandleFunc func(ctx context.Context, url string) (string, error)

type message struct {
	url    string
	chatID int64
}

type Result struct {
	ChatID   int64
	Filename string
	Err      error
}

type DownloadQueue struct {
	queue          chan *message
	doneWg         *sync.WaitGroup
	handler        HandleFunc
	maxProcessTime int64
}

func NewDownloadQueue(h HandleFunc, maxProcessTime int64) *DownloadQueue {
	return &DownloadQueue{
		queue:          make(chan *message),
		doneWg:         new(sync.WaitGroup),
		handler:        h,
		maxProcessTime: maxProcessTime,
	}
}

func (q *DownloadQueue) Start(result chan *Result) {
	go q.startProcess(result)
}

func (q *DownloadQueue) Stop() {
	q.doneWg.Wait()
	close(q.queue)
}

func (q *DownloadQueue) Enqueue(m *tgbotapi.Message) *message {
	return &message{
		chatID: m.Chat.ID,
		url:    m.Text,
	}
}
