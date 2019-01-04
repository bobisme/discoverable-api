package models

import "time"

type Widget struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Things    []*Thing  `json:"things"`
	CreatedAt time.Time `json:"createdAt"`
}

func NewWidget(name string) *Widget {
	return &Widget{
		ID:   newId(),
		Name: name,
	}
}

func (w *Widget) AddThings(things ...*Thing) {
	thingIds := make(map[string]struct{})
	for _, t := range w.Things {
		thingIds[t.Id] = struct{}{}
	}
	for _, thing := range things {
		if _, ok := thingIds[thing.Id]; ok {
			continue
		}
		w.Things = append(w.Things, thing)
		thingIds[thing.Id] = struct{}{}
	}
}
