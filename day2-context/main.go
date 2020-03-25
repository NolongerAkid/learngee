package main

import (
	"awesomeProject4/Gee/day2-context/gee"
	"net/http"
)

func main() {
	r := gee.New()

	r.GET("/",func(c *gee.Context){
		c.HTML(http.StatusOK,"<h1>Hello Gee </h1>")
	})

	r.GET("/hello",func(c *gee.Context){
		//expect /hello?name=jjz
		c.String(http.StatusOK,"hello %s,you are at %s\n",c.Query("name"),c.Path)
	})

	r.POST("/login",func(c *gee.Context){
		c.JSON(http.StatusOK,gee.H{
			"username" : c.PostForm("username"),
			"password" : c.PostForm("password"),
		})
	})

	r.Run(":9999")

}
