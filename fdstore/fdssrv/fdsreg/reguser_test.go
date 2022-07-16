/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package fdsreg

import(
    "testing"
    "github.com/stretchr/testify/require"

    "fdump/dscomm/dsdescr"
    "fdump/dscomm/dskvdb"
)

func TestUser01(t *testing.T) {
    var err error
    var has bool

    dataDir := t.TempDir()

    db, err := dskvdb.OpenDB(dataDir, "tmp.db")
    defer db.Close()
    require.NoError(t, err)


    reg, err := NewReg(db)
    require.NoError(t, err)
    require.NotEqual(t, reg, nil)

    descr0 := dsdescr.NewUser()
    descr0.Login    = "qwerty"
    descr0.Pass     = "123456"
    descr0.Role     = dsdescr.URoleUser
    descr0.State    = dsdescr.UStateEnabled
    descr0.CreatedAt = 1657645101
    descr0.UpdatedAt = 1657645102

    err = reg.PutUser(descr0)
    require.NoError(t, err)

    has, err = reg.HasUser(descr0.Login)
    require.NoError(t, err)
    require.Equal(t, has, true)

    descr1, err := reg.GetUser(descr0.Login)
    require.NoError(t, err)
    require.Equal(t, descr0, descr1)

    descrs, err := reg.ListUsers()
    require.NoError(t, err)
    require.Equal(t, len(descrs), 1)

    err = reg.DeleteUser(descr0.Login)
    require.NoError(t, err)

    has, err = reg.HasUser(descr0.Login)
    require.NoError(t, err)
    require.Equal(t, has, false)
}
