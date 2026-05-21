
package wal

import (
    "bufio"
    "fmt"
    "os"
)

type Log struct {
    file *os.File
    w    *bufio.Writer
}

func Open(
    path string,
) (*Log, error) {

    f, err := os.OpenFile(
        path,
        os.O_CREATE|
            os.O_WRONLY|
            os.O_APPEND,
        0644,
    )

    if err != nil {
        return nil, err
    }

    return &Log{
        file: f,
        w: bufio.NewWriterSize(
            f,
            1<<20,
        ),
    }, nil
}

func (l *Log) Append(
    id string,
) error {

    _, err := fmt.Fprintf(
        l.w,
        "%s\n",
        id,
    )

    return err
}

func (l *Log) Sync() error {

    if err := l.w.Flush(); err != nil {
        return err
    }

    return l.file.Sync()
}

func (l *Log) Close() error {

    if err := l.Sync(); err != nil {
        return err
    }

    return l.file.Close()
}
