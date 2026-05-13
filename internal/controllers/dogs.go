package controllers

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/innovative-io/go-api-sample/internal/models"
	"github.com/innovative-io/go-api-sample/internal/services"
)

// @Summary Deletes a dog by ID
// @Description deletes a dog
// @Produce  json
// @Param        dog_id    path      string     true  "Dog ID"
// @Success      200   {object}   interface{}	"ok"
// @Failure      400   {string}   string  "ok"
// @Failure      404   {string}   string  "ok"
// @Failure      500   {string}   string  "ok"
// @Router /dogs/{dog_id} [delete]
func (r *Router) DogsDelete(w http.ResponseWriter, req *http.Request) {
	id, err := uuid.Parse(req.PathValue("id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"message": "invalid id"})
		return
	}

	if err := r.dogs.Delete(req.Context(), id); err != nil {
		if errors.Is(err, services.ErrNotFound) {
			writeJSON(w, http.StatusNotFound, map[string]any{"message": "not found"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]any{"message": "there was an error"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"deleted": id.String()})
}

// @Summary Gets the total count of dogs
// @Description returns the total number of dogs
// @Produce  json
// @Success 200 {object} map[string]int64 "ok"
// @Failure 500 {string} string "error"
// @Router /dogs/count [get]
func (r *Router) DogsCount(w http.ResponseWriter, req *http.Request) {
	count, err := r.dogs.Count(req.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"message": "there was an error"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"count": count})
}

// @Summary Gets all the dogs in the database
// @Description get a list of dogs
// @Produce  json
// @Success 200 {object} []models.Dog	"ok"
// @Failure      400   {string}   string  "ok"
// @Failure      404   {string}   string  "ok"
// @Failure      500   {string}   string  "ok"
// @Router /dogs [get]
func (r *Router) DogsGet(w http.ResponseWriter, req *http.Request) {
	dogs, err := r.dogs.Get(req.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"message": "there was an error"})
		return
	}
	writeJSON(w, http.StatusOK, dogs)
}

// @Summary Gets a dog by ID
// @Description get a dog
// @Produce  json
// @Param        dog_id    path      string     true  "Dog ID"
// @Success 200 {object} models.Dog	"ok"
// @Failure      400   {string}   string  "ok"
// @Failure      404   {string}   string  "ok"
// @Failure      500   {string}   string  "ok"
// @Router /dogs/{dog_id} [get]
func (r *Router) DogsGetOne(w http.ResponseWriter, req *http.Request) {
	id, err := uuid.Parse(req.PathValue("id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"message": "invalid id"})
		return
	}

	dog, err := r.dogs.GetOne(req.Context(), id)
	if err != nil {
		if errors.Is(err, services.ErrNotFound) {
			writeJSON(w, http.StatusNotFound, map[string]any{"message": "not found"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]any{"message": "there was an error"})
		return
	}
	writeJSON(w, http.StatusOK, dog)
}

// @Summary Adds a dog
// @Description adds a dog
// @Accept   json
// @Produce  json
// @Param        message  body      models.Dog  true  "Dog"
// @Success      201   {string}  string  "answer"
// @Failure      400   {string}   string  "ok"
// @Failure      500   {string}   string  "ok"
// @Router /dogs [post]
func (r *Router) DogsPost(w http.ResponseWriter, req *http.Request) {
	dog := new(models.Dog)
	if err := r.bind(w, req, dog); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"message": "Bad request body"})
		return
	}

	if dog.ID == uuid.Nil {
		dog.ID = uuid.New()
	}

	id, err := r.dogs.Add(req.Context(), dog)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"message": "there was an error"})
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{"id": id.String()})
}

// @Summary Updates a dog by ID
// @Description updates a dog
// @Accept   json
// @Produce  json
// @Param        dog_id    path      string     true  "Dog ID"
// @Param        message  body      models.Dog  true  "Dog"
// @Success      202   {string}  string  "answer"
// @Failure      400   {string}   string  "ok"
// @Failure      404   {string}   string  "ok"
// @Failure      500   {string}   string  "ok"
// @Router /dogs/{dog_id} [put]
func (r *Router) DogsPut(w http.ResponseWriter, req *http.Request) {
	dog := new(models.Dog)
	if err := r.bind(w, req, dog); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"message": "Bad request body"})
		return
	}

	id, err := uuid.Parse(req.PathValue("id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"message": "invalid id"})
		return
	}

	if err := r.dogs.Update(req.Context(), id, dog); err != nil {
		if errors.Is(err, services.ErrNotFound) {
			writeJSON(w, http.StatusNotFound, map[string]any{"message": "not found"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]any{"message": "there was an error"})
		return
	}
	writeJSON(w, http.StatusAccepted, map[string]any{"id": id.String()})
}
