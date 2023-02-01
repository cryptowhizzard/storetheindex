package httpfinderserver

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/ipfs/go-cid"
	indexer "github.com/ipni/go-indexer-core"
	coremetrics "github.com/ipni/go-indexer-core/metrics"
	"github.com/ipni/storetheindex/api/v0/finder/model"
	"github.com/ipni/storetheindex/internal/counter"
	"github.com/ipni/storetheindex/internal/httpserver"
	"github.com/ipni/storetheindex/internal/metrics"
	"github.com/ipni/storetheindex/internal/registry"
	"github.com/ipni/storetheindex/server/finder/handler"
	"github.com/ipni/storetheindex/version"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multihash"
	"go.opencensus.io/stats"
	"go.opencensus.io/tag"
)

var (
	versionData []byte
	newline     = []byte("\n")
)

func init() {
	versionData, _ = json.Marshal(version.String())
}

// handler handles requests for the finder resource
type httpHandler struct {
	finderHandler *handler.FinderHandler
}

func newHandler(indexer indexer.Interface, registry *registry.Registry, indexCounts *counter.IndexCounts) *httpHandler {
	return &httpHandler{
		finderHandler: handler.NewFinderHandler(indexer, registry, indexCounts),
	}
}

func (h *httpHandler) find(w http.ResponseWriter, r *http.Request) {
	stream, err := explicitlyAcceptsNDJson(r)
	if err != nil {
		http.Error(w, "invalid Accept header: "+err.Error(), http.StatusBadRequest)
		return
	}
	vars := mux.Vars(r)
	mhVar := vars["multihash"]
	m, err := multihash.FromB58String(mhVar)
	if err != nil {
		log.Errorw("error decoding multihash", "multihash", mhVar, "err", err)
		httpserver.HandleError(w, err, "find")
		return
	}
	h.getIndexes(w, []multihash.Multihash{m}, stream)
}

func (h *httpHandler) findCid(w http.ResponseWriter, r *http.Request) {
	stream, err := explicitlyAcceptsNDJson(r)
	if err != nil {
		http.Error(w, "invalid Accept header: "+err.Error(), http.StatusBadRequest)
		return
	}
	vars := mux.Vars(r)
	cidVar := vars["cid"]
	c, err := cid.Decode(cidVar)
	if err != nil {
		log.Errorw("error decoding cid", "cid", cidVar, "err", err)
		httpserver.HandleError(w, err, "find")
		return
	}
	h.getIndexes(w, []multihash.Multihash{c.Hash()}, stream)
}

func (h *httpHandler) findBatch(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Errorw("failed reading get batch request", "err", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	req, err := model.UnmarshalFindRequest(body)
	if err != nil {
		log.Errorw("error unmarshalling get batch request", "err", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	h.getIndexes(w, req.Multihashes, false)
}

func (h *httpHandler) getIndexes(w http.ResponseWriter, mhs []multihash.Multihash, stream bool) {
	if len(mhs) != 1 && stream {
		log.Errorw("Streaming response is not supported for batch find")
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	startTime := time.Now()
	var found bool
	defer func() {
		msecPerMh := coremetrics.MsecSince(startTime) / float64(len(mhs))
		_ = stats.RecordWithOptions(context.Background(),
			stats.WithTags(tag.Insert(metrics.Method, "http"), tag.Insert(metrics.Found, fmt.Sprintf("%v", found))),
			stats.WithMeasurements(metrics.FindLatency.M(msecPerMh)))
	}()

	response, err := h.finderHandler.Find(mhs)
	if err != nil {
		httpserver.HandleError(w, err, "get")
		return
	}

	// If no info for any multihashes, then 404
	if len(response.MultihashResults) == 0 {
		http.Error(w, "no results for query", http.StatusNotFound)
		return
	}

	if stream {
		log := log.With("mh", mhs[0].B58String())
		pr := response.MultihashResults[0].ProviderResults
		if len(pr) == 0 {
			http.Error(w, "no results for query", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", mediaTypeNDJson)
		w.Header().Set("Connection", "Keep-Alive")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		flusher, flushable := w.(http.Flusher)
		encoder := json.NewEncoder(w)
		var count int
		for _, result := range pr {
			if err := encoder.Encode(result); err != nil {
				log.Errorw("Failed to encode streaming response", "err", err)
				break
			}
			if _, err := w.Write(newline); err != nil {
				log.Errorw("failed to write newline while streaming results", "err", err)
				break
			}
			// TODO: optimise the number of time we call flush based on some time-based or result
			//       count heuristic.
			if flushable {
				flusher.Flush()
			}
			count++
		}
		if count == 0 {
			log.Errorw("Failed to encode results; falling back on not found", "resultsCount", len(pr))
			http.Error(w, "no results for query", http.StatusNotFound)
			return
		}
		return
	}

	rb, err := model.MarshalFindResponse(response)
	if err != nil {
		log.Errorw("failed marshalling query response", "err", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	found = true
	httpserver.WriteJsonResponse(w, http.StatusOK, rb)
}

// ----- provider handlers -----

// GET /providers",
func (h *httpHandler) listProviders(w http.ResponseWriter, r *http.Request) {
	data, err := h.finderHandler.ListProviders()
	if err != nil {
		log.Errorw("cannot list providers", "err", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	httpserver.WriteJsonResponse(w, http.StatusOK, data)
}

// GET /providers/{providerid}
func (h *httpHandler) getProvider(w http.ResponseWriter, r *http.Request) {
	providerID, err := getProviderID(r)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	data, err := h.finderHandler.GetProvider(providerID)
	if err != nil {
		log.Error("cannot get provider", "err", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	if len(data) == 0 {
		http.Error(w, "provider not found", http.StatusNotFound)
		return
	}

	httpserver.WriteJsonResponse(w, http.StatusOK, data)
}

// GET /stats",
func (h *httpHandler) getStats(w http.ResponseWriter, r *http.Request) {
	data, err := h.finderHandler.GetStats()
	switch {
	case err != nil:
		log.Errorw("cannot get stats", "err", err)
		http.Error(w, "", http.StatusInternalServerError)
	case len(data) == 0:
		log.Warn("processing stats")
		http.Error(w, "processing", http.StatusTeapot)
	default:
		httpserver.WriteJsonResponse(w, http.StatusOK, data)
	}
}

func getProviderID(r *http.Request) (peer.ID, error) {
	vars := mux.Vars(r)
	pid := vars["providerid"]
	providerID, err := peer.Decode(pid)
	if err != nil {
		return providerID, fmt.Errorf("cannot decode provider id: %s", err)
	}
	return providerID, nil
}

func (h *httpHandler) health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache")
	httpserver.WriteJsonResponse(w, http.StatusOK, versionData)
}
