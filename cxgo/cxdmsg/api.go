package cxdmsg

import (
	"net/http"

	"github.com/skycoin/dmsg/cipher"
	"github.com/skycoin/dmsg/httputil"

	"github.com/SkycoinProject/cx-chains/src/api"
	"github.com/SkycoinProject/cx-chains/src/skycoin"

	"github.com/SkycoinProject/cx/cxgo/cxspec"
)

// ChainStats contain chain stats.
type NodeStats struct {
	CXChainVersion     string        `json:"cx_chain_version"`
	CXChainSpecEra     string        `json:"cx_chain_spec_era"`
	GenesisBlockHash   cipher.SHA256 `json:"genesis_block_hash"`
	PrevBlockHash      cipher.SHA256 `json:"prev_block_hash"`
	HeadBlockHash      cipher.SHA256 `json:"head_block_hash"`
	HeadBlockTimestamp uint64        `json:"head_block_timestamp"`
	HeadBlockHeight    uint64        `json:"head_block_height"`
}

// API represents the api being served with dmsg.
type API struct {
	Version   string
	NodeConf  skycoin.NodeConfig
	ChainSpec cxspec.ChainSpec
	Gateway   api.Gatewayer
}

// ServeHTTP implements http.Handler.
func (api *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/stats", api.handleNodeStats())
	mux.ServeHTTP(w, r)
}

func (api *API) obtainNodeStats() (NodeStats, error) {
	chainMeta, err := api.Gateway.GetBlockchainMetadata()
	if err != nil {
		return NodeStats{}, err
	}

	genBlock, err := api.ChainSpec.GenerateGenesisBlock()
	if err != nil {
		return NodeStats{}, err
	}

	stats := NodeStats{
		CXChainVersion:     api.Version,
		CXChainSpecEra:     api.ChainSpec.SpecEra,
		GenesisBlockHash:   cipher.SHA256(genBlock.HashHeader()),
		PrevBlockHash:      cipher.SHA256(chainMeta.HeadBlock.Head.PrevHash),
		HeadBlockHash:      cipher.SHA256(chainMeta.HeadBlock.HashHeader()),
		HeadBlockTimestamp: chainMeta.HeadBlock.Head.Time,
		HeadBlockHeight:    chainMeta.HeadBlock.Head.BkSeq,
	}
	return stats, nil
}

func (api *API) handleNodeStats() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		stats, err := api.obtainNodeStats()
		if err != nil {
			httputil.WriteJSON(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		httputil.WriteJSON(w, r, http.StatusOK, stats)
	}
}
