/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package dscron

import(
    "testing"
    "time"
    "github.com/stretchr/testify/require"
)


func TestExpand01(t *testing.T) {
    start := 2
    end := 31

    resMap := Expander("*", start, end)
    for i := start; i <= end ; i++ {
        require.Equal(t, resMap[i], true)
    }

    resMap = Expander("3-6,8,*/15", start, end)
    require.Equal(t, resMap[2], false)
    require.Equal(t, resMap[3], true)
    require.Equal(t, resMap[4], true)
    require.Equal(t, resMap[5], true)
    require.Equal(t, resMap[6], true)
    require.Equal(t, resMap[8], true)
    require.Equal(t, resMap[9], false)

    require.Equal(t, resMap[15], true)
    require.Equal(t, resMap[30], true)
    require.Equal(t, resMap[31], false)
}


func TestMatch01(t *testing.T) {

    ts, err := time.Parse(time.RFC3339, "2006-01-02T15:04:05+07:00")
    require.NoError(t, err)

    var match bool

    match = Match("*", "*", "10,18,23", "*", ts)
    require.Equal(t, match, false)

    match = Match("*", "*", "10,18,23,15", "*", ts)
    require.Equal(t, match, true)

    match = Match("*", "*", "10,18,23,15", "1-5", ts)
    require.Equal(t, match, true)

    match = Match("*", "*", "*", "*", ts)
    require.Equal(t, match, true)

    match = Match("*", "*", "*", "*/3", ts)
    require.Equal(t, match, false)


}
