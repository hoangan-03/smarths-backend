package controllers

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"time"

	"backend/database"
	"backend/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	_ "github.com/lib/pq"
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
	}
}

// temperature //light //humidity //humandetect
func GetRecordByCategory(category string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var recordList []models.Record
		var deviceID string
		err := db.QueryRow("SELECT dev_id FROM device WHERE category = $1", category).Scan(&deviceID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong. Please try again later"})
			return
		}
		rows, err := db.Query("SELECT * FROM record WHERE dev_id = $1", deviceID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong. Please try again later"})
			return
		}
		defer rows.Close()
		for rows.Next() {
			var record models.Record
			err = rows.Scan(&record.Rec_id, &record.Dev_id, &record.Value, &record.Status, &record.Timestamp)
			if err != nil {
				log.Println(err)
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
			recordList = append(recordList, record)
		}
		if err = rows.Err(); err != nil {
			log.Println(err)
			c.JSON(400, gin.H{"error": "Invalid"})
			return
		}
		c.JSON(200, recordList)
	}
}
