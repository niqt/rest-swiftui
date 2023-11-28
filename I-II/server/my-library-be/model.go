package main

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
	"log"
	"os"
)

// save the book
func SaveBook(book Book) Response {

	if book.ID.IsZero() { // if it's new
		result, err := GetDB().Collection(BooksCollection).InsertOne(
			context.Background(), book)

		if err != nil {
			return Response{
				Books:    []Book{},
				ErrorMsg: err.Error(),
				Total:    0,
			}
		}

		oid, _ := result.InsertedID.(primitive.ObjectID)
		book.ID = oid
		var books []Book
		books = append(books, book)
		return Response{
			Books:    books,
			ErrorMsg: "",
			Total:    1,
		}
	} else { // if it's an update
		_, err := GetDB().Collection(BooksCollection).ReplaceOne(context.Background(), bson.M{"_id": book.ID}, book)
		if err != nil {
			return Response{
				Books:    []Book{},
				ErrorMsg: err.Error(),
				Total:    0,
			}
		}
		var books []Book
		books = append(books, book)
		return Response{
			Books:    books,
			ErrorMsg: "",
			Total:    1,
		}
	}
}

// get the books by the query and user
func GetBooks(query, email string) Response {
	filter := bson.D{{"user", email}} // default query, if it's empty all the books for the user

	if len(query) > 0 { // if some text is specfied
		filter = bson.D{{"$and",
			bson.A{bson.D{{Key: "user", Value: email}},
				bson.D{{"$or",
					bson.A{
						bson.D{{Key: "title", Value: bson.M{"$regex": ".*" + query + ".*", "$options": "i"}}},
						bson.D{{Key: "isbn13", Value: bson.M{"$regex": ".*" + query + ".*", "$options": "i"}}},
						bson.D{{Key: "authors", Value: bson.M{"$regex": ".*" + query + ".*", "$options": "i"}}},
					}}}},
		}}
	}
	cur, err := GetDB().Collection(BooksCollection).Find(ctx, filter)

	if err != nil { // error getting the books
		return Response{
			Books:    []Book{},
			ErrorMsg: err.Error(),
			Total:    0,
		}
	} else { // get all the book
		var books []Book
		for cur.Next(ctx) {
			bookFromMongo := Book{}

			if err = cur.Decode(&bookFromMongo); err != nil {
				continue
			}
			books = append(books, bookFromMongo)
		}
		return Response{
			Books:    books,
			ErrorMsg: "",
			Total:    len(books),
		}
	}
}

// get the book by id
func GetBookById(id string) Response {
	var book Book
	if len(id) == 0 {
		return Response{
			Books:    []Book{},
			ErrorMsg: "wrong id",
			Total:    0,
		}
	}
	objID, err := primitive.ObjectIDFromHex(id)

	if err != nil { // wrong id
		return Response{
			Books:    []Book{},
			ErrorMsg: err.Error(),
			Total:    0,
		}
	}

	err = GetDB().Collection(BooksCollection).FindOne(ctx, bson.M{"_id": objID}).Decode(&book)

	if err != nil { // error getting the book
		return Response{
			Books:    []Book{},
			ErrorMsg: err.Error(),
			Total:    0,
		}
	}

	var books []Book
	books = append(books, book)
	return Response{
		Books:    books,
		ErrorMsg: "",
		Total:    1,
	}
}

// delete the book by id
func DeleteBookById(id string) error {
	if len(id) == 0 {
		return errors.New("wrong id")
	}
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("wrong id")
	}
	_, err = GetDB().Collection(BooksCollection).DeleteOne(context.TODO(), bson.D{{"_id", objID}})
	return err
}

// create the user
func CreateUser(account User) (map[string]User, error) {
	// hash the password
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)
	account.Password = string(hashedPassword)

	result, err := GetDB().Collection(UserCollection).InsertOne(
		context.Background(), account)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	oid, _ := result.InsertedID.(primitive.ObjectID)

	account.ID = oid

	if account.ID.IsZero() {
		return nil, errors.New("wrong user id")
	}

	// create the jwt token
	tk := &Token{UserId: account.ID, Username: account.Email}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))

	account.Token = tokenString
	account.Password = "" //delete password, no returned in the response due security reason

	var response = make(map[string]User)
	response["user"] = account

	return response, nil
}

// login
func Login(email, password string) (map[string]interface{}, error) {

	account := &User{}

	err := GetDB().Collection(UserCollection).FindOne(ctx, bson.M{"email": email}).Decode(&account)

	if err != nil {
		return nil, err
	}

	// verify the password
	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
		return nil, err
	}

	//Worked! Logged In
	account.Password = ""

	//Create a new JWT token
	tk := &Token{UserId: account.ID, Username: account.Email}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	account.Token = tokenString //Store the token in the response

	var response = make(map[string]interface{})
	response["user"] = *account
	return response, nil
}

func UserExists(email string) (bool, error) {
	account := &User{}

	err := GetDB().Collection(UserCollection).FindOne(ctx, bson.M{"email": email}).Decode(&account)

	if err != nil {
		return false, err
	}

	return true, nil
}
