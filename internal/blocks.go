package internal

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"time"
)

type Block struct {
	Index     uint64    `json:"index"`
	Timestamp time.Time `json:"time"`
	Data      string    `json:"data"`
	PrevHash  []byte    `json:"prevHash"`
	Hash      []byte    `json:"hash"`
}

func hash(b *Block) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	var err error
	encode := func(data any) {
		if err == nil {
			err = enc.Encode(data)
		}
	}

	encode(b.Index)
	encode(b.Timestamp)
	encode(b.Data)
	encode(b.PrevHash)

	if err != nil {
		return nil, err
	}

	return Sha256(buf.Bytes()), nil
}

func NewBlock(index uint64, data string, prevHash []byte) (*Block, error) {
	b := &Block{
		Index:     index,
		Timestamp: time.Now().UTC(),
		Data:      data,
		PrevHash:  prevHash,
	}

	// @TODO: Hashing should not return an error and
	// should always succeed.
	h, err := hash(b)
	if err != nil {
		return nil, err
	}

	b.Hash = h
	return b, nil
}

type Blockchain struct {
	Chain []*Block `json:"chain"`
}

func Sha256(b []byte) []byte {
	sha := sha256.New()
	sha.Write(b)
	return sha.Sum(nil)
}

func NewBlockchain() *Blockchain {
	b, _ := NewBlock(0, "genesis", []byte("0"))

	bc := &Blockchain{
		Chain: []*Block{b},
	}

	return bc
}

func (bc *Blockchain) AddBlock(data string) error {
	prevBlock := bc.getLastBlock()
	b, err := NewBlock(uint64(len(bc.Chain)), data, prevBlock.Hash)
	if err != nil {
		return err
	}

	bc.Chain = append(bc.Chain, b)
	return nil
}

func (bc *Blockchain) getLastBlock() *Block {
	return bc.Chain[len(bc.Chain)-1]
}
