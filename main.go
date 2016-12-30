//
package main

import (
	"github.com/gin-gonic/gin"
	"fmt"
	"net/http"
	// database drivers
	_ "github.com/lib/pq"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/types"

	//"github.com/jinzhu/gorm"
	//_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/dgrijalva/jwt-go"
	"time"
	"./models"
	//"./controllers"
)

//json:- ignores field


const schema = `
CREATE TABLE users (
    name varchar,
    age integer,
    data json
);`


const insert_query = `INSERT INTO events VALUES (:data)`


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

type UserShort struct{
	Name string `json:"name"`
	Age int `json:"age"`
	Data types.JSONText `json:"data"`
}

func main() {
	routes := gin.Default()

	db, err := sqlx.Connect("postgres", "host=localhost user=postgres dbname=postgres password=postgres sslmode=disable")
	if err != nil {
		fmt.Println("error", err)
	}else{
		fmt.Println("success",db)
	}
/*
	db, err := gorm.Open("postgres", "host=localhost user=postgres dbname=postgres sslmode=disable password=password")
	if err == nil {
		fmt.Println("Ok connected")
	}

	db.AutoMigrate(&models.Item{}, &models.User{})

	controller := controllers.NewController(db)


	var iController controllers.IController

	iController= controller


	routes.GET("/users/:userid", iController.Get)
	routes.POST("/users", iController.Create)
*/

	// Prepared named, util para usar objetos en query
	//all_users, _ := db.PrepareNamed(`SELECT * FROM users where name=:name`)

	// En cambio prepared se le pasa los argumentos sin estructura
	all_users, _ := db.Preparex(`SELECT * FROM users where name=$1`)

	routes.GET("/users", func(c *gin.Context) {
		var users [] UserShort
		//err := db.Select(&users, "SELECT * FROM users")
		err := all_users.Select(&users, "leo") //map[string]interface{}{"name": "leo"})

		if err != nil{
			fmt.Println("error with rows", err)
		}
		c.JSON(200, users)
	})

	routes.POST("/users", func(c *gin.Context) {
		var user UserShort
		c.BindJSON(&user)
		_, err = db.NamedExec(`INSERT INTO users(name,age,data) VALUES (:name,:age,:data)`, &user) 
		
		if err != nil{
			fmt.Println("error insert", err)
		}
		c.JSON(200, user)
	})


	routes.POST("/data", func(c *gin.Context) {
		var user types.JSONText
		c.BindJSON(&user)
		_, err = db.NamedExec(insert_query, map[string]interface{}{"data": user}) 
		c.JSON(200, user)
	})

/*	
	type Event struct{
		Data types.JSONText 
	}*/

	routes.GET("/data", func(c *gin.Context) {
		var events [] types.JSONText
		_ = db.Select(&events, "select data::JSON #>'{data,ok}' from events") 
		c.JSON(200, events)
	})


	routes.GET("/ping", func(c *gin.Context) {
		createToken()
		c.String(200, "pong")
	})


	routes.Run()
}
