package limiter_test

import (
	"github.com/techfront/core/src/lib/limiter"
	"gopkg.in/throttled/throttled.v2/store/memstore"
	"gopkg.in/throttled/throttled.v2/store/storetest"
	"log"
	"testing"
)

var ls *limiter.Store
var ms *memstore.MemStore

func init() {
	ldb, err := limiter.NewStore("./tmp")
	if err != nil {
		log.Fatal(err)
	}

	mdb, err := memstore.New(65536)
	if err != nil {
		log.Fatal(err)
	}

	ls = ldb
	ms = mdb
}

func TestLeveldbStore(t *testing.T) {
	log.Print("Testing leveldb store...")
	storetest.TestGCRAStore(t, ls)
}

func TestMemStore(t *testing.T) {
	log.Print("Testing mem store...")
	storetest.TestGCRAStore(t, ms)
}

func BenchmarkMemStoreLRU(b *testing.B) {
	log.Print("Benchmarking mem store...")
	storetest.BenchmarkGCRAStore(b, ms)
}

func BenchmarkLeveldbStore(b *testing.B) {
	log.Print("Benchmarking leveldb store...")
	storetest.BenchmarkGCRAStore(b, ls)
}
