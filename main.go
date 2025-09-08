package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yanaga11/todo/models"
)

func main() {
	if err := models.InitDB(); err != nil {
		log.Fatalln(err)
	}
	r := setupRouter()
	r.Run(":8080")
}

func setupRouter() *gin.Engine {

	//templateディレクトリ設定
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/todos/list")
	})

	// todo create
	r.POST("/todos/save", func(c *gin.Context) {

		models.CreateTodo(c.PostForm("content"))

		c.Redirect(http.StatusFound, "/todos/list")
	})

	// todo create
	r.POST("/todos/update", func(c *gin.Context) {
		id, _ := strconv.Atoi(c.PostForm("id"))
		content := c.PostForm("content")
		todo, _ := models.ListTodo(id)
		todo.Content = content
		models.UpdateTodo(todo)

		c.Redirect(http.StatusFound, "/todos/list")
	})

	//todo edit
	r.GET("/todos/edit", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Query("id"))
		if err != nil {
			log.Fatalln(err)
		}
		todo, _ := models.ListTodo(id)

		c.HTML(http.StatusOK, "edit.html", gin.H{
			"title": "Todo",
			"todo":  todo,
		})
	})

	//todo delete
	r.GET("/todos/destroy", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Query("id"))
		if err != nil {
			log.Fatalln(err)
		}
		models.DeleteTodo(id)

		c.Redirect(http.StatusFound, "/todos/list")
	})

	//todo list
    r.GET("/todos/list", func(c *gin.Context) {
        todos, err := models.GetAllTodos()
        if err != nil {
            c.String(http.StatusInternalServerError, "DB error")
            return
        }
        c.HTML(http.StatusOK, "list.html", gin.H{
            "title": "Todo",
            "todos": todos,
        })
    })

	//todo Search
    r.GET("/todos/search", func(c *gin.Context) {
        content := c.Query("content")
        todos, err := models.SearchTodos(content)
        if err != nil {
            c.String(http.StatusInternalServerError, "DB error")
            return
        }
        c.HTML(http.StatusOK, "list.html", gin.H{
            "title": "Todo",
            "todos": todos,
        })
    })

	return r
}
