/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package dstar

import (
    "archive/tar"
    "errors"
    "io"
    "io/fs"
    "os"
    "path/filepath"
    "runtime"
    "syscall"
    "time"
)


func Atime(fileInfo os.FileInfo) time.Time {
    var aTime time.Time
    switch runtime.GOOS {
        case "freebsd":
            ts := fileInfo.Sys().(*syscall.Stat_t).Atimespec
            aTime = time.Unix(int64(ts.Sec), int64(ts.Nsec))
        default:
    }
    return aTime
}


func Tar(baseDir, tarPath string) error {
    var err error

    tarFile, err := os.OpenFile(tarPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
    defer tarFile.Close()
    if err != nil {
        return err
    }

    tarWriter := tar.NewWriter(tarFile)

    tarFunc := func(filePath string, fileInfo os.FileInfo, walkErr error) error {
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

        file, err := os.OpenFile(tarPath, os.O_RDONLY, 0)
        defer file.Close()
        if err != nil {
            return err
        }

        fileSize    := fileInfo.Size()
        fileModtime := fileInfo.ModTime()
        fileAtime   := Atime(fileInfo)

        tarHeader := tar.Header{
            Format:     tar.FormatGNU,
            Name:       filePath,
            Size:       fileSize,
            Mode:       int64(fileMode),
            AccessTime: fileAtime,
            ChangeTime: fileModtime,
            ModTime:    fileModtime,
        }

        err = tarWriter.WriteHeader(&tarHeader)
        if err != nil {
            return err
        }

        _, err = copyData(file, tarWriter, fileSize)
        if err != nil {
            return err
        }

        return err
    }

    err = filepath.Walk(baseDir, tarFunc)

    return err
}


func copyData(reader io.Reader, writer io.Writer, size int64) (int64, error) {
    var err error
    var bufSize int64 = 1024 * 16
    var total   int64 = 0
    var remains int64 = size
    buffer := make([]byte, bufSize)

    for {
        if remains == 0 {
            return total, err
        }
        if remains < bufSize {
            bufSize = remains
        }
        received, err := reader.Read(buffer[0:bufSize])
        if err != nil {
            return total, err
        }
        written, err := writer.Write(buffer[0:received])
        if err != nil {
            return total, err
        }
        if written != received {
            err = errors.New("write error")
            return total, err
        }
        total += int64(written)
        remains -= int64(written)
    }
    return total, err
}
