package models

type Item struct {
	ID int `gorm:"AUTO_INCREMENT" json:"-"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Price       int `json:"-"`
}
