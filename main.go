package main

import (
	"github.com/gin-gonic/gin"
	"fmt"
	"net/http"
	"encoding/json"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/dgrijalva/jwt-go"
	"time"
	"./models"
	"./controllers"
)

func hasRole(userRoles []string, roles []string) bool{
	for _, u := range userRoles{
		for _, r :=range roles{
			if r == u{
				return true
			}
		}
	}
	return false
}

func authMiddleware(roles []string) gin.HandlerFunc {


	return func(c *gin.Context) {
		var token = c.Request.Header.Get("Authorization")
		fmt.Println("Header ", token)

		decodedToken, _ := jwt.ParseWithClaims(token, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte("here"), nil
		})

		if claims, ok := decodedToken.Claims.(*CustomClaims); ok && decodedToken.Valid {
			if hasRole(claims.Roles, []string{"admin"}){
				fmt.Println("OK HAS ADMIN ROLE")
			}
			//fmt.Printf("%v %v %v \n", claims.UserId, claims.StandardClaims.Issuer, claims.ExpiresAt)
		} else {
			//fmt.Println(err)
			c.JSON(http.StatusUnauthorized, "unauthorized")
			c.Abort()

		}
	}
}

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

type CustomClaims struct {
	UserId string
	Roles [] string
	jwt.StandardClaims
}

func createToken() string {

	claims := &CustomClaims{"yeah", []string{ "user"},  jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Minute * 5).Unix(),
		Issuer:    "test",
	}}
	mySigningKey := []byte("here")
	claim := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, _ := claim.SignedString(mySigningKey)
	fmt.Println(tokenString)
	return tokenString

}

func testMaps(){
	hashmap := make(map[string] string)
	hashmap["hola"]= "hola"
	hashmap["esto"]="esta bueno"
	jsonString, _:= json.Marshal(hashmap)
	fmt.Println(string(jsonString))
}


func main() {

	testMaps()
	routes := gin.Default()

	db, err := gorm.Open("postgres", "host=localhost user=postgres dbname=postgres sslmode=disable password=password")
	if err == nil {
		fmt.Println("Ok connected")
	}

	db.AutoMigrate(&models.Item{}, &models.User{})

	controller := controllers.NewController(db)
	itemController := controllers.NewItemController(db)

	var iController controllers.IController

	iController = controller

	//routes.Use(authMiddleware([]string{"admin"}))

	routes.GET("/token", func(c * gin.Context) {
		var token string = createToken()
		c.JSON(200, gin.H{"token": token})
	})

	routes.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	routes.GET("/users/:userid", iController.Get)
	routes.POST("/users", iController.Create)

	routes.GET("/items", itemController.Get)
	routes.POST("/items", itemController.Create)


	routes.Run()
}
