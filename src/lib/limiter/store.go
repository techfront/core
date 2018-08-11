package limiter

import (
	"encoding/binary"
	"github.com/techfront/core/src/kernel/store"
	"log"
	"time"
)

type Store struct {
	db *store.LocalStore
}

/**
* Функция инициализирует новое локальное хранилище.
*
* @param path string путь к локальному хранилищу.
*
* @return *Store
* @return error
 */
func NewStore(path string) (*Store, error) {

	db, err := store.NewLocalStore(path, true, false)
	if err != nil {
		log.Fatal(err)
	}

	s := &Store{
		db: db,
	}

	return s, nil
}

// GetWithTime returns the value of the key if it is in the store
// or -1 if it does not exist. It also returns the current time at
// the redis server to microsecond precision.
func (s *Store) GetWithTime(key string) (int64, time.Time, error) {
	now := time.Now()

	valP, err := s.db.Get([]byte(key))
	if err != nil {
		return -1, now, nil
	}

	return int64(binary.LittleEndian.Uint64(valP)), now, nil
}

// SetIfNotExistsWithTTL sets the value of key only if it is not
// already set in the store it returns whether a new value was set.
// If a new value was set, the ttl in the key is also set, though this
// operation is not performed atomically.
func (s *Store) SetIfNotExistsWithTTL(key string, value int64, _ time.Duration) (bool, error) {

	// Проверяем существует ли блок
	ok, _ := s.db.Has([]byte(key))
	if ok {
		return false, nil
	}

	// Сериализация Int64
	bsValue := make([]byte, 8)
	binary.LittleEndian.PutUint64(bsValue, uint64(value))

	// Запись в БД
	err := s.db.Put([]byte(key), bsValue)
	if err != nil {
		return false, err
	}

	return true, nil
}

// CompareAndSwapWithTTL atomically compares the value at key to the
// old value. If it matches, it sets it to the new value and returns
// true. Otherwise, it returns false. If the key does not exist in the
// store, it returns false with no error. If the swap succeeds, the
// ttl for the key is updated atomically.
func (s *Store) CompareAndSwapWithTTL(key string, old, new int64, _ time.Duration) (bool, error) {
	valP, err := s.db.Get([]byte(key))
	if err != nil {
		return false, nil
	}

	if int64(binary.LittleEndian.Uint64(valP)) == old {

		bsNew := make([]byte, 8)
		binary.LittleEndian.PutUint64(bsNew, uint64(new))

		err := s.db.Put([]byte(key), bsNew)
		if err != nil {
			return false, err
		}

		return true, nil
	}

	return false, nil
}
