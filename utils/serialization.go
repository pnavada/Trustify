package utils

import (
	"bytes"
	"encoding/gob"
	"trustify/blockchain"
)

func SerializeTransaction(tx *blockchain.UTXOTransaction) []byte {
	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	enc.Encode(tx)
	return buff.Bytes()
}

func DeserializeTransaction(data []byte) *blockchain.UTXOTransaction {
	var tx blockchain.UTXOTransaction
	dec := gob.NewDecoder(bytes.NewReader(data))
	dec.Decode(&tx)
	return &tx
}

func SerializeBlock(b *blockchain.Block) []byte {
	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	enc.Encode(b)
	return buff.Bytes()
}

func DeserializeBlock(data []byte) *blockchain.Block {
	var b blockchain.Block
	dec := gob.NewDecoder(bytes.NewReader(data))
	dec.Decode(&b)
	return &b
}
