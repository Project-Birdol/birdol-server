package model

import "time"

type LinkPassword struct {
	Password	string
	ExpireDate	time.Time
}