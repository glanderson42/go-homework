package main

import (
	"log"
	"net/http"
	"time"

	"github.com/appleboy/gin-jwt"
	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

// User is represents a user data structure
type User struct {
	ID       int    `gorm:"AUTO_INCREMENT" form:"id" json:"id"`
	Username string `gorm:"not null" form:"username" json:"username"`
	Password string `gorm:"not null" form:"password" json:"password"`
}

const (
	API_KEY = "blogengine"
)

// InitDb is for database init
func InitDb() *gorm.DB {
	db, err := gorm.Open("sqlite3", "data.db")
	db.LogMode(true)

	if err != nil {
		panic(err)
	}

	if !db.HasTable(&User{}) {
		db.CreateTable(&User{})
		db.Set("gorm:table_options", "ENGINE=InnoDB").CreateTable(&User{})
	}

	return db
}

// Cors is for CORS
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		c.Next()
	}
}

func indexGetMethod(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{ "message": "home get" }`))
}

// PostUser is...
func PostUser(c *gin.Context) {
	db := InitDb()
	defer db.Close()

	var user User
	c.Bind(&user)

	if user.Username != "" && user.Password != "" {
		// INSERT INTO "users" (name) VALUES (user.Name);
		db.Create(&user)
		// Display error
		c.JSON(201, gin.H{"success": user})

	} else {
		// Display error
		c.JSON(422, gin.H{"error": "Fields are empty"})
	}

	// curl -i -X POST -H "Content-Type: application/json" -d "{ \"firstname\": \"Thea\", \"lastname\": \"Queen\" }" http://localhost:8080/api/v1/users
}

// GetUsers is...
func GetUsers(c *gin.Context) {
	// Connection to the database
	db := InitDb()
	// Close connection database
	defer db.Close()

	var users []User
	// SELECT * FROM users
	db.Find(&users)

	// Display JSON result
	c.JSON(200, users)

	// curl -i http://localhost:8080/api/v1/users
}

// AuthUser is for authing user
func AuthUser(c *gin.Context) {
	db := InitDb()
	defer db.Close()

	var user User
	c.BindJSON(&user)

	var dbUser User
	db.Find(&dbUser, user)
	if dbUser.Username != "" {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user": dbUser.Username,
			"exp":  time.Now().Add(time.Hour * time.Duration(1)).Unix(),
			"iat":  time.Now().Unix(),
		})

		tokenString, err := token.SignedString([]byte(API_KEY))
		if err != nil {
			c.JSON(200, gin.H{"token": "generation failed"})
		}

		c.JSON(200, gin.H{"token": tokenString})
	} else {
		c.JSON(200, gin.H{"login": "not success"})
	}
}

func AuthMiddleware(next gin.HandlerFunc) http.Handler {
	if len(API_KEY) == 0 {
		log.Fatal("HTTP SERVER UNABLE TO START")
	}

	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte(API_KEY), nil
		},
		SigningMethod: jwt.SigningMethodHS256,
	})

	return jwtMiddleware.Handler(next)
}

func main() {

	router := gin.Default()
	router.Use(Cors())
	v1 := router.Group("api/v1")
	{
		v1.POST("/users", PostUser)
		v1.GET("/users", GetUsers)
		v1.POST("/auth", AuthUser)
	}
	log.Fatal(router.Run(":8080"))
}
