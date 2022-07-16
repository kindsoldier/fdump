/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package main

import (
    "encoding/json"
    "fmt"
    "flag"
    "os"
    "path/filepath"
    "errors"

    "fdump/fdstore/fdsapi"
    "fdump/dscomm/dsrpc"
)

type any = interface{}

func main() {
    var err error
    util := NewUtil()
    err = util.Exec()
    if err != nil {
        fmt.Printf("Exec error: %s\n", err)
    }
}

type Util struct {
    aLogin      string
    aPass       string

    Port        string
    Address     string
    Message     string
    URI         string
    SubCmd      string

    Login       string
    Pass        string

    bPort       string
    bAddress    string

    FileId      int64
    BatchId     int64
    BlockId     int64
    BlockType   int64

    FilePath   string
}

func NewUtil() *Util {
    var util Util
    util.Port       = "5101"
    util.Address    = "127.0.0.1"
    util.Message    = "hello"
    util.aLogin     = "admin"
    util.aPass      = "admin"
    return &util
}

const getStatusCmd      string = "getStatus"

const saveBlockCmd      string = "saveBlock"
const loadBlockCmd      string = "loadBlock"
const listBlocksCmd     string = "listBlocks"
const deleteBlockCmd    string = "deleteBlock"

const addUserCmd        string = "addUser"
const checkUserCmd      string = "checkUser"
const updateUserCmd     string = "updateUser"
const deleteUserCmd     string = "deleteUser"
const listUsersCmd      string = "listUsers"


const helpCmd           string = "help"


func (util *Util) GetOpt() error {
    var err error

    exeName := filepath.Base(os.Args[0])

    flag.StringVar(&util.Port, "port", util.Port, "service port")
    flag.StringVar(&util.Address, "address", util.Address, "service address")
    flag.StringVar(&util.aLogin, "aLogin", util.aLogin, "access login")
    flag.StringVar(&util.aPass, "aPass", util.aPass, "access password")

    help := func() {
        fmt.Println("")
        fmt.Printf("Usage: %s [option] command [command option]\n", exeName)
        fmt.Printf("\n")
        fmt.Printf("Command list: help, getStatus, \n")
        fmt.Printf("    addUser, checkUser, updateUser, listUsers, deleteUser \n")

        fmt.Printf("\n")
        fmt.Printf("Global options:\n")
        flag.PrintDefaults()
        fmt.Printf("\n")
    }
    flag.Usage = help
    flag.Parse()

    args := flag.Args()

    //if len(args) == 0 {
    //    args = append(args, getStatusCmd)
    //}

    var subCmd string
    var subArgs []string
    if len(args) > 0 {
        subCmd = args[0]
        subArgs = args[1:]
    }
    switch subCmd {
        case helpCmd:
            help()
            return errors.New("unknown command")
        case getStatusCmd:
            flagSet := flag.NewFlagSet(getStatusCmd, flag.ExitOnError)
            flagSet.Usage = func() {
                fmt.Printf("\n")
                fmt.Printf("Usage: %s [global options] %s [command options]\n", exeName, subCmd)
                fmt.Printf("\n")
                fmt.Printf("The command options: none\n")
                flagSet.PrintDefaults()
                fmt.Printf("\n")
            }
            flagSet.Parse(subArgs)
            util.SubCmd = subCmd

        case addUserCmd, checkUserCmd, updateUserCmd:
            flagSet := flag.NewFlagSet(addUserCmd, flag.ExitOnError)
            flagSet.StringVar(&util.Login, "login", util.Login, "login")
            flagSet.StringVar(&util.Pass, "pass", util.Pass, "pass")
            flagSet.Usage = func() {
                fmt.Printf("\n")
                fmt.Printf("Usage: %s [global options] %s [command options]\n", exeName, subCmd)
                fmt.Printf("\n")
                fmt.Printf("The command options:\n")
                flagSet.PrintDefaults()
                fmt.Printf("\n")
            }
            flagSet.Parse(subArgs)
            util.SubCmd = subCmd
        case deleteUserCmd:
            flagSet := flag.NewFlagSet(deleteUserCmd, flag.ExitOnError)
            flagSet.StringVar(&util.Login, "login", util.Login, "login")
            flagSet.Usage = func() {
                fmt.Printf("\n")
                fmt.Printf("Usage: %s [global options] %s [command options]\n", exeName, subCmd)
                fmt.Printf("\n")
                fmt.Printf("The command options:\n")
                flagSet.PrintDefaults()
                fmt.Printf("\n")
            }
            flagSet.Parse(subArgs)
            util.SubCmd = subCmd
        case listUsersCmd:
            flagSet := flag.NewFlagSet(deleteUserCmd, flag.ExitOnError)
            flagSet.Usage = func() {
                fmt.Printf("\n")
                fmt.Printf("Usage: %s [global options] %s [command options]\n", exeName, subCmd)
                fmt.Printf("\n")
                fmt.Printf("The command options: none\n")
                flagSet.PrintDefaults()
                fmt.Printf("\n")
            }
            flagSet.Parse(subArgs)
            util.SubCmd = subCmd

        default:
            help()
            return errors.New("unknown command")
    }
    return err
}

type Response struct {
    Error       bool       `json:"error"`
    ErrorMsg    string     `json:"errorMsg,omitempty"`
    Result      any        `json:"result,omitempty"`
}

func NewResponse(result any, err error) *Response {
    var errMsg string
    var errBool bool
    if err != nil {
        errMsg = err.Error()
        errBool = true
    }
    return &Response{
        Result:     result,
        Error:      errBool,
        ErrorMsg:   errMsg,
    }
}

func (util *Util) Exec() error {
    var err error
    err = util.GetOpt()
    if err != nil {
        return err
    }
    util.URI = fmt.Sprintf("%s:%s", util.Address, util.Port)
    auth := dsrpc.CreateAuth([]byte(util.aLogin), []byte(util.aPass))

    resp := NewResponse(nil, nil)
    var result interface{}

    switch util.SubCmd {
        case getStatusCmd:
            result, err = util.GetStatusCmd(auth)

        case addUserCmd:
            result, err = util.AddUserCmd(auth)
        case checkUserCmd:
            result, err = util.CheckUserCmd(auth)
        case updateUserCmd:
            result, err = util.UpdateUserCmd(auth)
        case deleteUserCmd:
            result, err = util.DeleteUserCmd(auth)
        case listUsersCmd:
            result, err = util.ListUsersCmd(auth)

        default:
            err = errors.New("unknown cli command")
    }
    resp = NewResponse(result, err)
    respJSON, _ := json.MarshalIndent(resp, "", "  ")
    fmt.Printf("%s\n", string(respJSON))
    err = nil
    return err
}

func (util *Util) GetStatusCmd(auth *dsrpc.Auth) (*fdsapi.GetStatusResult, error) {
    var err error
    params := fdsapi.NewGetStatusParams()
    result := fdsapi.NewGetStatusResult()
    err = dsrpc.Exec(util.URI, fdsapi.GetStatusMethod, params, result, auth)
    if err != nil {
        return result, err
    }
    return result, err
}

func (util *Util) AddUserCmd(auth *dsrpc.Auth) (*fdsapi.AddUserResult, error) {
    var err error
    params := fdsapi.NewAddUserParams()
    params.Login    = util.Login
    params.Pass     = util.Pass
    result := fdsapi.NewAddUserResult()
    err = dsrpc.Exec(util.URI, fdsapi.AddUserMethod, params, result, auth)
    if err != nil {
        return result, err
    }
    return result, err
}

func (util *Util) CheckUserCmd(auth *dsrpc.Auth) (*fdsapi.CheckUserResult, error) {
    var err error
    params := fdsapi.NewCheckUserParams()
    params.Login    = util.Login
    params.Pass     = util.Pass
    result := fdsapi.NewCheckUserResult()
    err = dsrpc.Exec(util.URI, fdsapi.CheckUserMethod, params, result, auth)
    if err != nil {
        return result, err
    }
    return result, err
}

func (util *Util) UpdateUserCmd(auth *dsrpc.Auth) (*fdsapi.UpdateUserResult, error) {
    var err error
    params := fdsapi.NewUpdateUserParams()
    params.Login    = util.Login
    params.Pass     = util.Pass
    result := fdsapi.NewUpdateUserResult()
    err = dsrpc.Exec(util.URI, fdsapi.UpdateUserMethod, params, result, auth)
    if err != nil {
        return result, err
    }
    return result, err
}

func (util *Util) DeleteUserCmd(auth *dsrpc.Auth) (*fdsapi.DeleteUserResult, error) {
    var err error
    params := fdsapi.NewDeleteUserParams()
    params.Login    = util.Login
    result := fdsapi.NewDeleteUserResult()
    err = dsrpc.Exec(util.URI, fdsapi.DeleteUserMethod, params, result, auth)
    if err != nil {
        return result, err
    }
    return result, err
}

func (util *Util) ListUsersCmd(auth *dsrpc.Auth) (*fdsapi.ListUsersResult, error) {
    var err error
    params := fdsapi.NewListUsersParams()
    result := fdsapi.NewListUsersResult()
    err = dsrpc.Exec(util.URI, fdsapi.ListUsersMethod, params, result, auth)
    if err != nil {
        return result, err
    }
    return result, err
}
