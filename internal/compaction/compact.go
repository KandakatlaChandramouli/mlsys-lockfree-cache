package compaction

import (
    "io"
    "os"
    "path/filepath"
    "sort"
)

func Merge(
    dir string,
    out string,
) error {

    entries, err := os.ReadDir(
        dir,
    )

    if err != nil {
        return err
    }

    files := make(
        []string,
        0,
        len(entries),
    )

    for _, e := range entries {

        if e.IsDir() {
            continue
        }

        files = append(
            files,
            filepath.Join(
                dir,
                e.Name(),
            ),
        )
    }

    sort.Strings(
        files,
    )

    dst, err := os.Create(
        out,
    )

    if err != nil {
        return err
    }

    defer dst.Close()

    for _, file := range files {

        src, err := os.Open(
            file,
        )

        if err != nil {
            return err
        }

        _, err = io.Copy(
            dst,
            src,
        )

        src.Close()

        if err != nil {
            return err
        }
    }

    return dst.Sync()
}
