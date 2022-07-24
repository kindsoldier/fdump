/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package dspack

import (
    "errors"
    "encoding/json"
    "encoding/binary"
    "bytes"
    "hash"
    "math/rand"
    "io"

    "github.com/minio/highwayhash"
)


const headerSize    int64   = 8 * 8
const sizeOfInt64   int     = 8
const magicCodeA    int64   = 0xEE00ABBA
const magicCodeB    int64   = 0xEE44ABBA

const HWHashInitSize  int64 = 32
const HWHashSumSize   int64 = 32


type Header struct {
    MagicCodeA      int64   `json:"magicCodeA"`
    HeaderVersion   int64   `json:"headerVersion"`
    DescrVersion    int64   `json:"descrVersion"`
    DescrSize       int64   `json:"descrSize"`
    BinSize         int64   `json:"binSize"`
    BinHashinitSize int64   `json:"binHashinitSize"`
    BinHashsumSize  int64   `json:"binHashsumSize"`
    MagicCodeB      int64   `json:"magicCodeB"`
}

func NewHeader() *Header {
    var header Header
    header.MagicCodeA       = magicCodeA
    header.HeaderVersion    = 1
    header.DescrVersion     = 1
    header.DescrSize        = 0
    header.BinSize          = 0
    header.BinHashinitSize  = HWHashInitSize
    header.BinHashsumSize   = HWHashSumSize
    header.MagicCodeB       = magicCodeB
    return &header
}

func (header *Header) Pack() ([]byte, error) {
    var err error

    headerBytes := make([]byte, 0, headerSize)
    headerBuffer := bytes.NewBuffer(headerBytes)

    magicCodeABytes := encoderI64(header.MagicCodeA)
    headerBuffer.Write(magicCodeABytes)

    headerVersionBytes := encoderI64(header.HeaderVersion)
    headerBuffer.Write(headerVersionBytes)

    descrVersionBytes := encoderI64(header.DescrVersion)
    headerBuffer.Write(descrVersionBytes)

    descrSizeBytes := encoderI64(header.DescrSize)
    headerBuffer.Write(descrSizeBytes)

    binSizeBytes := encoderI64(header.BinSize)
    headerBuffer.Write(binSizeBytes)

    binHashinitSizeBytes := encoderI64(header.BinHashinitSize)
    headerBuffer.Write(binHashinitSizeBytes)

    binHashsumSizeBytes := encoderI64(header.BinHashsumSize)
    headerBuffer.Write(binHashsumSizeBytes)

    magicCodeBBytes := encoderI64(header.MagicCodeB)
    headerBuffer.Write(magicCodeBBytes)

    return headerBuffer.Bytes(), err
}

func UnpackHeader(headerBytes []byte) (*Header, error) {
    var err error
    header := NewHeader()
    headerReader := bytes.NewReader(headerBytes)

    magicCodeABytes := make([]byte, sizeOfInt64)
    headerReader.Read(magicCodeABytes)
    header.MagicCodeA = decoderI64(magicCodeABytes)

    headerVersionBytes := make([]byte, sizeOfInt64)
    headerReader.Read(headerVersionBytes)
    header.HeaderVersion = decoderI64(headerVersionBytes)

    descrVersionBytes := make([]byte, sizeOfInt64)
    headerReader.Read(descrVersionBytes)
    header.HeaderVersion = decoderI64(descrVersionBytes)

    descrSizeBytes := make([]byte, sizeOfInt64)
    headerReader.Read(descrSizeBytes)
    header.DescrSize = decoderI64(descrSizeBytes)

    binSizeBytes := make([]byte, sizeOfInt64)
    headerReader.Read(binSizeBytes)
    header.BinSize = decoderI64(binSizeBytes)

    binHashinitSizeBytes := make([]byte, sizeOfInt64)
    headerReader.Read(binHashinitSizeBytes)
    header.BinHashinitSize = decoderI64(binHashinitSizeBytes)

    binHashsumSizeBytes := make([]byte, sizeOfInt64)
    headerReader.Read(binHashsumSizeBytes)
    header.BinHashsumSize = decoderI64(binHashsumSizeBytes)

    magicCodeBBytes := make([]byte, sizeOfInt64)
    headerReader.Read(magicCodeBBytes)
    header.MagicCodeB = decoderI64(magicCodeBBytes)

    if header.MagicCodeA != magicCodeA || header.MagicCodeB != magicCodeB {
        err = errors.New("wrong protocol magic code")
        return header, err
    }
    return header, err
}


const DTypeFile     int64 = 1 << 0
const DTypeSlink    int64 = 1 << 1
const DTypeDir      int64 = 1 << 4

type Descr struct {
    Path    string          `json:"path"`
    Mtime   int64           `json:"modTime"`
    Size    int64           `json:"size"`
    Mode    int64           `json:"mode"`
    Type    int64           `json:"type"`
    SLink   string          `json:"sLink,omitempty"`
    Match   bool            `json:"match"`
}

func NewDescr() *Descr {
    var descr Descr
    return &descr
}

func UnpackDescr(descrBin []byte) (*Descr, error) {
    var err error
    var descr Descr
    err = json.Unmarshal(descrBin, &descr)
    return &descr, err
}

func (descr *Descr) Pack() ([]byte, error) {
    var err error
    descrBin, err := json.Marshal(descr)
    return descrBin, err
}


type Writer struct {
    byteWriter  io.Writer
    hashInit    []byte
    hashSum     []byte
    hasher      hash.Hash
    mWriter     io.Writer
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


func (writer *Writer) WriteDescr(descr *Descr) error {
    var err error
    descrBin, err := descr.Pack()
    if err != nil {
        return err
    }
    header := NewHeader()
    header.DescrSize = int64(len(descrBin))
    header.BinSize = descr.Size
    headerBin, err := header.Pack()
    if err != nil {
        return err
    }
    _, err = writer.byteWriter.Write(headerBin)
    if err != nil {
        return err
    }
    _, err = writer.byteWriter.Write(descrBin)
    if err != nil {
        return err
    }
    return err
}

func (writer *Writer) WriteBin(buffer []byte) (int, error) {
    mWriter := io.MultiWriter(writer.byteWriter, writer.hasher)
    return mWriter.Write(buffer)
}


func (writer *Writer) WriteBinFrom(reader io.Reader, binSize int64) (int64, error) {
    var err error

    mWriter := io.MultiWriter(writer.byteWriter, writer.hasher)
    written, err := copy(reader, mWriter, binSize)
    if err != nil {
        return written, err
    }
    return written, err
}

func (writer *Writer) WriteHashInit() error {
    var err error

    writer.hashInit  = make([]byte, HWHashInitSize)
    rand.Read(writer.hashInit)

    writer.hasher, err = highwayhash.New(writer.hashInit)
    if err != nil {
        return err
    }
    _, err = writer.byteWriter.Write(writer.hashInit)
    if err != nil {
        return err
    }
    return err
}


func (writer *Writer) WriteHashSum() error {
    var err error
    hashSum := writer.hasher.Sum(nil)
    _, err = writer.byteWriter.Write(hashSum)
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
    mWriter     io.Writer
}


func NewReader(byteReader io.Reader) *Reader {
    var reader Reader
    reader.byteReader = byteReader

    reader.hashInit  = make([]byte, HWHashInitSize)
    reader.hashSum   = make([]byte, HWHashSumSize)
    reader.hasher, _ = highwayhash.New(reader.hashInit)

    return &reader
}

func (reader *Reader) NextDescr() (*Descr, error) {
    var err error
    var file *Descr

    headerBin := make([]byte, headerSize)
    _, err = reader.byteReader.Read(headerBin)
    if err != nil {
        return file, err
    }

    header, err := UnpackHeader(headerBin)
    if err != nil {
        return file, err
    }
    fileBin := make([]byte, header.DescrSize)
    _, err = reader.byteReader.Read(fileBin)
    if err != nil {
        return file, err
    }
    file, err = UnpackDescr(fileBin)
    if err != nil {
        return file, err
    }
    return file, err
}

func (reader *Reader) ReadBin(buffer []byte) (int, error) {
    var err error
    read, err := reader.byteReader.Read(buffer)
    reader.hasher.Write(buffer[0:read])
    if err != nil {
        return read, err
    }
    return read, err
}


func (reader *Reader) ReadBinTo(writer io.Writer, binSize int64) (int64, error) {
    var err error

    mWriter := io.MultiWriter(writer, reader.hasher)
    read, err := copy(reader.byteReader, mWriter, binSize)
    if err != nil {
        return read, err
    }
    return read, err
}

func (reader *Reader) ReadHashInit() error {
    var err error
    _, err = reader.byteReader.Read(reader.hashInit)
    if err != nil {
        return err
    }
    reader.hasher, err = highwayhash.New(reader.hashInit)
    if err != nil {
        return err
    }
    return err
}

func (reader *Reader) ReadHashSum() (bool, error) {
    var err error
    var match bool
    _, err = reader.byteReader.Read(reader.hashSum)
    if err != nil {
        return match, err
    }

    hashSum := reader.hasher.Sum(nil)
    if bytes.Compare(hashSum, reader.hashSum) == 0 {
        match = true
        return match, err
    }
    return match, err
}


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
