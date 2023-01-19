package main

import (
	"html/template"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type Todo struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

var todos []Todo

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Serve the HTML template
	t := &Template{
		templates: template.Must(template.ParseGlob("index.html")),
	}
	e.Renderer = t

	// Define routes
	e.GET("/", index)
	e.POST("/todos", createTodo)
	e.DELETE("/todos/:id", deleteTodo)

	// Start the server
	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))
}

// Index page
func index(c echo.Context) error {
	return c.Render(http.StatusOK, "index.html", todos)
}

// Create a new todo
func createTodo(c echo.Context) error {
	t := Todo{}
	if err := c.Bind(&t); err != nil {
		return err
	}

	t.ID = len(todos) + 1
	todos = append(todos, t)

	return c.Redirect(http.StatusMovedPermanently, "/")
}

// Delete a todo
func deleteTodo(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	todos = append(todos[:id-1], todos[id:]...)

	return c.Redirect(http.StatusMovedPermanently, "/")
}

// HTML template
type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
