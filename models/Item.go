package models

import "../utils"

type Object struct {
	Amazing bool
}

type Item struct {
	ID int `gorm:"AUTO_INCREMENT" json:"-"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Price       int `json:"-"`
	Object  utils.JSONRaw `sql:"type:JSONB NOT NULL DEFAULT '{}'::JSONB" json:"object"`
	Userid       string `json:"userid"`

}
