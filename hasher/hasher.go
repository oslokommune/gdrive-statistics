package hasher

import (
	"crypto/sha1"
	"encoding/hex"
)

type Hash struct {
	data []byte
}

func NewHash(data string) *Hash {
	bData := []byte(data)

	h := sha1.New()
	h.Write(bData)

	sha := h.Sum(nil) // "sha" is uint8 type, encoded in base16

	return &Hash{data: sha}
}

func (h *Hash) String() string {
	return hex.EncodeToString(h.data)
}
