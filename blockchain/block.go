package blockchain

import (
	"bytes"
	"encoding/gob"
	"log"
)

type Block struct {
	Hash     []byte
	Data     []byte
	PrevHash []byte
	Nonce    int
}

func NewBlock(data string, prevHash []byte) *Block {
	block := &Block{[]byte{}, []byte(data), prevHash, 0}
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

func (b *Block) Serialize() []byte {
	var res bytes.Buffer

	err := gob.NewEncoder(&res).Encode(b)
	if err != nil {
		log.Fatal("Could not Encode data: ", err)
	}

	return res.Bytes()
}

func (b *Block) Deserialize(data []byte) *Block {
	var block Block

	err := gob.NewDecoder(bytes.NewReader(data)).Decode(&block)
	if err != nil {
		log.Fatal("Could not Decode data: ", err)
	}

	return &block
}
