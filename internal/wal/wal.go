
package wal

import (
    "bufio"
    "fmt"
    "os"

    "fluxruntime/internal/segment"
)

const MaxEntries = 10000

type Log struct {
    dir     string
    segment int
    entries int

    file *os.File
    w    *bufio.Writer
}

func Open(
    dir string,
) (*Log, error) {

    if err := os.MkdirAll(
        dir,
        0755,
    ); err != nil {
        return nil, err
    }

    l := &Log{
        dir: dir,
    }

    if err := l.rotate(); err != nil {
        return nil, err
    }

    return l, nil
}

func (l *Log) rotate() error {

    if l.file != nil {

        if err := l.Sync(); err != nil {
            return err
        }

        if err := l.file.Close(); err != nil {
            return err
        }
    }

    path := segment.File(
        l.dir,
        l.segment,
    )

    f, err := os.OpenFile(
        path,
        os.O_CREATE|
            os.O_WRONLY|
            os.O_APPEND,
        0644,
    )

    if err != nil {
        return err
    }

    l.file = f

    l.w = bufio.NewWriterSize(
        f,
        1<<20,
    )

    l.entries = 0
    l.segment++

    return nil
}

func (l *Log) Append(
    id string,
) error {

    if l.entries >= MaxEntries {

        if err := l.rotate(); err != nil {
            return err
        }
    }

    _, err := fmt.Fprintf(
        l.w,
        "%s\n",
        id,
    )

    if err != nil {
        return err
    }

    l.entries++

    return nil
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
