package main

import (
	"awesomeProject4/Gee/day6-template/gee"
	"fmt"
	"html/template"
	"net/http"
	"time"
)

type student struct {
	Name string
	Age int8
}
func formAsDate(t time.Time) string{
	year,month,day := t.Date()
	return fmt.Sprintf("%d-%02d-%02d",year,month,day)
}
func main() {
	r := gee.New()
	r.Use(gee.Logger())
	r.SetFuncMap(template.FuncMap{
		"formatAsDate":formAsDate,
	})

	r.LoadHTMLGlob("C:/Users/dellp/go/src/awesomeProject4/Gee/day6-template/templates/*")
	r.Static("/assets","./static")
	stu1 := &student{
		Name:"jjz",
		Age:25,
	}
	stu2 := &student{
		Name: "Jack",
		Age:  20,
	}
	r.GET("/", func(c *gee.Context) {
		c.HTML(http.StatusOK,"css.tmpl",nil)
	})
	r.GET("/students",func(c *gee.Context){
		c.HTML(http.StatusOK,"arr.tmpl",gee.H{
			"title": "gee",
			"stuArr":[2]*student{stu1,stu2},
		})
	})
	r.GET("/date", func(c *gee.Context) {
		c.HTML(http.StatusOK,"custom_func.tmpl",gee.H{
			"title":"gee",
			"now":time.Date(2020,3,31,0,0,0,0,time.UTC),
		})
	})
	r.Run(":9999")

}
