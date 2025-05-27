package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"math/big"
)

// Difficulty представляет сложность алгоритма Proof of Work.
const Difficulty = 10

// ProofOfWork представляет Proof of Work для блока.
type ProofOfWork struct {
	Block  *Block
	Target *big.Int
}

// NewProofOfWork создает и возвращает новый ProofOfWork для блока.
func NewProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-Difficulty))

	return &ProofOfWork{b, target}
}

// InItData подготавливает данные для хеширования.
func (pow *ProofOfWork) InItData(nonce int) []byte {
	return bytes.Join(
		[][]byte{
			pow.Block.PrevHash,
			pow.Block.Data,
			ToHex(int64(nonce)),
			ToHex(int64(Difficulty)),
		},
		[]byte{},
	)
}

// Run выполняет алгоритм Proof of Work и возвращает валидный nonce и хэш блока.
func (pow *ProofOfWork) Run() (int, []byte) {
	var intHash big.Int
	var hash [32]byte

	nonce := 0
	for nonce <= math.MaxInt {
		data := pow.InItData(nonce)
		hash = sha256.Sum256(data)

		fmt.Printf("\r%x", hash)
		intHash.SetBytes(hash[:])

		// if hash < Target -> block signed
		if intHash.Cmp(pow.Target) == -1 {
			break
		} else {
			nonce++
		}
	}
	fmt.Println()

	return nonce, hash[:]
}

// Validate проверяет Proof of Work блока.
func (pow *ProofOfWork) Validate() bool {
	var intHash big.Int

	data := pow.InItData(pow.Block.Nonce)
	hash := sha256.Sum256(data)
	intHash.SetBytes(hash[:])

	return intHash.Cmp(pow.Target) == -1
}

// ToHex преобразует int64 в срез байтов.
func ToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Fatal(err)
	}

	return buff.Bytes()
}
