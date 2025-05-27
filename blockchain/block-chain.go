package blockchain

// BlockChain представляет цепочку блоков.
type BlockChain struct {
	Blocks []*Block
}

// InItBlockChain создает и возвращает новую цепочку блоков с генезис-блоком.
func InItBlockChain() *BlockChain {
	return &BlockChain{[]*Block{Genesis()}}
}

// Genesis создает и возвращает генезис-блок.
func Genesis() *Block {
	return NewBlock("Genesis", []byte{})
}

// AddBlock добавляет новый блок в цепочку блоков.
func (bc *BlockChain) AddBlock(data string) {
	prevBlock := bc.Blocks[len(bc.Blocks)-1]
	block := NewBlock(data, prevBlock.Hash)
	bc.Blocks = append(bc.Blocks, block)
}
