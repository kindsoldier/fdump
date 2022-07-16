/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package main

import (
    "flag"
    "fmt"
    "io/fs"
    "os"
    "os/signal"
    "os/user"
    "path/filepath"
    "strconv"
    "syscall"
    "io"

    "fdump/fdstore/fdsapi"
    "fdump/fdstore/fdssrv/fdscont"
    "fdump/fdstore/fdssrv/fdsreg"
    "fdump/fdstore/fdssrv/fdstore"

    "fdump/dscomm/dskvdb"
    "fdump/dscomm/dslog"
    "fdump/dscomm/dsrpc"
    "fdump/dscomm/dserr"
)

const successExit   int = 0
const errorExit     int = 1

func main() {
    var err error
    server := NewServer()

    dserr.SetDevelMode(false)
    dserr.SetDebugMode(false)

    err = server.Execute()
    if err != nil {
        dslog.LogError("config error:", err)
        os.Exit(errorExit)
    }
}

type Server struct {
    Params  *Config
    Backgr  bool
}

func (server *Server) Execute() error {
    var err error

    err = server.ReadConf()
    if err != nil {
        return err
    }
    err = server.GetOptions()
    if err != nil {
        return err
    }

    err = server.PrepareEnv()
    if err != nil {
        return err
    }
    if server.Backgr {
        err = server.ForkCmd()
        if err != nil {
            return err
        }
        err = server.CloseIO()
        if err != nil {
            return err
        }

        err = server.ChangeUid()
        if err != nil {
            return err
        }
    }
    err = server.RedirLog()
    if err != nil {
        return err
    }
    err = server.SavePid()
    if err != nil {
        return err
    }
    err = server.SetSHandler()
    if err != nil {
        return err
    }

    err = server.RunService()
    if err != nil {
        return err
    }
    return err
}


func NewServer() *Server {
    var server Server
    server.Params = NewConfig()
    server.Backgr = false
    return &server
}

func (server *Server) ReadConf() error {
    var err error
    err = server.Params.Read()
    if err != nil {
        return err
    }
    return err
}

func (server *Server) GetOptions() error {
    var err error
    exeName := filepath.Base(os.Args[0])

    flag.StringVar(&server.Params.RunDir, "runDir", server.Params.RunDir, "run direcory")
    flag.StringVar(&server.Params.LogDir, "logDir", server.Params.LogDir, "log direcory")
    flag.StringVar(&server.Params.DataDir, "dataDir", server.Params.DataDir, "data directory")

    flag.StringVar(&server.Params.Port, "port", server.Params.Port, "listen port")
    flag.BoolVar(&server.Backgr, "daemon", server.Backgr, "run as daemon")

    help := func() {
        fmt.Println("")
        fmt.Printf("usage: %s [option]\n", exeName)
        fmt.Println("")
        fmt.Println("options:")
        flag.PrintDefaults()
        fmt.Println("")
    }
    flag.Usage = help
    flag.Parse()

    return err
}


func (server *Server) ForkCmd() error {
    const successExit int = 0
    var keyEnv string = "IMX0LTSELMRF8K"
    var err error

    _, isChild := os.LookupEnv(keyEnv)
    switch  {
        case !isChild:
            os.Setenv(keyEnv, "TRUE")

            procAttr := syscall.ProcAttr{}
            cwd, err := os.Getwd()
            if err != nil {
                    return err
            }
            var sysFiles = make([]uintptr, 3)
            sysFiles[0] = uintptr(syscall.Stdin)
            sysFiles[1] = uintptr(syscall.Stdout)
            sysFiles[2] = uintptr(syscall.Stderr)

            procAttr.Files = sysFiles
            procAttr.Env = os.Environ()
            procAttr.Dir = cwd

            _, err = syscall.ForkExec(os.Args[0], os.Args, &procAttr)
            if err != nil {
                return err
            }
            os.Exit(successExit)
        case isChild:
            _, err = syscall.Setsid()
            if err != nil {
                    return err
            }
    }
    os.Unsetenv(keyEnv)
    return err
}


func (server *Server) ChangeUid() error {
    var err error

    username := server.Params.SrvUser

    currUid := syscall.Getuid()
    if currUid != 0 {
        err = fmt.Errorf("impossible to change uid for non-root")
        return err
    }

    userDescr, err := user.Lookup(username)
    if err != nil {
        err = fmt.Errorf("no username %s found, err: %v", username, err)
        return err
    }

    newGid, err := strconv.Atoi(userDescr.Gid)
    if err != nil {
        err = fmt.Errorf("cannot convert gid, err: %v", err)
        return err
    }
    err = syscall.Setgid(newGid)
    if err != nil {
        err = fmt.Errorf("cannot change gid, err: %v", err)
        return err
    }
    currGid := syscall.Getgid()
    if currGid != newGid {
        err = fmt.Errorf("unable to change gid for unknown reason")
        return err
    }

    newUid, err := strconv.Atoi(userDescr.Uid)
    if err != nil {
        err = fmt.Errorf("cannot convert uid, err: %v", err)
        return err
    }

    runDir := server.Params.RunDir
    err = os.Chown(runDir, newUid, newGid)
    if err != nil {
            return err
    }

    logDir := server.Params.LogDir
    err = os.Chown(logDir, newUid, newGid)
    if err != nil {
            return err
    }

    dataDir := server.Params.DataDir
    err = os.Chown(dataDir, newUid, newGid)
    if err != nil {
            return err
    }

    err = syscall.Setuid(newUid)
    if err != nil {
        err = fmt.Errorf("cannot change uid, err: %v", err)
        return err
    }

    currUid = syscall.Getuid()
    if currUid != newUid {
        err = fmt.Errorf("unable to change uid for unknown reason")
        return err
    }
    return err
}



func (server *Server) PrepareEnv() error {
    var err error

    var runDirPerm fs.FileMode = server.Params.DirPerm
    var logDirPerm fs.FileMode = server.Params.DirPerm
    var dataDirPerm fs.FileMode = server.Params.DirPerm

    runDir := server.Params.RunDir
    err = os.MkdirAll(runDir, runDirPerm)
    if err != nil {
            return err
    }
    err = os.Chmod(runDir, runDirPerm)
    if err != nil {
            return err
    }

    logDir := server.Params.LogDir
    err = os.MkdirAll(logDir, logDirPerm)
    if err != nil {
            return err
    }
    err = os.Chmod(logDir, logDirPerm)
    if err != nil {
            return err
    }

    dataDir := server.Params.DataDir
    err = os.MkdirAll(dataDir, dataDirPerm)
    if err != nil {
            return err
    }
    err = os.Chmod(dataDir, dataDirPerm)
    if err != nil {
            return err
    }
    return err
}

func (server *Server) SavePid() error {
    var err error

    var pidFilePerm fs.FileMode = server.Params.DirPerm

    pidFile := filepath.Join(server.Params.RunDir, server.Params.PidName)

    openMode := os.O_WRONLY | os.O_CREATE | os.O_TRUNC
    file, err := os.OpenFile(pidFile, openMode, pidFilePerm)
    if err != nil {
            return err
    }
    defer file.Close()

    pid := os.Getpid()
    _, err = file.WriteString(strconv.Itoa(pid))
    if err != nil {
            return err
    }
    err = os.Chmod(pidFile, pidFilePerm)
    if err != nil {
            return err
    }
    file.Sync()
    return err
}

func (server *Server) RedirLog() error {
    var err error

    var logFilePerm fs.FileMode = server.Params.FilePerm

    logOpenMode := os.O_WRONLY|os.O_CREATE|os.O_APPEND
    msgFileName := filepath.Join(server.Params.LogDir, server.Params.MsgName)
    msgFile, err := os.OpenFile(msgFileName, logOpenMode, logFilePerm)
    if err != nil {
            return err
    }

    logWriter := io.MultiWriter(os.Stdout, msgFile)
    dslog.SetOutput(logWriter)
    dsrpc.SetMessageWriter(logWriter)

    accFileName := filepath.Join(server.Params.LogDir, server.Params.AccName)
    accFile, err := os.OpenFile(accFileName, logOpenMode, logFilePerm)
    if err != nil {
            return err
    }

    accWriter := io.MultiWriter(os.Stdout, accFile)
    dsrpc.SetAccessWriter(accWriter)
    return err
}

func (server *Server) CloseIO() error {
    var err error
    file, err := os.OpenFile("/dev/null", os.O_RDWR, 0)
    if err != nil {
            return err
    }
    err = syscall.Dup2(int(file.Fd()), int(os.Stdin.Fd()))
    if err != nil {
            return err
    }
    err = syscall.Dup2(int(file.Fd()), int(os.Stdout.Fd()))
    if err != nil {
            return err
    }
    err = syscall.Dup2(int(file.Fd()), int(os.Stderr.Fd()))
    if err != nil {
            return err
    }
    return err
}


func (server *Server) SetSHandler() error {
    var err error
    sigs := make(chan os.Signal, 1)
    signal.Notify(sigs, syscall.SIGHUP, syscall.SIGINT, syscall.SIGSTOP,
                                    syscall.SIGTERM, syscall.SIGQUIT)

    handler := func() {
        var err error
        for {
            dslog.LogInfo("signal handler start")
            sig := <-sigs
            dslog.LogInfo("received signal", sig.String())

            switch sig {
                case syscall.SIGINT, syscall.SIGTERM, syscall.SIGSTOP:
                    dslog.LogInfo("exit by signal", sig.String())
                    server.StopAll()
                    os.Exit(successExit)

                case syscall.SIGHUP:
                    switch {
                        case server.Backgr:
                            dslog.LogInfo("restart server")

                            err = server.StopAll()
                            if err != nil {
                                dslog.LogError("stop all error:", err)
                            }
                            err = server.ForkCmd()
                            if err != nil {
                                dslog.LogError("fork error:", err)
                            }
                        default:
                            server.StopAll()
                            os.Exit(successExit)
                    }
            }
        }
    }
    go handler()
    return err
}

func (server *Server) StopAll() error {
    var err error
    dslog.LogInfo("stop processes")
    return err
}

func (server *Server) RunService() error {
    var err error

    filePerm    := server.Params.FilePerm
    dirPerm     := server.Params.DirPerm
    dataDir     := server.Params.DataDir

    develMode   := false //server.Params.DevelMode
    debugMode   := false //server.Params.DebugMode

    dslog.SetDebugMode(debugMode)

    dserr.SetDevelMode(develMode)
    dserr.SetDebugMode(debugMode)

    dsrpc.SetDevelMode(develMode)
    dsrpc.SetDebugMode(debugMode)

    db, err := dskvdb.OpenDB(dataDir, "storedb")
    if err != nil {
        return err
    }
    reg, err := fdsreg.NewReg(db)
    if err != nil {
        return err
    }
    store, err := fdstore.NewStore(dataDir, reg)
    if err != nil {
        return err
    }
    store.SetFilePerm(filePerm)
    store.SetDirPerm(dirPerm)

    err = store.SeedUsers()
    if err != nil {
        return err
    }
    contr, err := fdscont.NewContr(store)
    if err != nil {
        return err
    }

    dslog.LogInfof("dataDir is %s", server.Params.DataDir)
    dslog.LogInfof("logDir is %s", server.Params.LogDir)
    dslog.LogInfof("runDir is %s", server.Params.RunDir)

    serv := dsrpc.NewService()

    if debugMode || develMode {
        serv.PreMiddleware(dsrpc.LogRequest)
    }
    serv.PreMiddleware(contr.AuthMidware(debugMode))

    serv.Handler(fdsapi.AddUserMethod, contr.AddUserHandler)
    serv.Handler(fdsapi.CheckUserMethod, contr.CheckUserHandler)
    serv.Handler(fdsapi.UpdateUserMethod, contr.UpdateUserHandler)
    serv.Handler(fdsapi.ListUsersMethod, contr.ListUsersHandler)
    serv.Handler(fdsapi.DeleteUserMethod, contr.DeleteUserHandler)

    serv.Handler(fdsapi.GetStatusMethod, contr.GetStatusHandler)


    if debugMode || develMode {
        serv.PostMiddleware(dsrpc.LogResponse)
    }
    serv.PostMiddleware(dsrpc.LogAccess)

    listenParam := fmt.Sprintf(":%s", server.Params.Port)
    err = serv.Listen(listenParam)
    if err != nil {
        return err
    }
    return err
}
