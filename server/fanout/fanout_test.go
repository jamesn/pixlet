package fanout

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewFanout(t *testing.T) {
	fo := NewFanout()
	assert.NotNil(t, fo)
	assert.NotNil(t, fo.broadcast)
	assert.NotNil(t, fo.register)
	assert.NotNil(t, fo.unregister)
	assert.NotNil(t, fo.quit)
}

func TestWebsocketEventStructure(t *testing.T) {
	// Test that WebsocketEvent has the correct fields
	event := WebsocketEvent{
		Type:      EventTypeImage,
		Message:   "test message",
		ImageType: "webp",
	}

	assert.Equal(t, EventTypeImage, event.Type)
	assert.Equal(t, "test message", event.Message)
	assert.Equal(t, "webp", event.ImageType)
}

func TestEventTypeConstants(t *testing.T) {
	// Verify the event type constants are defined
	assert.Equal(t, "img", EventTypeImage)
	assert.Equal(t, "schema", EventTypeSchema)
	assert.Equal(t, "error", EventTypeErr)
}

func TestClientSendChannel(t *testing.T) {
	// Test that we can send events through a client's channel
	client := &Client{
		send: make(chan WebsocketEvent, 10),
		quit: make(chan bool, 1),
	}

	event := WebsocketEvent{
		Type:    EventTypeImage,
		Message: "test message",
	}

	// Send event
	client.send <- event

	// Verify event was sent
	received := <-client.send
	assert.Equal(t, event.Type, received.Type)
	assert.Equal(t, event.Message, received.Message)
}
