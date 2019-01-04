package storage

import "github.com/bobisme/discoverable-api/models"

type WidgetPage struct {
	Items      []*models.Widget `json:"items"`
	Total      int              `json:"total"`
	PrevCursor string           `json:"prevCursor,omitempty"`
	NextCursor string           `json:"nextCursor,omitempty"`
}

type WidgetStorage interface {
	Save(*models.Widget) (*models.Widget, error)
	Get(id string) (*models.Widget, error)
	List(cursor string, limit int) (WidgetPage, error)
	Delete(id string) error
}
