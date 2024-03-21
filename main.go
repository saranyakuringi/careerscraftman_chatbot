package main

import (
	"bytes"
	"careercraftsman_chatbot/app"
	"careercraftsman_chatbot/controller"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/openai/openai-go/v2/openai"
)

var db *sql.DB
var dbMutex sync.Mutex
var gpt3Client *openai.Client

// Initialize GPT-3 API credentials
func initGPT3Client() {
	// Initialize the OpenAI GPT-3 client with your API key
	apiKey := "YOUR_API_KEY"
	gpt3Client = openai.NewClient(apiKey)
}

func initDB() {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	// Check if the database is already initialized
	if db != nil {
		return
	}

	var err error
	db, err = controller.Connecting_database()
	if err != nil {
		fmt.Println("Error in connecting to database", err)
		panic("Failed to connect to the database")
	}
}

func main() {

	fmt.Println("Initializing the bot....")

	// Initialize database connection
	initDB()

	//router.Router()

	router := gin.Default()

	//Signup routes
	router.GET("/signup", GET_Signup)
	router.POST("/signup", POST_Signup)
	router.GET("/signup_success", GET_signup_success)

	//signup routers
	router.GET("/signin", GET_signin)
	router.POST("/signin", POST_Signin)
	router.GET("/chatbot", GET_Chatbot)
	router.GET("/signin_fail", GET_signin_fail)

	// query routes
	router.POST("/query", handleQuery)

	// Serve static files
	router.Static("/static", "./static")

	// Start server
	router.Run("localhost:8081")
}

func GET_Signup(c *gin.Context) {
	filePath := "./static/index.html"
	fmt.Println("Serving signup page from:", filePath)
	c.File(filePath)
}

func POST_Signup(c *gin.Context) {
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
}

func GET_signup_success(c *gin.Context) {
	c.File("./static/signup_success.html")
}

func GET_signin(c *gin.Context) {
	filePath := "./static/signin.html"
	fmt.Println("Serving signin page from:", filePath)
	c.File(filePath)
}

func POST_Signin(c *gin.Context) {
	var user app.Validuser
	// Bind JSON data to User struct
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Printf("username:%v,password:%v\n", user.Username, user.Password)

	rows, err := db.Query("SELECT * FROM users WHERE username = $1 AND password = $2", user.Username, user.Password)
	if err != nil {
		log.Println("Error in users search query:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	defer rows.Close()

	var authenticated bool
	for rows.Next() {
		authenticated = true
		break // Found a matching user, no need to continue checking
	}
	fmt.Println("authenticated:", authenticated)

	if authenticated {
		c.Redirect(http.StatusSeeOther, "/chatbot")
		return
	}
	c.Redirect(http.StatusSeeOther, "/signin_fail")
}

func GET_Chatbot(c *gin.Context) {
	c.File("./static/chatbot.html")
}

func GET_signin_fail(c *gin.Context) {
	c.File("./static/signin_fail.html")
}

// Process user query using GPT-3 API
func processQuery(query string) (string, error) {
	gpt3Mutex.Lock()
	defer gpt3Mutex.Unlock()

	// Prepare GPT-3 request data
	requestData := GPT3Request{
		Prompt:      query,
		MaxTokens:   50,
		Temperature: 0.7,
		APIKey:      apiKey,
	}

	// Convert request data to JSON
	requestDataJSON, err := json.Marshal(requestData)
	if err != nil {
		return "", err
	}

	// Send request to GPT-3 API
	resp, err := http.Post(apiUrl, "application/json", bytes.NewBuffer(requestDataJSON))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Read and parse GPT-3 API response
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var responseData GPT3Response
	if err := json.Unmarshal(responseBody, &responseData); err != nil {
		return "", err
	}

	// Extract and return GPT-3 response text
	if len(responseData.Choices) > 0 {
		return responseData.Choices[0].Text, nil
	}
	return "", fmt.Errorf("no response from GPT-3")
}

// Define a function to insert a chat message into the database
func insertChatMessage(sender, message string, db *sql.DB) error {
	_, err := db.Exec("INSERT INTO chat_messages (sender, message) VALUES ($1, $2)", sender, message)
	return err
}

func handleQuery(c *gin.Context) {
	// Parse user query from request body
	var requestData map[string]string
	if err := c.BindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}
	userQuery := requestData["query"]

	// Generate bot response using GPT-3 API
	response, err := generateGPT3Response(userQuery)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate bot response"})
		return
	}

	// Store user query and bot response in database
	if err := storeChatMessage("user", userQuery); err != nil {
		log.Println("Error storing user query in database:", err)
	}
	if err := storeChatMessage("bot", response); err != nil {
		log.Println("Error storing bot response in database:", err)
	}

	// Send bot response to client
	c.JSON(http.StatusOK, gin.H{"response": response})
}

func generateGPT3Response(userQuery string) (string, error) {
	// Use the GPT-3 API to generate a response to the user query
	prompt := fmt.Sprintf("User: %s\nBot:", userQuery)
	completion, err := gpt3Client.Completions.Create(openai.ChatCompletion{
		Model: "text-davinci-002",
		Messages: []openai.ChatMessage{
			{Role: "system", Content: "/start"},
			{Role: "user", Content: userQuery},
		},
	})
	if err != nil {
		return "", err
	}

	// Extract and return the bot response from the completion
	return completion.Choices[0].Message.Content, nil
}

func storeChatMessage(sender, message string) error {
	// Insert the chat message into the database
	_, err := db.Exec("INSERT INTO chat_messages (sender, message, timestamp) VALUES ($1, $2, $3)",
		sender, message, time.Now())
	return err
}
