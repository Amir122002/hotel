package models

import (
	"time"
)

type Config struct {
	Server Server
	DB     DB
}

type DB struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
	Sslmode  string `json:"sslmode"`
}

type Server struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

type Workers struct {
	Id       int       `json:"id"`
	FullName string    `json:"full_name"`
	Login    string    `json:"login"`
	Password string    `json:"password"`
	JobTitle int       `json:"job_title"`
	Active   bool      `json:"active"`
	CreateAt time.Time `json:"create_at"`
	UpdateAt time.Time `json:"update_at"`
	DeleteAt time.Time `json:"delete_at"`
}

type WorkersTokens struct {
	Id        int       `json:"id"`
	Token     string    `json:"token"`
	UserId    int       `json:"user_id"`
	Active    bool      `json:"active"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

type WorkingHours struct {
	Id         int       `json:"id"`
	UserId     int       `json:"user_id"`
	Active     bool      `json:"active"`
	StartWork  time.Time `json:"start_work"`
	FinishWork time.Time `json:"finish_work"`
}

type Clients struct {
	Id         int       `json:"id"`
	FullName   string    `json:"fill_name"`
	NumberRoom int       `json:"number_room"`
	Active     bool      `json:"active"`
	CreateAt   time.Time `json:"create_at"`
	UpdateAt   time.Time `json:"update_at"`
	DeleteAt   time.Time `json:"delete_at"`
}
type HotelRoom struct {
	Id             int    `json:"id"`
	NumberRoom     string `json:"number_room"`
	HotelRoomTypes int    `json:"hotel_room_types"`
	Active         bool   `json:"active"`
	PriceRoom      int    `json:"price_room"`
}
type Reservations struct {
	Id                int       `json:"id"`
	NumberRoom        int       `json:"number_room"`
	TableId           int       `json:"table_id"`
	TimeOfReservation time.Time `json:"time_of_reservation"`
}

type RestaurantTables struct {
	Id          int  `json:"id"`
	TableType   int  `json:"table_type"`
	NumberTable int  `json:"number_table"`
	Active      bool `json:"active"`
}
type TaxiOrdering struct {
	Id                 int       `json:"id"`
	NumberRoom         int       `json:"number_room"`
	TimeOfTaxiOrdering time.Time `json:"time_of_taxi_ordering"`
}

type Bill struct {
	Id                      int `json:"id"`
	UserId                  int `json:"user_id"`
	NumberRoom              int `json:"number_room"`
	PaymentForAccommodation int `json:"payment_for_accommodation"`
}
