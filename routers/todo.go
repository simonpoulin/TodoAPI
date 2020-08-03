package routers

import (
	"net/http"
	m "todoapi/models"
	u "todoapi/utils"

	"github.com/go-playground/validator"
	"github.com/labstack/echo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateTodo(c echo.Context) error {
	collection, ctxt := u.ConnectDB()
	v := validator.New()
	var t m.Todo
	if err := c.Bind(&t); err != nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}
	if err := v.Struct(t); err != nil {
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
	collection, ctxt := u.ConnectDB()
	cur, err := collection.Find(ctxt, bson.M{})
	if err != nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}
	defer cur.Close(ctxt)
	var todos []m.Todo
	for cur.Next(ctxt) {
		var result m.Todo
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
	collection, ctxt := u.ConnectDB()
	id := c.Param("id")
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}
	filter := bson.M{"_id": _id}
	var result m.Todo
	err = collection.FindOne(ctxt, filter).Decode(&result)
	if err != nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}
	return c.JSON(http.StatusOK, result)
}

func GetActiveTodos(c echo.Context) error {
	collection, ctxt := u.ConnectDB()
	cur, err := collection.Find(ctxt, bson.M{})
	if err != nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}
	defer cur.Close(ctxt)
	var todos []m.Todo
	for cur.Next(ctxt) {
		var result m.Todo
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
	collection, ctxt := u.ConnectDB()
	id := c.Param("id")
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}
	filter := bson.M{"_id": _id}
	var result m.Todo
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
	collection, ctxt := u.ConnectDB()
	v := validator.New()
	var t m.Todo
	if err := c.Bind(&t); err != nil {
		return err
	}
	if err := v.Struct(t); err != nil {
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
	var result m.Todo
	err = collection.FindOne(ctxt, filter).Decode(&result)
	if err != nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}
	return c.JSON(http.StatusOK, result)
}

func SelectTodos(c echo.Context) error {
	collection, ctxt := u.ConnectDB()
	cur, err := collection.Find(ctxt, bson.M{})
	if err != nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}
	defer cur.Close(ctxt)
	var todos []m.Todo
	for cur.Next(ctxt) {
		var result m.Todo
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
	var _todos []m.Todo
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
	collection, ctxt := u.ConnectDB()
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
	var todos []m.Todo
	for cur.Next(ctxt) {
		var result m.Todo
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
	collection, ctxt := u.ConnectDB()
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
	var todos []m.Todo
	for cur.Next(ctxt) {
		var result m.Todo
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
