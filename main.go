package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()

	ch := make(chan []gameDeal)
	r.SetTrustedProxies([]string{"127.0. 0.1"})
	r.GET("/", func(c *gin.Context) {
		minPrice := c.Query("minPrice")
		maxPrice := c.Query("maxPrice")
		gameTitle := c.Query("title")
		if minPrice == "" || maxPrice == "" {
			c.Status(http.StatusBadRequest)
			return
		}
		_, errMin := strconv.Atoi(minPrice)
		_, errMax := strconv.Atoi(maxPrice)
		if errMin != nil || errMax != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		go parseGameDeals(ch, maxPrice, minPrice, gameTitle)
		var result []gameDeal = <-ch
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST,HEAD,PATCH, OPTIONS, GET, PUT")
		c.JSON(http.StatusOK, result)
	})
	r.Run()
}
