package controllers

import "github.com/gin-gonic/gin"

type IController interface {
	Get(*gin.Context)
	Create(*gin.Context)
}
