
/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package fdsapi

const GetStatusMethod string = "getStatus"

type GetStatusParams struct {
}

type GetStatusResult struct {
    SrvUptime   int64       `json:"srvUptime" msgpack:"srvUptime"`
    DiskFree    uint64      `json:"diskFree"  msgpack:"diskFree"`
    DiskUsed    uint64      `json:"diskUsed"  msgpack:"diskUsed"`
    DiskAll     uint64      `json:"diskAll"   msgpack:"diskAll"`
}

func NewGetStatusResult() *GetStatusResult {
    return &GetStatusResult{}
}
func NewGetStatusParams() *GetStatusParams {
    return &GetStatusParams{}
}
