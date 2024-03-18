package routes

import (
	"context"
	"fmt"
	"time"

	"github.com/betonomochalka/go-react-calorie-tracker/models"

	"github.com/gin-gonic/gin" //its web app framwork
	"go.mongodb.org/mongo-driver/mongo" // driver package
	"go.mongodb.org/mongo-driver/mongo/primitive" // primitive package
)

var entryCollection *mongo.Collection = OpenCollection(Client, "calories") // open collection

func CreateEntries(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second) //
	var entry models.Entry

	err := c.BindJSON(&entry) //binds the json to the entry struct
	if err!= nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmr.Println(err)
        return
    }

	validationErr := validate.Struct(entry) //validates the entry struct
	if validationErr!= nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
        fmr.Println(validationErr)
        return
    }

	entry.ID = primitive.NewObjectID() //creates a new object id
	result, isertErr := entryCollection.InsertOne(ctx, entry) //inserts the entry into the collection
	if insertErr!= nil {
		msg := fmt.Sprintf("Error inserting entry")
        c.JSON(http.StatusInternalServerError, gin.H{"error": insertErr.Error()})
        fmr.Println(insertErr)
        return
    }
	defer cancel()
	c.JSON(http.StatusCreated, result}
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
	var EntryID := c.Params.ByName("id")
	docID, _ := primitive.ObjectIDFromHex(EntryID) // convert the id to a mongo id

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second
	var entry bson.M
	if err := FindOne(ctx, bson.M{"_id": docID}).Decode(&entry) /* find the entry */ ; err!= nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        fmt.Println(err)
        return
	}
	defer cancel() // close the cursor
	fmt.Println(entry)
	c.JSON(http.StatusOK, entry) // return the entry
}

func GetEntriesByIngredient(c *gin.Context) {
	ingridient := c.Params.ByName("id")
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	var entries []bson.M
	cursor, err := entryCollection.Find(ctx, bson.M{"ingrrdients": ingridient})
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
			"dish" entry.Dish,
			"fat": entry.Fat,
			"ingridients": entry.Ingridients,
			"calories": entry.Calories,
		},
	)
	if err!= nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        fmt.Println(err)
        return
    }
	defer cancel()
	c.JSON(http.StatusCreated, result.ModifiedCount)}
}

func UpdateIngridient(c *gin.Context) {
	entryID := c.Params.ByName("id") // get the id from the url
	docID, _ := primitive.ObjectIDFromHex(entryID) // convert the id to a mongo id

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second

	type Ingridient struct {
		Ingridients *string	`json:"ingridients`
	}
	var ingridient Ingridient

	if err := c.BindJSON(&Ingridient); err!= nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        fmt.Println(err)
        return
	}

	result, err := entryCollection.UpdateOne(ctx, bson.M{"_id": docID},
		bson.D{{"$set": bson.D{"ingridients": ingridient.Ingridients}}}},
	)
	if err!= nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        fmt.Println(err)
        return
    }
	defer cancel()
	c.JSON(http.StatusCreated, result.ModifiedCount)}
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