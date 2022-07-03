package model

type Navbar []*NavbarItem

type NavbarItem struct {
	Title       string `json:"title"`
	Path        string `json:"path"`
	Folder      string `json:"folder"`
	Description string `json:"description"`
	IsOutLink   bool   `json:"isOutLink"`
}
