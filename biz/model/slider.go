package model

type Link struct {
	Title       string `json:"title"`
	Url         string `json:"url"`
	Icon        string `json:"icon"`
	Description string `json:"description"`
}

type Slider struct {
	Author      string  `json:"author"`
	Avatar      string  `json:"avatar"`
	Description string  `json:"description"`
	Email       string  `json:"email"`
	Github      string  `json:"github"`
	OutLink     []*Link `json:"out_link"`
}
