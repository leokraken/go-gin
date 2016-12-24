package controllers

import (
	"github.com/gin-gonic/gin"
	"../models"
	"../services"
	"github.com/jinzhu/gorm"
	"fmt"
)

type ItemController struct {
	ItemService  services.ItemService
}


func  NewItemController(db *gorm.DB) *ItemController{
	return &ItemController{services.ItemService{db}}
}


func (ctrl *ItemController) Get(c *gin.Context) {
	var items * []models.Item = ctrl.ItemService.Get()
	if items == nil{
		c.JSON(404, "User not found ")
	}else {
		c.JSON(200, items)
	}
}


func (ctrl *ItemController) Create(c *gin.Context) {
	var model models.Item
	c.BindJSON(&model)
	fmt.Println(model.Description)
	fmt.Println(model.Title)

	//json.Unmarshal(model.Object, str)
	fmt.Println(string(model.Object))


	fmt.Println("printing")
	var user *models.Item = ctrl.ItemService.Create(model)
	if user!=nil{
		c.JSON(201, gin.H{"status": ""})
	}

}