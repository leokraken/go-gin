package controllers

import (
	"github.com/gin-gonic/gin"
	"../models"
	"../services"
	"github.com/jinzhu/gorm"
)

type Controller struct {
  UserService  services.UserService
}

//var UserService = &services.UserService{}


func NewController(db *gorm.DB) *Controller{
	return &Controller{services.UserService{db}}
}


func (ctrl *Controller) Get(c *gin.Context) {
	var userId = c.Params.ByName("userid")
	var user *models.User = ctrl.UserService.Get(userId)
	if user == nil{
		c.JSON(404, "User not found ")
	}else {
		c.JSON(200, user)
	}

}


func (ctrl *Controller) Create(c *gin.Context) {
	var model models.User
	c.BindJSON(&model)
	var user *models.User = ctrl.UserService.Create(model)
	if user!=nil{
		c.JSON(201, gin.H{"status": ""})
	}

}