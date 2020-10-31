package cxspec

type TrustedNodes struct {
	ChainPubKey    string   `json:"chain_pubkey"`    // public key of associated chain
	PublisherNodes []string `json:"publisher_nodes"` // addresses of publisher nodes
	AppNodes       []string `json:"app_nodes"`       // addresses of app nodes
}

