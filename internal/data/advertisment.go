package data

import (
	"time"
)

type Advertisment struct {
	Id          int
	Link        string
	Price       int
	PublishDate time.Time
	Flat        Flat
	IsValid     bool
}
