
package protocol

import (
	"encoding/binary"
	"io"
)

func WriteString(
	w io.Writer,
	s string,
) error {

	size := uint32(len(s))

	err := binary.Write(
		w,
		binary.LittleEndian,
		size,
	)

	if err != nil {
		return err
	}

	_, err = w.Write(
		[]byte(s),
	)

	return err
}

func ReadString(
	r io.Reader,
) (string, error) {

	var size uint32

	err := binary.Read(
		r,
		binary.LittleEndian,
		&size,
	)

	if err != nil {
		return "", err
	}

	buf := make(
		[]byte,
		size,
	)

	_, err = io.ReadFull(
		r,
		buf,
	)

	if err != nil {
		return "", err
	}

	return string(buf), nil
}

func WriteEmbedding(
	w io.Writer,
	v []float32,
) error {

	size := uint32(len(v))

	err := binary.Write(
		w,
		binary.LittleEndian,
		size,
	)

	if err != nil {
		return err
	}

	for _, x := range v {

		err = binary.Write(
			w,
			binary.LittleEndian,
			x,
		)

		if err != nil {
			return err
		}
	}

	return nil
}

func ReadEmbedding(
	r io.Reader,
) ([]float32, error) {

	var size uint32

	err := binary.Read(
		r,
		binary.LittleEndian,
		&size,
	)

	if err != nil {
		return nil, err
	}

	out := make(
		[]float32,
		size,
	)

	for i := range out {

		err = binary.Read(
			r,
			binary.LittleEndian,
			&out[i],
		)

		if err != nil {
			return nil, err
		}
	}

	return out, nil
}
