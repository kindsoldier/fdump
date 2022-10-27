
/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package dspack

import (
    "encoding/json"
)

type HeadDescr struct {
    Path    string          `json:"path"`
    Mtime   int64           `json:"mtime"`
    Atime   int64           `json:"atime"`
    Ctime   int64           `json:"atime"`
    Size    int64           `json:"size"`
    Mode    uint32          `json:"mode"`
    Type    int64           `json:"type"`
    SLink   string          `json:"sLink,omitempty"`
    Match   bool            `json:"match"`
    Uid     uint32          `json:"uid"`
    Gid     uint32          `json:"gid"`
    User    string          `json:"user"`
    Group   string          `json:"group"`
    HType   string          `json:"hType"`
    HInit   []byte          `json:"hInit"`
}

func NewHeadDescr() *HeadDescr {
    var descr HeadDescr
    return &descr
}

func UnpackHeadDescr(descrBin []byte) (*HeadDescr, error) {
    var err error
    var descr HeadDescr
    err = json.Unmarshal(descrBin, &descr)
    return &descr, err
}

func (descr *HeadDescr) Pack() ([]byte, error) {
    var err error
    descrBin, err := json.Marshal(descr)
    return descrBin, err
}
