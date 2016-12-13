package services

import (
	"../models"
	"github.com/jinzhu/gorm"
)

type UserService struct{
	DB *gorm.DB
}


func (srv * UserService) Get(userId string) * models.User{
	var user models.User
	proy := []string{"username"}
	if err := srv.DB.Where("username = ?", userId).Select(proy).First(&user).Error; err != nil {
		return nil
	}else{
		return &user
	}
}

func (srv * UserService) Create(model models.User) * models.User {
	srv.DB.Create(model)
	return &model
}


