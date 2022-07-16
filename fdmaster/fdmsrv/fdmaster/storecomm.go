/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package fdmaster

import (
    "io/fs"
    "time"
    "syscall"
    "fdump/dscomm/dsinter"
)

type Store struct {
    dataDir     string
    reg         dsinter.BStoreReg
    dirPerm     fs.FileMode
    filePerm    fs.FileMode
    startTime   int64
}

func NewStore(dataDir string, reg dsinter.BStoreReg) (*Store, error) {
    var err error
    var store Store
    store.dataDir   = dataDir
    store.reg       = reg
    store.dirPerm   = 0755
    store.filePerm  = 0644
    store.startTime = time.Now().Unix()
    return &store, err
}

func (store *Store) SetDirPerm(dirPerm fs.FileMode) {
    store.dirPerm = dirPerm
}

func (store *Store) SetFilePerm(filePerm fs.FileMode) {
    store.filePerm = filePerm
}

func (store *Store) GetUptime() (int64, error) {
    var err error
    uptime := time.Now().Unix() - store.startTime
    return uptime, err
}

func (store *Store) GetUsage() (uint64, uint64, uint64, error) {
    var err error
    var all, free, used uint64
    path := store.dataDir
    fs := syscall.Statfs_t{}
    err = syscall.Statfs(path, &fs)
    if err != nil {
        return all, free, used, err
    }
    all = fs.Blocks * uint64(fs.Bsize)
    free = fs.Bfree * uint64(fs.Bsize)
    used = all - free
    return all, free, used, err
}
