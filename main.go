package main

import (
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
	// router.Run(":8787") 指定port
	router.Run()

}
func getting(c *gin.Context) {
	p1 := c.Param("param1")
	p2 := c.Param("param2")
	c.JSON(http.StatusOK, gin.H{"param1": p1, "param2": p2})
}
