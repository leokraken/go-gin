package main

import (
	"github.com/gin-gonic/gin"
	"fmt"
	"net/http"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/dgrijalva/jwt-go"
	"time"
	"./models"
	"./controllers"
)

//json:- ignores field


func functest(c *gin.Context) {
	var l models.User
	c.BindJSON(&l)
	fmt.Printf("%s %s", l.Username, l.Password)
	if l.Password != "pass" {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "Fail"})
	}else {
		c.JSON(http.StatusOK, gin.H{
			"message": "auth",
		})
	}
}

type CustomClaims struct{
	UserId string
	jwt.StandardClaims
}

func createToken(){

	claims := &CustomClaims{ "yeah", jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Minute*5).Unix(),
		Issuer:    "test",
	}}
	mySigningKey := []byte("here")
	claim := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)


	tokenString, err := claim.SignedString(mySigningKey)
	fmt.Println(tokenString)

	decodedToken, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("here"), nil
	})

	if claims, ok := decodedToken.Claims.(*CustomClaims); ok && decodedToken.Valid {
		fmt.Printf("%v %v %v \n", claims.UserId, claims.StandardClaims.Issuer, claims.ExpiresAt)
	} else {
		fmt.Println(err)
	}

}


func main() {
	routes := gin.Default()

	db, err := gorm.Open("postgres", "host=localhost user=postgres dbname=postgres sslmode=disable password=password")
	if err == nil {
		fmt.Println("Ok connected")
	}

	db.AutoMigrate(&models.Item{}, &models.User{})

	controller := controllers.NewController(db)


	var iController controllers.IController

	iController= controller

	routes.GET("/ping", func(c *gin.Context) {
		createToken()
		c.String(200, "pong")
	})

	routes.GET("/users/:userid", iController.Get)
	routes.POST("/users", iController.Create)

	routes.Run()
}
