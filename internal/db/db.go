package db

import (
	"TestVK/internal/config"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Movie struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	ReleaseDate time.Time `json:"release_date"`
	Rating      float64   `json:"rating"`
}

type Actor struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	Gender    string    `json:"gender"`
	Birthdate time.Time `json:"birthdate"`
}

func Connection(config config.DBConfig) (*sql.DB, error) {
	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s sslmode=%s", config.Host, config.Port, config.User, config.Password, config.SSLMode)

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func AddActor(db *sql.DB, actor Actor) error {

	query := `
        INSERT INTO actors (name, gender, birthdate)
        VALUES ($1, $2, $3)
    `

	_, err := db.Exec(query, actor.Name, actor.Gender, actor.Birthdate)
	if err != nil {
		return err
	}

	return nil
}

func UpdateActor(db *sql.DB, actor Actor) error {
	query := `UPDATE actors SET `
	args := make([]interface{}, 0)
	argCounter := 1

	if actor.Name != "" {
		query += "name = $" + strconv.Itoa(argCounter) + ", "
		args = append(args, actor.Name)
		argCounter++
	}
	if actor.Gender != "" {
		query += "gender = $" + strconv.Itoa(argCounter) + ", "
		args = append(args, actor.Gender)
		argCounter++
	}
	if !actor.Birthdate.IsZero() {
		query += "birthdate = $" + strconv.Itoa(argCounter) + ", "
		args = append(args, actor.Birthdate)
		argCounter++
	}

	query = strings.TrimSuffix(query, ", ")

	query += " WHERE id = $1"

	args = append([]interface{}{actor.Id}, args...)

	_, err := db.Exec(query, args...)
	if err != nil {
		return err
	}

	return nil
}

func DeleteActor(db *sql.DB, actorID int) error {
	query := `DELETE FROM actors WHERE id = $1`

	_, err := db.Exec(query, actorID)
	if err != nil {
		return err
	}

	return nil
}

func AddMovie(db *sql.DB, movie Movie) error {
	query := `
        INSERT INTO movies (title, description, release_date, rating)
        VALUES ($1, $2, $3, $4)
    `

	_, err := db.Exec(query, movie.Title, movie.Description, movie.ReleaseDate, movie.Rating)
	if err != nil {
		return err
	}

	return nil
}

func UpdateMovie(db *sql.DB, movie Movie) error {
	query := "UPDATE movies SET "
	var args []interface{}

	if movie.Title != "" {
		query += "title = $2, "
		args = append(args, movie.Title)
	}
	if movie.Description != "" {
		query += "description = $3, "
		args = append(args, movie.Description)
	}
	if !movie.ReleaseDate.IsZero() {
		query += "release_date = $4, "
		args = append(args, movie.ReleaseDate)
	}
	if movie.Rating != 0 {
		query += "rating = $5, "
		args = append(args, movie.Rating)
	}

	query = strings.TrimSuffix(query, ", ")

	query += " WHERE id = $1"

	args = append([]interface{}{movie.ID}, args...)

	// Выполнение SQL запроса
	_, err := db.Exec(query, args...)
	if err != nil {
		return err
	}

	return nil
}

func DeleteMovie(db *sql.DB, movieID int) error {
	query := `DELETE FROM movies WHERE id = $1`

	_, err := db.Exec(query, movieID)
	if err != nil {
		return err
	}

	return nil
}

func AddMovieActor(db *sql.DB, movieID, actorID int) error {
	query := `INSERT INTO movie_actors (movie_id, actor_id) VALUES ($1, $2)`

	_, err := db.Exec(query, movieID, actorID)
	if err != nil {
		return err
	}

	return nil
}

func SearchMoviesByActorName(db *sql.DB, actorName string) ([]Movie, error) {
	query := `
        SELECT m.id, m.title, m.description, m.release_date, m.rating
        FROM movies m
        INNER JOIN movie_actors ma ON m.id = ma.movie_id
        INNER JOIN actors a ON ma.actor_id = a.id
        WHERE a.name = $1
    `

	rows, err := db.Query(query, actorName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []Movie
	for rows.Next() {
		var movie Movie
		err := rows.Scan(&movie.ID, &movie.Title, &movie.Description, &movie.ReleaseDate, &movie.Rating)
		if err != nil {
			return nil, err
		}
		movies = append(movies, movie)
	}

	return movies, nil
}

func GetMoviesWithSorting(db *sql.DB, orderBy, sortOrder string) ([]Movie, error) {
	query := fmt.Sprintf("SELECT id, title, description, release_date, rating FROM movies ORDER BY %s %s", orderBy, sortOrder)

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []Movie
	for rows.Next() {
		var movie Movie
		if err := rows.Scan(&movie.ID, &movie.Title, &movie.Description, &movie.ReleaseDate, &movie.Rating); err != nil {
			return nil, err
		}
		movies = append(movies, movie)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return movies, nil
}

func SearchMoviesByTitleOrActorName(db *sql.DB, titleFragment, actorNameFragment string) ([]Movie, error) {
	query := `
		SELECT DISTINCT m.id, m.title, m.description, m.release_date, m.rating
		FROM movies m
		LEFT JOIN movie_actors ma ON m.id = ma.movie_id
		LEFT JOIN actors a ON ma.actor_id = a.id
		WHERE m.title LIKE $1 OR a.name LIKE $2
    `

	rows, err := db.Query(query, "%"+titleFragment+"%", "%"+actorNameFragment+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []Movie
	for rows.Next() {
		var movie Movie
		if err := rows.Scan(&movie.ID, &movie.Title, &movie.Description, &movie.ReleaseDate, &movie.Rating); err != nil {
			return nil, err
		}
		movies = append(movies, movie)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return movies, nil
}

func GetActors(db *sql.DB) ([]Actor, error) {
	query := `
        SELECT id, name, gender, birthdate
        FROM actors
    `

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var actors []Actor
	for rows.Next() {
		var actor Actor
		if err := rows.Scan(&actor.Id, &actor.Name, &actor.Gender, &actor.Birthdate); err != nil {
			return nil, err
		}
		actors = append(actors, actor)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return actors, nil
}

func GetMoviesByActorName(db *sql.DB, actorName string) ([]Movie, error) {
	query := `
        SELECT m.id, m.title, m.description, m.release_date, m.rating
        FROM movies m
        INNER JOIN movie_actors ma ON m.id = ma.movie_id
        INNER JOIN actors a ON ma.actor_id = a.id
        WHERE a.name = $1
    `

	rows, err := db.Query(query, actorName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []Movie
	for rows.Next() {
		var movie Movie
		if err := rows.Scan(&movie.ID, &movie.Title, &movie.Description, &movie.ReleaseDate, &movie.Rating); err != nil {
			return nil, err
		}
		movies = append(movies, movie)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return movies, nil
}
