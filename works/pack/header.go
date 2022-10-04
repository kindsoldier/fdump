
/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package dspack

import (
    "errors"
    "bytes"
)

const magicCodeA    int64   = 0xEE00ABBA
const magicCodeB    int64   = 0xEE44ABBA

type Header struct {
    MagicCodeA      int64   `json:"magicCodeA"`
    HeaderVersion   int64   `json:"headerVersion"`
    HDescrVersion   int64   `json:"hDescrVersion"`
    HDescrSize      int64   `json:"hDescrSize"`
    BinSize         int64   `json:"binSize"`
    MagicCodeB      int64   `json:"magicCodeB"`
}

func NewHeader() *Header {
    var header Header
    header.MagicCodeA       = magicCodeA
    header.HeaderVersion    = 1
    header.HDescrVersion    = 1
    header.HDescrSize       = 0
    header.BinSize          = 0
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

    descrVersionBytes := encoderI64(header.HDescrVersion)
    headerBuffer.Write(descrVersionBytes)

    descrSizeBytes := encoderI64(header.HDescrSize)
    headerBuffer.Write(descrSizeBytes)

    binSizeBytes := encoderI64(header.BinSize)
    headerBuffer.Write(binSizeBytes)

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
    header.HDescrSize = decoderI64(descrSizeBytes)

    binSizeBytes := make([]byte, sizeOfInt64)
    headerReader.Read(binSizeBytes)
    header.BinSize = decoderI64(binSizeBytes)

    magicCodeBBytes := make([]byte, sizeOfInt64)
    headerReader.Read(magicCodeBBytes)
    header.MagicCodeB = decoderI64(magicCodeBBytes)

    if header.MagicCodeA != magicCodeA || header.MagicCodeB != magicCodeB {
        err = errors.New("wrong header magic code")
        return header, err
    }
    return header, err
}
