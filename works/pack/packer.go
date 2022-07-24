/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package dspack

import (
    "io/fs"
    "os"
    "path/filepath"
    "io"
    "strings"
)

func List(packPath string) ([]*Descr, error) {
    var err error
    descrs := make([]*Descr, 0)

    packFile, err := os.OpenFile(packPath, os.O_RDONLY, 0)
    defer packFile.Close()
    if err != nil {
        return descrs, err
    }

    reader := NewReader(packFile)

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


func Unpack(packPath, baseDir string) ([]*Descr, error) {
    var err error
    descrs := make([]*Descr, 0)

    packFile, err := os.OpenFile(packPath, os.O_RDONLY, 0)
    defer packFile.Close()
    if err != nil {
        return descrs, err
    }

    reader := NewReader(packFile)

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
        descrs = append(descrs, descr)

        switch descr.Type {
            case DTypeFile:
                filePath := strings.TrimLeft("/", descr.Path)
                unpackPath := filepath.Join(baseDir, filePath)

                file, err := os.OpenFile(unpackPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
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
            default:
        }
    }
    return descrs, err
}


func Pack(baseDir, packPath string) error {
    var err error

    packFile, err := os.OpenFile(packPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
    defer packFile.Close()
    if err != nil {
        return err
    }

    writer := NewWriter(packFile)

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

        filePath = strings.TrimLeft(filePath, "/")
        filePath = filepath.Clean(filePath)
        if filePath == "." {
            return err
        }

        fileMtime := fileInfo.ModTime().Unix()
        fileSize  := fileInfo.Size()

        descr := NewDescr()
        descr.Path  = filePath
        descr.Mtime = fileMtime

        switch {
            case fileMode & fs.ModeDir != 0:

                descr.Type  = DTypeDir
                descr.Size  = 0
                descr.Mode  = int64(fileMode)

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

                descr.Type  = DTypeSlink
                descr.Size  = 0
                descr.Mode  = int64(fileMode)
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
                descr.Type  = DTypeFile
                descr.Size  = fileSize
                descr.Mode  = int64(fileMode)

                err = writer.WriteDescr(descr)
                if err != nil {
                    return err
                }
                err = writer.WriteHashInit()
                if err != nil {
                    return err
                }
                _, err = writer.WriteBinFrom(file, fileSize)
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
    err = filepath.Walk(baseDir, packFunc)
    return err
}
