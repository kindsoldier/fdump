/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package dspack

import(
    "encoding/json"
    "fmt"
    "path/filepath"
    "testing"
    "github.com/stretchr/testify/require"
)


func TestPacker01(t *testing.T) {
    var err error

    baseDir := t.TempDir()
    packPath := filepath.Join(baseDir, "test.pack")
    err = Pack("./", packPath)
    require.NoError(t, err)

    descrs, err := List(packPath)
    require.NoError(t, err)
    for _, descr := range descrs {
        jsonBin, _ := json.MarshalIndent(descr, "", "    ")
        fmt.Println(string(jsonBin))
    }
}
