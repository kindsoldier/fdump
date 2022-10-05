
/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package dspack

import (
    "encoding/json"
)

type TDescr struct {
    HSum    []byte          `json:"hSum"`
}

func NewTDescr() *TDescr {
    var descr TDescr
    return &descr
}

func UnpackTDescr(descrBin []byte) (*TDescr, error) {
    var err error
    var descr TDescr
    err = json.Unmarshal(descrBin, &descr)
    return &descr, err
}

func (descr *TDescr) Pack() ([]byte, error) {
    var err error
    descrBin, err := json.Marshal(descr)
    return descrBin, err
}
