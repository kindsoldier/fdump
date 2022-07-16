/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package dsinter

import (
    "fdump/dscomm/dsdescr"
)

type IterFunc = func(key []byte, val []byte) (bool, error)
type DB interface {
    Put(key, val []byte) error
    Get(key []byte) ([]byte, error)
    Has(key []byte) (bool, error)
    Delete(key []byte) error
    Iter(prefix []byte, cb IterFunc) error
}

type Alloc interface {
    NewId() (int64, error)
    FreeId(id int64) error
    JSON() ([]byte, error)
}

type BStoreReg interface {
    PutUser(descr *dsdescr.User) error
    HasUser(login string) (bool, error)
    GetUser(login string) (*dsdescr.User, error)
    ListUsers() ([]*dsdescr.User, error)
    DeleteUser(login string) error

}
