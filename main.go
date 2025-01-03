package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Todo struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Completed bool   `json: "completed"`
	Body      string `json:"body"`
}


var collection *mongo.Collection
func main() {

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error in Loading env file")
	}

	PORT := os.Getenv("PORT")

	MONGODB_URL := os.Getenv("MONGODB_URI");
	clientOptions := options.Client().ApplyURI(MONGODB_URL);
	client,err := mongo.Connect(context.Background(), clientOptions);

	if err != nil {
		log.Fatal(err)
	}

    // close the mongodb connection after execution of main.
	defer client.Disconnect(context.Background())

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to Mongodb!!")

	collection = client.Database("golang").Collection("todos")

	app := fiber.New()

	app.Get("/api/todos", getTodos);
	app.Post("/api/todos", createTodo);
	app.Put("/api/todos/:id", updateTodo);
	app.Delete("/api/todos/:id", deleteTodo);

	log.Fatal(app.Listen(":" + PORT))
}

func getTodos(c *fiber.Ctx) error {
	var todos []Todo

	cursor, err := collection.Find(context.Background(),bson.M{})

	if err!= nil {
		return err;
	}

	//closing the db connection once function execution is done.
	defer cursor.Close(context.Background());

	for cursor.Next(context.Background()) {
		var todo Todo

		if err := cursor.Decode(&todo); err != nil {
			return err;
		}
		todos = append(todos, todo)
	}

	return c.JSON(todos)
}

func createTodo(c *fiber.Ctx) error {
	todo:= new(Todo)

	if err:= c.BodyParser(todo); err != nil {
		return err;
	}

	if todo.Body == "" {
		return c.Status(400).JSON(fiber.Map {"error": "Todo body cannot be empty"})
	}

	insertResult,err := collection.InsertOne(context.Background(), todo);
	if err != nil {
		return err;
	}

	todo.ID = insertResult.InsertedID.(primitive.ObjectID)

	return c.Status(201).JSON(todo)
}

func updateTodo(c *fiber.Ctx) error {
	id := c.Params("id");

	objectId , err := primitive.ObjectIDFromHex(id);

	if err != nil {
		return c.Status(400).JSON(fiber.Map {"error": "Invalid todo Id"})
	}

	filter := bson.M { "_id" : objectId}
	update := bson.M { "$set": bson.M{ "completed": true }};

	_ , err = collection.UpdateOne(context.Background(),filter, update)

	if err != nil {
		return err
	}

	return c.Status(200).JSON(fiber.Map{ "success" :  true })
}

func deleteTodo(c *fiber.Ctx) error {
	id := c.Params("id");

	objectId , err := primitive.ObjectIDFromHex(id);

	if err != nil {
		return c.Status(400).JSON(fiber.Map {"error": "Invalid todo Id"})
	}

	filter := bson.M{ "_id": objectId }
	_, err = collection.DeleteOne(context.Background(),filter)

	if err != nil {
		return nil
	}

	return c.Status(200).JSON(fiber.Map{ "success" : true });
}
