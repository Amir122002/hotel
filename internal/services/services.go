package services

import (
	"fmt"
	"github.com/Amir122002/hotel/internal/repositories"
	"github.com/Amir122002/hotel/pkg/models"
	"github.com/sirupsen/logrus"
	"os"
)

type Service struct {
	Repository *repositories.Repository
	logger     *logrus.Logger
}

func NewService(Repository *repositories.Repository, logger *logrus.Logger) *Service {
	return &Service{
		Repository: Repository,
		logger:     logger,
	}
}

//func NewService(Repository *repositories.Repository) *Service {
//	return &Service{Repository: Repository}
//}

func (s *Service) Registration(newWorker *models.Workers) error {
	return s.Repository.Registration(newWorker)
}

func (s *Service) GetToken(login, password string) (string, error) {
	return s.Repository.GetToken(login, password)
}

func (s *Service) CheckToken(token string) (int, error) {
	return s.Repository.CheckToken(token)
}

func (s *Service) EnterSystem(userID int) error {
	return s.Repository.EnterSystem(userID)
}

func (s *Service) ExitSystem(userID int) error {
	return s.Repository.ExitSystem(userID)
}

func (s *Service) ReadHotelRoom(userID, page int) (*[]models.HotelRoom, error) {
	return s.Repository.ReadHotelRoom(userID, page)
}

func (s *Service) CreateClient(userID int, client *models.Clients) error {
	return s.Repository.CreateClient(userID, client)
}

func (s *Service) ReadClient(userID int, client string) (*models.Clients, error) {
	return s.Repository.ReadClient(userID, client)
}

func (s *Service) DeleteClient(userID int, clientID int) error {
	clientTime, err := s.Repository.DeleteClient(userID, clientID)
	if err != nil {
		return err
	}

	// Parse the timestamps
	duration := clientTime.DeleteAt.Sub(clientTime.CreateAt)

	// Преобразовать длительность в дни
	daysDifference := duration.Hours() / 24

	bill := &models.Bill{
		UserId:                  clientTime.Id,
		NumberRoom:              clientTime.NumberRoom,
		PaymentForAccommodation: int(daysDifference) + 1,
	}

	totalAmount, err := s.Repository.Bill(bill)
	if err != nil {
		return err
	}

	file, err := os.Create("bill.txt")
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = fmt.Fprintf(file, "Счет для комнаты номер %d\n", clientTime.NumberRoom)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(file, "Общая сумма: %.d рублей\n", totalAmount)
	if err != nil {
		return err
	}

	return nil
}
func (s *Service) Restaurant(reservation *models.Reservations) error {
	return s.Repository.Restaurant(reservation)
}

func (s *Service) ReadTableReservation(page, userID int) (*[]models.Reservations, error) {
	return s.Repository.ReadTableReservation(page, userID)
}

func (s *Service) TaxiOrdering(taxi *models.TaxiOrdering) error {
	return s.Repository.TaxiOrdering(taxi)
}

func (s *Service) ReadTaxiOrdering(page, userID int) (*models.TaxiOrdering, error) {
	return s.Repository.ReadTaxiOrdering(page, userID)
}
