/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package fdacont

import (
    "fdump/fdagent/fdaapi"
    "fdump/dscomm/dsrpc"
    "fdump/dscomm/dserr"
)

func (contr *Contr) GetStatusHandler(context *dsrpc.Context) error {
    var err error
    params := fdaapi.NewGetStatusParams()
    err = context.BindParams(params)
    if err != nil {
        context.SendError(err)
        return dserr.Err(err)
    }
    result := fdaapi.NewGetStatusResult()
    uptime, err := contr.store.GetUptime()
    if err != nil {
        context.SendError(err)
        return dserr.Err(err)
    }
    result.SrvUptime = uptime

    diskAll, diskFree, diskUsed, err := contr.store.GetUsage()
    if err != nil {
        context.SendError(err)
        return dserr.Err(err)
    }
    result.DiskAll  = diskAll
    result.DiskFree = diskFree
    result.DiskUsed = diskUsed

    err = context.SendResult(result, 0)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}
