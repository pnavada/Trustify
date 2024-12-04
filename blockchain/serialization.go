package blockchain

import (
	"bytes"
	"encoding/gob"
)

func SerializeTransaction(tx *UTXOTransaction) []byte {
	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	enc.Encode(tx)
	return buff.Bytes()
}

func DeserializeTransaction(data []byte) *UTXOTransaction {
	var tx UTXOTransaction
	dec := gob.NewDecoder(bytes.NewReader(data))
	dec.Decode(&tx)
	return &tx
}

func SerializeBlock(b *Block) []byte {
	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	enc.Encode(b)
	return buff.Bytes()
}

func DeserializeBlock(data []byte) *Block {
	var b Block
	dec := gob.NewDecoder(bytes.NewReader(data))
	dec.Decode(&b)
	return &b
}
