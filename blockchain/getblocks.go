package blockchain

// import (
// 	"bytes"
// 	"context"
// 	"errors"
// 	"fmt"
// 	"sort"
// 	"time"
// 	"trustify/logger"

// 	"github.com/libp2p/go-libp2p/core/host"
// 	"github.com/libp2p/go-libp2p/core/network"
// 	"github.com/libp2p/go-libp2p/core/peer"
// 	"github.com/libp2p/go-libp2p/core/protocol"
// )

// // Protocol ID for GetBlocks communication
// const GetBlocksProtocolID = protocol.ID("/trustify/getblocks/1.0.0")

// // GetBlocksProtocol manages block synchronization between peers
// type GetBlocksProtocol struct {
// 	Timeout           time.Duration
// 	Blockchain        *Blockchain
// 	Host              host.Host
// 	BestBlocksChannel chan *GetBlocksResponse
// }

// // GetBlocksRequest represents a request for missing blocks
// type GetBlocksRequest struct {
// 	LastKnownHash []byte
// }

// // GetBlocksResponse contains blocks returned by a peer
// type GetBlocksResponse struct {
// 	Blocks  []*Block
// 	Success bool
// }

// // NewGetBlocksProtocol initializes the protocol with a timeout and associated node
// func NewGetBlocksProtocol(host host.Host, timeout int, blockchain *Blockchain, bestBlocksChannel chan *GetBlocksResponse) *GetBlocksProtocol {
// 	protocol := &GetBlocksProtocol{
// 		Timeout:           time.Duration(timeout) * time.Second,
// 		Host:              host,
// 		Blockchain:        blockchain,
// 		BestBlocksChannel: bestBlocksChannel,
// 	}
// 	host.SetStreamHandler("/trustify/getblocks/1.0.0", protocol.HandleGetBlocksRequest)
// 	logger.InfoLogger.Println("GetBlocksProtocol initialized and stream handler registered.")
// 	return protocol
// }

// // GetBlocks requests missing blocks from peers and resolves forks
// func (p *GetBlocksProtocol) GetBlocks(lastKnownHash []byte) error {
// 	logger.InfoLogger.Println("Starting GetBlocks protocol")

// 	// Retrieve the list of peers from the peerstore
// 	peerList := p.Host.Peerstore().Peers()
// 	if len(peerList) == 0 {
// 		logger.ErrorLogger.Println("No peers available for synchronization")
// 		return errors.New("no peers available")
// 	}

// 	// Send block requests to peers and gather responses
// 	responses, err := p.sendGetBlocksRequests(peerList, lastKnownHash)
// 	if err != nil {
// 		return fmt.Errorf("failed to get responses: %w", err)
// 	}

// 	// Validate and process responses to extract blocks
// 	validBlocks := p.processResponses(responses)

// 	// Resolve potential forks using the received blocks
// 	return p.resolveFork(validBlocks)
// }

// // sendGetBlocksRequests sends requests to peers concurrently and collects responses
// func (p *GetBlocksProtocol) sendGetBlocksRequests(peerList []peer.ID, lastKnownHash []byte) ([]*GetBlocksResponse, error) {
// 	responses := make(chan *GetBlocksResponse, len(peerList))

// 	// Send requests to each peer in parallel
// 	for _, peerID := range peerList {
// 		go func(peerID peer.ID) {
// 			logger.InfoLogger.Printf("Requesting blocks from peer %s", peerID)
// 			response, err := p.requestBlocks(peerID, lastKnownHash)
// 			if err != nil {
// 				logger.ErrorLogger.Printf("Failed to get blocks from peer %s: %v", peerID, err)
// 				return
// 			}
// 			responses <- response
// 		}(peerID)
// 	}

// 	// Wait for responses within the timeout period
// 	ctx, cancel := context.WithTimeout(context.Background(), p.Timeout)
// 	defer cancel()

// 	var allResponses []*GetBlocksResponse
// 	for {
// 		select {
// 		case response := <-responses:
// 			allResponses = append(allResponses, response)
// 		case <-ctx.Done():
// 			logger.InfoLogger.Println("Timeout reached while waiting for peer responses")
// 			if len(allResponses) == 0 {
// 				return nil, errors.New("no valid responses received")
// 			}
// 			return allResponses, nil
// 		}
// 	}
// }

// // processResponses validates and extracts blocks from received responses
// func (p *GetBlocksProtocol) processResponses(responses []*GetBlocksResponse) [][]*Block {
// 	var allValidBlocks [][]*Block

// 	for _, response := range responses {
// 		if response.Success {
// 			var validBlocks []*Block
// 			for _, block := range response.Blocks {
// 				// Validate each block individually
// 				if err := p.Blockchain.ValidateBlock(block); err == nil {
// 					validBlocks = append(validBlocks, block)
// 				} else {
// 					logger.ErrorLogger.Printf("Invalid block received: %v", err)
// 				}
// 			}
// 			allValidBlocks = append(allValidBlocks, validBlocks)
// 		} else {
// 			logger.ErrorLogger.Println("Received unsuccessful response from a peer")
// 		}
// 	}

// 	logger.InfoLogger.Printf("Processed %d sets of valid blocks", len(allValidBlocks))
// 	return allValidBlocks
// }

// // resolveFork chooses the best chain from available options based on chain length and transaction fees
// func (p *GetBlocksProtocol) resolveFork(blocks [][]*Block) error {
// 	if len(blocks) == 0 {
// 		return errors.New("no valid blocks to resolve fork")
// 	}

// 	// Sort chains by length and total transaction fees
// 	sort.Slice(blocks, func(i, j int) bool {
// 		if len(blocks[i]) != len(blocks[j]) {
// 			return len(blocks[i]) > len(blocks[j])
// 		}
// 		var feeI, feeJ int64
// 		for _, block := range blocks[i] {
// 			feeI += block.GetTransactionFee()
// 		}
// 		for _, block := range blocks[j] {
// 			feeJ += block.GetTransactionFee()
// 		}
// 		return feeI > feeJ
// 	})

// 	// Select the best chain and add its blocks to the blockchain
// 	bestBlocks := blocks[0]
// 	p.BestBlocksChannel <- &GetBlocksResponse{Blocks: bestBlocks, Success: true}

// 	logger.InfoLogger.Println("Fork resolved successfully")
// 	return nil
// }

// // requestBlocks sends a request to a single peer and returns its response
// func (p *GetBlocksProtocol) requestBlocks(peerID peer.ID, lastKnownHash []byte) (*GetBlocksResponse, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), p.Timeout)
// 	defer cancel()

// 	stream, err := p.Host.NewStream(ctx, peerID, GetBlocksProtocolID)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to create stream: %w", err)
// 	}
// 	defer stream.Close()

// 	// Send the request
// 	request := GetBlocksRequest{LastKnownHash: lastKnownHash}
// 	if err := Send(stream, request); err != nil {
// 		return nil, fmt.Errorf("failed to send request: %w", err)
// 	}

// 	// Receive the response
// 	var response GetBlocksResponse
// 	if err := Receive(stream, &response); err != nil {
// 		return nil, fmt.Errorf("failed to receive response: %w", err)
// 	}

// 	return &response, nil
// }

// // HandleGetBlocksRequest handles incoming GetBlocks requests from peers
// func (p *GetBlocksProtocol) HandleGetBlocksRequest(s network.Stream) {
// 	defer s.Close()

// 	var request GetBlocksRequest
// 	if err := Receive(s, &request); err != nil {
// 		logger.ErrorLogger.Printf("Failed to receive request: %v", err)
// 		return
// 	}

// 	blocks, err := p.getBlocksSinceHash(request.LastKnownHash)
// 	if err != nil {
// 		logger.ErrorLogger.Printf("Error retrieving blocks since hash: %v", err)
// 		Send(s, GetBlocksResponse{Success: false})
// 		return
// 	}

// 	response := GetBlocksResponse{
// 		Blocks:  blocks,
// 		Success: true,
// 	}

// 	if err := Send(s, response); err != nil {
// 		logger.ErrorLogger.Printf("Failed to send response: %v", err)
// 	}
// }

// // getBlocksSinceHash retrieves all blocks from the ledger after a given hash
// func (p *GetBlocksProtocol) getBlocksSinceHash(lastKnownHash []byte) ([]*Block, error) {
// 	blocks := p.Blockchain.Ledger
// 	startIndex := -1

// 	// Find the index of the last known hash
// 	for i, block := range blocks {
// 		if bytes.Equal(block.Header.BlockHash, lastKnownHash) {
// 			startIndex = i + 1
// 			break
// 		}
// 	}

// 	if startIndex == -1 {
// 		return nil, errors.New("last known hash not found")
// 	}

// 	return blocks[startIndex:], nil
// }

// // Send serializes and sends data over a LibP2P stream
// func Send(s network.Stream, v interface{}) error {
// 	data := Serialize(v)
// 	_, err := s.Write(data)
// 	return err
// }

// // Receive deserializes data received over a LibP2P stream
// func Receive(s network.Stream, v interface{}) error {
// 	data := make([]byte, 1024*1024) // 1 MB buffer
// 	n, err := s.Read(data)
// 	if err != nil {
// 		return err
// 	}
// 	Deserialize(data[:n], v)
// 	return nil
// }
