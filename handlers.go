package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func (a *App) GetAllPost() http.Handler {
	db := a.Broker.GetPostgres()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		posts := make([]*Post, 0)
		defer r.Body.Close()
		err := db.Table("posts").Find(&posts).Error
		if err != nil {
			log.Printf("get all posts %v", err)
			JSONResponse(w, http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
			return
		}

		JSONResponse(w, http.StatusOK, posts)
	})
}

func (a *App) GetSinglePost() http.Handler {
	db := a.Broker.GetPostgres()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		var post Post
		defer r.Body.Close()
		err := db.Table("posts").Where("id = ?", vars["post_id"]).First(&post).Error
		if err != nil {
			log.Printf("get single post %v", err)
			JSONResponse(w, http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
			return
		}

		JSONResponse(w, http.StatusOK, post)
	})
}

func (a *App) CreatePost() http.Handler {
	db := a.Broker.GetPostgres()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var post Post
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&post)
		if err != nil {
			JSONResponse(w, http.StatusBadRequest, map[string]string{"error": "invalid payload"})
			return
		}
		defer r.Body.Close()

		err = db.Create(&post).Error
		if err != nil {
			log.Printf("create post error %v", err)
			JSONResponse(w, http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
			return
		}

		JSONResponse(w, http.StatusCreated, nil)
	})
}

func (a *App) UpdatePost() http.Handler {
	db := a.Broker.GetPostgres()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		var post Post
		var newPost Post
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&newPost)
		if err != nil {
			JSONResponse(w, http.StatusBadRequest, map[string]string{"error": "invalid payload"})
			return
		}
		defer r.Body.Close()

		err = db.Table("posts").Where("id = ?", vars["post_id"]).First(&post).Error
		if err != nil {
			log.Printf("update post fetch error %v", err)
			JSONResponse(w, http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
			return
		}

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
