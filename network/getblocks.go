package network

import "trustify/blockchain"

type GetBlocksProtocol struct {
	timeout int
}

type GetBlocksRequest struct {
	LastKnownHash []byte
	Peer          string
}

type GetBlocksResponse struct {
	Blocks  []blockchain.Block
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
