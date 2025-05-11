package subpub

import (
	cerr "VKTest/pkg/customErrors"
	"VKTest/pkg/pubsub"
	"context"
	"fmt"
	"log/slog"
	"sync"
)

type hub struct {
	Subjects       map[string]map[*subscriber]struct{}
	mutex          sync.RWMutex
	isClosing      bool
	isClosed       bool
	MessagesSening int
}

type subscriber struct {
	Hub      *hub
	Topic    string
	Handler  pubsub.MessageHandler
	Messages chan interface{}
}

func (h *hub) Subscribe(subject string, cb pubsub.MessageHandler) (pubsub.Subscription, error) {
	slog.Info(fmt.Sprintf("Got new subscription request on topic `%v`", subject))
	if h.isClosing {
		return nil, cerr.ErrSubClosed
	}
	sub := &subscriber{Hub: h, Handler: cb, Topic: subject}
	sub.Messages = make(chan interface{}, 1)
	h.mutex.Lock()
	if len(h.Subjects[subject]) == 0 {
		h.Subjects[subject] = make(map[*subscriber]struct{})
	}
	h.Subjects[subject][sub] = struct{}{}
	h.mutex.Unlock()
	return sub, nil
}

func (h *hub) Publish(subject string, msg interface{}) error {
	slog.Info(fmt.Sprintf("Got new publish request on topic %v", subject))
	if h.isClosing {
		return cerr.ErrSubClosed
	}
	subscribers := h.Subjects[subject]
	if len(subscribers) == 0 {
		return cerr.ErrNoTopic
	}
	h.mutex.RLock()
	var w sync.WaitGroup
	w.Add(len(subscribers))
	h.MessagesSening += len(subscribers)
	for sub := range subscribers {
		go func() {
			defer w.Done()
			sub.Handler(msg)
			for {
				select {
				case sub.Messages <- msg:
					return
				default:
					if h.isClosed {
						return
					}
				}
			}
		}()
	}
	h.mutex.RUnlock()
	w.Wait()
	h.MessagesSening -= len(subscribers)
	return nil
}

func (h *hub) Close(ctx context.Context) error {
	slog.Info("Trying to close sub")
	if h.isClosed {
		return cerr.ErrSubClosed
	}
	h.isClosing = true
	slog.Info(fmt.Sprintf("Closing sub, number of remaining messages: %v", h.MessagesSening))
	for {
		select {
		case <-ctx.Done():
			h.isClosed = true
			return nil
		default:
			if h.MessagesSening == 0 {
				h.isClosed = true
				return nil
			}
		}
	}
}

func (sub *subscriber) Unsubscribe() {
	h := sub.Hub
	h.mutex.Lock()
	close(sub.Messages)
	delete(h.Subjects[sub.Topic], sub)
	h.mutex.Unlock()
}

func (sub *subscriber) GetMessages() <-chan any {
	return sub.Messages
}

func NewSubPub() pubsub.SubPub {
	hub := hub{Subjects: make(map[string]map[*subscriber]struct{}), isClosing: false, isClosed: false}
	return &hub
}
