package domain_test

import (
	"UnpakSiamida/common/domain"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type dummyEvent struct {
	id       string
	occurred time.Time
}

func (d dummyEvent) ID() string {
	return d.id
}

func (d dummyEvent) OccurredOnUTC() time.Time {
	return d.occurred
}

func TestEntity_DomainEvents(t *testing.T) {
	entity := domain.NewEntity()
	assert.Empty(t, entity.DomainEvents())

	evt1 := dummyEvent{id: "evt-1", occurred: time.Now()}
	evt2 := dummyEvent{id: "evt-2", occurred: time.Now()}

	entity.Raise(evt1)
	entity.Raise(evt2)

	events := entity.DomainEvents()
	assert.Len(t, events, 2)
	assert.Equal(t, "evt-1", events[0].ID())

	entity.ClearDomainEvents()
	assert.Empty(t, entity.DomainEvents())
}
