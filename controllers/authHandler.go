package controllers

import (
	"fmt"
	"net/http"
	"shines/middlewares"
	"shines/models"
	"shines/repositories"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func ViewRegisterHandler(c *gin.Context) {
	context := gin.H{
		"title": "Sign Up",
	}
	c.HTML(
		http.StatusOK,
		"register.html",
		context,
	)
}

func RegisterHandler(c *gin.Context) {
	var user models.User

	var username, email, phone, password string
	var usernameErr, phoneErr, passwordErr, emailErr string

	username = c.PostForm("name")
	email = c.PostForm("email")
	phone = c.PostForm("phone")
	password = c.PostForm("password")

	if len(username) < 5 {
		usernameErr = "Minimum Username is 5 Characters!"
	}

	if !strings.Contains(email, "@") {
		emailErr = "Email must included @"
	}

	if len(email) < 10 {
		emailErr = "Email must be at least 10 Characters and included @"
	}

	if len(phone) < 8 {
		phoneErr = "Minimum phone is 8 Characters!"
	}

	if !repositories.IsNumber(phone) {
		phoneErr = "Phone must be a number"
	}

	if len(password) < 5 {
		passwordErr = "Minimum Password is 5 Characters!"
	}

	if usernameErr == "" && phoneErr == "" && passwordErr == "" && emailErr == "" {
		hashedPassword, err := bcrypt.GenerateFromPassword(
			[]byte(password),
			bcrypt.DefaultCost,
		)
		if err != nil {
			panic(err)
		}

		user = models.User{
			Username: username,
			Email:    email,
			Phone:    phone,
			Password: string(hashedPassword),
			Role:     "Customer",
		}
		err = models.DB.Create(&user).Error
		if err != nil {

			ErrorHandler2("Error Create", "Failed to Create Data", "/shines/main/register-page", c)
			return
		}
		c.Redirect(
			http.StatusFound,
			"/shines/main/login-page",
		)
		return
	}

	context := gin.H{
		"title":       "Sign Up",
		"email":       email,
		"phone":       phone,
		"username":    username,
		"usernameErr": usernameErr,
		"emailErr":    emailErr,
		"phoneErr":    phoneErr,
		"passwordErr": passwordErr,
	}

	c.HTML(
		http.StatusOK,
		"register.html",
		context,
	)
}

func ViewLoginHandler(c *gin.Context) {
	isLogged := middlewares.CheckSession(c)
	if isLogged {
		c.Redirect(
			http.StatusFound,
			"/shines/main/home-page",
		)
		return
	}
	context := gin.H{
		"title": "Login",
	}
	c.HTML(
		http.StatusOK,
		"login.html",
		context,
	)
}

func LoginHandler(c *gin.Context) {
	var user models.User
	var username, password string
	var usernameErr, passwordErr string

	username = c.PostForm("username")
	password = c.PostForm("password")

	if len(username) < 5 {
		usernameErr = "Minimum Username is 5 Characters!"
	}

	if len(password) < 5 {
		passwordErr = "Minimum Password is 5 Characters!"
	}

	err := models.DB.Where("Username = ?", username).First(&user).Error
	if err != nil {
		usernameErr = "Invalid Username"
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(password),
	)
	if err != nil {
		passwordErr = "Invalid Username or Password"
	}

	if usernameErr == "" && passwordErr == "" {
		middlewares.SaveSession(c, username)
		err := repositories.CreateProfile(c)
		if err != nil {

			ErrorHandler1("Failed to Create Data", "/shines/main/personal-information-page", c)
			return
		}

		err = repositories.CreateShop(c)
		if err != nil {

			ErrorHandler1("Failed to Create Data", "/shines/main/personal-information-page", c)
			return
		}

		c.Redirect(
			http.StatusFound,
			"/shines/main/home-page",
		)
		return
	}
	context := gin.H{
		"title":       "Login",
		"username":    username,
		"usernameErr": usernameErr,
		"passwordErr": passwordErr,
	}
	c.HTML(
		http.StatusOK,
		"login.html",
		context,
	)
}

func LogoutHandler(c *gin.Context) {
	middlewares.ClearSession(c)
	c.Redirect(
		http.StatusFound,
		"/shines/main/login-page",
	)
}

func RootHandler(c *gin.Context) {
	isLogged := middlewares.CheckSession(c)
	if isLogged {
		c.Redirect(
			http.StatusFound,
			"/shines/main/home-page",
		)
	} else {
		c.Redirect(
			http.StatusFound,
			"/shines/main/login-page",
		)
	}
}

func ViewHomeHandler(c *gin.Context) {
	isLogged := middlewares.CheckSession(c)
	if !isLogged {
		c.Redirect(
			http.StatusFound,
			"shines/main/login-page",
		)
		return
	}
	userID := repositories.GetuserId(c)
	products := []models.Product{}
	err := models.DB.Model(&models.Product{}).
		Select("*").
		Where("Shop_id != ? AND product_stock != ?", userID, 0).
		Find(&products).Error
	if err != nil {

		ErrorHandler1("Failed to Get Data", "/shines/main/home-page", c)
		return
	}
	context := gin.H{
		"title":    "Home",
		"products": products,
		"isSeller": repositories.IsSeller(c),
		"isAdmin":  repositories.IsAdmin(c),
	}
	fmt.Println(products)
	c.HTML(
		http.StatusOK,
		"home.html",
		context,
	)
}
