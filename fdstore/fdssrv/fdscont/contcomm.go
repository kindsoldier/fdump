/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package fdscont

import (
    "fdump/fdstore/fdssrv/fdstore"
)


type Contr struct {
    store  *fdstore.Store
}

func NewContr(store *fdstore.Store) (*Contr, error) {
    var err error
    var contr Contr
    contr.store = store
    return &contr, err
}
