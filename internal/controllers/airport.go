package controllers

import (
	"net/http"

	"github.com/acme-sky/airline-api/internal/models"
	"github.com/acme-sky/airline-api/pkg/db"
	"github.com/gin-gonic/gin"
)

func AirportGet(ctx *gin.Context) {
	db := db.GetDb()
	var airports []models.Airport
	db.Find(&airports)

	ctx.JSON(http.StatusOK, gin.H{
		"count": len(airports),
		"data":  &airports,
	})
}

func AirportPost(ctx *gin.Context) {
	db := db.GetDb()
	var input models.AirportInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	aiport := models.NewAirport(input)
	db.Create(&aiport)

	ctx.JSON(http.StatusCreated, gin.H{"data": aiport})
}
