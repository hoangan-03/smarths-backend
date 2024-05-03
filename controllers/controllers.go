package controllers

import (
	"backend/database"
	"backend/models"
	"context"
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"time"

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
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
			return
		}

		if account.Key != "randomkey" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid key", "details": "Key must be 'randomkey'"})
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

		_, err = db.ExecContext(ctx, "INSERT INTO account (username, password,acc_id, key) VALUES ($1, $2, $3, $4)",
			account.Username, account.Password, account.Acc_id, account.Key)
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
		err := row.Scan(&returnedAccount.Username, &returnedAccount.Password, &returnedAccount.Acc_id, &returnedAccount.Key)
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
func Controlling() gin.HandlerFunc {
	return func(c *gin.Context) {
		var controlling models.Controlling
		if err := c.ShouldBindJSON(&controlling); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		query := `INSERT INTO Controlling (Dev_id, Room_id, Action, Ctrl_mode, Timestamp, Isviewed) VALUES ($1, $2, $3, $4, $5, $6)`
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		_, err := db.ExecContext(ctx, query, controlling.Dev_id, controlling.Room_id, controlling.Action, controlling.Ctrl_mode, controlling.Timestamp, controlling.Isviewed)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": controlling})
	}
}
func GetNofications() gin.HandlerFunc {
	return func(c *gin.Context) {
		query := `SELECT Ctrl_id, Action, Ctrl_mode, Timestamp, Isviewed FROM Controlling`
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		rows, err := db.QueryContext(ctx, query)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var result []models.Controlling
		for rows.Next() {
			var controlling models.Controlling
			err = rows.Scan(&controlling.Ctrl_id, &controlling.Action, &controlling.Ctrl_mode, &controlling.Timestamp, &controlling.Isviewed)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			result = append(result, controlling)
		}

		if err = rows.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": result})
	}
}

func UpdateIsViewed() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctrl_id, err := strconv.ParseInt(c.Param("ctrl_id"), 10, 64) // convert ctrl_id to int64
		if err != nil {
			log.Println("Error parsing ctrl_id:", err) // log the error
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ctrl_id"})
			return
		}

		query := `UPDATE Controlling SET Isviewed = true WHERE Ctrl_id = $1`
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		_, err = db.ExecContext(ctx, query, ctrl_id)
		if err != nil {
			log.Println("Error executing query:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Row updated successfully"})
	}
}
func AddBooking() gin.HandlerFunc {
	return func(c *gin.Context) {
		var booking models.Booking
		if err := c.ShouldBindJSON(&booking); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		query := `INSERT INTO booking (book_id, room_id, start_time, notes, remind_time, end_time) VALUES ($1, $2, $3, $4, $5, $6)`
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		_, err := db.ExecContext(ctx, query, booking.Book_id, booking.Room_id, booking.Start_time, booking.Notes, booking.Remind_time, booking.End_time)
		if err != nil {
			log.Println("Error executing query:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Booking added successfully", "booking": booking})
	}
}

// DeleteBooking deletes a booking with the given book_id
// DeleteBooking deletes a booking with the given book_id
func DeleteBooking() gin.HandlerFunc {
	return func(c *gin.Context) {
		book_id := c.Param("book_id")

		query := `DELETE FROM booking WHERE book_id = $1`
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		_, err := db.ExecContext(ctx, query, book_id)
		if err != nil {
			log.Println("Error executing query:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Booking deleted successfully"})
	}
}

// ModifyBooking modifies a booking with the given book_id
func ModifyBooking() gin.HandlerFunc {
	return func(c *gin.Context) {
		var booking models.Booking
		if err := c.ShouldBindJSON(&booking); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		query := `UPDATE booking SET room_id = $1, start_time = $2, notes = $3, remind_time = $4, end_time = $5 WHERE book_id = $6`
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		_, err := db.ExecContext(ctx, query, booking.Room_id, booking.Start_time, booking.Notes, booking.Remind_time, booking.End_time, booking.Book_id)
		if err != nil {
			log.Println("Error executing query:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Booking modified successfully", "booking": booking})
	}
}
func GetBookings() gin.HandlerFunc {
	return func(c *gin.Context) {
		query := `SELECT * FROM booking`
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		rows, err := db.QueryContext(ctx, query)
		if err != nil {
			log.Println("Error executing query:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var bookings []models.Booking
		for rows.Next() {
			var booking models.Booking
			if err := rows.Scan(&booking.Book_id, &booking.Room_id, &booking.Start_time, &booking.Notes, &booking.Remind_time, &booking.End_time); err != nil {
				log.Println("Error scanning row:", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			bookings = append(bookings, booking)
		}

		if err := rows.Err(); err != nil {
			log.Println("Error with rows:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"bookings": bookings})
	}
}
