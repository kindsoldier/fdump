/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package dspack

import(
    "context"
    "path/filepath"
    "testing"
    "os"
    "sync"
    "github.com/stretchr/testify/require"
)


func TestPacker01(t *testing.T) {
    var err error

    baseDir := t.TempDir()

    packPath := filepath.Join(baseDir, "test.pack")
    dirs := []string{ "/usr/bin" }

    packFile, err := os.OpenFile(packPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
    defer packFile.Close()
    require.NoError(t, err)

    err = Pack(dirs, packFile)
    require.NoError(t, err)

    packFile, err = os.OpenFile(packPath, os.O_RDONLY, 0)
    defer packFile.Close()
    require.NoError(t, err)

    var wg sync.WaitGroup
    ctx, _ := context.WithCancel(context.Background())
    descrChan := make(chan *HeadDescr, 1000)
    errChan := make(chan error, 10)

    //reporter := func() {
    //}

    //descrs, err := List(packFile, os.Stdout)
    wg.Add(1)
    go ListBG(ctx, &wg, packFile, descrChan, errChan)
    wg.Wait()
    err = <- errChan
    require.NoError(t, err)

    //for _, descr := range descrs {
    //    jsonBin, _ := json.MarshalIndent(descr, "", "    ")
    //    fmt.Println(string(jsonBin))
    //    require.Equal(t, descr.Match, true)
    //}

    packFile, err = os.OpenFile(packPath, os.O_RDONLY, 0)
    defer packFile.Close()
    require.NoError(t, err)

    _, err = Unpack(packFile, "./xxx")
    require.NoError(t, err)

    os.RemoveAll("./xxx")
}
