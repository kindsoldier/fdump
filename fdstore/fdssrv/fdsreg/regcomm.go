
package fdsreg

import (
    "fdump/dscomm/dsinter"
)

type Reg struct {
    db      dsinter.DB
    sep         string
    entryBase   string
    userBase    string
    blockBase   string
    batchBase   string
    fileBase    string
}

func NewReg(db dsinter.DB) (*Reg, error) {
    var err error
    var reg Reg
    reg.db          = db
    reg.sep         = ":"
    reg.entryBase   = "entry"
    reg.userBase    = "user"
    reg.blockBase   = "block"
    reg.batchBase   = "batch"
    reg.fileBase    = "file"
    return &reg, err
}
