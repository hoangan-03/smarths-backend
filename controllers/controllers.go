package controllers

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"time"

	"backend/database"
	"backend/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	_ "github.com/lib/pq"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB = database.DBSet()
var Validate = validator.New()

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		var account models.Account
		if err := c.BindJSON(&account); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		row := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM account WHERE username = $1", account.Username)
		var count int
		err := row.Scan(&count)
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Username already exists"})
			return
		}

		password := HashPassword(account.Password)
		account.Password = password

		_, err = db.ExecContext(ctx, "INSERT INTO account (username, password) VALUES ($1, $2)",
			account.Username, account.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong. Account not created"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Successfully signed up"})
		return
	}
}
func SignIn() gin.HandlerFunc {
	return func(c *gin.Context) {
		var account models.Account
		if err := c.BindJSON(&account); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		row := db.QueryRowContext(ctx, "SELECT password FROM account WHERE username = $1", account.Username)
		var storedPassword string
		err := row.Scan(&storedPassword)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Username or password is incorrect"})
			return
		}

		passwordIsValid := CheckPasswordHash(account.Password, storedPassword)
		if !passwordIsValid {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Username or password is incorrect"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Successfully signed in"})
		return
	}
}

func GetAllProducts() gin.HandlerFunc {
	return func(c *gin.Context) {
		var response models.Response
		var productList []models.Product

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		cursor, err := ProductCollection.Find(ctx, bson.D{{}})
		if err != nil {
			response.Status = "Failed"
			response.Code = http.StatusInternalServerError
			response.Msg = "Something went wrong. Please try again later"
			c.IndentedJSON(http.StatusInternalServerError, response)
			return
		}
		err = cursor.All(ctx, &productList)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		defer cursor.Close(ctx)
		if err := cursor.Err(); err != nil {
			log.Println(err)
			response.Status = "Failed"
			response.Code = 400
			response.Msg = "Invalid"
			c.IndentedJSON(400, response)
			return
		}

		response.Status = "OK"
		response.Code = 200
		response.Msg = "Successfully"
		response.Data = productList
		c.IndentedJSON(200, response)
		return
	}
}

func SearchProductByQuery() gin.HandlerFunc {
	return func(c *gin.Context) {
		var response models.Response
		var searchedProducts []models.Product
		queryParam := c.Query("name")
		if queryParam == "" {
			log.Println("query is empty")
			c.Header("Content-Type", "application/json")

			response.Status = "Failed"
			response.Code = http.StatusNotFound
			response.Msg = "Invalid searched index"
			c.IndentedJSON(http.StatusNotFound, response)
			c.Abort()
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		searchQueryDB, err := ProductCollection.Find(ctx, bson.M{"product_name": bson.M{"$regex": queryParam}})
		if err != nil {
			response.Status = "Failed"
			response.Code = 404
			response.Msg = "Something went wrong in fetching the db queries"
			c.IndentedJSON(404, response)
			return
		}
		err = searchQueryDB.All(ctx, &searchedProducts)
		if err != nil {
			log.Println(err)

			response.Status = "Failed"
			response.Code = 400
			response.Msg = "Invalid"
			c.IndentedJSON(400, response)
			return
		}
		defer searchQueryDB.Close(ctx)
		if err := searchQueryDB.Err(); err != nil {
			log.Println(err)

			response.Status = "Failed"
			response.Code = 400
			response.Msg = "Invalid request"
			c.IndentedJSON(400, response)
			return
		}

		response.Status = "OK"
		response.Code = 200
		response.Msg = "Successfully"
		response.Data = searchedProducts
		c.IndentedJSON(200, response)
		return
	}
}

func GetAllOrders() gin.HandlerFunc {
	return func(c *gin.Context) {
		var response models.Response
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var orderList []models.Order
		cursor, err := OrderCollection.Find(ctx, bson.D{{}})
		if err != nil {
			response.Status = "Failed"
			response.Code = http.StatusInternalServerError
			response.Msg = "Something went wrong. Please try again later"
			c.IndentedJSON(http.StatusInternalServerError, response)
			return
		}
		err = cursor.All(ctx, &orderList)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		defer cursor.Close(ctx)
		if err := cursor.Err(); err != nil {
			log.Println(err)
			response.Status = "Failed"
			response.Code = 400
			response.Msg = "Invalid"
			c.IndentedJSON(400, response)
			return
		}

		response.Status = "OK"
		response.Code = 200
		response.Msg = "Successfully"
		response.Data = orderList
		c.IndentedJSON(200, response)
		return
	}
}

func ProductAdderAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		var response models.Response
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var product models.Product
		if err := c.BindJSON(&product); err != nil {
			response.Status = "Failed"
			response.Code = http.StatusBadRequest
			response.Msg = err.Error()
			c.IndentedJSON(http.StatusBadRequest, response)
			return
		}
		product.ProductId = primitive.NewObjectID()
		product.Comments = make([]models.Comment, 0)
		_, anyErr := ProductCollection.InsertOne(ctx, product)
		if anyErr != nil {
			response.Status = "Failed"
			response.Code = http.StatusInternalServerError
			response.Msg = "Not created"
			c.IndentedJSON(http.StatusInternalServerError, response)
			return
		}

		response.Status = "OK"
		response.Code = http.StatusOK
		response.Msg = "New product has been successfully added by an admin"
		c.IndentedJSON(http.StatusOK, response)
		return
	}
}

func ProductUpdaterAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		var response models.Response
		productQueryId := c.Query("productId")
		if productQueryId == "" {
			response.Status = "Failed"
			response.Code = http.StatusBadRequest
			response.Msg = "Missing product id"
			c.IndentedJSON(http.StatusBadRequest, response)
			return
		}

		productId, err := primitive.ObjectIDFromHex(productQueryId)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		var name string = c.PostForm("name")
		if name == "" {
			response.Status = "Failed"
			response.Code = http.StatusBadRequest
			response.Msg = "Missing name"
			c.IndentedJSON(http.StatusBadRequest, response)
			return
		}
		var priceString string = c.PostForm("price")
		if priceString == "" {
			response.Status = "Failed"
			response.Code = http.StatusBadRequest
			response.Msg = "Missing price"
			c.IndentedJSON(http.StatusBadRequest, response)
			return
		}

		price, err := strconv.Atoi(priceString)
		if err != nil {
			response.Status = "Failed"
			response.Code = http.StatusBadRequest
			response.Msg = "Missing price"
			c.IndentedJSON(http.StatusBadRequest, response)
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		filter := bson.D{primitive.E{Key: "_id", Value: productId}}
		update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "product_name", Value: name}, {Key: "price", Value: price}}}}
		_, err = ProductCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			response.Status = "Failed"
			response.Code = 500
			response.Msg = "Something went wrong"
			c.IndentedJSON(500, response)
			return
		}

		ctx.Done()

		response.Status = "OK"
		response.Code = 200
		response.Msg = "Successfully updated the product"
		c.IndentedJSON(200, response)
		return
	}
}
