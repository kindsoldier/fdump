/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package dstar

import(
    "path/filepath"
    "testing"
    "github.com/stretchr/testify/require"

)

func TestUser01(t *testing.T) {
    var err error

    baseDir := t.TempDir()
    tarPath := filepath.Join(baseDir, "test.tar")
    err = Tar("/usr/share", tarPath)
    require.NoError(t, err)
}
