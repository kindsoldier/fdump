/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package dspack

import (
    "io/fs"
    "os"
    "path/filepath"
    "io"
)

func List(packPath string) ([]*Descr, error) {
    var err error
    descrs := make([]*Descr, 0)


    packFile, err := os.OpenFile(packPath, os.O_RDONLY, 0644)
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

        _, readErr := Copy(reader, io.Discard, descr.Size)
        if readErr == io.EOF {
            return descrs, err
        }
        if err != nil {
            return descrs, readErr
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
        if fileMode & (fs.ModeDir|fs.ModeDevice|fs.ModeCharDevice) != 0  {
            return err
        }
        if fileMode & (fs.ModeNamedPipe|fs.ModeSocket|fs.ModeIrregular) != 0 {
            return err
        }

        if fileMode & fs.ModeSymlink != 0 {
            return err
        }

        file, err := os.OpenFile(packPath, os.O_RDONLY, 0)
        defer file.Close()
        if err != nil {
            return err
        }

        fileSize    := fileInfo.Size()

        descr := NewDescr()
        descr.Path = filePath
        descr.Size = fileSize

        err = writer.WriteDescr(descr)
        if err != nil {
            return err
        }

        _, err = Copy(file, writer, fileSize)
        if err != nil {
            return err
        }

        return err
    }
    err = filepath.Walk(baseDir, packFunc)
    return err
}
