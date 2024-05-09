package handler

import (
	"TZ/internal/service"
	"database/sql"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type CarHandler struct {
	DB     *sql.DB
	Logger *log.Logger
}

func (h *CarHandler) GetCar(c *gin.Context) {
	carData, err := service.GetCar(h.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching internal data"})
		return
	}
	c.JSON(200, gin.H{
		"info": carData,
	})
}

func (h *CarHandler) AddCar(c *gin.Context) {
	carData, err := service.AddNewCars(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error add new car to database"})
		return
	}
	c.JSON(200, gin.H{
		"Info": carData,
	})
}

func (h *CarHandler) UpdateCar(c *gin.Context) {
	carData, ownerInfo, err := service.PutCar(c, h.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error update car"})
		return
	}

	c.JSON(200, gin.H{
		"Old_information": carData,
		"New_information": ownerInfo,
	})
}

func (h *CarHandler) DeleteCar(c *gin.Context) {
	carData, carOwner, err := service.DeleteCar(c, h.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error delete internal data"})
		return
	}

	c.JSON(200, gin.H{
		"deleted_car":          carData,
		"Owner_of_deleted_car": carOwner,
	})
}
