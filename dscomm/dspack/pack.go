/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package dspack

import (
    "bytes"
    "encoding/json"
    "io"
    "io/fs"
    "os"
    "path/filepath"
    "strings"
    "syscall"
    "strconv"
    "os/user"
    "time"
    "sync"
    "context"
)


func Pack(dirs []string, outWriter io.Writer) error {
    var err error

    writer := NewWriter(outWriter)

    packFunc := func(filePath string, fileInfo os.FileInfo, walkErr error) error {
        var err error
        if walkErr != nil {
            return err
        }
        fileMode := fileInfo.Mode()

        if fileMode & (fs.ModeDevice|fs.ModeCharDevice) != 0  {
            return err
        }
        if fileMode & (fs.ModeNamedPipe|fs.ModeSocket|fs.ModeIrregular) != 0 {
            return err
        }

        filePath = filepath.Clean(filePath)
        if filePath == "." {
            return err
        }

        headDescr := NewHeadDescr()
        headDescr.Path  = strings.TrimLeft(filePath, "/")
        headDescr.HInit = writer.hashInit


        tailDescr := NewTailDescr()

        switch {
            case fileMode & fs.ModeDir != 0:

                var sysStat syscall.Stat_t
                err = syscall.Stat(filePath, &sysStat)

                headDescr.Mtime = sysStat.Mtimespec.Sec
                headDescr.Atime = sysStat.Atimespec.Sec
                headDescr.Ctime = sysStat.Ctimespec.Sec

                headDescr.Uid = sysStat.Uid
                headDescr.Gid = sysStat.Gid

                headDescr.HType = HashTypeNone

                uid := strconv.FormatUint(uint64(headDescr.Uid), 10)
                gid := strconv.FormatUint(uint64(headDescr.Gid), 10)

                iUser, err := user.LookupId(uid)
                if err == nil && iUser != nil {
                    headDescr.User = iUser.Username
                }
                iGroup, _ := user.LookupGroupId(gid)
                if err == nil && iGroup != nil {
                    headDescr.Group = iGroup.Name
                }

                headDescr.Type  = DTypeDir
                headDescr.Size  = 0
                headDescr.Mode  = uint32(sysStat.Mode) & 0777

                err = writer.WriteHeadDescr(headDescr)
                if err != nil {
                    return err
                }

                err = writer.WriteTailDescr(tailDescr)
                if err != nil {
                    return err
                }

            case fileMode & fs.ModeSymlink != 0:
                sLink, err := os.Readlink(filePath)
                if err != nil {
                    return err
                }

                var sysStat syscall.Stat_t
                err = syscall.Lstat(filePath, &sysStat)

                headDescr.Mtime = sysStat.Mtimespec.Sec
                headDescr.Atime = sysStat.Atimespec.Sec
                headDescr.Ctime = sysStat.Ctimespec.Sec

                headDescr.Uid = sysStat.Uid
                headDescr.Gid = sysStat.Gid

                uid := strconv.FormatUint(uint64(headDescr.Uid), 10)
                gid := strconv.FormatUint(uint64(headDescr.Gid), 10)



                iUser, err := user.LookupId(uid)
                if err == nil && iUser != nil {
                    headDescr.User = iUser.Username
                }
                iGroup, _ := user.LookupGroupId(gid)
                if err == nil && iGroup != nil {
                    headDescr.Group = iGroup.Name
                }

                headDescr.Type  = DTypeSlink
                headDescr.Size  = 0
                headDescr.Mode  = uint32(sysStat.Mode)
                headDescr.SLink = sLink

                headDescr.HType = HashTypeNone

                err = writer.WriteHeadDescr(headDescr)
                if err != nil {
                    return err
                }

                err = writer.WriteTailDescr(tailDescr)
                if err != nil {
                    return err
                }
            default:

                file, openErr := os.OpenFile(filePath, os.O_RDONLY, 0)
                defer file.Close()
                if openErr != nil {
                    return err
                }

                var sysStat syscall.Stat_t
                err = syscall.Stat(filePath, &sysStat)

                headDescr.Mtime = sysStat.Mtimespec.Sec
                headDescr.Atime = sysStat.Atimespec.Sec
                headDescr.Ctime = sysStat.Ctimespec.Sec

                headDescr.Uid = sysStat.Uid
                headDescr.Gid = sysStat.Gid

                uid := strconv.FormatUint(uint64(headDescr.Uid), 10)
                gid := strconv.FormatUint(uint64(headDescr.Gid), 10)

                iUser, err := user.LookupId(uid)
                if err == nil && iUser != nil {
                    headDescr.User = iUser.Username
                }
                iGroup, _ := user.LookupGroupId(gid)
                if err == nil && iGroup != nil {
                    headDescr.Group = iGroup.Name
                }

                headDescr.Type  = DTypeFile
                headDescr.Size  = sysStat.Size
                headDescr.Mode  = uint32(fileMode)

                headDescr.HType = HashTypeHW

                err = writer.WriteHeadDescr(headDescr)
                if err != nil {
                    return err
                }
                _, err = writer.WriteBin(file, headDescr.Size)
                if err != nil {
                    return err
                }

                tailDescr.HSum = writer.hashSum

                err = writer.WriteTailDescr(tailDescr)
                if err != nil {
                    return err
                }
        }
        return err
    }

    for _, dir := range dirs {
        err = filepath.Walk(dir, packFunc)
        if err != nil {
            return err
        }
    }
    return err
}

func Unpack(ioReader io.Reader, baseDir string) ([]*HeadDescr, error) {
    var err error
    descrs := make([]*HeadDescr, 0)

    reader := NewReader(ioReader)

    for {
        headDescr, readerErr := reader.ReadHeadDescr()

        if err == io.EOF {
            return descrs, err
        }
        if err != nil {
            return descrs, readerErr
        }
        if headDescr == nil {
            return descrs, err
        }

        filePath := strings.TrimLeft(headDescr.Path, "/")
        unpackPath := filepath.Join(baseDir, filePath)

        file, err := os.OpenFile(unpackPath, os.O_RDONLY, 0)
        if os.IsExist(err) {
            //continue
        }
        file.Close()

        switch headDescr.Type {
            case DTypeFile:

                dir := filepath.Dir(unpackPath)
                os.MkdirAll(dir, 0700)

                tmpPath := filepath.Join(unpackPath + ".tmp")

                file, err := os.OpenFile(tmpPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0640)
                defer file.Close()
                if err != nil {
                    return descrs, err
                }

                _, readErr := reader.ReadBin(file, headDescr.Size)
                if readErr == io.EOF {
                    return descrs, err
                }
                if err != nil {
                    return descrs, readErr
                }

                tailDescr, readErr := reader.ReadTailDescr()
                if readErr == io.EOF {
                    return descrs, err
                }
                if err != nil {
                    return descrs, readErr
                }
                if tailDescr == nil {
                    return descrs, err
                }

                if bytes.Compare(tailDescr.HSum, reader.hashSum) == 0 {
                    headDescr.Match = true
                }

                mTime := time.Unix(headDescr.Mtime, 0)
                aTime := time.Now()
                err = os.Chtimes(tmpPath, mTime, aTime)
                if err != nil {
                    return descrs, err
                }

                if os.Getuid() == 0 {
                    err = os.Chown(tmpPath, int(headDescr.Uid), int(headDescr.Gid))
                    if err != nil {
                        return descrs, err
                    }
                }

                err = os.Chmod(tmpPath, fs.FileMode(headDescr.Mode))
                if err != nil {
                    return descrs, err
                }

                err = os.Rename(tmpPath, unpackPath)
                if err != nil {
                    return descrs, err
                }

            case DTypeSlink:
                dir := filepath.Dir(unpackPath)
                os.MkdirAll(dir, 0750)

                err = os.Symlink(headDescr.SLink, unpackPath)
                if err != nil {
                    return descrs, err
                }

                tailDescr, readErr := reader.ReadTailDescr()
                if readErr == io.EOF {
                    return descrs, err
                }
                if err != nil {
                    return descrs, readErr
                }
                if tailDescr == nil {
                    return descrs, err
                }

                //err = syscall.Chmod(unpackPath, fs.FileMode(headDescr.Mode))
                //if err != nil {
                //    return descrs, err
                //}

                if os.Getuid() == 0 {
                    err = syscall.Lchown(unpackPath, int(headDescr.Uid), int(headDescr.Gid))
                    if err != nil {
                        return descrs, err
                    }
                }

                //mTime := time.Unix(headDescr.Mtime, 0)
                //aTime := time.Now()
                //err = os.Chtimes(unpackPath, mTime, aTime)
                //if err != nil {
                //    return descrs, err
                //}

            case DTypeDir:
                err = os.MkdirAll(unpackPath, 0750)
                if err != nil {
                    return descrs, err
                }

                tailDescr, readErr := reader.ReadTailDescr()
                if readErr == io.EOF {
                    return descrs, err
                }
                if err != nil {
                    return descrs, readErr
                }
                if tailDescr == nil {
                    return descrs, err
                }

                //err = os.Chmod(unpackPath, fs.FileMode(headDescr.Mode))
                //if err != nil {
                //    return descrs, err
                //}

                if os.Getuid() == 0 {
                    err = syscall.Chown(unpackPath, int(headDescr.Uid), int(headDescr.Gid))
                    if err != nil {
                        return descrs, err
                    }
                }
                mTime := time.Unix(headDescr.Mtime, 0)
                aTime := time.Now()
                err = os.Chtimes(unpackPath, mTime, aTime)
                if err != nil {
                    return descrs, err
                }

            default:
                _, readErr := reader.ReadBin(io.Discard, headDescr.Size)
                if readErr == io.EOF {
                    return descrs, err
                }
                if err != nil {
                    return descrs, readErr
                }

                tailDescr, readErr := reader.ReadTailDescr()
                if readErr == io.EOF {
                    return descrs, err
                }
                if err != nil {
                    return descrs, readErr
                }
                if tailDescr == nil {
                    return descrs, err
                }
                headDescr.Match = true
        }
        descrs = append(descrs, headDescr)
    }
    return descrs, err
}

func ListBG(ctx context.Context, wg *sync.WaitGroup, ioReader io.Reader, descrChan chan *HeadDescr, errChan chan error) {
    var err error
    descrs := make([]*HeadDescr, 0)

    reader := NewReader(ioReader)

    exitFunc := func() {
        errChan <- err
        wg.Done()
    }
    defer exitFunc()

    for {
        select {
            case <-ctx.Done():
                return
            default:
        }

        headDescr, err := reader.ReadHeadDescr()
        if err == io.EOF {
            return
        }
        if err != nil {
            return
        }
        if headDescr == nil {
            return
        }
        reader.hashInit = headDescr.HInit

        switch headDescr.Type {
            case DTypeFile:

                _, readErr := reader.ReadBin(io.Discard, headDescr.Size)
                if readErr == io.EOF {
                    return
                }
                if err != nil {
                    return
                }

                tailDescr, readErr := reader.ReadTailDescr()
                if readErr == io.EOF {
                    return
                }
                if err != nil {
                    return
                }
                if tailDescr == nil {
                    return
                }


                if bytes.Compare(tailDescr.HSum, reader.hashSum) == 0 {
                    headDescr.Match = true
                }

            case DTypeSlink, DTypeDir:

                headDescr.Match = true

                _, readErr := reader.ReadBin(io.Discard, headDescr.Size)
                if readErr == io.EOF {
                    return
                }
                if err != nil {
                    return
                }

                tailDescr, readErr := reader.ReadTailDescr()
                if readErr == io.EOF {
                    return
                }
                if err != nil {
                    return
                }
                if tailDescr == nil {
                    return
                }
            default:

                headDescr.Match = true

                _, readErr := reader.ReadBin(io.Discard, headDescr.Size)
                if readErr == io.EOF {
                    return
                }
                if err != nil {
                    return
                }

                tailDescr, readErr := reader.ReadTailDescr()
                if readErr == io.EOF {
                    return
                }
                if err != nil {
                    return
                }
                if tailDescr == nil {
                    return
                }

        }
        descrs = append(descrs, headDescr)
        descrChan <- headDescr
    }
    return
}



func List(ioReader io.Reader, outWriter io.Writer) ([]*HeadDescr, error) {
    var err error
    descrs := make([]*HeadDescr, 0)

    reader := NewReader(ioReader)

    for {
        headDescr, readerErr := reader.ReadHeadDescr()
        if err == io.EOF {
            return descrs, err
        }
        if err != nil {
            return descrs, readerErr
        }
        if headDescr == nil {
            return descrs, err
        }
        reader.hashInit = headDescr.HInit

        switch headDescr.Type {
            case DTypeFile:

                _, readErr := reader.ReadBin(io.Discard, headDescr.Size)
                if readErr == io.EOF {
                    return descrs, err
                }
                if err != nil {
                    return descrs, readErr
                }

                tailDescr, readErr := reader.ReadTailDescr()
                if readErr == io.EOF {
                    return descrs, err
                }
                if err != nil {
                    return descrs, readErr
                }
                if tailDescr == nil {
                    return descrs, err
                }


                if bytes.Compare(tailDescr.HSum, reader.hashSum) == 0 {
                    headDescr.Match = true
                }

            case DTypeSlink, DTypeDir:

                headDescr.Match = true

                _, readErr := reader.ReadBin(io.Discard, headDescr.Size)
                if readErr == io.EOF {
                    return descrs, err
                }
                if err != nil {
                    return descrs, readErr
                }

                tailDescr, readErr := reader.ReadTailDescr()
                if readErr == io.EOF {
                    return descrs, err
                }
                if err != nil {
                    return descrs, readErr
                }
                if tailDescr == nil {
                    return descrs, err
                }
            default:

                headDescr.Match = true

                _, readErr := reader.ReadBin(io.Discard, headDescr.Size)
                if readErr == io.EOF {
                    return descrs, err
                }
                if err != nil {
                    return descrs, readErr
                }

                tailDescr, readErr := reader.ReadTailDescr()
                if readErr == io.EOF {
                    return descrs, err
                }
                if err != nil {
                    return descrs, readErr
                }
                if tailDescr == nil {
                    return descrs, err
                }

        }
        descrs = append(descrs, headDescr)
        headDescrJson, err := json.Marshal(headDescr)
        if err != nil {
            return descrs, err
        }
        outWriter.Write(headDescrJson)
    }
    return descrs, err
}
