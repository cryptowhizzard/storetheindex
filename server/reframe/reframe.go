package reframe

import (
	"bytes"
	"net/http"

	"github.com/filecoin-project/go-indexer-core"
	"github.com/filecoin-project/storetheindex/internal/registry"
	"github.com/filecoin-project/storetheindex/server/finder/handler"
	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-delegated-routing/client"
	"github.com/ipfs/go-delegated-routing/server"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/routing"
	"github.com/multiformats/go-multicodec"
	"github.com/multiformats/go-multihash"
	"github.com/multiformats/go-varint"
)

func NewReframeHTTPHandler(indexer indexer.Interface, registry *registry.Registry) http.HandlerFunc {
	return server.DelegatedRoutingAsyncHandler(NewReframeService(handler.NewFinderHandler(indexer, registry)))
}

func NewReframeService(fh *handler.FinderHandler) *ReframeService {
	return &ReframeService{finderHandler: fh}
}

type ReframeService struct {
	finderHandler *handler.FinderHandler
}

func (x *ReframeService) FindProviders(key cid.Cid) (<-chan client.FindProvidersAsyncResult, error) {
	mh := key.Hash()
	fr, err := x.finderHandler.Find([]multihash.Multihash{mh})
	if err != nil {
		return nil, err
	}
	ch := make(chan client.FindProvidersAsyncResult)
	peerAddrs := []peer.AddrInfo{}
	for _, mhr := range fr.MultihashResults {
		if !bytes.Equal(mhr.Multihash, mh) {
			continue
		}
		for _, pr := range mhr.ProviderResults {
			if !isBitswapMetadata(pr.Metadata) {
				continue
			}
			peerAddrs = append(peerAddrs, pr.Provider)
		}
	}
	go func() {
		defer close(ch)
		ch <- client.FindProvidersAsyncResult{AddrInfo: peerAddrs}
	}()
	return ch, nil
}

func (x *ReframeService) GetIPNS(id []byte) (<-chan client.GetIPNSAsyncResult, error) {
	return nil, routing.ErrNotSupported
}

func (x *ReframeService) PutIPNS(id []byte, record []byte) (<-chan client.PutIPNSAsyncResult, error) {
	return nil, routing.ErrNotSupported
}

var BitswapMetadataBytes = varint.ToUvarint(uint64(multicodec.TransportBitswap))

func isBitswapMetadata(meta []byte) bool {
	return bytes.Equal(meta, BitswapMetadataBytes)
}
