/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package dsrpc

import (
    "encoding/json"
)

type Packet struct {
    header      []byte
    rcpPayload  []byte
}

func NewPacket() *Packet {
    return &Packet{}
}

func (pkt *Packet) JSON() []byte {
    jBytes, _ := json.Marshal(pkt)
    return jBytes
}
