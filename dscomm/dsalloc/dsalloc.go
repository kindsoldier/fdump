/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package dsalloc

import (
    "encoding/json"
    "sync"
    "dstore/dscomm/dsinter"
)

type Alloc struct {
    db      dsinter.DB
    topId   int64
    freeIds []int64
    key     []byte
    giantMtx      sync.Mutex
}

func OpenAlloc(db dsinter.DB, key []byte) (*Alloc, error) {
    var err error
    var alloc Alloc

    alloc.db        = db
    alloc.freeIds   = make([]int64, 0)
    alloc.key       = key
    alloc.topId     = 0

    has, err := alloc.db.Has(alloc.key)
    if err != nil {
        return &alloc, err
    }
    if has {
        descrBin, err := alloc.db.Get(alloc.key)
        if err != nil {
            return &alloc, err
        }
        descr, err := UnpackAllocDescr(descrBin)
        if err != nil {
            return &alloc, err
        }
        alloc.freeIds = descr.FreeIds
        alloc.topId   = descr.TopId
    }
    return &alloc, err
}

func (alloc *Alloc) NewId() (int64, error) {
    var err error
    var newId int64

    alloc.giantMtx.Lock()
    defer alloc.giantMtx.Unlock()

    freeIds := len(alloc.freeIds)
    if freeIds > 0 {
        newId = alloc.freeIds[freeIds - 1]
        alloc.freeIds = alloc.freeIds[0:freeIds - 1]
        return newId, err
    }

    newId = alloc.topId + 1
    err = alloc.storeState()
    if err != nil {
        alloc.freeIds = append(alloc.freeIds, newId)
        newId = -1
        return newId, err
    }
    alloc.topId = newId
    return newId, err
}

func (alloc *Alloc) FreeId(id int64) error {
    var err error

    alloc.giantMtx.Lock()
    defer alloc.giantMtx.Unlock()

    switch {
        case id == alloc.topId:
            alloc.topId--
        case id > alloc.topId:  // todo: ???
            return err
        default:
            alloc.freeIds = append(alloc.freeIds, id)
    }
    err = alloc.storeState()
    if err != nil {
        return err
    }
    return err
}

func (alloc *Alloc) JSON() ([]byte, error) {
    var err error
    descr := alloc.toDescr()
    descrBin, err := descr.Pack()
    if err != nil {
        return descrBin, err
    }
    return descrBin, err
}


func (alloc *Alloc) toDescr() *AllocDescr {
    descr := NewAllocDescr()
    descr.TopId     = alloc.topId
    descr.FreeIds   = alloc.freeIds
    return descr
}

func (alloc *Alloc) storeState() error {
    var err error
    descr := alloc.toDescr()
    descrBin, err := descr.Pack()
    if err != nil {
        return err
    }
    err = alloc.db.Put(alloc.key, descrBin)
    if err != nil {
        return err
    }
    return err
}


type AllocDescr struct {
    TopId   int64           `json:"topId"`
    FreeIds []int64         `json:"freeIds"`
}

func NewAllocDescr() *AllocDescr {
    var descr AllocDescr
    return &descr
}

func UnpackAllocDescr(descrBin []byte) (*AllocDescr, error) {
    var err error
    var descr AllocDescr
    err = json.Unmarshal(descrBin, &descr)
    return &descr, err
}

func (descr *AllocDescr) Pack() ([]byte, error) {
    var err error
    descrBin, err := json.Marshal(descr)
    return descrBin, err
}
