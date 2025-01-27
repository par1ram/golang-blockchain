package blockchain

type BlockChain struct {
	Blocks []*Block
}

func InItBlockChain() *BlockChain {
	return &BlockChain{[]*Block{Genesis()}}
}

func Genesis() *Block {
	return NewBlock("Genesis", []byte{})
}

func (bc *BlockChain) AddBlock(data string) {
	prevBlock := bc.Blocks[len(bc.Blocks)-1]
	block := NewBlock(data, prevBlock.Hash)
	bc.Blocks = append(bc.Blocks, block)
}
