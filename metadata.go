package nfsm

import "sync"

type MetadataImpl struct {
	metadata map[string]any

	metadataMu sync.RWMutex
}

func NewMetadata() *MetadataImpl {
	return &MetadataImpl{}
}

func (m *MetadataImpl) GetAll() map[string]any {
	m.metadataMu.RLock()
	defer m.metadataMu.RUnlock()

	return m.metadata
}

func (m *MetadataImpl) Get(k string) any {
	m.metadataMu.RLock()
	defer m.metadataMu.RUnlock()

	return m.metadata[k]
}

func (m *MetadataImpl) Set(k string, v any) {
	m.metadataMu.Lock()
	defer m.metadataMu.Unlock()

	m.metadata[k] = v
}
