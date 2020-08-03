package main

import (
	"net/http"
	r "todoapi/routers"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	//Echo
	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
	}))

	e.POST("/todos", r.CreateTodo) //Tạo Todo

	e.GET("/todos", r.GetTodos) //Lấy danh sách Todo

	e.GET("/todos/:id", r.GetTodo) //Lấy Todo theo ID

	e.GET("/todos/active", r.GetActiveTodos) //Lấy danh sách Todo trạng thái false

	e.PATCH("/todos/select/:id", r.SelectTodo) //Chọn Todo theo ID

	e.PATCH("/todos/:id", r.UpdateTodo) //Sửa Todo theo ID

	e.PUT("/todos", r.SelectTodos) //Chọn hết (bỏ hết) Todo

	e.DELETE("/todos", r.DeleteTodos) //Xóa hết Todo trạng thái true

	e.DELETE("/todos/:id", r.DeleteTodo) //Xóa Todo theo ID

	e.Start(":9999")
}
