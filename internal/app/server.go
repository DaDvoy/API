package server

import (
	postgres "API/internal/DB"
	"API/internal/config"
	"fmt"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/exp/slog"
	"log"
	"net/http"
)

type Server struct {
	db   *postgres.DB
	logs *slog.Logger
	cfg  *config.Config
}

func New(logs *slog.Logger, conf *config.Config) *Server {
	return &Server{logs: logs, cfg: conf}
}

func (s *Server) RequestIdMiddleware(c *gin.Context) {
	uuidStr := uuid.NewV4().String()
	c.Writer.Header().Set("X-Request-ID", uuidStr)

	c.Next()
	defer func() {
		s.logs.Info(
			"", slog.Attr{Key: "UUID", Value: slog.StringValue(uuidStr)},
		)
	}()
}

func (s *Server) getPlants(c *gin.Context) {
	// Context.IndentedJSON to serialize the struct into JSON and add it to the response
	ch := make(chan *postgres.Plant)
	go postgres.ExecuteQuery(s.db, "SELECT * FROM plants", ch)
	for req := range ch {
		c.IndentedJSON(http.StatusOK, req)
	}
}

func (s *Server) getPlantsByID(c *gin.Context) {
	// Context.Param получает id параметра из пути URL-адреса
	id := c.Param("id")
	ch := make(chan *postgres.Plant)

	req := fmt.Sprintf("SELECT * FROM plants WHERE id=%s", id)
	go postgres.ExecuteQuery(s.db, req, ch)
	for res := range ch {
		c.IndentedJSON(http.StatusOK, res)
		return
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "plant not found"})
}

func (s *Server) postPlants(c *gin.Context) {
	var newPlant postgres.Plant

	// использую Context.BindJSON для привязки тела запроса к newPlant
	if err := c.BindJSON(&newPlant); err != nil {
		log.Fatal("Error: Unable to bindJSON:", err)
		return
	}
	req := fmt.Sprintf("INSERT INTO plants(id, product, amount, price)"+
		" values (%s , '%s', %s, %.2f)", newPlant.ID, newPlant.Product, newPlant.Amount, newPlant.Price)
	if err := postgres.ExecuteQuery(s.db, req, nil); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, newPlant)
	}
	// добавляет 201 в ответ код состояния, а также JSON, представляющий добавленный вами plant
	c.IndentedJSON(http.StatusCreated, newPlant)
}

func (s *Server) deletePlants(c *gin.Context) {
	// Context.Param получает id параметра из пути URL-адреса
	id := c.Param("id")

	req := fmt.Sprintf("DELETE FROM plants WHERE id=%s", id)
	if err := postgres.ExecuteQuery(s.db, req, nil); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "unable to delete statement"})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"message": "The plant deleted successfully"})
}

func (s *Server) putPlants(c *gin.Context) {
	id := c.Param("id")
	var newData postgres.Plant
	if err := c.BindJSON(&newData); err != nil {
		log.Fatal(err)
	}
	formReq := func(db *postgres.DB, specifiers ...string) {
		fmt.Println(specifiers[0])
		fmt.Println(specifiers[1])
		req := fmt.Sprintf("UPDATE plants SET %s WHERE id='%s'", specifiers[0], specifiers[1])
		if err := postgres.ExecuteQuery(db, req, nil); err != nil {
			c.IndentedJSON(http.StatusInternalServerError, newData)
			return
		}
	}

	if len(newData.Product) != 0 {
		req := fmt.Sprintf("product='%s'", newData.Product)
		formReq(s.db, req, id)
	}
	if len(newData.Amount) != 0 {
		req := fmt.Sprintf("amount='%s'", newData.Amount)
		formReq(s.db, req, id)
	}
	if newData.Price != 0 {
		req := fmt.Sprintf("price=%.2f", newData.Price)
		formReq(s.db, req, id)
	}
	c.IndentedJSON(http.StatusOK, newData)
}

func (s *Server) Start(_db *postgres.DB) error {
	router := gin.Default()
	router.Use(s.RequestIdMiddleware)
	s.db = _db
	router.GET("/plants", s.getPlants)
	// В Gin двоеточие перед элементом пути означает, что элемент является параметром пути
	router.GET("/plants/:id", s.getPlantsByID)
	router.POST("/plants", s.postPlants)
	router.DELETE("/plants/:id", s.deletePlants)
	router.PUT("/plants/:id", s.putPlants)

	srvListen := &http.Server{
		Addr:         s.cfg.HTTPServer.Address,
		Handler:      router,
		ReadTimeout:  s.cfg.HTTPServer.Timeout,
		WriteTimeout: s.cfg.HTTPServer.Timeout,
		IdleTimeout:  s.cfg.HTTPServer.IdleTimeout,
	}
	srvListen.ListenAndServe()
	return nil
}
