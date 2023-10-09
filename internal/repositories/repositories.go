package repositories

import (
	"errors"
	"fmt"
	"github.com/Amir122002/hotel/internal/database"
	"github.com/Amir122002/hotel/pkg/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
)

type Repository struct {
	db     *database.DB
	logger *logrus.Logger
}

func NewRepository(db *database.DB, logger *logrus.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

func (r *Repository) Registration(newWorker *models.Workers) error {
	var worker models.Workers
	err := r.db.Pool.Raw("select * from workers where login=$1", newWorker.Login).Scan(&worker).Error
	if err != nil {
		return err
	} else if worker.Login != "" {
		return errors.New("Такой логин уже существует")
	}
	err = r.db.Pool.
		Exec("INSERT INTO Workers (full_name,login,password,job_title) values ($1,$2,$3,$4)", newWorker.FullName, newWorker.Login, newWorker.Password, newWorker.JobTitle).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetToken(login, password string) (string, error) {
	var workers models.Workers
	err := r.db.Pool.Raw("select * from workers where login=$1 and active=true", login).Scan(&workers).Error
	if err != nil {
		return "", err
	}
	err = bcrypt.CompareHashAndPassword([]byte(workers.Password), []byte(password))
	if err != nil {
		return "", err
	}

	EndTime := time.Now().Add(time.Hour * 3).Unix()

	tokenClaims := jwt.MapClaims{
		"worker_id": workers.Id,
		"exp":       EndTime,
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims)
	Token, err := jwtToken.SignedString([]byte("secret"))
	if err != nil {
		return "", err
	}
	//Token, EndTime, err := GenerateToken(workers)
	err = r.db.Pool.Exec("INSERT into workers_tokens(token,user_id,end_time) values ($1,$2,$3)", Token, workers.Id, time.Unix(EndTime, 0)).Error
	if err != nil {
		return "", err
	}

	return Token, nil
}

func (r *Repository) CheckToken(token string) (int, error) {
	var workerToken models.WorkersTokens
	err := r.db.Pool.Raw("SELECT user_id,end_time from workers_tokens where token=$1 and end_time>$2", token, time.Now()).Scan(&workerToken).Error
	if err != nil {
		return 0, err
	}
	//&& workerToken.EndTime.IsZero()
	if workerToken.UserId == 0 {
		log.Println("Токен не найден или истек!")
		return 0, errors.New("Токен не найден или истек")
	}

	return workerToken.UserId, nil
}

func (r *Repository) EnterSystem(userID int) error {
	var worker models.Workers
	fmt.Println("id", userID)
	err := r.db.Pool.Raw("select * from workers where id=$1", userID).Scan(&worker).Error
	if err != nil {
		return err
	}
	var workingHours models.WorkingHours
	err = r.db.Pool.Raw("SELECT * FROM working_hours WHERE user_id=$1 and finish_work IS NULL;", userID).Scan(&workingHours).Error
	if err != nil {
		return err
	}
	fmt.Println("id", workingHours.UserId)
	if workingHours.UserId != 0 {
		log.Println("Вы не закончили работу!")
		return errors.New("Вы не закончили работу!")
	}

	err = r.db.Pool.Exec("insert into working_hours(user_id) values($1)", worker.Id).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) ExitSystem(userID int) error {
	var workingHours models.WorkingHours
	fmt.Println("id", userID)
	err := r.db.Pool.Raw("SELECT * FROM working_hours WHERE date(start_work) = CURRENT_DATE and user_id=$1 and finish_work IS NULL;", userID).Scan(&workingHours).Error
	if err != nil {
		log.Println(err)
		return err
	}
	if workingHours.UserId == 0 {
		return errors.New("Вы не начили работу! ")
	}
	err = r.db.Pool.Exec("UPDATE working_hours set finish_work=current_timestamp,active=false where id=$1", workingHours.Id).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) ReadHotelRoom(userID, page int) (*[]models.HotelRoom, error) {
	var working models.WorkingHours
	err := r.db.Pool.Raw("SELECT * from working_hours where user_id=$1 and active=true", userID).Scan(&working).Error
	if err != nil {

		return nil, err
	} else if working.Id == 0 {

		return nil, errors.New("Работник не начил смену ")
	}

	var workers models.Workers
	err = r.db.Pool.Raw("SELECT * from workers where id=$1 and active=true", userID).Scan(&workers).Error
	if err != nil {
		return nil, err
	}
	if workers.JobTitle != 1 {
		return nil, errors.New("У вас не доступа! ")
	}

	perPage := 10
	offset := (page - 1) * perPage

	var hotelRoom *[]models.HotelRoom
	err = r.db.Pool.Raw("SELECT * FROM hotel_rooms WHERE active = true ORDER BY id LIMIT $1 OFFSET $2", perPage, offset).Scan(&hotelRoom).Error
	if err != nil {

		return nil, err
	}

	return hotelRoom, nil
}

func (r *Repository) CreateClient(userID int, client *models.Clients) error {
	var working models.WorkingHours
	err := r.db.Pool.Raw("SELECT * from working_hours where user_id=$1 and active=true", userID).Scan(&working).Error
	if err != nil {
		return err
	} else if working.Id == 0 {
		return errors.New("Работник не начил смену ")
	}

	var workers models.Workers
	err = r.db.Pool.Raw("SELECT * from workers where id=$1 and active=true", userID).Scan(&workers).Error
	if err != nil {
		return err
	}
	if workers.JobTitle != 2 {
		return errors.New("У вас не доступа! ")
	}
	var room models.HotelRoom
	err = r.db.Pool.Raw("SELECT * from hotel_rooms where number_room=$1 and active=true", client.NumberRoom).Scan(&room).Error
	if err != nil {
		return err
	}
	if room.Id == 0 {
		log.Println("ERROR 2")
		return errors.New("ERROR 2")
	}

	err = r.db.Pool.Exec("INSERT into clients(full_name,number_room) values ($1,$2)", client.FullName, client.NumberRoom).Error
	if err != nil {
		return err
	}

	err = r.db.Pool.Exec("UPDATE hotel_rooms set active=false where number_room=$1", client.NumberRoom).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) ReadClient(userID int, client string) (*models.Clients, error) {
	var working models.WorkingHours
	err := r.db.Pool.Raw("SELECT * from working_hours where user_id=$1 and active=true", userID).Scan(&working).Error
	if err != nil {
		return nil, err
	} else if working.Id == 0 {
		return nil, errors.New("Работник не начил смену ")
	}

	var workers models.Workers
	err = r.db.Pool.Raw("SELECT * from workers where id=$1 and active=true", userID).Scan(&workers).Error
	if err != nil {
		return nil, err
	}
	if workers.JobTitle != 2 {
		return nil, errors.New("У вас не доступа! ")
	}
	var person models.Clients
	err = r.db.Pool.Raw("SELECT * from clients where id=$1 ", client).Scan(&person).Error
	if err != nil {
		return nil, err
	} else if person.Id == 0 {
		return nil, errors.New("Такого клиента нет! ")
	}

	return &person, nil
}

func (r *Repository) DeleteClient(userID int, clientID int) (*models.Clients, error) {

	var working models.WorkingHours
	err := r.db.Pool.Raw("SELECT * from working_hours where user_id=$1 and active=true", userID).Scan(&working).Error
	if err != nil {
		return nil, err
	} else if working.Id == 0 {
		return nil, errors.New("Работник не начил смену ")
	}

	var workers models.Workers
	err = r.db.Pool.Raw("SELECT * from workers where id=$1 and active=true", userID).Scan(&workers).Error
	if err != nil {
		return nil, err
	}
	if workers.JobTitle != 2 {
		return nil, errors.New("У вас не доступа! ")
	}

	var client models.Clients
	err = r.db.Pool.Raw("select * from clients where id=$1 and active=true", clientID).Scan(&client).Error
	if err != nil {
		return nil, err
	}
	if client.Id == 0 {
		return nil, errors.New("Клиент не найден ")
	}
	err = r.db.Pool.Exec("UPDATE hotel_rooms set active=true where number_room=$1", client.NumberRoom).Error
	if err != nil {
		return nil, err
	}

	err = r.db.Pool.Exec("UPDATE clients set active=false, delete_at=current_timestamp where id=$1", client.Id).Error
	if err != nil {
		return nil, err
	}

	var clientTime models.Clients
	err = r.db.Pool.Raw("select * from clients where id=$1", clientID).Scan(&clientTime).Error
	if err != nil {
		return nil, err
	}

	return &clientTime, nil
}

func (r *Repository) Bill(bill *models.Bill) (int, error) {
	var room models.HotelRoom
	err := r.db.Pool.Raw("select * from hotel_rooms where number_room=$1", bill.NumberRoom).Scan(&room).Error
	if err != nil {
		return 0, err
	}
	totalAmount := bill.PaymentForAccommodation * room.PriceRoom

	err = r.db.Pool.Exec("insert into bill(user_id,number_room,payment_for_accommodation) values($1,$2,$3)", bill.UserId, bill.NumberRoom, totalAmount).Error
	if err != nil {
		return 0, err
	}
	return totalAmount, nil
}

func (r *Repository) Restaurant(reservation *models.Reservations) error {
	var restaurantTables models.RestaurantTables
	err := r.db.Pool.Raw("SELECT * from restaurant_tables where id=$1 and active=true", reservation.TableId).Scan(&restaurantTables).Error
	if err != nil {
		return err
	} else if restaurantTables.Active == false {

		return errors.New("Этот стол уже забронирован! ")
	}

	err = r.db.Pool.Exec("Insert into reservations(table_id,time_of_reservation,number_room) values ($1,$2,$3)", reservation.TableId, reservation.TimeOfReservation, reservation.NumberRoom).Error
	if err != nil {
		return err
	}

	err = r.db.Pool.Exec("Update restaurant_tables set active=false where id=$1", reservation.TableId).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) ReadTableReservation(page, userID int) (*[]models.Reservations, error) {
	var workers models.Workers
	err := r.db.Pool.Raw("SELECT * from workers where id=$1 and active=true", userID).Scan(&workers).Error
	if err != nil {
		return nil, err
	}
	if workers.JobTitle != 1 {
		return nil, errors.New("У вас не доступа! ")
	}

	perPage := 10
	offset := (page - 1) * perPage

	var reservation []models.Reservations
	err = r.db.Pool.Raw("SELECT * from reservations LIMIT $1 OFFSET $2", perPage, offset).Scan(&reservation).Error
	if err != nil {
		return nil, err
	}

	return &reservation, nil
}

func (r *Repository) TaxiOrdering(taxi *models.TaxiOrdering) error {
	err := r.db.Pool.Exec("Insert into taxi_ordering(number_room,time_of_taxi_ordering) values ($1,$2)", taxi.NumberRoom, taxi.TimeOfTaxiOrdering).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) ReadTaxiOrdering(page, userID int) (*models.TaxiOrdering, error) {

	err := r.db.Pool.Raw("SELECT * from workers where id=$1 and active=true", userID).Error
	if err != nil {
		return nil, err
	}

	perPage := 10
	offset := (page - 1) * perPage

	var taxi models.TaxiOrdering
	err = r.db.Pool.Raw("SELECT * from taxi_ordering where date(time_of_taxi_ordering) = CURRENT_DATE  LIMIT $1 OFFSET $2", perPage, offset).Scan(&taxi).Error
	if err != nil {
		return nil, err
	}
	log.Println("eee", taxi)
	return &taxi, nil
}
