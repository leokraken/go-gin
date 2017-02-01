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
	"runtime"
	_ "github.com/lib/pq"
	"github.com/jmoiron/sqlx"
	"./utils"
)
type RolesMap map[string]bool

func hasRole(userRoles RolesMap, roles RolesMap) bool {
	// Order n userRoles length
	for key, _ := range userRoles {
		_, ok := roles[key]
		if ok {
			return true
		}
	}
	return false
}

func authMiddleware(roles RolesMap) gin.HandlerFunc {
	return func(c *gin.Context) {
		var token = c.Request.Header.Get("Authorization")
		decodedToken, _ := jwt.ParseWithClaims(token, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte("here"), nil
		})


		claims, ok := decodedToken.Claims.(*CustomClaims);

		if !ok || !decodedToken.Valid || !hasRole(claims.Roles, roles) {
			c.String(http.StatusUnauthorized, "unauthorized")
			c.Abort()
		}
	}
}

type CustomClaims struct {
	UserId string
	Roles  RolesMap
	jwt.StandardClaims
}

func createToken() string {

	claims := &CustomClaims{"yeah", RolesMap{"user": true}, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Minute * 10).Unix(),
		Issuer:    "test",
	}}
	mySigningKey := []byte("here")
	claim := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, _ := claim.SignedString(mySigningKey)
	fmt.Println(tokenString)
	return tokenString

}


/*
func rawsql (c * gin.Context){
	id := 1
	u := models.User{}
	err := db2.QueryRow("select id, name from users where id = $1", id).Scan(&u.Id, &u.Name)
	if err!= nil{
		fmt.Println(err)
	}
	//users = append(users, u)

	c.JSON(200, u)
}*/


func pong2(c *gin.Context) {
	fmt.Println("Ejecuta esto")
	c.String(200, "pong handler")
}



func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())
	gin.SetMode(gin.ReleaseMode)

	routes := gin.New()

	db, err := gorm.Open("postgres", "host=localhost user=postgres dbname=postgres sslmode=disable password=password")
	db.DB().SetMaxOpenConns(20)
	if err == nil {
		fmt.Println("Ok connected")
	}

	db2, _ := sqlx.Connect("postgres", "user=postgres dbname=postgres sslmode=disable")

	db.AutoMigrate(&models.Item{}, &models.User{}, &models.Place{}, &models.Specialism{})

	controller := controllers.NewController(db)
	itemController := controllers.NewItemController(db)

	var iController controllers.IController

	iController = controller


	gin.SetMode(gin.ReleaseMode)
	db2.SetMaxOpenConns(20)



	//routes.GET("/users/:userid", controller.Get)
	routes.GET("/users/:userid", func(c *gin.Context) {
		id := 1
		u := models.User{}
		db.First(&u, "id=?", id)
		//err := db2.Get(&u, "select * from users where id = $1", id)//"select users.id, json_agg(items.*) as items from users join items on(items.userid=users.id) where users.id = $1 group by users.id", id)
		if err != nil {
			fmt.Println(err)
		}

		//items := [] models.Item{}
		//json.Unmarshal(u.Items, &items)
		//fmt.Println(items)
		c.JSON(200, u)
	})

	routes.POST("/users", iController.Create)
	routes.GET("/items", itemController.Get)
	routes.POST("/items", itemController.Create)

	routes.GET("/ping2", authMiddleware(RolesMap{"user":true}), pong2)

	routes.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"says":"pong"})
	})


	/**
	Specialties
	 */
	routes.GET("/specialties", func(c *gin.Context){
		var specialties = []models.Specialism{}
		db.Find(&specialties)
		c.JSON(200, specialties)
	})


	/**
		Supply token
	 */
	routes.GET("/token", func(c *gin.Context) {

		var token string = createToken()
		c.JSON(200, gin.H{"token": token})
	})


	routes.GET("/parallel", func(c *gin.Context) {
		var opts utils.RequestOptions
		opts.Host = "api.sarem.uy"
		opts.Path = "/api/especialidades"

		ch1 := make(chan utils.Response)
		ch2 := make(chan utils.Response)

		go func(){
			ch1 <- utils.Get(opts)
		}()

		go func(){
			ch2<- utils.Get(opts)
		}()
		res1 :=<- ch1
		res2 :=<- ch2

		results := make(map[string] utils.JSONResponse)
		if res1.Error == nil{
			results["r1"]= res1.Message
		}

		if res2.Error == nil{
			results["r2"]= res2.Message
		}
		fmt.Println(results)
		c.JSON(200, results)
	})



	routes.Run(":8080")
}
