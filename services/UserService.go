package services

import (
	"../models"
	"github.com/jinzhu/gorm"
)

type UserService struct{
	DB *gorm.DB
}


func (srv * UserService) Get(id int) * models.User{
	var user models.User
	proy := []string{"id", "name"}
	if err := srv.DB.Where("id = ?", id).Select(proy).First(&user).Error; err != nil {
		return nil
	}else{
		return &user
	}
}

func (srv * UserService) Create(model models.User) * models.User {
	srv.DB.Create(model)
	return &model
}


