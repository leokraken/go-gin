package models

type Place struct {
	ID int `json:"id" gorm:"AUTO_INCREMENT"`
	Name string `json:"name"`

}
