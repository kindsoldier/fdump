/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package dskvdb

import (
    "path/filepath"
    "github.com/syndtr/goleveldb/leveldb"
    "github.com/syndtr/goleveldb/leveldb/util"

    "fdump/dscomm/dsinter"
)

type DB struct {
    ldb *leveldb.DB
}

func OpenDB(dataDir, name string) (*DB, error) {
    var err error
    var db DB
    dbPath  := filepath.Join(dataDir, name)
    ldb, err := leveldb.OpenFile(dbPath, nil)
    db.ldb = ldb
    return &db, err
}

func (db *DB) Put(key, val []byte) error {
    return db.ldb.Put(key, val, nil)
}

func (db *DB) Get(key []byte) ([]byte, error) {
    return db.ldb.Get(key, nil)
}

func (db *DB) Has(key []byte) (bool, error) {
    return db.ldb.Has(key, nil)
}

func (db *DB) Delete(key []byte) error {
    return db.ldb.Delete(key, nil)
}

func (db *DB) Close() error {
    return db.ldb.Close()
}

func (db *DB) Iter(prefix []byte, cb dsinter.IterFunc) error {
    var err error
    bPrefix := util.BytesPrefix(prefix)
    iter := db.ldb.NewIterator(bPrefix, nil)
    defer iter.Release()
    for iter.Next() {
        stop, _ := cb(iter.Key(), iter.Value())
        if stop {
            break
        }
    }
    err = iter.Error()
    return err
}
