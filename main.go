package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
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

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/todos/list")
	})

	// todo create
	r.POST("/todos/save", func(c *gin.Context) {
		todo := models.Todo{Content: c.PostForm("content")}
		validate := validator.New()
		if err := validate.Struct(todo); err != nil {
			c.String(http.StatusBadRequest, "Todo内容が未入力です")
			return
		}
		models.CreateTodo(todo.Content)
		c.Redirect(http.StatusSeeOther, "/todos/list")
	})

	// todo update
	r.POST("/todos/update", func(c *gin.Context) {
		id, err := strconv.Atoi(c.PostForm("id"))
		if err != nil {
			c.String(http.StatusBadRequest, "ID内容が未入力です")
			return
		}
		content := c.PostForm("content")
		todo, err := models.ListTodo(id)
		if err != nil {
			c.String(http.StatusNotFound, "Todoが見つかりません")
			return
		}
		todo.Content = content

		validate := validator.New()
		if err := validate.Struct(todo); err != nil {
			c.String(http.StatusBadRequest, "Todo内容が未入力です")
			return
		}

		models.UpdateTodo(todo)
		c.Redirect(http.StatusSeeOther, "/todos/list")
	})

	//todo edit
	r.GET("/todos/edit", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Query("id"))
		if err != nil {
			c.String(http.StatusBadRequest, "Invalid ID")
			return
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
			c.Redirect(http.StatusSeeOther, "/todos/list")
			return
		}
		models.DeleteTodo(id)
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
