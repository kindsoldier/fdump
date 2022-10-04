/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package dspack

import(
    "encoding/json"
    "fmt"
    "path/filepath"
    "testing"
    "os"
    "github.com/stretchr/testify/require"
)


func TestPacker01(t *testing.T) {
    var err error

    baseDir := t.TempDir()

    packPath := filepath.Join(baseDir, "test.pack")
    dirs := []string{ "/usr/bin" }

    packFile, err := os.OpenFile(packPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
    defer packFile.Close()
    require.NoError(t, err)


    err = Pack(dirs, packFile)
    require.NoError(t, err)

    packFile, err = os.OpenFile(packPath, os.O_RDONLY, 0)
    defer packFile.Close()
    require.NoError(t, err)


    descrs, err := List(packFile)
    require.NoError(t, err)

    for _, descr := range descrs {
        jsonBin, _ := json.MarshalIndent(descr, "", "    ")
        fmt.Println(string(jsonBin))
        require.Equal(t, descr.Match, false)
    }

    packFile, err = os.OpenFile(packPath, os.O_RDONLY, 0)
    defer packFile.Close()
    require.NoError(t, err)

    descrs, err = Unpack(packFile, "./xxx")
    require.NoError(t, err)
}
