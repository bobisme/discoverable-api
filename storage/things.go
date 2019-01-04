package storage

import (
	"github.com/bobisme/discoverable-api/models"
)

type ThingPage struct {
	Items      []*models.Thing `json:"items"`
	Total      int             `json:"total"`
	PrevCursor string          `json:"prevCursor,omitempty"`
	NextCursor string          `json:"nextCursor,omitempty"`
}

type ThingStorage interface {
	Get(id string) *models.Thing
	List(cursor string, limit int) ThingPage
}
