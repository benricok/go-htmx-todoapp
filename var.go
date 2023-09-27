package main

import (
	"time"

	"github.com/google/uuid"
)

type CustomTime struct {
	time.Time
}

type Item struct {
	ID        uuid.UUID  `json:"id"`
	Owner     string     `json:"owner"`
	Name      string     `json:"name,omitempty"`
	Done      bool       `json:"done"`
	Created   CustomTime `json:"datecreated,omitempty"`
	Completed CustomTime `json:"datecompleted,omitempty"`
}

