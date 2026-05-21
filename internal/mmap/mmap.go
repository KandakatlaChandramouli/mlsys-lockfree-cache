package mmap

import (
    "os"
    "syscall"
)

type Mapping struct {
    Data []byte
}

func Open(
    path string,
) (*Mapping, error) {

    f, err := os.Open(
        path,
    )

    if err != nil {
        return nil, err
    }

    stat, err := f.Stat()

    if err != nil {
        return nil, err
    }

    data, err := syscall.Mmap(
        int(f.Fd()),
        0,
        int(stat.Size()),
        syscall.PROT_READ,
        syscall.MAP_SHARED,
    )

    if err != nil {
        return nil, err
    }

    return &Mapping{
        Data: data,
    }, nil
}

func (m *Mapping) Close() error {

    return syscall.Munmap(
        m.Data,
    )
}
