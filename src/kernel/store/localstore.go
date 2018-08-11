package store

import (
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"os"
)

type LocalStore struct {
	path string
	db *leveldb.DB
	options *opt.Options
	writeOptions *opt.WriteOptions
	readOptions *opt.ReadOptions
}

/**
* Функция NewLocalStore инициализирует новое локальное хранилище.
*
* @param path string путь в файловой системе.
* @param reset bool предварительный сброс хранилища.
* @param compression bool сжатие блоков (по умолчанию Snappy).
*
* @return ls *LocalStore
* @return err error
 */
func NewLocalStore(path string, reset bool, compression bool) (*LocalStore, error) {

	// Определяем базовые опции
	options := &opt.Options{
		Filter:      filter.NewBloomFilter(10),
		Compression: opt.NoCompression,
	}

	if compression {
		options.Compression = opt.SnappyCompression
	}

	// Определяем опции записи
	writeOptions := &opt.WriteOptions{
		Sync: false,
	}

	if reset {
		os.RemoveAll(path)
		writeOptions.Sync = true
	}

	// Определяем опции чтения
	readOptions := &opt.ReadOptions{
		Strict: opt.DefaultStrict,
	}

	// Инициализация LevelDB
	db, err := leveldb.OpenFile(path, options)
	if err != nil {
		return nil, err
	}

	ls := &LocalStore{
		path: path,
		db: db,
		options: options,
		writeOptions: writeOptions,
		readOptions: readOptions,
	}

	return ls, nil
}

/**
* Функция добавляет ключ и значение в локальную БД.
*
* @param key []byte ключ (байт-массив)
* @param value []byte значение (байт-массив)
* @return error
 */
func (ls *LocalStore) Put(key, value []byte) error {
	return ls.db.Put(key, value, ls.writeOptions)
}

/**
* Функция получает блок из локальной БД.
*
* @param key []byte ключ (байт-массив)
* @return []byte
* @return error
 */
func (ls *LocalStore) Get(key []byte) ([]byte, error) {
	return ls.db.Get(key, ls.readOptions)
}

/**
* Функция проверяет есть ли блок в БД.
*
* @param key []byte ключ (байт-массив)
* @return bool
* @return error
 */
func (ls *LocalStore) Has(key []byte) (bool, error) {
	return ls.db.Has(key, ls.readOptions)
}

/**
* Функция удаляет блок из локальной БД.
*
* @param key []byte ключ (байт-массив)
* @return error
 */
func (ls *LocalStore) Delete(key []byte) error {
	return ls.db.Delete(key, ls.writeOptions)
}

/**
* Функция HotReset выполняет горячую очистку инициализированной локальной БД.
*
* @return error
 */
func (ls *LocalStore) HotReset() error {
	if err := ls.db.Close(); err != nil {
		return err
	}

	if err := os.RemoveAll(ls.path); err != nil {
		return err
	}

	db, err := leveldb.OpenFile(ls.path, ls.options)
	if err != nil {
		return err
	}

	ls.db = db

	return nil
}
