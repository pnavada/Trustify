package network

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"time"
	"trustify/blockchain"
	"trustify/logger"
)

// Global port for broadcast and listening
var Port = 9000

// Message type constants
const (
	MessageTypeTransaction = iota
	MessageTypeBlock
)

// BroadcastMessage sends a UDP message to all networks the peer is part of.
func BroadcastMessage(message []byte) error {
	// Retrieve all network interfaces
	interfaces, err := net.Interfaces()
	if err != nil {
		logger.ErrorLogger.Printf("Failed to get network interfaces: %v\n", err)
		return fmt.Errorf("failed to get network interfaces: %w", err)
	}

	successCount := 0

	// Iterate over each interface to broadcast the message
	for _, iface := range interfaces {
		// Skip interfaces that are not up or do not support broadcasting
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagBroadcast == 0 {
			continue
		}

		// Get the broadcast address for the interface
		addrs, err := iface.Addrs()
		if err != nil {
			logger.ErrorLogger.Printf("Failed to get addresses for interface %s: %v\n", iface.Name, err)
			continue
		}

		for _, addr := range addrs {
			ipNet, ok := addr.(*net.IPNet)
			if !ok || ipNet.IP.To4() == nil {
				continue
			}

			// Compute the broadcast address
			broadcastAddr := computeBroadcastAddress(ipNet)
			if broadcastAddr == nil {
				continue
			}

			// Add the port to the broadcast address
			broadcastWithPort := fmt.Sprintf("%s:%d", broadcastAddr.String(), Port)

			// Resolve the UDP address
			udpAddr, err := net.ResolveUDPAddr("udp", broadcastWithPort)
			if err != nil {
				logger.ErrorLogger.Printf("Failed to resolve UDP address for interface %s: %v\n", iface.Name, err)
				continue
			}

			// Create a UDP connection
			conn, err := net.DialUDP("udp", nil, udpAddr)
			if err != nil {
				logger.ErrorLogger.Printf("Failed to establish UDP connection on interface %s: %v\n", iface.Name, err)
				continue
			}
			defer conn.Close()

			// Set the write deadline
			if err := conn.SetWriteDeadline(time.Now().Add(5 * time.Second)); err != nil {
				logger.ErrorLogger.Printf("Failed to set write deadline on interface %s: %v\n", iface.Name, err)
				continue
			}

			// Send the message
			_, err = conn.Write(message)
			if err != nil {
				logger.ErrorLogger.Printf("Failed to send broadcast message on interface %s: %v\n", iface.Name, err)
			} else {
				logger.InfoLogger.Printf("Broadcast message sent on interface %s to %s\n", iface.Name, broadcastWithPort)
				successCount++
			}
		}
	}

	if successCount == 0 {
		return fmt.Errorf("broadcast failed on all interfaces")
	}

	logger.InfoLogger.Printf("Broadcast message sent on %d interface(s)\n", successCount)
	return nil
}

// computeBroadcastAddress calculates the broadcast address for a given network.
func computeBroadcastAddress(ipNet *net.IPNet) net.IP {
	ip := ipNet.IP.To4()
	mask := ipNet.Mask
	if ip == nil || mask == nil {
		return nil
	}

	// Perform bitwise OR on the inverted subnet mask and the IP
	broadcast := make(net.IP, len(ip))
	for i := 0; i < len(ip); i++ {
		broadcast[i] = ip[i] | ^mask[i]
	}

	return broadcast
}

// ReceiveMessages listens for UDP broadcast messages on the global port.
func ReceiveMessages(readChannel chan InboundMessage) error {
	// Listen on all interfaces for UDP messages on the specified port
	addr := fmt.Sprintf(":%d", Port)
	conn, err := net.ListenPacket("udp", addr)
	if err != nil {
		logger.ErrorLogger.Printf("Failed to listen for UDP messages: %v\n", err)
		return fmt.Errorf("failed to listen for UDP messages: %w", err)
	}
	defer conn.Close()

	logger.InfoLogger.Printf("Listening for UDP messages on port %d...\n", Port)

	buffer := make([]byte, 1024*1024) // 1 MB buffer size
	for {
		n, remoteAddr, err := conn.ReadFrom(buffer)

		if err != nil {
			logger.ErrorLogger.Printf("Error reading from UDP connection: %v\n", err)
			continue
		}

		if n == 0 {
			continue
		}

		readChannel <- InboundMessage{Data: buffer[:n], Sender: remoteAddr}
	}
}

// SendTransaction serializes and broadcasts a transaction along with its signature and public key.
func SendTransaction(tx *blockchain.Transaction, signature []byte, publicKey []byte) error {
	txData := blockchain.SerializeTransaction(tx)
	if txData == nil {
		return fmt.Errorf("failed to serialize transaction")
	}

	// Construct the message with the type header
	var buf bytes.Buffer
	buf.WriteByte(byte(MessageTypeTransaction)) // Message type

	// Write Transaction Length and Data
	err := writeVarBytes(&buf, txData)
	if err != nil {
		return fmt.Errorf("failed to write transaction data: %v", err)
	}

	// Write Signature Length and Data
	err = writeVarBytes(&buf, signature)
	if err != nil {
		return fmt.Errorf("failed to write signature data: %v", err)
	}

	// Write PublicKey Length and Data
	err = writeVarBytes(&buf, publicKey)
	if err != nil {
		return fmt.Errorf("failed to write public key data: %v", err)
	}

	// Broadcast the message
	return BroadcastMessage(buf.Bytes())
}

// deserializeTransactionMessage deserializes a Transaction, signature, and public key from the payload.
func deserializeTransactionMessage(data []byte) (*blockchain.Transaction, []byte, []byte, error) {
	buf := bytes.NewReader(data)

	// Read Transaction Data
	txData, err := readVarBytes(buf)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to read transaction data: %v", err)
	}
	tx, err := blockchain.DeserializeTransaction(txData)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to deserialize transaction: %v", err)
	}

	// Read Signature Data
	signature, err := readVarBytes(buf)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to read signature data: %v", err)
	}

	// Read PublicKey Data
	publicKey, err := readVarBytes(buf)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to read public key data: %v", err)
	}

	return tx, signature, publicKey, nil
}

// SendBlock serializes and broadcasts a block.
func SendBlock(block *blockchain.Block) error {
	serializedBlock := blockchain.SerializeBlock(block)
	if serializedBlock == nil {
		return fmt.Errorf("failed to serialize block")
	}

	// if len(serializedBlock) > 1200 { // Safe threshold below typical MTU
	// 	logger.ErrorLogger.Printf("Serialized block size %d exceeds 1200 bytes limit", len(serializedBlock))
	// 	return fmt.Errorf("serialized block size %d exceeds limit", len(serializedBlock))
	// }

	// Construct the message with the type header and length
	var buf bytes.Buffer
	buf.WriteByte(byte(MessageTypeBlock)) // Message type

	// Write the length of the serialized block
	// err := binary.Write(&buf, binary.BigEndian, uint32(len(serializedBlock)))
	// if err != nil {
	// return fmt.Errorf("failed to write block length: %v", err)
	// }

	// Write the serialized block
	buf.Write(serializedBlock)

	// Log serializedBlock
	logger.InfoLogger.Printf("Serialized block: %v\n", serializedBlock)

	// Broadcast the message
	return BroadcastMessage(buf.Bytes())
}

// Utility functions for reading and writing variable length bytes
func writeVarBytes(w *bytes.Buffer, data []byte) error {
	length := uint64(len(data))
	err := binary.Write(w, binary.BigEndian, length)
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}

// Define a maximum allowed length (e.g., 10MB)
const MaxVarBytesLength = 10 * 1024 * 1024 // 10 MB

func readVarBytes(r *bytes.Reader) ([]byte, error) {
	var length uint64
	err := binary.Read(r, binary.BigEndian, &length)
	if err != nil {
		return nil, fmt.Errorf("failed to read length: %v", err)
	}

	if length > MaxVarBytesLength {
		return nil, fmt.Errorf("readVarBytes: length %d exceeds maximum allowed %d", length, MaxVarBytesLength)
	}

	data := make([]byte, length)
	n, err := r.Read(data)
	if err != nil {
		return nil, fmt.Errorf("failed to read data: %v", err)
	}
	if uint64(n) != length {
		return nil, fmt.Errorf("expected %d bytes, read %d bytes", length, n)
	}
	return data, nil
}
