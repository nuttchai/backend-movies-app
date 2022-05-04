package models

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type DBModel struct {
	DB *sql.DB
}

// Get returns one movie and error, if any
func (m *DBModel) GetMovie(id int) (*Movie, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel() // NOTE: if thing go wrong, it will be canceled

	// NOTE: $1 is placeholder for id
	query := `select id, title, description, year, release_date, 
		runtime, rating, mpaa_rating, created_at, updated_at from movies where id = $1`

	/* NOTE:
	ctx tells the database to timeout the query if it takes more than 3 seconds to query
	if you don't use ctx, the database will wait forever
	id will be put in the placeholder where we define in query parameter */
	row := m.DB.QueryRowContext(ctx, query, id)

	var movie Movie

	// NOTE: row.Scan() will scan the query result to the struct and return error if any
	err := row.Scan(
		&movie.ID,
		&movie.Title,
		&movie.Description,
		&movie.Year,
		&movie.ReleaseDate,
		&movie.Runtime,
		&movie.Rating,
		&movie.MPAARating,
		&movie.CreatedAt,
		&movie.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	// Get the genres of the movie
	query = `select 
						mg.id, mg.movie_id, mg.genre_id, g.genre_name
					from 
						movies_genres mg
						left join genres g on (g.id = mg.genre_id)
					where mg.movie_id = $1`

	rows, _ := m.DB.QueryContext(ctx, query, id)
	defer rows.Close() // NOTE: close the rows after we finish using it

	genres := make(map[int]string)

	// NOTE: rows.Next() will move the cursor to the next row
	for rows.Next() {
		var mg MovieGenre
		err := rows.Scan(
			&mg.ID,
			&mg.MovieID,
			&mg.GenreID,
			&mg.Genre.GenreName,
		)
		if err != nil {
			return nil, err
		}
		genres[mg.ID] = mg.Genre.GenreName
	}

	movie.MovieGenre = genres

	return &movie, nil
}

// All returns all movies and error, if any
func (m *DBModel) GetAllMovies(genre ...int) ([]*Movie, error) {
	// NOTE: genre ...int is a variadic parameter (allows a function to accept any number of extra arguments)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel() // NOTE: if thing go wrong, it will be canceled

	where := ""
	if len(genre) > 0 {
		where = fmt.Sprintf("where id in (select movie_id from movies_genres where genre_id = %d)", genre[0])
	}

	query := fmt.Sprintf(`select id, title, description, year, release_date, 
		runtime, rating, mpaa_rating, created_at, updated_at from movies %s order by title`, where)

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close() // NOTE: close the rows after we finish using it

	var movies []*Movie

	for rows.Next() {
		var movie Movie
		err := rows.Scan(
			&movie.ID,
			&movie.Title,
			&movie.Description,
			&movie.Year,
			&movie.ReleaseDate,
			&movie.Runtime,
			&movie.Rating,
			&movie.MPAARating,
			&movie.CreatedAt,
			&movie.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Get the genres of the movie
		genreQuery := `select 
										mg.id, mg.movie_id, mg.genre_id, g.genre_name
									from 
										movies_genres mg
										left join genres g on (g.id = mg.genre_id)
									where mg.movie_id = $1`

		genreRows, _ := m.DB.QueryContext(ctx, genreQuery, movie.ID)

		genres := make(map[int]string)
		for genreRows.Next() {
			var mg MovieGenre
			err := genreRows.Scan(
				&mg.ID,
				&mg.MovieID,
				&mg.GenreID,
				&mg.Genre.GenreName,
			)
			if err != nil {
				return nil, err
			}
			genres[mg.ID] = mg.Genre.GenreName
		}
		defer genreRows.Close() // NOTE: close the rows after we finish using it

		movie.MovieGenre = genres
		movies = append(movies, &movie)
	}

	return movies, nil
}

func (m *DBModel) GetAllGenres() ([]*Genre, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel() // NOTE: if thing go wrong, it will be canceled

	query := `select id, genre_name, created_at, updated_at from genres order by genre_name`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var genres []*Genre

	for rows.Next() {
		var g Genre
		err := rows.Scan(
			&g.ID,
			&g.GenreName,
			&g.CreatedAt,
			&g.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		genres = append(genres, &g)
	}

	return genres, nil
}
