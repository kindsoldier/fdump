/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package fdmcont

import (
    "fdump/fdmaster/fdmsrv/fdmaster"
)


type Contr struct {
    store  *fdmaster.Store
}

func NewContr(store *fdmaster.Store) (*Contr, error) {
    var err error
    var contr Contr
    contr.store = store
    return &contr, err
}
