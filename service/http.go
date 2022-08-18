package service

import (
	"net/http"

	delivemux "github.com/aditya37/geofence-service/delivery/mux"
	"github.com/gorilla/mux"
)

type httpServer struct {
	// Deprecated....
	eventstate  *delivemux.EventstateDelivery
	geofenecing *delivemux.GeofenceDelivery
	muxrouter   *mux.Router
}

func NewHttpServer(geofenecing *delivemux.GeofenceDelivery) (*httpServer, error) {
	router := mux.NewRouter()
	return &httpServer{
		muxrouter:   router,
		geofenecing: geofenecing,
	}, nil
}

// handler
func (h *httpServer) handler() http.Handler {

	// Health check
	h.muxrouter.Methods(http.MethodGet).Path("/").HandlerFunc(h.healthCheckHandler)

	// prefix for event
	// Deprecated....
	eventStateRoute := h.muxrouter.PathPrefix("/event")
	eventStateRoute.Methods(http.MethodGet).Path("/state/{service_name}/{event_id}").HandlerFunc(h.eventstate.GetServiceEventState)

	// geofencing route...
	goefencingRoute := h.muxrouter.PathPrefix("/geofencing")
	goefencingRoute.Methods(http.MethodPost).Path("/").HandlerFunc(h.geofenecing.AddGeofence)

	return h.muxrouter
}

func (h *httpServer) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("up"))
}
