package routes

import (
	"backend/controllers"

	"github.com/gin-gonic/gin"
)

func Routes(router *gin.Engine) {
	router.GET("/temperature_records", controllers.GetRecordByCategory("temperature"))
	router.GET("/light_level_records", controllers.GetRecordByCategory("light"))
	router.GET("/humidity_records", controllers.GetRecordByCategory("humidity"))
	router.GET("/camera_records", controllers.GetRecordByCategory("humandetect"))
	router.POST("/register", controllers.Register())
	router.POST("/signin", controllers.SignIn())
	router.POST("/controlling", controllers.Controlling())
	router.GET("/get_nofications", controllers.GetNofications())
	router.PUT("/update_isviewed/:ctrl_id", controllers.UpdateIsViewed())
	router.POST("/add_booking", controllers.AddBooking())
	router.DELETE("/delete_booking/:book_id", controllers.DeleteBooking())
	router.PUT("/modify_booking/:book_id", controllers.ModifyBooking())
	router.GET("/get_bookings", controllers.GetBookings())

}
