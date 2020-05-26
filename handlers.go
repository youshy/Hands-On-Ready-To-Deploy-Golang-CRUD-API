package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
)

// For the sake of not flooding this file with comments,
// I'll heavy comment this one heavily,
// And add whatever new in each handler.

// Each handler follows the same design pattern.
// I'm using the decorator pattern here, meaning
// that each handler returns a handler func.
// In other words - anything that might be reused by any consequent handler
// is created before. This way the garbage collector won't clean the memory for this variable.
func (a *App) GetAllPost() http.Handler {
	// Get the DB connection
	db := a.Broker.GetPostgres()
	// Basically this is a closure here - we're creating a function
	// that wraps another function and shares the variables across them.
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create an array of Posts.
		// make is a built-in function that creates a dynamically-sized slice of any given type.
		// Also used with maps.
		posts := make([]*Post, 0)
		// Another keyword - defer means that this function will run after the function exits.
		// defer r.Body.Close() closes the stream from body - we're preventing here any memory leaks.
		defer r.Body.Close()
		// Query the database with no parameters and find all the posts
		err := db.Table("posts").Find(&posts).Error
		// Check if there's any error
		if err != nil {
			// Log the error
			log.Printf("get all posts %v", err)
			// Return the error to the user. We could use map[string]string here which is a representation of JSON files
			// but it makes a bit more sense to use interface{} here.
			// This way you can see how the implementation wouldn't have to change if we'd like to return something else.
			JSONResponse(w, http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
			return
		}

		// Return StatusOK (200) and array of posts.
		// JSONResponse will marshal the data into JSON automatically.
		JSONResponse(w, http.StatusOK, posts)
	})
}

// TODO: create a method for GetSinglePost

func (a *App) CreatePost() http.Handler {
	db := a.Broker.GetPostgres()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a variable to which the JSON payload will be unmarshaled.
		var post Post
		// Build the decoder, and decode the fields.
		// DisallowUnknownFields is a function that will return error,
		// if the payload has more fields than our required struct.
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		err := decoder.Decode(&post)
		if err != nil {
			JSONResponse(w, http.StatusBadRequest, map[string]interface{}{"error": err.Error()})
			return
		}
		defer r.Body.Close()

		// Create a new UUID V4.
		// Better way would be to do it automatically by the db, but this way you can see all the steps.
		uid, _ := uuid.NewV4()
		post.Id = uid
		// Create the post and check if there's any error
		err = db.Create(&post).Error
		if err != nil {
			log.Printf("create post error %v", err)
			JSONResponse(w, http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
			return
		}

		// Return StatusCreated (201) and nothing in the payload
		JSONResponse(w, http.StatusCreated, nil)
	})
}

func (a *App) UpdatePost() http.Handler {
	db := a.Broker.GetPostgres()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// We need some way to understand which post is to be updated.
		// mux.Vars returns the parameters from the handler route.
		vars := mux.Vars(r)
		var post Post
		var newPost Post
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		err := decoder.Decode(&newPost)
		if err != nil {
			JSONResponse(w, http.StatusBadRequest, map[string]interface{}{"error": err.Error()})
			return
		}
		defer r.Body.Close()

		// Find the post we need to update
		err = db.Table("posts").Where("id = ?", vars["post_id"]).First(&post).Error
		if err != nil {
			log.Printf("update post fetch error %v", err)
			JSONResponse(w, http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
			return
		}

		// Update the content of the post.
		// If we'd like to allow users to update more fields,
		// we could specify it here.
		post.Content = newPost.Content
		db.Save(&post)
		JSONResponse(w, http.StatusNoContent, nil)
	})
}

func (a *App) DeletePost() http.Handler {
	db := a.Broker.GetPostgres()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		defer r.Body.Close()
		err := db.Where("id = ?", vars["post_id"]).Delete(&Post{}).Error
		if err != nil {
			log.Printf("delete post etch error %v", err)
			JSONResponse(w, http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
			return
		}

		JSONResponse(w, http.StatusOK, nil)
	})
}

// JSONResponse is a helper function that wraps a marshaler, sets a header for response
// writes the code in the header and then returns the response.
// Note the type of output, which is interface{} - meaning, we could use anything here
// as a return.
func JSONResponse(w http.ResponseWriter, code int, output interface{}) {
	response, _ := json.Marshal(output)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
