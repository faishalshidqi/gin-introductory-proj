package models

import (
	"time"

	"github.com/google/uuid"
)

type Recipe struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Tags         []string  `json:"tags"`
	Ingredients  []string  `json:"ingredients"`
	Instructions []string  `json:"instructions"`
	PublishedAt  time.Time `json:"published_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
