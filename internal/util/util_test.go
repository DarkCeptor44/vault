package util_test

import (
	"testing"
	"time"

	"github.com/DarkCeptor44/vault/internal/util"
)

func TestCache(t *testing.T) {
	// janky way of checking if the first derivation is cached and the other times it's not
	for i := 0; i < 3; i++ {
		start := time.Now()
		key := util.DeriveKey([]byte("password"), []byte("salt"))
		elapsed := time.Since(start)

		if i == 0 && elapsed.String() == "0s" {
			t.Fatal("first run is not supposed to be 0s")
		} else if i > 0 && elapsed.String() != "0s" {
			t.Fatal("other runs are supposed to be 0s")
		}

		t.Logf("%x: %s\n", key[:5], elapsed)
	}
}

func BenchmarkUtils(b *testing.B) {
	c := util.NewCache()
	c2 := util.NewCache()

	b.Run("Cache.Load", wrap(func() {
		c.Load([]byte("password"))
	}))

	b.Run("Cache.Store", wrap(func() {
		c2.Store([]byte("password"), []byte("key"))
	}))

	b.Run("NewSalt", wrap(func() {
		util.NewSalt()
	}))

	b.Run("DeriveKey", wrap(func() {
		util.DeriveKey([]byte("password"), []byte("salt"))
	}))

	b.Run("EncryptData", wrap(func() {
		util.EncryptData([]byte("password"), []byte("data"))
	}))

	b.Run("DecryptData", wrap(func() {
		util.DecryptData([]byte("password"), []byte("data"))
	}))

	b.Run("ClearFilename", wrap(func() {
		x := util.ClearFilename("hello world")
		if x != "helloworld" {
			b.Fatal(x)
		}
	}))
}

// auxiliary functions

func wrap(f func()) func(b *testing.B) {
	return func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			f()
		}
	}
}
