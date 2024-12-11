package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"trustify/logger"
)

func init() {
	// Register all concrete types used in interfaces
	gob.Register(&PurchaseTransactionData{})
	gob.Register(&ReviewTransactionData{})
	gob.Register(&UTXOTransaction{})
	gob.Register(&UTXOTransaction{})
	gob.Register(&UTXOTransactionID{})
	gob.Register(&CoinbaseTransactionData{})
	gob.Register(&BlockHeader{})
	gob.Register(&Block{})
	gob.Register(&Transaction{})

}

func SerializeTransaction(tx *Transaction) []byte {
	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	enc.Encode(tx)
	return buff.Bytes()
}

func DeserializeTransaction(data []byte) (*Transaction, error) {
	var tx Transaction
	dec := gob.NewDecoder(bytes.NewReader(data))
	err := dec.Decode(&tx)
	if err != nil {
		logger.ErrorLogger.Printf("Failed to deserialize transaction: %v\n", err)
		return nil, err
	}
	return &tx, nil
}

func SerializeBlockHeader(b *BlockHeader) []byte {
	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	enc.Encode(b)
	return buff.Bytes()
}

func SerializeBlock(b *Block) []byte {
	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)

	enc.Encode(b)

	// serialized := buff.Bytes()
	// deserialized := DeserializeBlock(serialized)

	// // Compare deserialized and b with deep equality
	// if !bytes.Equal(SerializeBlockHeader(deserialized.Header), SerializeBlockHeader(b.Header)) {
	// 	logger.ErrorLogger.Println("Block header serialization failed")
	// } else {
	// 	logger.InfoLogger.Println("Block header serialization successful")
	// }

	// // Also deep print the block and deserialized block
	// logger.InfoLogger.Println("Block: ", b)
	// logger.InfoLogger.Println("Deserialized block: ", deserialized)
	// b.PrintToString()
	// deserialized.PrintToString()

	return buff.Bytes()
}

func DeserializeBlock(data []byte) (*Block, error) {
	var b Block
	dec := gob.NewDecoder(bytes.NewReader(data))
	err := dec.Decode(&b)
	if err != nil {
		logger.ErrorLogger.Printf("Failed to deserialize block: %v\n", err)
		return nil, err
	}
	if b.Header == nil {
		logger.ErrorLogger.Println("Deserialized block has nil Header")
		return nil, fmt.Errorf("block header is nil")
	}
	// Additional validation if necessary
	return &b, nil
}

func HashObject(serializedData []byte) []byte {
	hash := sha256.Sum256(serializedData)
	return hash[:]
}

func Serialize(data interface{}) []byte {
	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	enc.Encode(data)
	return buff.Bytes()
}

func Deserialize(data []byte, v interface{}) {
	dec := gob.NewDecoder(bytes.NewReader(data))
	dec.Decode(v)
}
