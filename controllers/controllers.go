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
)

var db *sql.DB = database.DBSet()
var Validate = validator.New()

func Register() gin.HandlerFunc {
	return func(c *gin.Context) {
		var account models.Account
		if err := c.BindJSON(&account); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if account.Key != "randomkey" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid key"})
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

		_, err = db.ExecContext(ctx, "INSERT INTO account (acc_id, username, password, key) VALUES ($1, $2, $3, $4)",
			account.Acc_id, account.Username, account.Password, account.Key)
		if err != nil {
			log.Printf("Error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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

		row := db.QueryRowContext(ctx, "SELECT * FROM account WHERE username = $1", account.Username)
		var returnedAccount models.Account
		err := row.Scan(&returnedAccount.Acc_id, &returnedAccount.Username, &returnedAccount.Password, &returnedAccount.Key)
		if err != nil {
			log.Printf("Error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		passwordIsValid := account.Password == returnedAccount.Password
		if !passwordIsValid {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Username or password is incorrect"})
			return
		}
		returnedAccount.Password = ""

		c.JSON(http.StatusOK, gin.H{"message": "Successfully signed in", "user": returnedAccount})
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
