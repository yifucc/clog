package model

type Categories []*Category

func (c Categories) Len() int {
	return len(c)
}

func (c Categories) Less(i, j int) bool {
	return c[i].Number > c[j].Number
}

func (c Categories) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

type Category struct {
	Name   string `json:"name"`
	Path   string `json:"path"`
	Number int    `json:"number"`
}
