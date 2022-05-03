package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (app *application) getOneMovie(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		app.logger.Print(errors.New("invalid movie id"))
		app.errorJSON(w, err)
		return
	}

	app.logger.Println("get movie with id:", id)

	movie, err := app.models.DB.Get(id)

	if err != nil {
		app.logger.Println(err)
		app.errorJSON(w, err)
		return
	}

	// movie := models.Movie{
	// 	ID:          id,
	// 	Title:       "The Shawshank Redemption",
	// 	Description: "Two imprisoned",
	// 	Year:        1994,
	// 	ReleaseDate: time.Date(1994, time.January, 14, 0, 0, 0, 0, time.Local),
	// 	Runtime:     142,
	// 	Rating:      9,
	// 	MPAARating:  "R",
	// 	CreatedAt:   time.Now(),
	// 	UpdatedAt:   time.Now(),
	// }

	err = app.writeJSON(w, http.StatusOK, movie, "movie")
	if err != nil {
		app.logger.Println(err)
	}
}

func (app *application) getAllMovies(w http.ResponseWriter, r *http.Request) {

}
