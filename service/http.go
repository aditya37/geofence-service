package service

import (
	"net/http"

	delivemux "github.com/aditya37/geofence-service/delivery/mux"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
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

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	// Health check
	h.muxrouter.Methods(http.MethodGet).Path("/").HandlerFunc(h.healthCheckHandler)

	// prefix for event
	// Deprecated....
	eventStateRoute := h.muxrouter.PathPrefix("/event")
	eventStateRoute.Methods(http.MethodGet).Path("/state/{service_name}/{event_id}").HandlerFunc(h.eventstate.GetServiceEventState)

	// geofencing route...
	goefencingRoute := h.muxrouter.PathPrefix("/geofencing").Subrouter()
	goefencingRoute.Methods(http.MethodPost).Path("/").HandlerFunc(h.geofenecing.AddGeofence)
	goefencingRoute.Methods(http.MethodGet).Path("/type").Queries("id", "", "name", "").HandlerFunc(h.geofenecing.GetGeofenceTypeDetail)
	goefencingRoute.Methods(http.MethodGet).Path("/counts").HandlerFunc(h.geofenecing.GetGeofenceCount)
	goefencingRoute.Methods(http.MethodGet).Path("/{id}").HandlerFunc(h.geofenecing.GetGeofenceById)
	goefencingRoute.Methods(http.MethodGet).Path("/location/{location_id}").HandlerFunc(h.geofenecing.GetGeofenceByLocationId)
	goefencingRoute.Methods(http.MethodGet).Path("/{type}/").Queries("page", "", "itemPerPage", "").HandlerFunc(h.geofenecing.GetGeofenceByType)

	// get mobility....
	goefencingRoute.Methods(http.MethodGet).Path("/mobility/{geofence_id}/avg").Queries("interval", "").HandlerFunc(h.geofenecing.GetAggregateMobilityByArea)
	// get geofenecing type...
	return corsHandler.Handler(h.muxrouter)
}

func (h *httpServer) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("up"))
}
