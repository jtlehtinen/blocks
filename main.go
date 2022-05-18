package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"time"
)

type Block struct {
	Index     uint64    `json:"index"`
	Timestamp time.Time `json:"time"`
	Data      int       `json:"data"`
	PrevHash  []byte    `json:"prevHash"`
	Hash      []byte    `json:"hash"`
}

type Blockchain struct {
	Chain []*Block `json:"chain"`
}

func New() *Blockchain {
	bc := new(Blockchain)
	bc.Chain = make([]*Block, 0)
	bc.addBlock(Sha256([]byte("genesis")), 0)
	return bc
}

func (bc *Blockchain) hash(b *Block) ([]byte, error) {
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

func (bc *Blockchain) getLastBlock() *Block {
	return bc.Chain[len(bc.Chain)-1]
}

func (bc *Blockchain) addBlock(prevHash []byte, data int) error {
	b := new(Block)
	b.Index = uint64(len(bc.Chain))
	b.Timestamp = time.Now().UTC()
	b.Data = data
	b.PrevHash = prevHash

	// @TODO: Hashing should not return an error and
	// should always succeed.
	hash, err := bc.hash(b)
	if err != nil {
		return err
	}

	b.Hash = hash
	bc.Chain = append(bc.Chain, b)

	return nil
}

func (bc *Blockchain) isValidBlock(prev, block *Block) bool {
	if !bytes.Equal(prev.Hash, block.PrevHash) || prev.Index+1 != block.Index {
		return false
	}

	hash, _ := bc.hash(block)
	return bytes.Equal(hash, block.Hash)
}

func Sha256(b []byte) []byte {
	sha := sha256.New()
	sha.Write(b)
	return sha.Sum(nil)
}

func main() {
	bc := New()
	bc.addBlock(bc.getLastBlock().Hash, 1)
	bc.addBlock(bc.getLastBlock().Hash, 2)
	bc.addBlock(bc.getLastBlock().Hash, 3)
	bc.addBlock(bc.getLastBlock().Hash, 4)
	bc.addBlock(bc.getLastBlock().Hash, 5)

	b, _ := json.MarshalIndent(bc, "", "  ")
	fmt.Println(string(b))
}
