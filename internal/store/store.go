package store

import (
	"sync"
)

type FileStore struct {
	store map[string]string
	mutex *sync.RWMutex
}

func NewFileStore() *FileStore {
	return &FileStore{
		store: make(map[string]string),
		mutex: &sync.RWMutex{},
	}
}

func (fs *FileStore) Add(id, filePath string) {
	fs.mutex.Lock()
	defer fs.mutex.Unlock()
	fs.store[id] = filePath
}

func (fs *FileStore) Get(id string) (string, bool) {
	fs.mutex.RLock()
	defer fs.mutex.RUnlock()
	filePath, ok := fs.store[id]
	return filePath, ok
}

