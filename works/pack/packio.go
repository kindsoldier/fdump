/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package dspack

import (
    "errors"
    "encoding/json"
    "encoding/binary"
    "bytes"
    "io"
)


const headerSize    int64   = 8 * 4
const sizeOfInt64   int     = 8
const magicCodeA    int64   = 0xEE00ABBA
const magicCodeB    int64   = 0xEE44ABBA


type Header struct {
    MagicCodeA  int64   `json:"magicCodeA"`
    DescrSize   int64   `json:"descrSize"`
    BinSize     int64   `json:"binSize"`
    MagicCodeB  int64   `json:"magicCodeB"`
}

func NewHeader() *Header {
    var descr Header
    descr.MagicCodeA = magicCodeA
    descr.MagicCodeB = magicCodeB
    return &descr
}

func (descr *Header) Pack() ([]byte, error) {
    var err error

    headerBytes := make([]byte, 0, headerSize)
    headerBuffer := bytes.NewBuffer(headerBytes)

    magicCodeABytes := encoderI64(descr.MagicCodeA)
    headerBuffer.Write(magicCodeABytes)

    descrSizeBytes := encoderI64(descr.DescrSize)
    headerBuffer.Write(descrSizeBytes)

    binSizeBytes := encoderI64(descr.BinSize)
    headerBuffer.Write(binSizeBytes)

    magicCodeBBytes := encoderI64(descr.MagicCodeB)
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

    descrSizeBytes := make([]byte, sizeOfInt64)
    headerReader.Read(descrSizeBytes)
    header.DescrSize = decoderI64(descrSizeBytes)

    binSizeBytes := make([]byte, sizeOfInt64)
    headerReader.Read(binSizeBytes)
    header.BinSize = decoderI64(binSizeBytes)

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
}


func NewWriter(byteWriter io.Writer) *Writer {
    var writer Writer
    writer.byteWriter = byteWriter
    return &writer
}


func (writer *Writer) Write(data []byte) (int, error) {
    return writer.byteWriter.Write(data)
}

func (writer *Writer) WriteDescr(file *Descr) error {
    var err error
    fileBin, err := file.Pack()
    if err != nil {
        return err
    }
    header := NewHeader()
    header.DescrSize = int64(len(fileBin))
    header.BinSize = file.Size
    headerBin, err := header.Pack()
    if err != nil {
        return err
    }
    _, err = writer.byteWriter.Write(headerBin)
    if err != nil {
        return err
    }
    _, err = writer.byteWriter.Write(fileBin)
    if err != nil {
        return err
    }
    return err
}



type Reader struct {
    byteReader  io.Reader
    pos         int64
}


func NewReader(byteReader io.Reader) *Reader {
    var reader Reader
    reader.byteReader = byteReader
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

func (reader *Reader) Read(data []byte) (int, error) {
    return reader.byteReader.Read(data)
}

func Copy(reader io.Reader, writer io.Writer, size int64) (int64, error) {
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
