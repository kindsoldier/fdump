/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package dspack

import (
    "errors"
    "encoding/binary"
    "hash"
    "math/rand"
    "io"

    "github.com/minio/highwayhash"
)


const headerSize    int64   = 8 * 6
const tailendSize   int64   = 8 * 5
const sizeOfInt64   int     = 8

const DTypeFile     int64 = 1 << 0
const DTypeSlink    int64 = 1 << 1
const DTypeDir      int64 = 1 << 2

const HWHashInitSize int64 = 32
const HWHashSumSize  int64 = 32

const HashTypeHW    string = "hw"
const HashTypeNone  string = "none"

type Writer struct {
    byteWriter  io.Writer
    hashInit    []byte
    hashSum     []byte
    hasher      hash.Hash
}

func NewWriter(byteWriter io.Writer) *Writer {
    var writer Writer
    writer.byteWriter = byteWriter

    writer.hashInit  = make([]byte, HWHashInitSize)
    rand.Read(writer.hashInit)
    writer.hashSum   = make([]byte, HWHashSumSize)
    writer.hasher, _ = highwayhash.New(writer.hashInit)

    return &writer
}

func (writer *Writer) WriteHDescr(headDescr *HDescr) error {
    var err error

    headDescrBin, err := headDescr.Pack()
    if err != nil {
        return err
    }

    header := NewHeader()
    header.HDescrSize = int64(len(headDescrBin))
    header.BinSize = headDescr.Size
    headerBin, err := header.Pack()
    if err != nil {
        return err
    }
    _, err = writer.byteWriter.Write(headerBin)
    if err != nil {
        return err
    }
    _, err = writer.byteWriter.Write(headDescrBin)
    if err != nil {
        return err
    }
    return err
}

func (writer *Writer) WriteBin(reader io.Reader, binSize int64) (int64, error) {
    var err error

    mWriter := io.MultiWriter(writer.byteWriter, writer.hasher)
    written, err := copy(reader, mWriter, binSize)
    if err != nil {
        return written, err
    }
    return written, err
}


func (writer *Writer) WriteTDescr(tailDescr *TDescr) error {
    var err error
    tailDescrBin, err := tailDescr.Pack()
    if err != nil {
        return err
    }
    tailend := NewTailend()
    tailend.TDescrSize = int64(len(tailDescrBin))

    tailendBin, err := tailend.Pack()
    if err != nil {
        return err
    }
    _, err = writer.byteWriter.Write(tailendBin)
    if err != nil {
        return err
    }
    _, err = writer.byteWriter.Write(tailDescrBin)
    if err != nil {
        return err
    }
    return err
}

type Reader struct {
    byteReader  io.Reader
    pos         int64

    hashInit    []byte
    hashSum     []byte
    hasher      hash.Hash
}

func NewReader(byteReader io.Reader) *Reader {
    var reader Reader
    reader.byteReader = byteReader

    reader.hashInit  = make([]byte, HWHashInitSize)
    reader.hashSum   = make([]byte, HWHashSumSize)
    reader.hasher, _ = highwayhash.New(reader.hashInit)

    return &reader
}

func (reader *Reader) ReadHDescr() (*HDescr, error) {
    var err error
    var headDescr *HDescr

    headerBin := make([]byte, headerSize)
    _, err = reader.byteReader.Read(headerBin)
    if err != nil {
        return headDescr, err
    }

    header, err := UnpackHeader(headerBin)
    if err != nil {
        return headDescr, err
    }
    headDescrBin := make([]byte, header.HDescrSize)
    _, err = reader.byteReader.Read(headDescrBin)
    if err != nil {
        return headDescr, err
    }
    headDescr, err = UnpackHDescr(headDescrBin)
    if err != nil {
        return headDescr, err
    }
    return headDescr, err
}

func (reader *Reader) ReadBin(writer io.Writer, binSize int64) (int64, error) {
    var err error

    mWriter := io.MultiWriter(writer, reader.hasher)
    read, err := copy(reader.byteReader, mWriter, binSize)
    if err != nil {
        return read, err
    }
    return read, err
}

func (reader *Reader) ReadTDescr() (*TDescr, error) {
    var err error
    var tailDescr *TDescr

    tailendBin := make([]byte, tailendSize)
    _, err = reader.byteReader.Read(tailendBin)
    if err != nil {
        return tailDescr, err
    }

    tailend, err := UnpackTailend(tailendBin)
    if err != nil {
        return tailDescr, err
    }
    tailDescrBin := make([]byte, tailend.TDescrSize)
    _, err = reader.byteReader.Read(tailDescrBin)
    if err != nil {
        return tailDescr, err
    }
    tailDescr, err = UnpackTDescr(tailDescrBin)
    if err != nil {
        return tailDescr, err
    }
    return tailDescr, err
}

//func (reader *Reader) ReadHashSum() (bool, error) {
//    var err error
//    var match bool
//    _, err = reader.byteReader.Read(reader.hashSum)
//    if err != nil {
//        return match, err
//    }
//
//    hashSum := reader.hasher.Sum(nil)
//    if bytes.Compare(hashSum, reader.hashSum) == 0 {
//        match = true
//        return match, err
//    }
//    return match, err
//}


func copy(reader io.Reader, writer io.Writer, size int64) (int64, error) {
    var err error
    var bufSize int64 = 1024 * 8
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

func encoderI64(i int64) []byte {
    buffer := make([]byte, sizeOfInt64)
    binary.BigEndian.PutUint64(buffer, uint64(i))
    return buffer
}

func decoderI64(b []byte) int64 {
    return int64(binary.BigEndian.Uint64(b))
}
