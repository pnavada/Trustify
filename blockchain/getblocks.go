package blockchain

import (
	"context"
	"errors"
	"fmt"
	"time"
	"trustify/logger"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
)

// Protocol ID for GetBlocks communication
const GetBlocksProtocolID = protocol.ID("/trustify/getblocks/1.0.0")

// GetBlocksProtocol manages block synchronization between peers
type GetBlocksProtocol struct {
	Timeout       time.Duration
	Host          host.Host
	BlocksChannel chan *GetBlocksResponse
}

// GetBlocksRequest represents a request for missing blocks
type GetBlocksRequest struct {
	LastKnownHash []byte
}

// GetBlocksResponse contains blocks returned by a peer
type GetBlocksResponse struct {
	Blocks  []*Block
	Success bool
}

// NewGetBlocksProtocol initializes the protocol with a timeout and associated node
func NewGetBlocksProtocol(host host.Host, timeout int) *GetBlocksProtocol {
	protocol := &GetBlocksProtocol{
		Timeout:       time.Duration(timeout) * time.Second,
		Host:          host,
		BlocksChannel: make(chan *GetBlocksResponse),
	}
	logger.InfoLogger.Println("GetBlocksProtocol initialized and stream handler registered.")
	return protocol
}

// GetBlocks requests missing blocks from peers
func (p *GetBlocksProtocol) GetBlocks(lastKnownHash []byte) error {
	logger.InfoLogger.Println("Starting GetBlocks protocol")

	// Retrieve the list of peers from the peerstore
	peerList := p.Host.Peerstore().Peers()
	if len(peerList) == 0 {
		logger.ErrorLogger.Println("No peers available for synchronization")
		return errors.New("no peers available")
	}

	// Send block requests to peers and gather responses
	err := p.sendGetBlocksRequests(peerList, lastKnownHash)
	if err != nil {
		return fmt.Errorf("failed to get responses: %w", err)
	}

	return nil
}

// sendGetBlocksRequests sends requests to peers concurrently and collects responses
func (p *GetBlocksProtocol) sendGetBlocksRequests(peerList []peer.ID, lastKnownHash []byte) error {

	// Send requests to each peer in parallel
	for _, peerID := range peerList {
		go func(peerID peer.ID) {
			logger.InfoLogger.Printf("Requesting blocks from peer %s", peerID)
			response, err := p.requestBlocks(peerID, lastKnownHash)
			if err != nil {
				logger.ErrorLogger.Printf("Failed to get blocks from peer %s: %v", peerID, err)
				return
			}
			p.BlocksChannel <- response
		}(peerID)
	}

	// Wait for responses within the timeout period
	ctx, cancel := context.WithTimeout(context.Background(), p.Timeout)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			logger.InfoLogger.Println("Timeout reached while waiting for peer responses")
			return nil
		}
	}
}

// requestBlocks sends a request to a single peer and returns its response
func (p *GetBlocksProtocol) requestBlocks(peerID peer.ID, lastKnownHash []byte) (*GetBlocksResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), p.Timeout)
	defer cancel()

	stream, err := p.Host.NewStream(ctx, peerID, GetBlocksProtocolID)
	if err != nil {
		return nil, fmt.Errorf("failed to create stream: %w", err)
	}
	defer stream.Close()

	// Send the request
	request := GetBlocksRequest{LastKnownHash: lastKnownHash}
	if err := Send(stream, request); err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	// Receive the response
	var response GetBlocksResponse
	if err := Receive(stream, &response); err != nil {
		return nil, fmt.Errorf("failed to receive response: %w", err)
	}

	return &response, nil
}

// Send serializes and sends data over a LibP2P stream
func Send(s network.Stream, v interface{}) error {
	data := Serialize(v)
	_, err := s.Write(data)
	return err
}

// Receive deserializes data received over a LibP2P stream
func Receive(s network.Stream, v interface{}) error {
	data := make([]byte, 1024*1024) // 1 MB buffer
	n, err := s.Read(data)
	if err != nil {
		return err
	}
	Deserialize(data[:n], v)
	return nil
}
