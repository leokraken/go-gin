package services

import (
	"../models"
	"github.com/jinzhu/gorm"
)

type ItemService struct{
	DB *gorm.DB
}


func (srv * ItemService) Get() * []models.Item{
	var items  []models.Item
	srv.DB.Find(&items)
	return &items
}

func (srv * ItemService) Create(model models.Item) * models.Item {
	srv.DB.Create(model)
	return &model
}


