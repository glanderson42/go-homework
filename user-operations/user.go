package user

import (
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	ID       uint   `gorm:"AUTO_INCREMENT" form:"id" json:"id"`
	Username string `gorm:"not null" form:"username" json:"username"`
	Password string `gorm:"not null" form:"password" json:"password"`
}

// PostUser is...
func PostUser(c *gin.Context) {
	db := InitDb()
	defer db.Close()

	var user User
	c.Bind(&user)

	if user.Username != "" && user.Password != "" {
		db.Create(&user)
		c.JSON(201, gin.H{"success": user})

	} else {
		c.JSON(422, gin.H{"error": "Fields are empty"})
	}
}

// GetUsers is...
func GetUsers(c *gin.Context) {
	db := InitDb()
	defer db.Close()

	var users []User
	db.Find(&users)

	c.JSON(200, users)
}
