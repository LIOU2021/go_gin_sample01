package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type book struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

var books = []book{
	{ID: "1", Title: "title01", Author: "mark01"},
	{ID: "2", Title: "title02", Author: "mark02"},
	{ID: "3", Title: "title03", Author: "mark03"},
	{ID: "4", Title: "title04", Author: "mark04"},
}

func main() {
	router := gin.Default()
	router.SetTrustedProxies([]string{"127.0.0.1"})

	//global middleware
	router.Use(GlobalMiddleware1(), CheckTokenMiddleware())

	router.GET("/DoGetByPath/:param1/*param2", getting)

	router.GET("/DoGetByQueryString", func(c *gin.Context) {
		// Parameters key值不存在時，回傳參數2的值
		p1 := c.DefaultQuery("param1", "Default")
		// 也可以寫成c.Request.URL.Query().Get("param2")
		//因為param2是使用*，所以可輸可不輸，但是如果使用*這前綴的話，取得的Parameters會有前綴"/"的符號存在
		p2 := c.Query("param2")
		c.JSON(http.StatusOK, gin.H{"param1": p1, "param2": p2})
	})

	router.GET("/", middleware1, func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": "hello world"})
	})

	g1 := router.Group("/v1").Use(middleware1)
	//By Group設定middleware
	g1.GET("/getting", func(c *gin.Context) {
		fmt.Println("doing v1 getting")
		c.JSON(http.StatusOK, gin.H{"data": "v1 getting"})
	})

	//By Group設定middleware
	g2 := router.Group("/v2").Use(middleware2)
	g2.GET("/getting", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": "v2 getting"})
	})

	router.GET("/api", func(c *gin.Context) {
		fmt.Println("First Middle Before Next")
		// c.Next()
		fmt.Println("First Middle After Next")

	}, func(c *gin.Context) {
		fmt.Println("Second Middle Before Next")
		// c.Next()
		// c.Abort()
		fmt.Println("Second Middle After Next")

	}, func(c *gin.Context) {

		fmt.Println("Third Middle Before Next")
		// c.Next()
		fmt.Println("Third Middle After Next")
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	router.GET("/api/1", func(c *gin.Context) {
		c.String(200, "ok")
	})

	g3 := router.Group("/books")
	g3.GET("/", IndexBooks)

	// router.Run(":8787") 指定port。預設8080
	router.Run("127.0.0.1:80") //指定127.0.0.1避免觸發win 防火牆
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

func middleware2(c *gin.Context) {
	fmt.Println("exec middleware2")
	//c.Abort()停止執行後面的hanlder，可以用來做auth
	c.Abort()
	c.JSON(200, gin.H{"msg": "i'm fail..."})
}

func CheckTokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//取得header的token值
		token := c.Query("token")
		if token == "" {
			//沒有token的話就不執行後面的handler
			c.Abort()
			//回傳錯誤值
			c.JSON(http.StatusOK, map[string]string{"msg": "token empty"})
		} else {
			//有token，繼續執行後面的handler
			fmt.Println("token checked")
			c.Next()
		}

	}
}

func GlobalMiddleware1() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("test GlobalMiddleware1")
	}
}

func IndexBooks(c *gin.Context) {
	c.JSON(http.StatusOK, books)
}
