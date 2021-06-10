package locker

import (
	"github.com/google/uuid"
	"sync"
)

var uuids = make(map[uuid.UUID]*sync.Mutex)

func For(u uuid.UUID) *sync.Mutex {
	if uuids[u] == nil {
		uuids[u] = &sync.Mutex{}
	}

	return uuids[u]
}