/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package fdstore

import (
    "testing"
    "github.com/stretchr/testify/require"

    "fdump/dscomm/dskvdb"
    "fdump/dscomm/dsdescr"
    "fdump/fdstore/fdssrv/fdsreg"
)


func TestUser01(t *testing.T) {
    var err error

    dataDir := t.TempDir()

    db, err := dskvdb.OpenDB(dataDir, "storedb")
    defer db.Close()
    require.NoError(t, err)

    reg, err := bsreg.NewReg(db)
    require.NoError(t, err)

    store, err := NewStore(dataDir, reg)
    require.NoError(t, err)

    err = store.SeedUsers()
    require.NoError(t, err)

    descr0 := dsdescr.NewUser()
    descr0.Login    = "qwerty"
    descr0.Pass     = "123456"

    adminLogin   := "admin"
    wrongLogin   := "wrong"

    err = store.AddUser(wrongLogin, descr0)
    require.Error(t, err)

    err = store.AddUser(adminLogin, descr0)
    require.NoError(t, err)

    has, descr1, err := store.GetUser(descr0.Login)
    require.NoError(t, err)
    require.Equal(t, has, true)
    require.Equal(t, descr0, descr1)

    var ok bool
    ok, err = store.CheckUser(adminLogin, descr0.Login, descr0.Pass)
    require.NoError(t, err)
    require.Equal(t, true, ok)

    ok, err = store.CheckUser(wrongLogin, descr0.Login, descr0.Pass)
    require.Error(t, err)
    require.Equal(t, false, ok)

    err = store.DeleteUser(wrongLogin, descr0.Login)
    require.Error(t, err)

    err = store.DeleteUser(adminLogin, descr0.Login)
    require.NoError(t, err)

    err = store.DeleteUser(descr0.Login, descr0.Login)
    require.NoError(t, err)

    _, err = store.ListUsers(wrongLogin)
    require.Error(t, err)

    descrs, err := store.ListUsers(adminLogin)
    require.NoError(t, err)
    require.Equal(t, len(descrs), 2)

}
