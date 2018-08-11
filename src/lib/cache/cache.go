package cache

import (
	"encoding/binary"
	"github.com/spaolacci/murmur3"
	"github.com/techfront/core/src/kernel/store"
	"gopkg.in/vmihailenco/msgpack.v2"
	"log"
	"net/http"
	"runtime"
	"strconv"
	"time"
)

var (
	cacheConfig *Cache
	defaultExpiration time.Duration
)

type Cache struct {
	LocalDB *store.LocalStore
	MemoryDB *store.MemoryStore
	janitor *janitor
}

/**
* Тип Context это интерфейс для обработки контекста.
 */
type Context interface {
	Path() string
	Request() *http.Request
}

/**
* Конфигурирование.
*
* @param config map[string]string
*/
func Setup(config map[string]string) {
	mdb := store.NewMemoryStore()

	ldb, err := store.NewLocalStore(config["cache_store"], true, false)
	if err != nil {
		log.Fatal(err)
	}

	cacheConfig = &Cache{
		LocalDB:  ldb,
		MemoryDB: mdb,
	}

	de, err := strconv.ParseInt(config["cache_default_expiration"], 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	ci, err := strconv.ParseInt(config["cache_cleanup_interval"], 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	defaultExpiration = secondsToDuration(de)
	cleanupInterval := secondsToDuration(ci)

	runJanitor(cacheConfig, cleanupInterval)
	runtime.SetFinalizer(cacheConfig, stopJanitor)
}

/**
* Функция Set записывает interface{} с ключем на заданую продолжительность.
*
* @param key string ключ.
* @param value interface{} значение.
* @param seconds int64 секунды.
*/
func Set(key string, value interface{}, seconds int64) error {

	// Определение времени истечения. По умолчанию 5 минут.
	expiration := time.Now().Add(defaultExpiration).UnixNano()

	// Но если переданы секунды, то ставим новое время.
	if seconds > 0 {
		duration := secondsToDuration(seconds)
		expiration = time.Now().Add(duration).UnixNano()
	}

	// Сериализация времени (int64).
	bsExpiration := make([]byte, 8)
	binary.LittleEndian.PutUint64(bsExpiration, uint64(expiration))

	// Сериализация значения (interface{}).
	bsValue, err := msgpack.Marshal(value)
	if err != nil {
		return err
	}

	// Записываем ключ и время в MemoryDB.
	if err := cacheConfig.MemoryDB.Put([]byte(key), bsExpiration); err != nil {
		return err
	}

	// Ключ и значение в LocalDB.
	return cacheConfig.LocalDB.Put([]byte(key), bsValue)
}

/**
* Функция Get получает значение по ключу, если время не истекло.
*
* @param key string ключ.
* @param dst interface{} указатель на память.
*/
func Get(key string, dist interface{}) error {

	bsValue, err := cacheConfig.LocalDB.Get([]byte(key))
	if err != nil {
		return err
	}

	if err := msgpack.Unmarshal(bsValue, dist); err != nil {
		return err
	}

	return nil
}

/**
* Функция Delete удаляет ключ из базы данных.
*
* @param key string ключ.
*/
func Delete(key string) error {
	err := cacheConfig.MemoryDB.Delete([]byte(key))
	if err != nil {
		return err
	}

	return cacheConfig.LocalDB.Delete([]byte(key))
}

/**
* Функция сlearExpired выполняет очистку кэша от просроченых блоков.
*/
func clearExpired() {
	iterator := cacheConfig.MemoryDB.Iterator()

	for iterator.Next() {
		key := string(iterator.Key())
		expiration := int64(binary.LittleEndian.Uint64(iterator.Value()))

		if isExpired(expiration) {
			Delete(key)
		}
	}

	iterator.Release()
}

/**
* Функция isExpired проверяет истекло ли время значения кэше.
*
* @param exp int64 время.
*/
func isExpired(expiration int64) bool {
	if expiration == 0 {
		return false
	}

	return time.Now().UnixNano() > expiration
}

/**
* Функция secondsToDuration преобразовывает секунды в time.Duration.
*/
func secondsToDuration(s int64) time.Duration {
	return time.Duration(s) * time.Second
}

/**
* Функция ContextToHash хеширует текущий контекст.
*
* @param context Context
* @return uint64, error
*/
func ContextToHash(context Context) uint64 {
	bs := []byte(context.Path() + context.Request().URL.RawQuery)

	return murmur3.Sum64(bs)
}
