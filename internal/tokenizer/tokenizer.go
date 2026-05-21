
package tokenizer

import (
	"strings"
)

type Tokenizer struct {
	vocab map[string]int32
}

func New() *Tokenizer {

	v := map[string]int32{
		"[PAD]": 0,
		"[UNK]": 100,
		"[CLS]": 101,
		"[SEP]": 102,
	}

	return &Tokenizer{
		vocab: v,
	}
}

func (t *Tokenizer) Encode(
	text string,
) []int32 {

	out := []int32{
		101,
	}

	parts := strings.Fields(
		strings.ToLower(text),
	)

	for _, p := range parts {

		id, ok := t.vocab[p]

		if !ok {

			id = int32(
				len(t.vocab) + 1000,
			)

			t.vocab[p] = id
		}

		out = append(
			out,
			id,
		)
	}

	out = append(
		out,
		102,
	)

	return out
}
