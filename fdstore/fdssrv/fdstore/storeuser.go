/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package fdstore

import (
    "errors"
    "fmt"
    "time"
    "fdump/dscomm/dsdescr"
    "fdump/dscomm/dserr"
)


const defaultAUser      string  = "admin"
const defaultAPass      string  = "admin"

const defaultUser       string  = "user"
const defaultPass       string  = "user"

func (store *Store) SeedUsers() error {
    var err error
    users, err := store.reg.ListUsers()
    if err != nil {
        return dserr.Err(err)
    }
    if len(users) < 1 {
        var user *dsdescr.User
        user = dsdescr.NewUser()
        user.Login  = defaultAUser
        user.Pass   = defaultAPass
        user.State  = dsdescr.UStateEnabled
        user.Role   = dsdescr.URoleAdmin
        user.CreatedAt = time.Now().Unix()
        user.UpdatedAt = user.CreatedAt

        err = store.reg.PutUser(user)
        if err != nil {
            return dserr.Err(err)
        }
        user = dsdescr.NewUser()
        user.Login  = defaultUser
        user.Pass   = defaultPass
        user.State  = dsdescr.UStateEnabled
        user.Role   = dsdescr.URoleUser
        user.CreatedAt = time.Now().Unix()
        user.UpdatedAt = user.CreatedAt


        err = store.reg.PutUser(user)
        if err != nil {
            return dserr.Err(err)
        }
    }
    return dserr.Err(err)
}

func (store *Store) AddUser(authLogin string, user *dsdescr.User) error {
    var err error
    var ok bool

    role, err := store.getUserRole(authLogin)
    if role != dsdescr.URoleAdmin {
        err = fmt.Errorf("insufficient rights for %s", authLogin)
        return dserr.Err(err)
    }
    ok, err = validateLogin(user.Login)
    if !ok {
        return dserr.Err(err)
    }
    ok, err = validatePass(user.Pass)
    if !ok {
        return dserr.Err(err)
    }

    has, err := store.reg.HasUser(user.Login)
    if err != nil {
        return dserr.Err(err)
    }
    if has {
        err = fmt.Errorf("login %s exist", user.Login)

    }
    user.State  = dsdescr.UStateEnabled
    user.Role   = dsdescr.URoleUser
    user.CreatedAt = time.Now().Unix()
    user.UpdatedAt = user.CreatedAt

    err = store.reg.PutUser(user)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}

func (store *Store) GetUser(login string) (bool, *dsdescr.User, error) {
    var err error
    var user *dsdescr.User
    has, err := store.reg.HasUser(login)
    if err != nil {
        return has, user, dserr.Err(err)
    }
    if !has {
        return has, user, dserr.Err(err)
    }
    user, err = store.reg.GetUser(login)
    if err != nil {
        return has, user, dserr.Err(err)
    }
    return has,user, dserr.Err(err)
}

func (store *Store) CheckUser(authLogin, login, passw string) (bool, error) {
    var err error
    var ok bool

    if len(login) == 0 {
        login = authLogin
    }
    userRole, err := store.getUserRole(authLogin)
    if authLogin != login && userRole != dsdescr.URoleAdmin {
        err = fmt.Errorf("user %s have insufficient rights", authLogin)
        return ok, dserr.Err(err)
    }
    has, err := store.reg.HasUser(login)
    if err != nil {
        return ok, dserr.Err(err)
    }
    if !has {
        err = fmt.Errorf("user %s not exist", login)
    }
    user, err := store.reg.GetUser(login)
    if err != nil {
        return ok, dserr.Err(err)
    }
    if passw == user.Pass {
        ok = true
    }
    return ok, dserr.Err(err)
}

func (store *Store) UpdateUser(authLogin string, user *dsdescr.User) error {
    var err error
    // Get current role
    userRole, err := store.getUserRole(authLogin)
    if err != nil {
        return dserr.Err(err)
    }
    // Set defaults
    if len(user.Login) < 1 {
        user.Login = authLogin
    }
    // Rigth control
    if  authLogin != user.Login && userRole != dsdescr.URoleAdmin {
        err = fmt.Errorf("user %s have insufficient rights", authLogin)
        return dserr.Err(err)
    }

    // Get old profile and copy to new
    oldUser, err := store.reg.GetUser(user.Login)
    if err != nil {
        return dserr.Err(err)
    }
    newUser := dsdescr.NewUser()
    newUser.Login       = oldUser.Login
    newUser.Pass        = oldUser.Pass
    newUser.Role        = oldUser.Role
    newUser.State       = oldUser.State
    newUser.CreatedAt   = oldUser.CreatedAt
    newUser.UpdatedAt   = time.Now().Unix()

    // Update property if exists
    if len(user.Pass) > 0 {
        newUser.Pass = user.Pass
    }
    if len(user.Role) > 0 {
        newUser.Role = user.Role
    }
    if len(user.State) > 0 {
        newUser.State = user.State
    }
    // Rigth control
    if newUser.Role != oldUser.Role && userRole != dsdescr.URoleAdmin {
        err = errors.New("insufficient rights for changing role")
        return dserr.Err(err)
    }
    if newUser.State != oldUser.State && userRole != dsdescr.URoleAdmin {
        err = errors.New("insufficient rights for changing state")
        return dserr.Err(err)
    }

    // Validation new property
    var ok bool
    ok, err = validateUState(newUser.State)
    if !ok {
        return dserr.Err(err)
    }
    ok, err = validateURole(newUser.Role)
    if !ok {
        return dserr.Err(err)
    }
    ok, err = validatePass(newUser.Pass)
    if !ok {
        return dserr.Err(err)
    }
    // Delete old user descr
    err = store.reg.DeleteUser(user.Login)
    if err != nil {
        return dserr.Err(err)
    }
    // Put new user descr
    err = store.reg.PutUser(newUser)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}

func (store *Store) ListUsers(authLogin string) ([]*dsdescr.User, error) {
    var err error
    users := make([]*dsdescr.User, 0)
    userRole, err := store.getUserRole(authLogin)
    if userRole != dsdescr.URoleAdmin {
        err = fmt.Errorf("user %s have insufficient rights", authLogin)
        return users, dserr.Err(err)
    }
    users, err = store.reg.ListUsers()
    if err != nil {
        return users, dserr.Err(err)
    }
    return users, dserr.Err(err)
}

func (store *Store) DeleteUser(authLogin string, login string) error {
    var err error

    userRole, err := store.getUserRole(authLogin)
    if authLogin != login && userRole != dsdescr.URoleAdmin {
        err = fmt.Errorf("user %s have insufficient rights", authLogin)
        return dserr.Err(err)
    }

    err = store.reg.DeleteUser(login)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}

func (store *Store) getUserRole(authLogin string) (string, error) {
    var err error
    var userRole string
    has, err := store.reg.HasUser(authLogin)
    if !has {
        err = fmt.Errorf("user %s not exists", authLogin)
        return userRole, dserr.Err(err)
    }
    user, err := store.reg.GetUser(authLogin)
    if err != nil {
        return userRole, dserr.Err(err)
    }
    userRole = user.Role
    if err != nil {
        return userRole, dserr.Err(err)
    }
    return userRole, dserr.Err(err)
}

func validateURole(role string) (bool, error) {
    var err error
    var ok bool = true
    if role == dsdescr.URoleAdmin  {
        return ok, dserr.Err(err)
    }
    if role == dsdescr.URoleUser  {
        return ok, dserr.Err(err)
    }
    err = errors.New("irrelevant role name")
    ok = false
    return ok, dserr.Err(err)
}

func validateUState(state string) (bool, error) {
    var err error
    var ok bool = true
    if state == dsdescr.UStateDisabled  {
        return ok, dserr.Err(err)
    }
    if state == dsdescr.UStateEnabled  {
        return ok, dserr.Err(err)
    }
    err = errors.New("irrelevant state name")
    ok = false
    return ok, dserr.Err(err)
}

func validateLogin(login string) (bool, error) {
    var err error
    var ok bool = true
    if len(login) == 0 {
        ok = false
        err = errors.New("zero len password")
    }
    return ok, dserr.Err(err)
}

func validatePass(passw string) (bool, error) {
    var err error
    var ok bool = true
    if len(passw) == 0 {
        ok = false
        err = errors.New("zero len password")
    }
    return ok, dserr.Err(err)
}
