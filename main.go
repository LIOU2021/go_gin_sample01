package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/DoGetByPath/:param1/*param2", getting)

	router.GET("/DoGetByQueryString", func(c *gin.Context) {
		// Parameters key值不存在時，回傳參數2的值
		p1 := c.DefaultQuery("param1", "Default")
		// 也可以寫成c.Request.URL.Query().Get("param2")
		//因為param2是使用*，所以可輸可不輸，但是如果使用*這前綴的話，取得的Parameters會有前綴"/"的符號存在
		p2 := c.Query("param2")
		c.JSON(http.StatusOK, gin.H{"param1": p1, "param2": p2})
	})

	router.Use(middleware1).GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": "hello world"})
	})

	g1 := router.Group("/v1").Use(middleware1)
	//By Group設定middleware
	g1.GET("/getting", func(c *gin.Context) {
		fmt.Println("doing v1 getting")
		c.JSON(http.StatusOK, gin.H{"data": "v1 getting"})
	})

	// router.Run(":8787") 指定port
	router.Run()

}

func getting(c *gin.Context) {
	p1 := c.Param("param1")
	p2 := c.Param("param2")
	c.JSON(http.StatusOK, gin.H{"param1": p1, "param2": p2})
}

func middleware1(c *gin.Context) {
	fmt.Println("exec middleware1")
	//c.Next() 執行middleware後面接的function，執行完後再回到middleware繼續執行下去
	// c.Next()
	fmt.Println("after exec middleware1")
}
