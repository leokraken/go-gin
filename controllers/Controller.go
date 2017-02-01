package controllers

import (
	"github.com/gin-gonic/gin"
	"../models"
	"../services"
	"github.com/jinzhu/gorm"
	"strconv"
)

type ServiceController struct {
  	UserService  services.UserService
}


func NewController(db *gorm.DB) *ServiceController{
	return &ServiceController{services.UserService{db}}
}


func (ctrl *ServiceController) Get(c *gin.Context) {
	var userId, _ = strconv.Atoi(c.Params.ByName("userid"))
	var user *models.User = ctrl.UserService.Get(userId)
	c.JSON(200, user)
}


func (ctrl *ServiceController) Create(c *gin.Context) {
	var model models.User
	c.BindJSON(&model)
	var user *models.User = ctrl.UserService.Create(model)
	if user!=nil{
		c.JSON(201, gin.H{"status": ""})
	}

}