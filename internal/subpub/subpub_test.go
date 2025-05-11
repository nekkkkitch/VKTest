package subpub

import (
	"VKTest/pkg/pubsub"
	"context"
	"testing"
)

var h pubsub.SubPub
var sub pubsub.Subscription
var err error

func TestCreation(t *testing.T) {
	h = NewSubPub()
	t.Log("Good hub:", h)
}

func TestSubscribe(t *testing.T) {
	sub, err = h.Subscribe("cool_topic", func(msg interface{}) { return })
	if err != nil {
		t.Error(err)
	}
	t.Log("Good sub:", sub)
}

func TestPublish(t *testing.T) {
	err := h.Publish("cool_topic", "ababa")
	if err != nil {
		t.Error(err)
	}
	result := <-sub.(*subscriber).Messages
	t.Log("Result:", result)
}

func TestUsubbing(t *testing.T) {
	sub.Unsubscribe()
	t.Log("Result subscribers after unsubbing:", h.(*hub).Subjects)
	t.Log("Works")
}

func TestClosingHub(t *testing.T) {
	h.Close(context.Background())
	t.Log("Closed")
}

func TestSubbingAfterClosed(t *testing.T) {
	sub, err = h.Subscribe("cool_topic", func(msg interface{}) {})
	if err == nil {
		t.Error("Error is nil(wrong)")
	}
	t.Log("Got error: " + err.Error())
}

func TestPublishingAfterClosed(t *testing.T) {
	err = h.Publish("cool_topic", "abab")
	if err == nil {
		t.Error("Error is nil(wrong)")
	}
	t.Log("Got error: " + err.Error())
}
