package service

import (
	"net/http"

	delivemux "github.com/aditya37/geofence-service/delivery/mux"
	"github.com/gorilla/mux"
)

type httpServer struct {
	eventstate *delivemux.EventstateDelivery
	muxrouter  *mux.Router
}

func NewHttpServer(evenstate *delivemux.EventstateDelivery) (*httpServer, error) {
	router := mux.NewRouter()
	return &httpServer{
		eventstate: evenstate,
		muxrouter:  router,
	}, nil
}

// handler
func (h *httpServer) handler() http.Handler {
	// Health check
	h.muxrouter.Methods(http.MethodGet).Path("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("up"))
	})

	// prefix for event
	eventStateRoute := h.muxrouter.PathPrefix("/event")
	eventStateRoute.Methods(http.MethodGet).Path("/state/{service_name}/{event_id}").HandlerFunc(h.eventstate.GetServiceEventState)
	return h.muxrouter
}
