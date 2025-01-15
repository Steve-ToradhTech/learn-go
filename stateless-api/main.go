package main

import (
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

type User struct {
	gorm.Model
	Name  string
	Email string `gorm:"uniqueIndex"`
}

type CreateUserRequest struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
}

type UserResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type CreateUserForm struct {
	Name  string `form:"name" binding:"required"`
	Email string `form:"email" binding:"required,email"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.ClaimStrings
}

func getUsers(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Fetching all users",
	})
}

func getProducts(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Fetcing all prodcuts"})
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Header.Get("Authorization") != "Bearer token" {
			c.AbortWithStatus(401)
			return
		}
		c.Next()
	}
}

func adminDashboard(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Admin Dashboard",
	})
}

func adminSettings(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Admin Settings"})
}

func getUser(c *gin.Context) {
	// In this example, when a request is made to /user/123, the userID variable will be 123
	userId := c.Param("id")
	c.JSON(200, gin.H{"user_id": userId})
}

func search(c *gin.Context) {
	// In this example, when a request is made to /search?q=golang, the query variable will be golang.
	query := c.Query("q")
	c.JSON(200, gin.H{"query": query})
}

func login(c *gin.Context) {
	var creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if creds.Username == "admin" && creds.Password == "password123" {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256,
			jwt.MapClaims{
				"username": creds.Username,
				"exp":      time.Now().Add(time.Hour * 24).Unix(),
			})
		tokenString, err := token.SignedString(c.Keys)
		if err != nil {
			log.Default()
		}
		c.JSON(http.StatusOK, gin.H{"token": tokenString})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	}
}

func authenticateJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := jwt.Parse(tokenString)
	}
}

func main() {
	// Initialize Gin Routes
	router := gin.Default()

	// Initalize PostgreSQL connection
	dsn := "user=postgres password=postgres dbname=postgres sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}

	db.AutoMigrate(&User{})

	// Define routes
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello Gin!",
		})
	})

	router.GET("/users/:id", func(c *gin.Context) {
		id := c.Param("d")
		c.JSON(http.StatusOK, gin.H{
			"id": id,
		})
	})

	// Group routes
	api := router.Group("/api")
	{
		api.GET("/users", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"users": []string{"Alice", "Bob", "Charlie"},
			})
		})
	}

	api2 := router.Group("/api")
	api2.Use(AuthMiddleware())
	{
		api2.GET("/users", getUsers)
		api2.GET("/products", getProducts)
	}

	router.POST("/users", func(c *gin.Context) {
		var req CreateUserRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		c.JSON(201, git.H{"message": "User created sucessfully"})
	})

	router.GET("/users/:id", func(c *gin.Context) {
		var user User
		if err := db.First(&user, c.Param("id")).Error; err != nil {
			c.JSON(404, gin.H{"error": "user not found"})
			return
		}
		userResp := UserResponse{ID: user.ID, Name: user.Name, Email: user.Email}
		c.JSON(200, userResp)
	})

	router.POST("/users", func(c *gin.Context) {
		var form CreateUserForm
		if err := c.ShouldBind(&form); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		c.JSON(201, gin.H{"message": "User created successfully"})
	})

	// In this example, any request to /assets will serve files from the ./assets directory.
	// For instance, a request to /assets/style.css will serve the style.css file from the ./assets directory.
	router.Static("/assets", "./assets")

	// In this example, a request to /index will render the index.tmpl template with the title “Main website”.
	router.LoadHTMLGlob("templates/*")
	router.GET("/index", func(c *gin.Context) {
		c.HTML(200, "index.tmpl", gin.H{"title": "Main website"})
	})

	authorized := router.Group("/", gin.BasicAuth(gin.Accounts{
		"admin": "password123",
	}))

	// In this example, the gin.BasicAuth middleware checks for a valid username and password.
	// If the credentials are correct, the request proceeds; otherwise, it returns a 401 Unauthorized status.
	authorized.GET("/protected", func(c *gin.Context) {
		user := c.MustGet(gin.AuthUserKey).(string)
		c.JSON(http.StatusOK, gin.H{"user": user, "message": "welcome to the protected site"})
	})

	// Start the server
	router.Run(":8080")
}
