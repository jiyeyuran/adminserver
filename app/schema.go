package app

import "time"

type User struct {
	Id       int64     `json:"id,omitempty"`
	Name     string    `json:"name,omitempty"`
	Mobile   string    `json:"mobile,omitempty"`
	Password string    `json:"password,omitempty"`
	Ctime    time.Time `json:"ctime,omitempty"`
}
