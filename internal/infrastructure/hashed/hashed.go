package hashed

import (
	"fmt"

	"github.com/Vanak/Backend/bootstrap"
	"github.com/speps/go-hashids/v2"
)

type HashID struct {
	encoder *hashids.HashID
}

func NewHashID(cfg bootstrap.Config) *HashID {
	hd := hashids.NewData()
	hd.Salt = cfg.Env.Hashed.Salt
	hd.MinLength = cfg.Env.Hashed.MinLength

	encoder, err := hashids.NewWithData(hd)
	if err != nil {
		panic(fmt.Sprintf("failed to create hashid: %v", err))
	}

	return &HashID{encoder: encoder}
}

func (h *HashID) Encode(id uint) (string, error) {
	hash, err := h.encoder.Encode([]int{int(id)})
	if err != nil {
		return "", err
	}
	return hash, nil
}

func (h *HashID) Decode(hash string) (uint, error) {
	ids, err := h.encoder.DecodeWithError(hash)
	if err != nil {
		return 0, fmt.Errorf("invalid hash: %w", err)
	}

	if len(ids) == 0 {
		return 0, fmt.Errorf("empty hash")
	}

	return uint(ids[0]), nil
}
