package store

import (
	"github.com/syndtr/goleveldb/leveldb/comparer"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/memdb"
	"sync"
)

type MemoryStore struct {
	mu sync.RWMutex
	db *memdb.DB
}

/**
* Функция NewMemoryStore инициализирует новое хранилище в памяти.
*
* @return ms *MemoryStore
* @return err error
 */
func NewMemoryStore() *MemoryStore {
	ms := &MemoryStore{
		db: memdb.New(comparer.DefaultComparer, 128),
	}

	return ms
}

/**
* Функция добавляет ключ и значение в память.
*
* @param key []byte ключ (байт-массив)
* @param value []byte значение (байт-массив)
* @return error
 */
func (ms *MemoryStore) Put(key, value []byte) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	return ms.db.Put(key, value)
}

/**
* Функция получает блок из в памяти.
*
* @param key []byte ключ (байт-массив)
* @return []byte
* @return error
 */
func (ms *MemoryStore) Get(key []byte) ([]byte, error) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	return ms.db.Get(key)
}

/**
* Функция удаляет блок из памяти.
*
* @param key []byte ключ (байт-массив)
* @return error
 */
func (ms *MemoryStore) Delete(key []byte) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	return ms.db.Delete(key)
}

/**
* Функция возвращает новый iterator для memdb.
*
* @return iterator.Iterator
 */
func (ms *MemoryStore) Iterator() iterator.Iterator {
	return ms.db.NewIterator(nil)
}
