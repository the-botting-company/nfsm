package nfsm

import "sync"

type Metadata struct {
	metadata map[string]any

	metadataMu sync.RWMutex
}

func NewMetadata() *Metadata {
	return &Metadata{
		metadata: make(map[string]any),
	}
}

func (m *Metadata) GetAll() map[string]any {
	m.metadataMu.RLock()
	defer m.metadataMu.RUnlock()

	return m.metadata
}

func (m *Metadata) Get(k string) any {
	m.metadataMu.RLock()
	defer m.metadataMu.RUnlock()

	return m.metadata[k]
}

func (m *Metadata) Set(k string, v any) {
	m.metadataMu.Lock()
	defer m.metadataMu.Unlock()

	m.metadata[k] = v
}
