package data
import (
	"context"
	"database/sql"
	"errors"
	"finalproject.arailym/internal/validator"
	"fmt"
	"github.com/lib/pq"
	"time"
)
type Anime struct {
	ID int64 `json:"id"`
	CreatedAt time.Time `json:"-"`
	Title string `json:"title"`
	Year int32 `json:"year,omitempty"`
	Runtime Runtime `json:"runtime,omitempty"`
	Genres []string `json:"genres,omitempty"`
	Version int32 `json:"version"`
}
func ValidateAnime(v *validator.Validator, movie *Anime) {
	v.Check(movie.Title != "", "title", "must be provided")
	v.Check(len(movie.Title) <= 500, "title", "must not be more than 500 bytes long")
	v.Check(movie.Year != 0, "year", "must be provided")
	v.Check(movie.Year >= 1888, "year", "must be greater than 1888")
	v.Check(movie.Year <= int32(time.Now().Year()), "year", "must not be in the future")
	v.Check(movie.Runtime != 0, "runtime", "must be provided")
	v.Check(movie.Runtime > 0, "runtime", "must be a positive integer")
	v.Check(movie.Genres != nil, "genres", "must be provided")
	v.Check(len(movie.Genres) >= 1, "genres", "must contain at least 1 genre")
	v.Check(len(movie.Genres) <= 5, "genres", "must not contain more than 5 genres")
	v.Check(validator.Unique(movie.Genres), "genres", "must not contain duplicate values")
}

type AnimeModel struct {
	DB *sql.DB
}
// Add a placeholder method for inserting a new record in the movies table.
func (m AnimeModel) Insert(anime *Anime) error {
	// Define the SQL query for inserting a new record in the movies table and returning
	// the system-generated data.
	query := `
INSERT INTO animes (title, year, runtime, genres)
VALUES ($1, $2, $3, $4)
RETURNING id, created_at, version`
	args := []interface{}{anime.Title, anime.Year, anime.Runtime, pq.Array(anime.Genres)}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return m.DB.QueryRowContext(ctx, query, args...).Scan(&anime.ID, &anime.CreatedAt, &anime.Version)
	}
func (m AnimeModel) Get(id int64) (*Anime, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	query := `
SELECT  id, created_at, title, year, runtime, genres, version
FROM animes
WHERE id = $1`
	var anime Anime

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	    err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&anime.ID,
		&anime.CreatedAt,
		&anime.Title,
		&anime.Year,
		&anime.Runtime,
		pq.Array(&anime.Genres),
		&anime.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &anime, nil

}

func (m AnimeModel) Update(anime *Anime) error {
	// Declare the SQL query for updating the record and returning the new version
	// number.
	query := `
UPDATE animes
SET title = $1, year = $2, runtime = $3, genres = $4, version = version + 1
WHERE id = $5 AND version=$6
RETURNING version`
	// Create an args slice containing the values for the placeholder parameters.
	args := []interface{}{
		anime.Title,
		anime.Year,
		anime.Runtime,
		pq.Array(anime.Genres),
		anime.ID,
		anime.Version,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&anime.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}

func (m AnimeModel) Delete(id int64) error {
	// Return an ErrRecordNotFound error if the movie ID is less than 1.
	if id < 1 {
		return ErrRecordNotFound
	}
	// Construct the SQL query to delete the record.
	query := `
DELETE FROM animes
WHERE id = $1`
	// Execute the SQL query using the Exec() method, passing in the id variable as
	// the value for the placeholder parameter. The Exec() method returns a sql.Result
	// object.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	result, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	// Call the RowsAffected() method on the sql.Result object to get the number of rows
	// affected by the query.
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	// If no rows were affected, we know that the movies table didn't contain a record
	// with the provided ID at the moment we tried to delete it. In that case we
	// return an ErrRecordNotFound error.
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}

// Create a new GetAll() method which returns a slice of movies. Although we're not
// using them right now, we've set this up to accept the various filter parameters as
// arguments.
func (m AnimeModel) GetAll(title string, genres []string, filters Filters) ([]*Anime, Metadata, error) {	// Construct the SQL query to retrieve all movie records.
	query := fmt.Sprintf(`
SELECT count(*) OVER(), id, created_at, title, year, runtime, genres, version
FROM animes
WHERE (to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) OR $1 = '')
AND (genres @> $2 OR $2 = '{}')
ORDER BY %s %s, id ASC
LIMIT $3 OFFSET $4`, filters.sortColumn(), filters.sortDirection())

	// Create a context with a 3-second timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	args := []interface{}{title, pq.Array(genres), filters.limit(), filters.offset()}
	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	// Use QueryContext() to execute the query. This returns a sql.Rows resultset
	// containing the result.

	defer rows.Close()
	totalRecords := 0
	// Initialize an empty slice to hold the movie data.
	animes := []*Anime{}
	// Use rows.Next to iterate through the rows in the resultset.
	for rows.Next() {
		// Initialize an empty Movie struct to hold the data for an individual movie.
		var anime Anime
		// Scan the values from the row into the Movie struct. Again, note that we're
		// using the pq.Array() adapter on the genres field here.
		err := rows.Scan(
			&totalRecords,
			&anime.ID,
			&anime.CreatedAt,
			&anime.Title,
			&anime.Year,
			&anime.Runtime,
			pq.Array(&anime.Genres),
			&anime.Version,
		)
		if err != nil {
			return nil, Metadata{}, err
		}
	// Add the Movie struct to the slice.
	animes = append(animes, &anime)
}
// When the rows.Next() loop has finished, call rows.Err() to retrieve any error
// that was encountered during the iteration.
if err = rows.Err(); err != nil {
	return nil, Metadata{}, err
}
	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)
// If everything went OK, then return the slice of movies.
	return animes,  metadata, nil
}






