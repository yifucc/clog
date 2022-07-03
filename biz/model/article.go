package model

import "time"

type Articles []*Article

func (a Articles) Len() int {
	return len(a)
}

func (a Articles) Less(i, j int) bool {
	return a[i].UpdatedTime.After(a[j].UpdatedTime) || a[i].CreatedTime.After(a[j].CreatedTime)
}

func (a Articles) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

type Article struct {
	Title       string    `json:"title"`
	Category    string    `json:"category"`
	Description string    `json:"description"`
	Body        string    `json:"body"`
	Path        string    `json:"path"`
	Type        int       `json:"type"`
	CreatedTime time.Time `json:"create_time"`
	UpdatedTime time.Time `json:"update_time"`
}
