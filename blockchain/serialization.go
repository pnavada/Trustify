package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"trustify/logger"
)

func SerializeTransaction(tx *Transaction) []byte {
	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	enc.Encode(tx)
	return buff.Bytes()
}

func HashObject(serializedData []byte) []byte {
	hash := sha256.Sum256(serializedData)
	logger.InfoLogger.Println("Computed hash for block:", hash[:])
	return hash[:]
}

func DeserializeTransaction(data []byte) *Transaction {
	var tx Transaction
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

func HashObject(serializedData []byte) []byte {
    hash := sha256.Sum256(serializedData)
    logger.InfoLogger.Println("Computed hash for block:", hash[:])
    return hash[:]
}