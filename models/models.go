package models

import (
	"time"

	_ "github.com/lib/pq"
)

type Account struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Acc_id   int    `json:"acc_id"`
	Key      string `json:"key"`
}

type Room struct {
	Room_id   int    `json:"room_id"`
	Acc_id    string `json:"acc_id"`
	Room_name string `json:"room_name"`
	Mode      string `json:"mode"`
}

type Device struct {
	Dev_id   string `json:"dev_id"`
	Room_id  int    `json:"room_id"`
	Category string `json:"category"`
}

type Record struct {
	Rec_id    int       `json:"rec_id"`
	Dev_id    string    `json:"dev_id"`
	Value     string    `json:"value"`
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
}

type Booking struct {
	Book_id     int       `json:"book_id"`
	Room_id     int       `json:"room_id"`
	Start_time  time.Time `json:"start_time"`
	Notes       string    `json:"notes"`
	Remind_time time.Time `json:"remind_time"`
	End_time    time.Time `json:"end_time"`
}

type Criteria struct {
	Crit_id   string `json:"crit_id"`
	Dev_id    string `json:"dev_id"`
	Threshold string `json:"threshold"`
	Action    string `json:"action"`
}

type Controlling struct {
	Ctrl_id   int       `json:"ctrl_id"`
	Dev_id    string    `json:"dev_id"`
	Room_id   int       `json:"room_id"`
	Action    string    `json:"action"`
	Ctrl_mode string    `json:"ctrl_mode"`
	Timestamp time.Time `json:"timestamp"`
}
