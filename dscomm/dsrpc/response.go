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


type Response struct {
    Error   string      `json:"error"   msgpack:"error"`
    Result  any         `json:"result"  msgpack:"result"`
}

func NewResponse() *Response {
    return &Response{}
}

func (this *Response) JSON() []byte {
    jBytes, _ := json.Marshal(this)
    return jBytes
}

func (this *Response) Pack() ([]byte, error) {
    rBytes, err := encoder.Marshal(this)
    return rBytes, Err(err)
}
