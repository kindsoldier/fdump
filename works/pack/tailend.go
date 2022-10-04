
/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package dspack

import (
    "errors"
    "bytes"
)


const magicCodeC    int64   = 0xAD55ACDC
const magicCodeD    int64   = 0xAD77ACDC


type Tailend struct {
    MagicCodeC      int64   `json:"magicCodeA"`
    TailendVersion  int64   `json:"tailendVersion"`
    TDescrVersion   int64   `json:"tDescrVersion"`
    TDescrSize      int64   `json:"tDescrSize"`
    MagicCodeD      int64   `json:"magicCodeB"`
}


func NewTailend() *Tailend {
    var tailend Tailend
    tailend.MagicCodeC       = magicCodeC
    tailend.TailendVersion   = 1
    tailend.TDescrVersion    = 1
    tailend.TDescrSize       = 0
    tailend.MagicCodeD       = magicCodeD
    return &tailend
}

func (tailend *Tailend) Pack() ([]byte, error) {
    var err error

    tailendBytes := make([]byte, 0, tailendSize)
    tailendBuffer := bytes.NewBuffer(tailendBytes)

    magicCodeCBytes := encoderI64(tailend.MagicCodeC)
    tailendBuffer.Write(magicCodeCBytes)

    tailendVersionBytes := encoderI64(tailend.TailendVersion)
    tailendBuffer.Write(tailendVersionBytes)

    descrVersionBytes := encoderI64(tailend.TDescrVersion)
    tailendBuffer.Write(descrVersionBytes)

    descrSizeBytes := encoderI64(tailend.TDescrSize)
    tailendBuffer.Write(descrSizeBytes)

    magicCodeDBytes := encoderI64(tailend.MagicCodeD)
    tailendBuffer.Write(magicCodeDBytes)

    return tailendBuffer.Bytes(), err
}

func UnpackTailend(tailendBytes []byte) (*Tailend, error) {
    var err error
    tailend := NewTailend()
    tailendReader := bytes.NewReader(tailendBytes)

    magicCodeCBytes := make([]byte, sizeOfInt64)
    tailendReader.Read(magicCodeCBytes)
    tailend.MagicCodeC = decoderI64(magicCodeCBytes)

    tailendVersionBytes := make([]byte, sizeOfInt64)
    tailendReader.Read(tailendVersionBytes)
    tailend.TailendVersion = decoderI64(tailendVersionBytes)

    descrVersionBytes := make([]byte, sizeOfInt64)
    tailendReader.Read(descrVersionBytes)
    tailend.TailendVersion = decoderI64(descrVersionBytes)

    descrSizeBytes := make([]byte, sizeOfInt64)
    tailendReader.Read(descrSizeBytes)
    tailend.TDescrSize = decoderI64(descrSizeBytes)

    magicCodeDBytes := make([]byte, sizeOfInt64)
    tailendReader.Read(magicCodeDBytes)
    tailend.MagicCodeD = decoderI64(magicCodeDBytes)

    if tailend.MagicCodeC != magicCodeC || tailend.MagicCodeD != magicCodeD {
        err = errors.New("wrong tailend magic code")
        return tailend, err
    }
    return tailend, err
}
