package network

import "trustify/blockchain"

type GetBlocksProtocol struct {
	timeout int
}

func (p *GetBlocksProtocol) GetBlocks(peer string, lastKnownHash string) error {
	// TODO: Implement the method
	return nil
}

func (p *GetBlocksProtocol) HandleGetBlocksRequest(request GetBlocksRequest) error {
	// TODO: Implement the method
	return nil
}

func (p *GetBlocksProtocol) ProcessGetBlocksResponse(response GetBlocksResponse) error {
	// TODO: Implement the method
	return nil
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
