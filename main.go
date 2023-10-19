package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

type plant struct {
	ID      string  `json:"id"`
	Product string  `json:"product"`
	Amount  int     `json:"amount"`
	Price   float64 `json:"price"`
}

var plants = []plant{
	{ID: "1", Product: "Syngonium green lime", Amount: 10, Price: 29.99},
	{ID: "2", Product: "Alocasia dragon skin", Amount: 10, Price: 39.99},
	{ID: "3", Product: "Monstera Alba", Amount: 10, Price: 50.00},
}

func getPlants(c *gin.Context) {
	// Context.IndentedJSON to serialize the struct into JSON and add it to the response
	c.IndentedJSON(http.StatusOK, plants)
}

func getPlantsByID(c *gin.Context) {
	// Context.Param получает id параметра из пути URL-адреса
	id := c.Param("id")

	for _, a := range plants {
		if a.ID == id {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "plant not found"})
}

func postPlants(c *gin.Context) {
	var newPlant plant

	// использую Context.BindJSON для привязки тела запроса к newPlant
	if err := c.BindJSON(&newPlant); err != nil {
		return
	}

	// добавляю plant структуру, инициализированную из JSON, в plants срез
	plants = append(plants, newPlant)
	// добавляет 201 в ответ код состояния, а также JSON, представляющий добавленный вами альбом
	c.IndentedJSON(http.StatusCreated, newPlant)
}

func deletePlants(c *gin.Context) {
	// Context.Param получает id параметра из пути URL-адреса
	id := c.Param("id")

	for _, a := range plants {
		if a.ID == id {
			idInt, err := strconv.Atoi(id)
			if err != nil {
				log.Fatal(err)
			}
			if idInt != len(plants) {
				plants[idInt-1] = plants[len(plants)-1]
				plants[idInt-1].ID = strconv.Itoa(idInt)
			}
			plants[len(plants)-1] = plant{}
			plants = plants[:len(plants)-1]
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "plant not found"})
}

func putPlats(c *gin.Context) {
	id := c.Param("id")
	var newData plant
	if err := c.BindJSON(&newData); err != nil {
		log.Fatal(err)
	}

	for _, a := range plants {
		if a.ID == id {
			itr, _ := strconv.Atoi(id)
			if newData.Product != "" {
				plants[itr-1].Product = newData.Product
			}
			if newData.Amount != 0 {
				plants[itr-1].Amount = newData.Amount
			}
			if newData.Price != 0 {
				plants[itr-1].Price = newData.Price
			}
			c.IndentedJSON(http.StatusOK, newData)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, newData)
}

func main() {
	router := gin.Default()
	router.GET("/plants", getPlants)
	// В Gin двоеточие перед элементом пути означает, что элемент является параметром пути
	router.GET("/plants/:id", getPlantsByID)
	router.POST("/plants", postPlants)
	router.DELETE("/plants/:id", deletePlants)
	router.PUT("/plants/:id", putPlats)

	router.Run("localhost:8080")

}
