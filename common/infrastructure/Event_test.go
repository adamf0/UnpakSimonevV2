package infrastructure_test

import (
	"context"
	"UnpakSiamida/common/infrastructure"
	"testing"
	"time"

	"github.com/mehdihadeli/go-mediatr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type MyTestEvent struct {
	EventID   string
	Timestamp time.Time
}

func (e MyTestEvent) ID() string { return e.EventID }
func (e MyTestEvent) OccurredOnUTC() time.Time { return e.Timestamp }

type MyTestEventHandler struct {
	Called bool
	Event  MyTestEvent
}

func (h *MyTestEventHandler) Handle(ctx context.Context, event MyTestEvent) error {
	h.Called = true
	h.Event = event
	return nil
}

func TestCanonicalTypeName(t *testing.T) {
	evt := MyTestEvent{}
	name := infrastructure.CanonicalTypeName(evt)
	assert.Contains(t, name, "infrastructure_test.MyTestEvent")

	namePtr := infrastructure.CanonicalTypeName(&evt)
	assert.Contains(t, namePtr, "infrastructure_test.MyTestEvent")
}

func TestEventRegistryAndDispatcher(t *testing.T) {
	// Register the handler in go-mediatr
	handler := &MyTestEventHandler{}
	errReg := mediatr.RegisterNotificationHandler[MyTestEvent](handler)
	require.NoError(t, errReg)

	// Register event in our registry
	infrastructure.RegisterDomainEvent(&MyTestEvent{})

	dispatcher := infrastructure.NewEventDispatcher()
	infrastructure.RegisterEvent[MyTestEvent](dispatcher)

	// Dispatch
	event := MyTestEvent{EventID: "123", Timestamp: time.Now()}
	err := dispatcher.Dispatch(context.Background(), event)
	require.NoError(t, err)

	assert.True(t, handler.Called)
	assert.Equal(t, "123", handler.Event.EventID)
}
