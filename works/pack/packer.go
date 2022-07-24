/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package dspack

import (
    "io"
    "io/fs"
    "os"
    "path/filepath"
    "strings"
    "syscall"
    "strconv"
    "os/user"
    "time"
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

        descr := NewDescr()
        descr.Path  = strings.TrimLeft(filePath, "/")

        switch {
            case fileMode & fs.ModeDir != 0:

                var sysStat syscall.Stat_t
                err = syscall.Stat(filePath, &sysStat)

                descr.Mtime = sysStat.Mtimespec.Sec
                descr.Atime = sysStat.Atimespec.Sec
                descr.Ctime = sysStat.Ctimespec.Sec

                descr.Uid = sysStat.Uid
                descr.Gid = sysStat.Gid

                uid := strconv.FormatUint(uint64(descr.Uid), 10)
                gid := strconv.FormatUint(uint64(descr.Gid), 10)

                iUser, err := user.LookupId(uid)
                if err == nil && iUser != nil {
                    descr.User = iUser.Username
                }
                iGroup, _ := user.LookupGroupId(gid)
                if err == nil && iGroup != nil {
                    descr.Group = iGroup.Name
                }

                descr.Type  = DTypeDir
                descr.Size  = 0
                descr.Mode  = uint32(sysStat.Mode) & 0777

                err = writer.WriteDescr(descr)
                if err != nil {
                    return err
                }
                err = writer.WriteHashInit()
                if err != nil {
                    return err
                }
                err = writer.WriteHashSum()
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

                descr.Mtime = sysStat.Mtimespec.Sec
                descr.Atime = sysStat.Atimespec.Sec
                descr.Ctime = sysStat.Ctimespec.Sec

                descr.Uid = sysStat.Uid
                descr.Gid = sysStat.Gid

                uid := strconv.FormatUint(uint64(descr.Uid), 10)
                gid := strconv.FormatUint(uint64(descr.Gid), 10)

                iUser, err := user.LookupId(uid)
                if err == nil && iUser != nil {
                    descr.User = iUser.Username
                }
                iGroup, _ := user.LookupGroupId(gid)
                if err == nil && iGroup != nil {
                    descr.Group = iGroup.Name
                }

                descr.Type  = DTypeSlink
                descr.Size  = 0
                descr.Mode  = uint32(sysStat.Mode)
                descr.SLink = sLink
                err = writer.WriteDescr(descr)
                if err != nil {
                    return err
                }
                err = writer.WriteHashInit()
                if err != nil {
                    return err
                }
                err = writer.WriteHashSum()
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

                descr.Mtime = sysStat.Mtimespec.Sec
                descr.Atime = sysStat.Atimespec.Sec
                descr.Ctime = sysStat.Ctimespec.Sec

                descr.Uid = sysStat.Uid
                descr.Gid = sysStat.Gid

                uid := strconv.FormatUint(uint64(descr.Uid), 10)
                gid := strconv.FormatUint(uint64(descr.Gid), 10)

                iUser, err := user.LookupId(uid)
                if err == nil && iUser != nil {
                    descr.User = iUser.Username
                }
                iGroup, _ := user.LookupGroupId(gid)
                if err == nil && iGroup != nil {
                    descr.Group = iGroup.Name
                }

                descr.Type  = DTypeFile
                descr.Size  = sysStat.Size
                descr.Mode  = uint32(fileMode)

                err = writer.WriteDescr(descr)
                if err != nil {
                    return err
                }
                err = writer.WriteHashInit()
                if err != nil {
                    return err
                }
                _, err = writer.WriteBinFrom(file, descr.Size)
                if err != nil {
                    return err
                }
                err = writer.WriteHashSum()
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


func Unpack(outReader io.Reader, baseDir string) ([]*Descr, error) {
    var err error
    descrs := make([]*Descr, 0)

    reader := NewReader(outReader)

    for {
        descr, readerErr := reader.NextDescr()
        if err == io.EOF {
            return descrs, err
        }
        if err != nil {
            return descrs, readerErr
        }
        if descr == nil {
            return descrs, err
        }

        filePath := strings.TrimLeft(descr.Path, "/")
        unpackPath := filepath.Join(baseDir, filePath)

        switch descr.Type {
            case DTypeFile:

                dir := filepath.Dir(unpackPath)
                os.MkdirAll(dir, 0700)

                file, err := os.OpenFile(unpackPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0640)
                defer file.Close()
                if err != nil {
                    return descrs, err
                }
                readErr := reader.ReadHashInit()
                if readErr == io.EOF {
                    return descrs, err
                }
                _, readErr = reader.ReadBinTo(file, descr.Size)
                if readErr == io.EOF {
                    return descrs, err
                }
                if err != nil {
                    return descrs, readErr
                }
                _, readErr = reader.ReadHashSum()
                if readErr == io.EOF {
                    return descrs, err
                }
                if err != nil {
                    return descrs, readErr
                }

                mTime := time.Unix(descr.Mtime, 0)
                aTime := time.Now()
                err = os.Chtimes(unpackPath, mTime, aTime)
                if err != nil {
                    return descrs, err
                }

                if os.Getuid() == 0 {
                    err = os.Chown(unpackPath, int(descr.Uid), int(descr.Gid))
                    if err != nil {
                        return descrs, err
                    }
                }

                err = os.Chmod(unpackPath, fs.FileMode(descr.Mode))
                if err != nil {
                    return descrs, err
                }

            case DTypeSlink:
                dir := filepath.Dir(unpackPath)
                os.MkdirAll(dir, 0750)

                err = os.Symlink(descr.SLink, unpackPath)
                if err != nil {
                    return descrs, err
                }
                readErr := reader.ReadHashInit()
                if readErr == io.EOF {
                    return descrs, err
                }
                _, readErr = reader.ReadHashSum()
                if readErr == io.EOF {
                    return descrs, err
                }
                if err != nil {
                    return descrs, readErr
                }

                //err = syscall.Chmod(unpackPath, fs.FileMode(descr.Mode))
                //if err != nil {
                //    return descrs, err
                //}

                if os.Getuid() == 0 {
                    err = syscall.Lchown(unpackPath, int(descr.Uid), int(descr.Gid))
                    if err != nil {
                        return descrs, err
                    }
                }

                //mTime := time.Unix(descr.Mtime, 0)
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

                readErr := reader.ReadHashInit()
                if readErr == io.EOF {
                    return descrs, err
                }
                _, readErr = reader.ReadHashSum()
                if readErr == io.EOF {
                    return descrs, err
                }
                if err != nil {
                    return descrs, readErr
                }

                //err = os.Chmod(unpackPath, fs.FileMode(descr.Mode))
                //if err != nil {
                //    return descrs, err
                //}

                if os.Getuid() == 0 {
                    err = syscall.Chown(unpackPath, int(descr.Uid), int(descr.Gid))
                    if err != nil {
                        return descrs, err
                    }
                }
                mTime := time.Unix(descr.Mtime, 0)
                aTime := time.Now()
                err = os.Chtimes(unpackPath, mTime, aTime)
                if err != nil {
                    return descrs, err
                }

            default:
                readErr := reader.ReadHashInit()
                if readErr == io.EOF {
                    return descrs, err
                }
                _, readErr = reader.ReadBinTo(io.Discard, descr.Size)
                if readErr == io.EOF {
                    return descrs, err
                }
                if err != nil {
                    return descrs, readErr
                }
                _, readErr = reader.ReadHashSum()
                if readErr == io.EOF {
                    return descrs, err
                }
                if err != nil {
                    return descrs, readErr
                }

        }
        descrs = append(descrs, descr)
    }
    return descrs, err
}

func List(outReader io.Reader) ([]*Descr, error) {
    var err error
    descrs := make([]*Descr, 0)

    reader := NewReader(outReader)

    for {
        descr, readerErr := reader.NextDescr()
        if err == io.EOF {
            return descrs, err
        }
        if err != nil {
            return descrs, readerErr
        }
        if descr == nil {
            return descrs, err
        }
        descr.Match = false

        readErr := reader.ReadHashInit()
        if readErr == io.EOF {
            return descrs, err
        }

        _, readErr = reader.ReadBinTo(io.Discard, descr.Size)
        if readErr == io.EOF {
            return descrs, err
        }
        if err != nil {
            return descrs, readErr
        }
        match, readErr := reader.ReadHashSum()
        if readErr == io.EOF {
            return descrs, err
        }
        if err != nil {
            return descrs, readErr
        }
        descr.Match = match
        descrs = append(descrs, descr)
    }
    return descrs, err
}
