package router

import (
	"github.com/Amir122002/hotel/internal/handlers"
	"github.com/gorilla/mux"
	"net/http"
)

func NewRouter(h *handlers.Handler) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/registration", h.Registration).Methods(http.MethodPost)
	router.HandleFunc("/get_token", h.GetToken).Methods(http.MethodGet)

	worker := router.PathPrefix("/worker").Subrouter()
	worker.Use(h.CheckToken)

	worker.HandleFunc("/enter", h.EnterSystem).Methods(http.MethodPost)
	worker.HandleFunc("/exit", h.ExitSystem).Methods(http.MethodPost)

	worker.HandleFunc("/read_hotel_room/{page}", h.ReadHotelRoom).Methods(http.MethodGet)

	worker.HandleFunc("/creat_client", h.CreatClient).Methods(http.MethodPost)
	worker.HandleFunc("/read_client", h.ReadClient).Methods(http.MethodGet)
	worker.HandleFunc("/delete_client", h.DeleteClient).Methods(http.MethodPost)

	client := router.PathPrefix("/client").Subrouter()
	client.HandleFunc("/table_reservation", h.TableReservation).Methods(http.MethodPost)
	client.HandleFunc("/taxi_ordering", h.TaxiOrdering).Methods(http.MethodPost)

	waiter := router.PathPrefix("/waiter").Subrouter()
	waiter.Use(h.CheckToken)
	waiter.HandleFunc("/read_table_reservation/{page}", h.ReadTableReservation).Methods(http.MethodGet)

	taxi := router.PathPrefix("").Subrouter()
	taxi.Use(h.CheckToken)
	taxi.HandleFunc("/read_taxi_ordering/{page}", h.ReadTaxiOrdering).Methods(http.MethodGet)

	return router
}
