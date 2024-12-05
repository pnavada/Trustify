package network

import (
	"bytes"
	"context"
	"errors"
	"time"
	"trustify/blockchain"
	"trustify/logger"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
)

const GetBlocksProtocolID = protocol.ID("/trustify/getblocks/1.0.0")

type GetBlocksProtocol struct {
	Timeout time.Duration
	Node    *Node
	Host    host.Host
}

type GetBlocksRequest struct {
	LastKnownHash []byte
}

type GetBlocksResponse struct {
	Blocks  []*blockchain.Block
	Success bool
}

func NewGetBlocksProtocol(node *Node, host host.Host, timeout int) *GetBlocksProtocol {
	return &GetBlocksProtocol{
		Timeout: time.Duration(timeout) * time.Second,
		Node:    node,
		Host:    host,
	}
}

// GetBlocks requests missing blocks from a peer.
func (p *GetBlocksProtocol) GetBlocks(peerName string, lastKnownHash []byte) error {

	logger.InfoLogger.Printf("Initiating GetBlocks request to peer %s\n", peerName)

	peerID, err := peer.Decode(peerName)
	if err != nil {
		logger.ErrorLogger.Printf("Invalid peer ID: %v", err)
		return err
	}

	peerAddr := p.Host.Peerstore().PeerInfo(peerID)

	ctx, cancel := context.WithTimeout(context.Background(), p.Timeout)
	defer cancel()

	stream, err := p.Host.NewStream(ctx, peerAddr.ID, GetBlocksProtocolID)
	if err != nil {
		logger.ErrorLogger.Printf("Failed to create stream to peer %s: %v", peerID, err)
		return err
	}
	defer stream.Close()

	request := GetBlocksRequest{LastKnownHash: lastKnownHash}
	err = Send(stream, request)
	if err != nil {
		logger.ErrorLogger.Printf("Failed to send GetBlocksRequest to peer %s: %v", peerID, err)
		return err
	}

	var response GetBlocksResponse
	err = Receive(stream, &response)
	if err != nil {
		logger.ErrorLogger.Printf("Failed to receive GetBlocksResponse from peer %s: %v", peerID, err)
		return err
	}

	return p.ProcessGetBlocksResponse(response)
}

// HandleGetBlocksRequest handles incoming GetBlocks requests.
func (p *GetBlocksProtocol) HandleGetBlocksRequest(s network.Stream) {
	defer s.Close()

	var request GetBlocksRequest
	err := Receive(s, &request)
	if err != nil {
		logger.ErrorLogger.Printf("Failed to receive GetBlocksRequest: %v", err)
		return
	}

	logger.InfoLogger.Printf("Received GetBlocksRequest from %s", s.Conn().RemotePeer().String())

	blocks, err := p.getBlocksSinceHash(request.LastKnownHash)
	if err != nil {
		logger.ErrorLogger.Printf("Failed to get blocks since hash: %v", err)
		Send(s, GetBlocksResponse{Success: false})
		return
	}

	response := GetBlocksResponse{
		Blocks:  blocks,
		Success: true,
	}

	err = Send(s, response)
	if err != nil {
		logger.ErrorLogger.Printf("Failed to send GetBlocksResponse: %v", err)
		return
	}

	logger.InfoLogger.Printf("Sent %d blocks to peer %s", len(blocks), s.Conn().RemotePeer().String())
}

// ProcessGetBlocksResponse processes the response from a GetBlocks request.
func (p *GetBlocksProtocol) ProcessGetBlocksResponse(response GetBlocksResponse) error {
	if !response.Success {
		logger.ErrorLogger.Println("GetBlocksResponse indicates failure")
		return errors.New("failed to get blocks from peer")
	}

	logger.InfoLogger.Printf("Processing %d received blocks", len(response.Blocks))

	// Validate and integrate received blocks
	for _, block := range response.Blocks {
		err := p.Node.Blockchain.AddBlock(block)
		if err != nil {
			logger.ErrorLogger.Printf("Failed to add block %x: %v", block.Header.BlockHash, err)
			// Handle invalid block: add transactions to mempool
			for _, tx := range block.Transactions {
				p.Node.Mempool.AddTransaction(tx)
			}
			continue
		}
	}

	logger.InfoLogger.Println("Successfully processed received blocks")
	return nil
}

// Helper method to get blocks since a given hash.
func (p *GetBlocksProtocol) getBlocksSinceHash(lastKnownHash []byte) ([]*blockchain.Block, error) {
	blocks := p.Node.Blockchain.Ledger
	var startIndex int = -1

	for i, block := range blocks {
		if bytes.Equal(block.Header.BlockHash, lastKnownHash) {
			startIndex = i + 1
			break
		}
	}

	if startIndex == -1 {
		logger.ErrorLogger.Println("Last known block hash not found")
		return nil, errors.New("last known block hash not found")
	}

	return blocks[startIndex:], nil
}

// Send is a helper function to send data over a LibP2P stream.
func Send(s network.Stream, v interface{}) error {
	data := blockchain.Serialize(v)
	_, err := s.Write(data)
	return err
}

// Receive is a helper function to receive data over a LibP2P stream.
func Receive(s network.Stream, v interface{}) error {
	data := make([]byte, 1024*1024) // 1 MB buffer
	n, err := s.Read(data)
	if err != nil {
		return err
	}
	blockchain.Deserialize(data[:n], v)
	return nil
}
