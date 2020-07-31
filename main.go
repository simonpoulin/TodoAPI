package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ctxt       context.Context
	collection *mongo.Collection
)

type todo struct {
	ID      primitive.ObjectID `json:"_id" bson:"_id"`
	Status  bool               `json:"isComplete" bson:"isComplete"`
	Content string             `json:"content" bson:"content"`
}

func hi(c echo.Context) error {
	return c.String(http.StatusOK, "Yass, I'm in!")
}

func CreateTodo(c echo.Context) error {
	var t todo
	if err := c.Bind(&t); err != nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}
	t.ID = primitive.NewObjectID()
	t.Status = false
	_, err := collection.InsertOne(ctxt, t)
	if err != nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}
	return c.JSON(http.StatusOK, t)
}

func GetTodos(c echo.Context) error {
	cur, err := collection.Find(ctxt, bson.M{})
	if err != nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}
	defer cur.Close(ctxt)
	var todos []todo
	for cur.Next(ctxt) {
		var result todo
		err := cur.Decode(&result)
		if err != nil {
			return c.JSON(http.StatusNotFound, err.Error())
		}
		todos = append(todos, result)
	}
	if err := cur.Err(); err != nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}
	return c.JSON(http.StatusOK, todos)
}

func GetTodo(c echo.Context) error {
	id := c.Param("id")
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}
	filter := bson.M{"_id": _id}
	var result todo
	err = collection.FindOne(ctxt, filter).Decode(&result)
	if err != nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}
	return c.JSON(http.StatusOK, result)
}

func GetActiveTodos(c echo.Context) error {
	cur, err := collection.Find(ctxt, bson.M{})
	if err != nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}
	defer cur.Close(ctxt)
	var todos []todo
	for cur.Next(ctxt) {
		var result todo
		err := cur.Decode(&result)
		if err != nil {
			return c.JSON(http.StatusNotFound, err.Error())
		}
		if !result.Status {
			todos = append(todos, result)
		}
	}
	if err := cur.Err(); err != nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}
	return c.JSON(http.StatusOK, todos)
}

func SelectTodo(c echo.Context) error {
	id := c.Param("id")
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}
	filter := bson.M{"_id": _id}
	var result todo
	err = collection.FindOne(ctxt, filter).Decode(&result)
	if err != nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}
	result.Status = !result.Status
	update := bson.M{"$set": bson.M{"isComplete": result.Status}}
	_, err = collection.UpdateOne(ctxt, filter, update)
	return c.JSON(http.StatusOK, result)
}

func UpdateTodo(c echo.Context) error {
	var t todo
	if err := c.Bind(&t); err != nil {
		return err
	}
	id := c.Param("id")
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}
	filter := bson.M{"_id": _id}
	update := bson.M{"$set": bson.M{"content": t.Content}}
	_, err = collection.UpdateOne(ctxt, filter, update)
	if err != nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}
	var result todo
	err = collection.FindOne(ctxt, filter).Decode(&result)
	if err != nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}
	return c.JSON(http.StatusOK, result)
}

func SelectTodos(c echo.Context) error {
	cur, err := collection.Find(ctxt, bson.M{})
	if err != nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}
	defer cur.Close(ctxt)
	var todos []todo
	for cur.Next(ctxt) {
		var result todo
		err := cur.Decode(&result)
		if err != nil {
			return c.JSON(http.StatusNotFound, err.Error())
		}
		todos = append(todos, result)
	}
	if err := cur.Err(); err != nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}
	count := 0
	for _, t := range todos {
		if !t.Status {
			count++
		} else {
			count--
		}
	}
	stt := false
	if count >= 0 {
		stt = true
	}
	var _todos []todo
	for _, i := range todos {
		i.Status = stt
		_todos = append(_todos, i)
	}
	filter := bson.M{"isComplete": !stt}
	update := bson.M{"$set": bson.M{"isComplete": stt}}
	_, err = collection.UpdateMany(ctxt, filter, update)
	if err != nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}
	return c.JSON(http.StatusOK, _todos)
}

func DeleteTodos(c echo.Context) error {
	filter := bson.M{"isComplete": true}
	_, err := collection.DeleteOne(ctxt, filter)
	collection.DeleteMany(ctxt, filter)
	if err != nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}
	cur, err := collection.Find(ctxt, bson.M{})
	if err != nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}
	defer cur.Close(ctxt)
	var todos []todo
	for cur.Next(ctxt) {
		var result todo
		err := cur.Decode(&result)
		if err != nil {
			return c.JSON(http.StatusNotFound, err.Error())
		}
		todos = append(todos, result)
	}
	if err := cur.Err(); err != nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}
	return c.JSON(http.StatusOK, todos)
}

func DeleteTodo(c echo.Context) error {
	id := c.Param("id")
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}
	filter := bson.M{"_id": _id}
	_, err = collection.DeleteOne(ctxt, filter)
	if err != nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}
	cur, err := collection.Find(ctxt, bson.M{})
	if err != nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}
	defer cur.Close(ctxt)
	var todos []todo
	for cur.Next(ctxt) {
		var result todo
		err := cur.Decode(&result)
		if err != nil {
			return c.JSON(http.StatusNotFound, err.Error())
		}
		todos = append(todos, result)
	}
	if err := cur.Err(); err != nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}
	return c.JSON(http.StatusOK, todos)
}

func main() {
	//MongoDB
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://simonpl:123@cluster0.cg68p.mongodb.net/todo?retryWrites=true&w=majority"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	collection = client.Database("todo").Collection("todo")
	ctxt = context.Background()
	//Echo
	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
	}))

	e.GET("/", hi)

	e.POST("/todos", CreateTodo) //Tạo Todo

	e.GET("/todos", GetTodos) //Lấy danh sách Todo

	e.GET("/todos/:id", GetTodo) //Lấy Todo theo ID

	e.GET("/todos/active", GetActiveTodos) //Lấy danh sách Todo trạng thái false

	e.PATCH("/todos/select/:id", SelectTodo) //Chọn Todo theo ID

	e.PATCH("/todos/:id", UpdateTodo) //Sửa Todo theo ID

	e.PUT("/todos", SelectTodos) //Chọn hết (bỏ hết) Todo

	e.DELETE("/todos", DeleteTodos) //Xóa hết Todo trạng thái true

	e.DELETE("/todos/:id", DeleteTodo) //Xóa Todo theo ID

	e.Start(":9999")
}
