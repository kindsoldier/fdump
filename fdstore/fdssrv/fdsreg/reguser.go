package fdsreg

import (
    "strings"
    "fdump/dscomm/dsdescr"
)

func (reg *Reg) PutUser(descr *dsdescr.User) error {
    var err error
    keyArr := []string{ reg.userBase, descr.Login }
    keyBin := []byte(strings.Join(keyArr, reg.sep))
    valBin, _ := descr.Pack()
    err = reg.db.Put(keyBin, valBin)
    return err
}

func (reg *Reg) HasUser(login string) (bool, error) {
    var err error
    keyArr := []string{ reg.userBase, login }
    keyBin := []byte(strings.Join(keyArr, reg.sep))
    has, err := reg.db.Has(keyBin)
    if err != nil {
        return has, err
    }
    return has, err
}

func (reg *Reg) GetUser(login string) (*dsdescr.User, error) {
    var err error
    var descr *dsdescr.User
    keyArr := []string{ reg.userBase, login }
    keyBin := []byte(strings.Join(keyArr, reg.sep))
    valBin, err := reg.db.Get(keyBin)
    if err != nil {
        return descr, err
    }
    descr, err = dsdescr.UnpackUser(valBin)
    if err != nil {
        return descr, err
    }
    return descr, err
}

func (reg *Reg) DeleteUser(login string) error {
    var err error
    keyArr := []string{ reg.userBase, login }
    keyBin := []byte(strings.Join(keyArr, reg.sep))
    err = reg.db.Delete(keyBin)
    if err != nil {
        return err
    }
    return err
}

func (reg *Reg) ListUsers() ([]*dsdescr.User, error) {
    var err error
    descrs := make([]*dsdescr.User, 0)
    cb := func(key []byte, val []byte) (bool, error) {
        var err error
        var interr bool
        descr, err := dsdescr.UnpackUser(val)
        if err != nil {
            return interr, err
        }
        descrs = append(descrs, descr)
        return interr, err
    }
    userKeyBaseBin := []byte(reg.userBase + reg.sep)
    err = reg.db.Iter(userKeyBaseBin, cb)
    if err != nil {
        return descrs, err
    }
    return descrs, err
}
