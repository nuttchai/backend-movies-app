package models

import (
	"context"
	"database/sql"
	"time"
)

type DBModel struct {
	DB *sql.DB
}

// Get returns one movie and error, if any
func (m *DBModel) Get(id int) (*Movie, error) {
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

	var genres []MovieGenre

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
		genres = append(genres, mg)
	}

	movie.MovieGenre = genres

	return &movie, nil
}

// All returns all movies and error, if any
func (m *DBModel) All(id int) ([]*Movie, error) {
	return nil, nil
}
