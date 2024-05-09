package route

import (
	"TZ/internal/handler"
	"github.com/gin-gonic/gin"
)

func InitializeRoutes(router *gin.Engine) {
	handler := handler.CarHandler{} // Предположим, что CarHandler определен в вашем пакете internal
	router.GET("/cars", handler.GetCar)
	router.POST("/cars", handler.AddCar)
	router.PUT("/cars/:id", handler.UpdateCar)
	router.DELETE("/cars", handler.DeleteCar)

}
