package deduplicator

import "sync"

type Deduplicator struct {
	seen sync.Map
}

func New() *Deduplicator {
	return &Deduplicator{}
}

func (d *Deduplicator) Seen(eventID string) bool {
	_, exists := d.seen.LoadOrStore(eventID, struct{}{})
	return exists
}
