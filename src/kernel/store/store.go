/**
* Пакет store представляет собой обёртку над клиентом LevelDB.
*
* LocalStore - локальное и основное key-value хранилище.
* MemoryStore - key-value хранилище в памяти. Удобно для некоторых промежуточных и мелких операций.
*
 */
package store

/**
*
* type Store interface {
*
*	Put(key, value []byte) error
*
*	Get(key []byte) ([]byte, error)
*
*	Delete(key []byte) error
*
* }
*
*/
