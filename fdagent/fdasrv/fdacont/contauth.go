/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package fdacont

import (
    "errors"

    "fdump/dscomm/dsrpc"
    "fdump/dscomm/dslog"
    "fdump/dscomm/dserr"
)


func (contr *Contr) AuthMidware(debugMode bool) dsrpc.HandlerFunc {
    return func(context *dsrpc.Context) error {

        var err error
        login := context.AuthIdent()
        salt := context.AuthSalt()
        hash := context.AuthHash()

        has, user, err := contr.store.GetUser(string(login))
        if err != nil {
            resErr := errors.New("auth mismatch")
            context.SendError(resErr)
            return dserr.Err(err)
        }
        if !has {
            resErr := errors.New("auth error")
            context.SendError(resErr)
            return dserr.Err(err)
        }

        if debugMode {
            auth := context.Auth()
            dslog.LogDebug("auth ", string(auth.JSON()))
        }

        pass := []byte(user.Pass)
        ok := dsrpc.CheckHash(login, pass, salt, hash)
        if debugMode {
            dslog.LogDebugf("auth for %s is %v", login, ok)
        }
        if !ok {
            resErr := errors.New("auth mismatch")
            context.SendError(resErr)
            return dserr.Err(err)
        }
        return dserr.Err(err)
    }
}
