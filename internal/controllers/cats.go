package controllers

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/innovative-io/go-api-sample/internal/models"
	"github.com/innovative-io/go-api-sample/internal/services"
)

// @Summary Deletes a cat by ID
// @Description deletes a cat
// @Produce  json
// @Success 200 {object} interface{}	"ok"
// @Failure      400   {string}   string  "ok"
// @Failure      404   {string}   string  "ok"
// @Failure      500   {string}   string  "ok"
// @Router /cats/{cat_id} [delete]
func (r *Router) CatsDelete(w http.ResponseWriter, req *http.Request) {
	id, err := uuid.Parse(req.PathValue("id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"message": "invalid id"})
		return
	}

	if err := r.cats.Delete(req.Context(), id); err != nil {
		if errors.Is(err, services.ErrNotFound) {
			writeJSON(w, http.StatusNotFound, map[string]any{"message": "not found"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]any{"message": "there was an error"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"deleted": id.String()})
}

// @Summary Gets the total count of cats
// @Description returns the total number of cats
// @Produce  json
// @Success 200 {object} map[string]int64 "ok"
// @Failure 500 {string} string "error"
// @Router /cats/count [get]
func (r *Router) CatsCount(w http.ResponseWriter, req *http.Request) {
	count, err := r.cats.Count(req.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"message": "there was an error"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"count": count})
}

// @Summary Gets all the cats in the database
// @Description get a list of cats
// @Produce  json
// @Success 200 {object} []models.Cat	"ok"
// @Failure      400   {string}   string  "ok"
// @Failure      404   {string}   string  "ok"
// @Failure      500   {string}   string  "ok"
// @Router /cats [get]
func (r *Router) CatsGet(w http.ResponseWriter, req *http.Request) {
	cats, err := r.cats.Get(req.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"message": "there was an error"})
		return
	}
	writeJSON(w, http.StatusOK, cats)
}

// @Summary Gets a cat by ID
// @Description get a cat
// @Produce  json
// @Param        cat_id    path      string     true  "Cat ID"
// @Success 200 {object} models.Cat	"ok"
// @Failure      400   {string}   string  "ok"
// @Failure      404   {string}   string  "ok"
// @Failure      500   {string}   string  "ok"
// @Router /cats/{cat_id} [get]
func (r *Router) CatsGetOne(w http.ResponseWriter, req *http.Request) {
	id, err := uuid.Parse(req.PathValue("id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"message": "invalid id"})
		return
	}

	cat, err := r.cats.GetOne(req.Context(), id)
	if err != nil {
		if errors.Is(err, services.ErrNotFound) {
			writeJSON(w, http.StatusNotFound, map[string]any{"message": "not found"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]any{"message": "there was an error"})
		return
	}
	writeJSON(w, http.StatusOK, cat)
}

// @Summary Adds a cat
// @Description adds a cat
// @Accept   json
// @Produce  json
// @Param        message  body      models.Cat  true  "Cat"
// @Success      201   {string}  string  "answer"
// @Failure      400   {string}   string  "ok"
// @Failure      500   {string}   string  "ok"
// @Router /cats [post]
func (r *Router) CatsPost(w http.ResponseWriter, req *http.Request) {
	cat := new(models.Cat)
	if err := r.bind(w, req, cat); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"message": "Bad request body"})
		return
	}

	if cat.ID == uuid.Nil {
		cat.ID = uuid.New()
	}

	id, err := r.cats.Add(req.Context(), cat)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"message": "there was an error"})
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{"id": id.String()})
}

// @Summary Updates a cat by ID
// @Description updates a cat
// @Accept   json
// @Produce  json
// @Param        cat_id    path      string     true  "Cat ID"
// @Param        message  body      models.Cat  true  "Cat"
// @Success      202   {string}  string  "answer"
// @Failure      400   {string}   string  "ok"
// @Failure      404   {string}   string  "ok"
// @Failure      500   {string}   string  "ok"
// @Router /cats/{cat_id} [put]
func (r *Router) CatsPut(w http.ResponseWriter, req *http.Request) {
	cat := new(models.Cat)
	if err := r.bind(w, req, cat); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"message": "Bad request body"})
		return
	}

	id, err := uuid.Parse(req.PathValue("id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"message": "invalid id"})
		return
	}

	if err := r.cats.Update(req.Context(), id, cat); err != nil {
		if errors.Is(err, services.ErrNotFound) {
			writeJSON(w, http.StatusNotFound, map[string]any{"message": "not found"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]any{"message": "there was an error"})
		return
	}
	writeJSON(w, http.StatusAccepted, map[string]any{"id": id.String()})
}
