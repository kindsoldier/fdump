/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package main

import (
    "fmt"
    "testing"

    "github.com/stretchr/testify/require"
    "fdump/dscomm/dsrpc"
)

func BenchmarkStatus(b *testing.B) {
    util := NewUtil()
    util.URI = fmt.Sprintf("%s:%s", util.Address, util.Port)

    auth := dsrpc.CreateAuth([]byte(util.aLogin), []byte(util.aPass))
    b.ResetTimer()

    pBench := func(pb *testing.PB) {
        for pb.Next() {
            _, err := util.GetStatusCmd(auth)
            require.NoError(b, err)
        }
    }
    b.SetParallelism(10)
    b.RunParallel(pBench)
}
