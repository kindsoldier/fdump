/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package main

import (
    "encoding/json"
    "fmt"
    "io/fs"
    "flag"
    "os"
    "path/filepath"
    "errors"
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
    SubCmd      string
    FilePath    string
}

func NewUtil() *Util {
    var util Util
    return &util
}

const packCmd      string = "pack"
const unpackCmd    string = "unpack"
const listCmd      string = "list"
const helpCmd      string = "help"


func (util *Util) GetOpt() error {
    var err error

    exeName := filepath.Base(os.Args[0])

    help := func() {
        fmt.Println("")
        fmt.Printf("Usage: %s [option] command [command option]\n", exeName)
        fmt.Printf("\n")
        fmt.Printf("Command list: help, pack, unpack, list\n")

        //fmt.Printf("\n")
        //fmt.Printf("Global options:\n")
        //flag.PrintDefaults()
        //fmt.Printf("\n")
    }
    flag.Usage = help
    flag.Parse()

    args := flag.Args()

    //if len(args) == 0 {
    //    args = append(args, packCmd)
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
            util.SubCmd = subCmd
        case packCmd:
            flagSet := flag.NewFlagSet(packCmd, flag.ExitOnError)
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

        case unpackCmd:
            flagSet := flag.NewFlagSet(unpackCmd, flag.ExitOnError)
            flagSet.StringVar(&util.FilePath, "file", util.FilePath, "pack file name")
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

        case listCmd:
            flagSet := flag.NewFlagSet(listCmd, flag.ExitOnError)

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

    resp := NewResponse(nil, nil)
    var result interface{}

    switch util.SubCmd {
        case packCmd:
            result, err = util.PackCmd()
        case unpackCmd:
            result, err = util.UnpackCmd()
        case listCmd:
            result, err = util.ListCmd()
        case helpCmd:
            return err
        default:
            err = errors.New("unknown cli command")
    }
    resp = NewResponse(result, err)
    respJSON, _ := json.MarshalIndent(resp, "", "  ")
    fmt.Printf("%s\n", string(respJSON))
    err = nil
    return err
}

const dirPerm   fs.FileMode = 0755
const filePerm  fs.FileMode = 0644

type PackResult struct {

}

type UnpackResult struct {
}

type ListResult struct {
}


func (util *Util) PackCmd() (*PackResult, error) {
    var err error
    var result PackResult
    return &result, err
}

func (util *Util) UnpackCmd() (*UnpackResult, error) {
    var err error
    var result UnpackResult
    return &result, err
}

func (util *Util) ListCmd() (*ListResult, error) {
    var err error
    var result ListResult
    return &result, err
}
