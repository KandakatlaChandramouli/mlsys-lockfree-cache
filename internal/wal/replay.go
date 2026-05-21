package wal

import (
    "bufio"
    "os"
)

func Replay(
    path string,
    fn func(id string),
) error {

    f, err := os.Open(
        path,
    )

    if err != nil {
        return err
    }

    defer f.Close()

    scanner := bufio.NewScanner(
        f,
    )

    for scanner.Scan() {

        fn(
            scanner.Text(),
        )
    }

    return scanner.Err()
}
