/*
 *
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 *
 */

package dsrpc

import (
    "encoding/json"
    encoder "github.com/vmihailenco/msgpack/v5"
)

type Request struct {
    Method  string      `json:"method"            msgpack:"method"`
    Params  any         `json:"params,omitempty"  msgpack:"params"`
    Auth    *Auth       `json:"auth,omitempty"    msgpack:"auth"`
}

func NewRequest() *Request {
    req := &Request{}
    req.Auth = &Auth{}
    return req
}

func (this *Request) Pack() ([]byte, error) {
    rBytes, err := encoder.Marshal(this)
    return rBytes, Err(err)
}

func (this *Request) JSON() []byte {
    jBytes, _ := json.Marshal(this)
    return jBytes
}
