package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Amir122002/hotel/internal/services"
	"github.com/Amir122002/hotel/pkg/models"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
)

type Handler struct {
	Service *services.Service
	logger  *logrus.Logger
}

func NewHandler(Service *services.Service, logger *logrus.Logger) *Handler {
	return &Handler{
		Service: Service,
		logger:  logger,
	}
}

//func NewHandler(Service *services.Service) *Handler {
//	return &Handler{Service: Service}
//}

func (h *Handler) Registration(w http.ResponseWriter, r *http.Request) {
	var worker models.Workers
	err := json.NewDecoder(r.Body).Decode(&worker)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.logger.Error(err)
		//fmt.Println(err)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(worker.Password), bcrypt.DefaultCost)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.logger.Error(err)
		return
	}

	newWorker := &models.Workers{
		FullName: worker.FullName,
		Login:    worker.Login,
		Password: string(hashedPassword),
		JobTitle: worker.JobTitle,
	}

	err = h.Service.Registration(newWorker)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.logger.Error(err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) GetToken(w http.ResponseWriter, r *http.Request) {
	login := r.Header.Get("login")
	password := r.Header.Get("password")

	token, err := h.Service.GetToken(login, password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.logger.Error(err)
		return
	}

	w.Header().Set("Token", token)
	fmt.Println(token)
	w.WriteHeader(http.StatusOK)

}

func (h *Handler) CheckToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("token")
		fmt.Println("sdasd")
		userId, err := h.Service.CheckToken(token)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			h.logger.Error(err)
			return
		}
		ctx := context.WithValue(r.Context(), "userID", userId)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func (h *Handler) EnterSystem(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)
	fmt.Println("2")
	err := h.Service.EnterSystem(userID)
	if err != nil {
		http.Error(w, "Enter Error", http.StatusUnauthorized)
		h.logger.Error(err)
		return
	}
	w.WriteHeader(http.StatusOK)
	return
}

func (h *Handler) ExitSystem(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)

	err := h.Service.ExitSystem(userID)
	if err != nil {
		http.Error(w, "Exit Error", http.StatusUnauthorized)
		h.logger.Error(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}

func (h *Handler) ReadHotelRoom(w http.ResponseWriter, r *http.Request) {

	userID := r.Context().Value("userID").(int)
	vars := mux.Vars(r)
	pageStr := vars["page"]
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.logger.Error(err)
		return
	}

	Rooms, err := h.Service.ReadHotelRoom(userID, page)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.logger.Error(err)
		return
	}

	data, err := json.Marshal(Rooms)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.logger.Error(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.logger.Error(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return
}

func (h *Handler) CreatClient(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("userID").(int)
	var client models.Clients
	err := json.NewDecoder(r.Body).Decode(&client)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.logger.Error(err)
		return
	}

	clients := &models.Clients{
		FullName:   client.FullName,
		NumberRoom: client.NumberRoom,
	}

	err = h.Service.CreateClient(userId, clients)
	if err != nil {
		http.Error(w, "CreateClient Error", http.StatusUnauthorized)
		h.logger.Error(err)
		return
	}
	w.WriteHeader(http.StatusOK)

	return
}

func (h *Handler) ReadClient(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)
	client := r.Header.Get("client_id")

	Client, err := h.Service.ReadClient(userID, client)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.logger.Error(err)
		return
	}

	data, err := json.Marshal(Client)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.logger.Error(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.logger.Error(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}

func (h *Handler) DeleteClient(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)
	clientStr := r.Header.Get("client_id")
	clientID, err := strconv.Atoi(clientStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.logger.Error(err)
		return
	}
	err = h.Service.DeleteClient(userID, clientID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.logger.Error(err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

//func client Restaurant{

func (h *Handler) TableReservation(w http.ResponseWriter, r *http.Request) {
	var reservations models.Reservations
	err := json.NewDecoder(r.Body).Decode(&reservations)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.logger.Error(err)
		return
	}
	reservation := &models.Reservations{
		NumberRoom:        reservations.NumberRoom,
		TableId:           reservations.TableId,
		TimeOfReservation: reservations.TimeOfReservation,
	}

	err = h.Service.Restaurant(reservation)
	if err != nil {
		h.logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}

func (h *Handler) ReadTableReservation(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)
	vars := mux.Vars(r)
	pageStr := vars["page"]
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.logger.Error(err)
		return
	}

	Table, err := h.Service.ReadTableReservation(page, userID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.logger.Error(err)
		return
	}

	data, err := json.Marshal(Table)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.logger.Error(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.logger.Error(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}

func (h *Handler) TaxiOrdering(w http.ResponseWriter, r *http.Request) {
	fmt.Println("q")
	var taxi models.TaxiOrdering
	err := json.NewDecoder(r.Body).Decode(&taxi)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.logger.Error(err)
		return
	}
	ordering := &models.TaxiOrdering{
		NumberRoom:         taxi.NumberRoom,
		TimeOfTaxiOrdering: taxi.TimeOfTaxiOrdering,
	}

	err = h.Service.TaxiOrdering(ordering)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.logger.Error(err)
		return
	}
	w.WriteHeader(http.StatusOK)
	return
}

func (h *Handler) ReadTaxiOrdering(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)
	vars := mux.Vars(r)
	pageStr := vars["page"]
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.logger.Error(err)
		return
	}

	Taxi, err := h.Service.ReadTaxiOrdering(page, userID)
	if err != nil {
		h.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Println("taxi", Taxi)
	data, err := json.Marshal(Taxi)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.logger.Error(err)
		return
	}

	//w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.logger.Error(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}

//func client Cleaning

func cleaning() {

}
