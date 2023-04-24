package main

import (
"bytes"
"encoding/json"
"net/http"
"net/http/httptest"
"testing"
"time"

"github.com/stretchr/testify/assert"
)

func TestHandlers(t *testing.T) {
	app := newTestApplication(t)
	defer app.close()

	t.Run("testPingHandler", func(t *testing.T) {
		rr := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/v1/ping", nil)
		assert.NoError(t, err)

		app.router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

		var response struct {
			Message string `json:"message"`
		}

		err = json.Unmarshal(rr.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "pong", response.Message)
	})

	t.Run("testShowMovieHandler", func(t *testing.T) {
		// Create a new movie
		movie := createRandomMovie(t, app)

		rr := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/v1/movies/"+movie.ID, nil)
		assert.NoError(t, err)

		app.router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

		var response struct {
			Movie *data.Movie `json:"movie"`
		}

		err = json.Unmarshal(rr.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, movie.ID, response.Movie.ID)
	})

	t.Run("testCreateMovieHandler", func(t *testing.T) {
		// Define the request body.
		jsonBody := []byte(`{
			"title": "Test movie",
			"year": 2022,
			"runtime": 90,
			"genres": ["comedy", "drama", "romance"],
			"version": 1
		}`)

		rr := httptest.NewRecorder()
		req, err := http.NewRequest("POST", "/v1/movies", bytes.NewBuffer(jsonBody))
		assert.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")

		app.router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusCreated, rr.Code)
		assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

		var response struct {
			ID string `json:"id"`
		}

		err = json.Unmarshal(rr.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotEmpty(t, response.ID)

		// Check that the movie was inserted into the database.
		movie, err := app.models.Movies.Get(response.ID)
		assert.NoError(t, err)
		assert.Equal(t, "Test movie", movie.Title)
		assert.Equal(t, 2022, movie.Year)
		assert.Equal(t, 90, movie.Runtime)
		assert.Equal(t, []string{"comedy", "drama", "romance"}, movie.Genres)
		assert.Equal(t, response.ID, movie.ID)
	})

	t.Run("testUpdateMovieHandler", func(t *testing.T) {
		// Create a new movie
		movie := createRandomMovie(t, app)

		// Define the request body.
		jsonBody := []byte(`{
			"title": "Updated test movie",
			"year": 2023,
			"runtime": 100,


