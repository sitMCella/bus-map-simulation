package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"hub/start/database"
)

type bus struct {
	Id        string `json:"id"`
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}

type busPosition struct {
	BusId         string `json:"bus_id"`
	Latitude      string `json:"latitude"`
	Longitude     string `json:"longitude"`
	NextBusStopId string `json:"next_bus_stop_id"`
	IsBusStop     bool   `json:"is_bus_stop"`
}

type Handler struct {
	DC *database.DatabaseConnection
}

// curl -X GET http://localhost:9090/hub/health
func (h *Handler) GetHealthStatus(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, "ok")
}

// curl -X GET http://localhost:9090/hub/bus_stop
func (h *Handler) GetBusStopEntries(c *gin.Context) {
	err, busStopEntries := h.DC.GetBusStopEntries()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error while retrieving the bus stop entries", "detail": err})
		return
	}
	c.IndentedJSON(http.StatusOK, busStopEntries)
}

// curl -X GET http://localhost:9090/hub/bus
func (h *Handler) GetBusEntries(c *gin.Context) {
	err, busEntries := h.DC.GetBusEntries()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error while retrieving the bus entries", "detail": err})
		return
	}
	c.IndentedJSON(http.StatusOK, busEntries)
}

// curl -X GET http://localhost:9090/hub/bus/492/time_table
func (h *Handler) GetBusTimeTableEntries(c *gin.Context) {
	busId := c.Param("bus_id")
	err, busTimeTableEntries := h.DC.GetBusTimeTableEntries(busId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error while retrieving the bus time table entries", "detail": err})
		return
	}
	c.IndentedJSON(http.StatusOK, busTimeTableEntries)
}

// curl -X POST http://localhost:9090/hub/bus/register --header "Content-Type: application/json" --data '{"id": "1","latitude": "0.34","longitude":"1.1"}'
func (h *Handler) BusRegister(c *gin.Context) {
	var newBus bus
	if err := c.BindJSON(&newBus); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "wrong bus parameters"})
		return
	}
	err, exists := h.DC.BusExists(newBus.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error while retrieving bus"})
		return
	}
	if exists {
		c.JSON(http.StatusConflict, gin.H{"error": "bus already exists"})
		return
	}
	if err := h.DC.CreateBus(newBus.Id, newBus.Latitude, newBus.Longitude); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error while creating bus", "detail": err})
		return
	}
	c.IndentedJSON(http.StatusCreated, newBus)
}

// curl -X POST http://localhost:9090/hub/bus/position --header "Content-Type: application/json" --data '{"bus_id": "492","latitude": "0.34","longitude":"1.1", "next_bus_stop_id": "1", "is_stop": "true"}'
func (h *Handler) InsertBusPosition(c *gin.Context) {
	var newBusPosition busPosition
	if err := c.BindJSON(&newBusPosition); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "wrong bus position parameters"})
		return
	}
	err, exists := h.DC.BusExists(newBusPosition.BusId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error while retrieving bus"})
		return
	}
	if exists == false {
		c.JSON(http.StatusConflict, gin.H{"error": "bus does not exist"})
		return
	}
	err, busPosition := h.DC.CreateBusPosition(newBusPosition.BusId, newBusPosition.Latitude, newBusPosition.Longitude, newBusPosition.NextBusStopId, newBusPosition.IsBusStop)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error while creating bus position", "detail": err})
		return
	}
	c.IndentedJSON(http.StatusCreated, busPosition)
}

func main() {
	var dc database.DatabaseConnection
	dc, err := database.NewDatabaseConnection()
	if err != nil {
		fmt.Println("Error while connecting to the database ", err)
		panic(err)
	}
	err = dc.InitDatabase()
	if err != nil {
		fmt.Println("Error while initializing the database ", err)
		panic(err)
	}

	h := &Handler{DC: &dc}

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "accepted"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.GET("/hub/health", h.GetHealthStatus)
	router.GET("/hub/bus_stop", h.GetBusStopEntries)
	router.GET("/hub/bus", h.GetBusEntries)
	router.GET("/hub/bus/:bus_id/time_table", h.GetBusTimeTableEntries)
	router.POST("/hub/bus/register", h.BusRegister)
	router.POST("/hub/bus/position", h.InsertBusPosition)

	srv := &http.Server{
		Addr:    ":9090",
		Handler: router.Handler(),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no params) by default sends syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fmt.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		fmt.Println("Server Shutdown:", err)
	}
	fmt.Println("Server Shutdown")
	if err := dc.Close(); err != nil {
		fmt.Println("Database Shutdown:", err)
	}
	fmt.Println("Database Shutdown")
}
