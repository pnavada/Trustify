package blockchain

type GetBlocksProtocol struct {
	timeout int
}

type GetBlocksRequest struct {
	LastKnownHash []byte
	Peer          string
}

type GetBlocksResponse struct {
	Blocks  []Block
	Peer    string
	Success bool
}

func NewGetBlocksProtocol(timeout int) *GetBlocksProtocol {
	return &GetBlocksProtocol{timeout: timeout}
}

func (p *GetBlocksProtocol) GetBlocks(peer string, lastKnownHash string) error {
	// This method is used to request blocks from a peer.
	// It sends a GetBlocksRequest to the peer and waits for a response.
	return nil
}

func (p *GetBlocksProtocol) HandleGetBlocksRequest(request GetBlocksRequest) error {
	// This method is used to handle a GetBlocksRequest.
	// It retrieves the requested blocks from the blockchain and sends them back to the requesting peer.
	// Send all the blocks A
	// If the requested blocks are not found, it sends an error response.
	return nil
}

func (p *GetBlocksProtocol) ProcessGetBlocksResponse(response GetBlocksResponse) error {
	// This handles the response for the protocol
	// If there are no response even after the timeout - go ahead and initiate getblocks from a block behind until it succeeds
	// And once succeeds and with valid blocks - drop the invalid blocks and attach the missing blocks and complete the protocol
	// Dropping the invalid blocks here means to add those transactions to the mempool
	// After the timeout if there are responses from multiple peers
	// First of all each received list of blocks needs to be validated
	// Then chooses the longest chain
	// If there are competing chains of same lenght, choose the chain with the higher fee
	// most likely you will need to use a max heap foe this usecase
	// Add additional methods or files as needed maintaining separation of concerns

	return nil

}

// // Version 1

// package network

// import (
//     "context"
//     "encoding/hex"
//     "errors"
//     "time"
//     "trustify/blockchain"
//     "trustify/utils/logger"
// )

// const GetBlocksProtocolID = protocol.ID("/trustify/getblocks/1.0.0")

// type GetBlocksProtocol struct {
//     Host    host.Host
//     Timeout time.Duration
//     Node    *Node
// }

// type GetBlocksRequest struct {
//     LastKnownHash []byte
// }

// type GetBlocksResponse struct {
//     Blocks  []*blockchain.Block
//     Success bool
// }

// // NewGetBlocksProtocol initializes the GetBlocksProtocol with a timeout and host.
// func NewGetBlocksProtocol(node *Node, timeout int) *GetBlocksProtocol {
//     return &GetBlocksProtocol{
//         Host:    node.Host,
//         Timeout: time.Duration(timeout) * time.Second,
//         Node:    node,
//     }
// }

// // GetBlocks requests missing blocks from a peer.
// func (p *GetBlocksProtocol) GetBlocks(peerID string, lastKnownHash []byte) error {
//     logger.InfoLogger.Printf("Initiating GetBlocks request to peer %s", peerID)

//     peerAddr, err := peer.Decode(peerID)
//     if err != nil {
//         logger.ErrorLogger.Printf("Invalid peer ID: %v", err)
//         return err
//     }

//     ctx, cancel := context.WithTimeout(context.Background(), p.Timeout)
//     defer cancel()

//     stream, err := p.Host.NewStream(ctx, peerAddr, GetBlocksProtocolID)
//     if err != nil {
//         logger.ErrorLogger.Printf("Failed to create stream to peer %s: %v", peerID, err)
//         return err
//     }
//     defer stream.Close()

//     request := GetBlocksRequest{LastKnownHash: lastKnownHash}
//     err = Send(stream, request)
//     if err != nil {
//         logger.ErrorLogger.Printf("Failed to send GetBlocksRequest to peer %s: %v", peerID, err)
//         return err
//     }

//     var response GetBlocksResponse
//     err = Receive(stream, &response)
//     if err != nil {
//         logger.ErrorLogger.Printf("Failed to receive GetBlocksResponse from peer %s: %v", peerID, err)
//         return err
//     }

//     return p.ProcessGetBlocksResponse(response)
// }

// // HandleGetBlocksRequest handles incoming GetBlocks requests.
// func (p *GetBlocksProtocol) HandleGetBlocksRequest(s network.Stream) {
//     defer s.Close()

//     var request GetBlocksRequest
//     err := Receive(s, &request)
//     if err != nil {
//         logger.ErrorLogger.Printf("Failed to receive GetBlocksRequest: %v", err)
//         return
//     }

//     logger.InfoLogger.Printf("Received GetBlocksRequest from %s", s.Conn().RemotePeer().Pretty())

//     blocks, err := p.getBlocksSinceHash(request.LastKnownHash)
//     if err != nil {
//         logger.ErrorLogger.Printf("Failed to get blocks since hash: %v", err)
//         Send(s, GetBlocksResponse{Success: false})
//         return
//     }

//     response := GetBlocksResponse{
//         Blocks:  blocks,
//         Success: true,
//     }

//     err = Send(s, response)
//     if err != nil {
//         logger.ErrorLogger.Printf("Failed to send GetBlocksResponse: %v", err)
//         return
//     }

//     logger.InfoLogger.Printf("Sent %d blocks to peer %s", len(blocks), s.Conn().RemotePeer().Pretty())
// }

// // ProcessGetBlocksResponse processes the response from a GetBlocks request.
// func (p *GetBlocksProtocol) ProcessGetBlocksResponse(response GetBlocksResponse) error {
//     if !response.Success {
//         logger.ErrorLogger.Println("GetBlocksResponse indicates failure")
//         return errors.New("failed to get blocks from peer")
//     }

//     logger.InfoLogger.Printf("Processing %d received blocks", len(response.Blocks))

// // Validate and integrate received blocks
//     for _, block := range response.Blocks {
//         err := p.Node.Blockchain.AddBlock(block)
//         if err != nil {
//             logger.ErrorLogger.Printf("Failed to add block %x: %v", block.Header.BlockHash, err)
//             // Handle invalid block: add transactions to mempool
//             for _, tx := range block.Transactions {
//                 p.Node.Mempool.AddTransaction(tx)
//             }
//             continue
//         }
//     }

//     logger.InfoLogger.Println("Successfully processed received blocks")
//     return nil
// }

// // Helper method to get blocks since a given hash.
// func (p *GetBlocksProtocol) getBlocksSinceHash(lastKnownHash []byte) ([]*blockchain.Block, error) {
//     blocks := p.Node.Blockchain.Ledger
//     var startIndex int = -1

//     for i, block := range blocks {
//         if bytes.Equal(block.Header.BlockHash, lastKnownHash) {
//             startIndex = i + 1
//             break
//         }
//     }

//     if startIndex == -1 {
//         logger.ErrorLogger.Println("Last known block hash not found")
//         return nil, errors.New("last known block hash not found")
//     }

//     return blocks[startIndex:], nil
// }

// // Send is a helper function to send data over a LibP2P stream.
// func Send(s network.Stream, v interface{}) error {
//     data, err := utils.Serialize(v)
//     if err != nil {
//         return err
//     }
//     _, err = s.Write(data)
//     return err
// }

// // Receive is a helper function to receive data over a LibP2P stream.
// func Receive(s network.Stream, v interface{}) error {
//     data := make([]byte, 1024*1024) // 1 MB buffer
//     n, err := s.Read(data)
//     if err != nil {
//         return err
//     }
//     return utils.Deserialize(data[:n], v)
// }
