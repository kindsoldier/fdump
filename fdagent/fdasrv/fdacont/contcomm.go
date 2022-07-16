/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package fdacont

import (
    "fdump/fdagent/fdasrv/fdagent"
)


type Contr struct {
    store  *fdagent.Store
}

func NewContr(store *fdagent.Store) (*Contr, error) {
    var err error
    var contr Contr
    contr.store = store
    return &contr, err
}
