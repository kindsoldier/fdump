/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package dsdescr

import (
    encoder "github.com/vmihailenco/msgpack/v5"
)

const UStateEnabled     string  = "enabled"
const UStateDisabled    string  = "disabled"

const URoleAdmin        string  = "admin"
const URoleUser         string  = "user"

type User struct {
    Login       string      `json:"login"	msgpack:"login"`
    Pass        string      `json:"pass"	msgpack:"pass"`
    Role        string      `json:"role"	msgpack:"role"`
    State       string      `json:"state"	msgpack:"state"`
    CreatedAt   int64       `json:"updatedAt"	msgpack:"updatedAt"`
    UpdatedAt   int64       `json:"createdAt"	msgpack:"createdAt"`
}

func NewUser() *User {
    var descr User
    return &descr
}

func UnpackUser(descrBin []byte) (*User, error) {
    var err error
    var descr User
    err = encoder.Unmarshal(descrBin, &descr)
    return &descr, err
}

func (descr *User) Pack() ([]byte, error) {
    var err error
    descrBin, err := encoder.Marshal(descr)
    return descrBin, err
}


type File struct {
    FilePath    string      `json:"filePath"	msgpack:"filePath"`
    Login       string      `json:"login"	msgpack:"login"`
    FileId      int64       `json:"fileId"	msgpack:"fileId"`
    BatchCount  int64       `json:"batchCount"	msgpack:"batchCount"`
    BatchSize   int64       `json:"batchSize"	msgpack:"batchSize"`
    BlockSize   int64       `json:"blockSize"	msgpack:"blockSize"`
    DataSize    int64       `json:"dataSize"	msgpack:"dataSize"`
    CreatedAt   int64       `json:"createdAt"	msgpack:"createdAt"`
    UpdatedAt   int64       `json:"updatedAt"	msgpack:"updatedAt"`
}

func NewFile() *File {
    var descr File
    return &descr
}

func UnpackFile(descrBin []byte) (*File, error) {
    var err error
    var descr File
    err = encoder.Unmarshal(descrBin, &descr)
    return &descr, err
}

func (descr *File) Pack() ([]byte, error) {
    var err error
    descrBin, err := encoder.Marshal(descr)
    return descrBin, err
}


type Batch struct {
    BatchId     int64       `json:"batchId"	msgpack:"batchId"`
    FileId      int64       `json:"fileId"	msgpack:"fileId"`
    BatchSize   int64       `json:"batchSize"	msgpack:"batchSize"`
    BlockSize   int64       `json:"blockSize"	msgpack:"blockSize"`
    CreatedAt   int64       `json:"createdAt"	msgpack:"createdAt"`
    UpdatedAt   int64       `json:"updatedAt"	msgpack:"updatedAt"`
}

func NewBatch() *Batch {
    var descr Batch
    return &descr
}

func UnpackBatch(descrBin []byte) (*Batch, error) {
    var err error
    var descr Batch
    err = encoder.Unmarshal(descrBin, &descr)
    return &descr, err
}

func (descr *Batch) Pack() ([]byte, error) {
    var err error
    descrBin, err := encoder.Marshal(descr)
    return descrBin, err
}

const BTData int64 = 1
const BTReco int64 = 2

type Block struct {
    FileId      int64       `json:"fileId"	msgpack:"fileId"`
    BatchId     int64       `json:"batchId"	msgpack:"batchId"`
    BlockType   int64       `json:"blockType"	msgpack:"blockType"`
    BlockId     int64       `json:"blockId"	msgpack:"blockId"`

    BlockSize   int64       `json:"blockSize"	msgpack:"blockSize"`
    DataSize    int64       `json:"dataSize"	msgpack:"dataSize"`
    CreatedAt   int64       `json:"createdAt"	msgpack:"createdAt"`
    UpdatedAt   int64       `json:"updatedAt"	msgpack:"updatedAt"`
    FilePath    string      `json:"filePath"	msgpack:"filePath"`

    HashInit    string      `json:"hashInit"	msgpack:"hashInit"`
    HashSum     string      `json:"hashSum"	msgpack:"hashSum"`

    HasLocal    bool        `json:"hasLocal"	msgpack:"hasLocal"`
    HasRemote   bool        `json:"hasRemote"	msgpack:"hasRemote"`
    BstoreId    int64       `json:"bstoreId"	msgpack:"bstoreId"`
}

func NewBlock() *Block {
    var descr Block
    return &descr
}

func UnpackBlock(descrBin []byte) (*Block, error) {
    var err error
    var descr Block
    err = encoder.Unmarshal(descrBin, &descr)
    return &descr, err
}

func (descr *Block) Pack() ([]byte, error) {
    var err error
    descrBin, err := encoder.Marshal(descr)
    return descrBin, err
}
