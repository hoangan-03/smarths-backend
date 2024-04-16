package models

import (
	"time"

	_ "github.com/lib/pq"
)

type Account struct {
	Acc_id   int    `json:"acc_id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Room struct {
	Room_id int `json:"room_id"`
	Acc_id  int `json:"acc_id"`
}

type Device struct {
	Sens_id int    `json:"aens_id"`
	Acc_id  int    `json:"acc_id"`
	Status  string `json:"status"`
}

type Record struct {
	Rec_id  int    `json:"rec_id"`
	Sens_id int    `json:"sens_id"`
	Value   string `json:"value"`
}

type Booking struct {
	Book_id int       `json:"book_id"`
	Acc_id  int       `json:"acc_id"`
	Day     int       `json:"day"`
	Month   int       `json:"month"`
	Time    time.Time `json:"time"`
}

type Reminder struct {
	Rem_id      int       `json:"rem_id"`
	Time_on     int       `json:"time_on"`
	Remain_time time.Time `json:"remain_time"`
	Acc_id      int       `json:"acc_id"`
	Book_id     int       `json:"book_id"`
	Notes       string    `json:"notes"`
}

type Controlling struct {
	Cmd_id int    `json:"cmd_id"`
	Acc_id int    `json:"acc_id"`
	Status string `json:"status"`
	Action string `json:"action"`
}
