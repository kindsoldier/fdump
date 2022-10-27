
/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package dspack

import (
    "encoding/json"
)

type TailDescr struct {
    HSum    []byte          `json:"hSum"`
}

func NewTailDescr() *TailDescr {
    var descr TailDescr
    return &descr
}

func UnpackTailDescr(descrBin []byte) (*TailDescr, error) {
    var err error
    var descr TailDescr
    err = json.Unmarshal(descrBin, &descr)
    return &descr, err
}

func (descr *TailDescr) Pack() ([]byte, error) {
    var err error
    descrBin, err := json.Marshal(descr)
    return descrBin, err
}
