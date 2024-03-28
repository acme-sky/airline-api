package handlers

import (
	"net/http"

	"github.com/acme-sky/airline-api/internal/models"
	"github.com/acme-sky/airline-api/pkg/db"
	"github.com/gin-gonic/gin"
)

func AirportHandlerGet(ctx *gin.Context) {
	db := db.GetDb()
	var airports []models.Airport
	db.Find(&airports)

	ctx.JSON(http.StatusOK, gin.H{
		"count": len(airports),
		"data":  &airports,
	})
}

func AirportHandlerPost(ctx *gin.Context) {
	db := db.GetDb()
	var input models.AirportInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	aiport := models.NewAirport(input)
	db.Create(&aiport)

	ctx.JSON(http.StatusCreated, aiport)
}

func AirportHandlerGetId(ctx *gin.Context) {
	db := db.GetDb()
	var airport models.Airport
	if err := db.Where("id = ?", ctx.Param("id")).First(&airport).Error; err != nil {
		ctx.JSON(http.StatusNotFound, map[string]string{})
		return
	}

	ctx.JSON(http.StatusOK, airport)
}

func AirportHandlerPut(ctx *gin.Context) {
	db := db.GetDb()
	var airport models.Airport
	if err := db.Where("id = ?", ctx.Param("id")).First(&airport).Error; err != nil {
		ctx.JSON(http.StatusNotFound, map[string]string{})
		return
	}

	var input models.AirportInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.Model(&airport).Updates(input)

	ctx.JSON(http.StatusOK, airport)
}
