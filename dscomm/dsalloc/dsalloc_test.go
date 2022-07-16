
package dsalloc

import (
    "testing"
    "github.com/stretchr/testify/require"
    "dstore/dscomm/dskvdb"
)

func TestAlloc01(t *testing.T) {
    var err error

    dataDir := t.TempDir()

    db, err := dskvdb.OpenDB(dataDir, "tmp.db")
    defer db.Close()
    require.NoError(t, err)

    key := []byte("blockids")
    alloc, err := OpenAlloc(db, key)
    require.NoError(t, err)

    id1, err := alloc.NewId()
    require.NoError(t, err)

    err = alloc.FreeId(id1)
    require.NoError(t, err)

    id2, err := alloc.NewId()
    require.NoError(t, err)
    require.Equal(t, id1, id2)

    err = alloc.FreeId(id2)
    require.NoError(t, err)
}

func BenchmarkIdAlloc(b *testing.B) {
    var err error

    dataDir := b.TempDir()

    db, err := dskvdb.OpenDB(dataDir, "tmp.db")
    defer db.Close()
    require.NoError(b, err)

    key := []byte("blockids")
    alloc, err := OpenAlloc(db, key)
    require.NoError(b, err)

    pBench := func(pb *testing.PB) {
        for pb.Next() {

            id, err := alloc.NewId()
            require.NoError(b, err)

            err = alloc.FreeId(id)
            require.NoError(b, err)
        }
    }
    b.SetParallelism(1000)
    b.RunParallel(pBench)
}
