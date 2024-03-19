package routes

import (
	"context"
	"fmt"
	"time"
	"net/http"

	"github.com/betonomochalka/go-react-calorie-tracker/models"

	"github.com/go-playground/validator/v10"
	"github.com/gin-gonic/gin" //its web app framwork
	"go.mongodb.org/mongo-driver/mongo" // driver package
	"go.mongodb.org/mongo-driver/bson/primitive" // primitive package
	"go.mongodb.org/mongo-driver/bson"
)

var validate =validator.New()
var entryCollection *mongo.Collection = OpenCollection(Client, "calories") // open collection

func CreateEntries(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second) //
	var entry models.Entry

	err := c.BindJSON(&entry) //binds the json to the entry struct
	if err!= nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
        return
    }

	validationErr := validate.Struct(entry) //validates the entry struct
	if validationErr!= nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
        fmt.Println(validationErr)
        return
    }

	entry.ID = primitive.NewObjectID() //creates a new object id
	result, insertErr := entryCollection.InsertOne(ctx, entry) //inserts the entry into the collection
	if insertErr!= nil {
		msg := fmt.Sprintf("Error inserting entry")
        c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
        fmt.Println(insertErr)
        return
    }
	defer cancel()
	c.JSON(http.StatusCreated, result)
}

func GetEntries(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	var entries []bson.M 
	cursor, err := entryCollection.Find(ctx, bson.M{}) // find all entries
	if err!= nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
    }

	if err = cursor.All(ctx, &entries); err!= nil { // get all entries
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        fmt.Println(err)
        return
	}

	defer cancel() // close the cursor
	fmt.Println(entries)
	c.JSON(http.StatusOK, entries) // return all entries
}

func GetEntriesByID(c *gin.Context) {
	var EntryID = c.Params.ByName("id")
	docID, _ := primitive.ObjectIDFromHex(EntryID) // convert the id to a mongo id

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	var entry bson.M
	if err := entryCollection.FindOne(ctx, bson.M{"_id": docID}).Decode(&entry) /* find the entry */ ; err!= nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        fmt.Println(err)
        return
	}
	defer cancel() // close the cursor
	fmt.Println(entry)
	c.JSON(http.StatusOK, entry) // return the entry
}

func GetEntriesByIngredients(c *gin.Context) {
	ingredient := c.Params.ByName("id")
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	var entries []bson.M
	cursor, err := entryCollection.Find(ctx, bson.M{"ingrrdients": ingredient})
	if err!= nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        fmt.Println(err)
        return
    }
	if err = cursor.All(ctx, &entries); err!= nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        fmt.Println(err)
        return
    }
	defer cancel() // close the cursor
	fmt.Println(entries)
	c.JSON(http.StatusOK, entries) // return all entries
}

func UpdateEntry(c *gin.Context) {
	entryID := c.Params.ByName("id") // get the id from the url
	docID, _ := primitive.ObjectIDFromHex(entryID) // convert the id to a mongo id

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	var entry models.Entry

	if err := c.BindJSON(&entry); err!= nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}

	validatenErr := validate.Struct(entry) // validates the entry struct
	if validatenErr!= nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": validatenErr.Error()})
        fmt.Println(validatenErr)
        return
    }

	result, err := entryCollection.ReplaceOne(
		ctx,
		bson.M{"_id": docID},
        bson.M{
			"dish": entry.Dish,
			"fat": entry.Fat,
			"ingredients": entry.Ingredients,
			"calories": entry.Calories,
		},
	)
	if err!= nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        fmt.Println(err)
        return
    }
	defer cancel()
	c.JSON(http.StatusCreated, result.ModifiedCount)
}

func UpdateIngredient(c *gin.Context) {
	entryID := c.Params.ByName("id") // get the id from the url
	docID, _ := primitive.ObjectIDFromHex(entryID) // convert the id to a mongo id

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	type Ingredient struct {
		Ingredients *string	`json:"ingredients`
	}
	var ingredient Ingredient

	if err := c.BindJSON(&ingredient); err!= nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        fmt.Println(err)
        return
	}

	result, err := entryCollection.UpdateOne(ctx, bson.M{"_id": docID},
		bson.D{{"$set", bson.D{{"ingredients", ingredient.Ingredients}}}},
	)
	if err!= nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        fmt.Println(err)
        return
    }
	defer cancel()
	c.JSON(http.StatusCreated, result.ModifiedCount)
}

func DeleteEntry(c *gin.Context) {
	entryID := c.Params.ByName("id") //gets the id from the url
	docID, _ := primitive.ObjectIDFromHex(entryID) //gets the docID from the url

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second) //cancel

	result, err := entryCollection.DeleteOne(ctx, bson.M{"_id": docID}) //deletes the entry
	if err!= nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
    }

	defer cancel()
	c.JSON(http.StatusOK, result.DeletedCount)
}