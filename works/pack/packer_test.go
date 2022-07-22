/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package dspack

import(
    "fmt"
    "path/filepath"
    "testing"
    "github.com/stretchr/testify/require"
)


func TestPacker01(t *testing.T) {
    var err error

    baseDir := t.TempDir()
    packPath := filepath.Join(baseDir, "test.pack")
    err = Pack("/usr/share", packPath)
    require.NoError(t, err)

    descrs, err := List(packPath)
    require.NoError(t, err)
    for _, descr := range descrs {
        fmt.Println(descr.Path)
    }
}
