package router

import (
	"careercraftsman_chatbot/app"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

var db *sql.DB

func Router() {
	router := gin.Default()

	//define route handlers

	router.GET("/signup", func(c *gin.Context) {
		// Serve the signup page HTML
		filePath := "./static/index.html"
		fmt.Println("Serving signup page from:", filePath)
		c.File(filePath)
	})

	//router.GET("/signup", app.GET_Signup)

	router.POST("/signup", func(c *gin.Context) {
		var user app.User
		// Bind JSON data to User struct
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// Insert user details into database
		_, err := db.Exec("INSERT INTO users (first_name, last_name, username, password) VALUES ($1, $2, $3, $4)",
			user.FirstName, user.LastName, user.Username, user.Password)
		if err != nil {

			fmt.Println("Error inserting user into database:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Redirect(http.StatusSeeOther, "/signup_success")
	})

	// Serve the sign-up success page
	router.GET("/signup_success", func(c *gin.Context) {
		c.File("./static/signup_success.html")
	})

	// Serve static files
	router.Static("/static", "./static")

	// Route to serve the sign-in page
	router.GET("/signin", func(c *gin.Context) {
		// Render the sign-in page
		filePath_signin := "./static/signin.html"
		fmt.Println("Serving signup page from:", filePath_signin)
		c.File(filePath_signin)
	})

	// Route to handle sign-in POST request
	router.POST("/signin", func(c *gin.Context) {

	})

	// Start server
	router.Run(":8081")

}
