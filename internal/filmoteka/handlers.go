package filmoteka

import (
	"TestVK/internal/db"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type MovieRequest struct {
	Id             int     `json:"id"`
	Title          string  `json:"title"`
	Description    string  `json:"description"`
	ReleaseDateStr string  `json:"release_date"`
	Rating         float64 `json:"rating"`
}

func (f *Filmoteka) handleAddActor(w http.ResponseWriter, r *http.Request) {
	var actor db.Actor
	err := json.NewDecoder(r.Body).Decode(&actor)
	if err != nil {
		f.Logger.Info("Response", slog.String("Body", err.Error()))
		http.Error(w, "Невозможно прочитать тело запроса", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if actor.Name == "" {
		http.Error(w, "Имя актера обязательно для заполнения", http.StatusBadRequest)
		return
	}

	if err := db.AddActor(f.Db, actor); err != nil {
		f.Logger.Warn("Error creating actor", http.StatusInternalServerError, err)
		http.Error(w, "Ошибка при добавлении актера в базу данных", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Актер успешно добавлен в базу данных"))
	f.Logger.Info("New Actor", actor.Name)
}

func (f *Filmoteka) handleUpdateActor(w http.ResponseWriter, r *http.Request) {
	var actor db.Actor
	err := json.NewDecoder(r.Body).Decode(&actor)
	if err != nil {
		f.Logger.Info("Response", slog.String("Body", err.Error()))
		http.Error(w, "Невозможно прочитать тело запроса", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err := db.UpdateActor(f.Db, actor); err != nil {
		f.Logger.Warn("Error updating actor", http.StatusInternalServerError, err)
		http.Error(w, "Ошибка при обновлении актера", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Актер успешно обновлен в базу данных"))
	f.Logger.Info("New Actor", actor.Id, actor.Name)
}

func (f *Filmoteka) handleDeleteActor(w http.ResponseWriter, r *http.Request) {
	idParam := r.URL.Query().Get("id")
	actorID, err := strconv.Atoi(idParam)
	if err != nil {
		f.Logger.Info("Can't get actor id", http.StatusBadRequest, err)
		http.Error(w, "Неверный идентификатор актера", http.StatusBadRequest)
		return
	}

	if err := db.DeleteActor(f.Db, actorID); err != nil {
		f.Logger.Warn("Error deleting actor", http.StatusInternalServerError, err)
		http.Error(w, "Ошибка при удалении актера", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Актер успешно удален из базы данных"))
	f.Logger.Info("Deleted actor", actorID)
}

func (f *Filmoteka) handleAddMovie(w http.ResponseWriter, r *http.Request) {
	var movieReq MovieRequest
	err := json.NewDecoder(r.Body).Decode(&movieReq)
	if err != nil {
		f.Logger.Info("Ошибка при декодировании запроса", err.Error(), http.StatusBadRequest)
		http.Error(w, "Ошибка при декодировании запроса", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if movieReq.ReleaseDateStr == "" {
		f.Logger.Info("Дата выхода фильма обязательна для заполнения", http.StatusBadRequest)
		http.Error(w, "Дата выхода фильма обязательна для заполнения", http.StatusBadRequest)
		return
	}

	releaseDate, err := parseDate(movieReq.ReleaseDateStr)
	if err != nil {
		f.Logger.Info("Wrong data format", err.Error(), http.StatusBadRequest)
		http.Error(w, "Неверный формат даты выхода фильма", http.StatusBadRequest)
		return
	}

	movie := db.Movie{
		Title:       movieReq.Title,
		Description: movieReq.Description,
		ReleaseDate: releaseDate,
		Rating:      movieReq.Rating,
	}
	defer r.Body.Close()

	if movie.Title == "" || movie.ReleaseDate.IsZero() {
		f.Logger.Info("Response", slog.String("Body", err.Error()))
		http.Error(w, "Название и дата выхода фильма обязательны для заполнения", http.StatusBadRequest)
		return
	}

	if err := db.AddMovie(f.Db, movie); err != nil {
		f.Logger.Warn("Error creating movie", http.StatusInternalServerError, err)
		http.Error(w, "Ошибка при добавлении фильма в базу данных", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Фильм успешно добавлен в базу данных"))
	f.Logger.Info("New Movie", movie.Title)
}

func (f *Filmoteka) handleUpdateMovie(w http.ResponseWriter, r *http.Request) {
	var movieReq MovieRequest
	err := json.NewDecoder(r.Body).Decode(&movieReq)
	if err != nil {
		f.Logger.Info("Response", slog.String("Body", err.Error()))
		http.Error(w, "Ошибка при декодировании запроса", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if movieReq.ReleaseDateStr == "" {
		f.Logger.Info("Дата выхода фильма обязательна для заполнения", http.StatusBadRequest)
		http.Error(w, "Дата выхода фильма обязательна для заполнения", http.StatusBadRequest)
		return
	}

	releaseDate, err := parseDate(movieReq.ReleaseDateStr)
	if err != nil {
		f.Logger.Info("Response", slog.String("Body", err.Error()))
		http.Error(w, "Неверный формат даты выхода фильма", http.StatusBadRequest)
		return
	}

	movie := db.Movie{
		Title:       movieReq.Title,
		Description: movieReq.Description,
		ReleaseDate: releaseDate,
		Rating:      movieReq.Rating,
	}

	if movie.ID == 0 {
		f.Logger.Info("Response", slog.String("Body", err.Error()))
		http.Error(w, "Идентификатор фильма обязателен для обновления", http.StatusBadRequest)
		return
	}

	if err := db.UpdateMovie(f.Db, movie); err != nil {
		f.Logger.Warn("Error updating movie", http.StatusInternalServerError, err)
		http.Error(w, "Ошибка при обновлении информации о фильме", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Информация о фильме успешно обновлена"))
	f.Logger.Info("Movie update", movie.ID, movie.Title)
}

func (f *Filmoteka) handleDeleteMovie(w http.ResponseWriter, r *http.Request) {
	idParam := r.URL.Query().Get("id")
	movieID, err := strconv.Atoi(idParam)
	if err != nil {
		f.Logger.Info("Response", slog.String("Body", err.Error()))
		http.Error(w, "Неверный идентификатор фильма", http.StatusBadRequest)
		return
	}

	if movieID == 0 {
		f.Logger.Info("Response", slog.String("Body", err.Error()))
		http.Error(w, "Идентификатор фильма обязателен для обновления", http.StatusBadRequest)
		return
	}

	if err := db.DeleteMovie(f.Db, movieID); err != nil {
		f.Logger.Warn("Error deleting movie", http.StatusInternalServerError, err)
		http.Error(w, "Ошибка при удалении фильма", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Фильм успешно удален из базы данных"))
	f.Logger.Info("Movie deleted", movieID)
}

func (f *Filmoteka) handleUpdateMovieActors(w http.ResponseWriter, r *http.Request) {
	movieIDParam := r.URL.Query().Get("movie_id")
	movieID, err := strconv.Atoi(movieIDParam)
	if err != nil {
		f.Logger.Info("Response", slog.String("Body", err.Error()))
		http.Error(w, "Неверный идентификатор фильма", http.StatusBadRequest)
		return
	}

	var actorID int
	err = json.NewDecoder(r.Body).Decode(&actorID)
	if err != nil {
		f.Logger.Info("Response", slog.String("Body", err.Error()))
		http.Error(w, "Невозможно прочитать тело запроса", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err := db.AddMovieActor(f.Db, movieID, actorID); err != nil {
		f.Logger.Warn("Error adding actor to movie", http.StatusInternalServerError, err)
		http.Error(w, "Ошибка при обновлении списка актеров для фильма", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Список актеров для фильма успешно обновлен"))
	f.Logger.Info("Movie actors update", movieID, actorID)
}

func (f *Filmoteka) handleSearchMoviesByActorName(w http.ResponseWriter, r *http.Request) {
	actorName := r.URL.Query().Get("actor_name")
	if actorName == "" {
		f.Logger.Info("Response", slog.String("Body", actorName))
		http.Error(w, "Не указано имя актера", http.StatusBadRequest)
		return
	}

	movies, err := db.SearchMoviesByActorName(f.Db, actorName)
	if err != nil {
		f.Logger.Warn("Error searching movie", http.StatusInternalServerError, err)
		http.Error(w, "Ошибка при поиске фильмов по имени актера", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)
	f.Logger.Info("Actors movies", movies)
}

func (f *Filmoteka) handleGetMovies(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()
	sortBy := queryValues.Get("sort_by")
	sortOrder := queryValues.Get("sort_order")

	var orderBy string
	switch sortBy {
	case "title":
		orderBy = "title"
	case "rating":
		orderBy = "rating"
	case "release_date":
		orderBy = "release_date"
	default:
		orderBy = "rating"
	}

	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "desc"
	}

	movies, err := db.GetMoviesWithSorting(f.Db, orderBy, sortOrder)
	if err != nil {
		f.Logger.Warn("Error searching movies", http.StatusInternalServerError, err)
		http.Error(w, "Ошибка при получении списка фильмов", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)
	f.Logger.Info("Movies", movies)
}

func (f *Filmoteka) handleSearchMoviesByTitleOrActor(w http.ResponseWriter, r *http.Request) {
	titleFragment := r.URL.Query().Get("title_fragment")
	actorNameFragment := r.URL.Query().Get("actor_name_fragment")

	if titleFragment == "" && actorNameFragment == "" {
		f.Logger.Info("Response", slog.String("Body", titleFragment))
		http.Error(w, "Не указан фрагмент названия или фрагмент имени актёра", http.StatusBadRequest)
		return
	}

	movies, err := db.SearchMoviesByTitleOrActorName(f.Db, titleFragment, actorNameFragment)
	if err != nil {
		f.Logger.Warn("Error searching movie by title or actor", http.StatusInternalServerError, err)
		http.Error(w, "Ошибка при поиске фильмов", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)
	f.Logger.Info("Movies by title or actor", movies)
}

func (f *Filmoteka) handleGetActors(w http.ResponseWriter, r *http.Request) {
	actors, err := db.GetActors(f.Db)
	if err != nil {
		f.Logger.Warn("Error getting actors", http.StatusInternalServerError, err)
		http.Error(w, "Ошибка при получении списка актёров", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(actors)
	f.Logger.Info("actors", actors)
}

func (f *Filmoteka) handleGetActorMovies(w http.ResponseWriter, r *http.Request) {
	actorName := r.URL.Query().Get("actor_name")
	if actorName == "" {
		f.Logger.Info("Response", slog.String("Body", actorName))
		http.Error(w, "Не указано имя актёра", http.StatusBadRequest)
		return
	}

	movies, err := db.GetMoviesByActorName(f.Db, actorName)
	if err != nil {
		f.Logger.Warn("Error searching movie by actor", http.StatusInternalServerError, err)
		http.Error(w, "Ошибка при получении списка фильмов по имени актёра", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)
	f.Logger.Info("movies by actor", actorName, movies)
}

func parseDate(dateString string) (time.Time, error) {
	dateParts := strings.Split(dateString, ".")
	if len(dateParts) != 3 {
		return time.Time{}, errors.New("неверный формат даты")
	}

	year, err := strconv.Atoi(dateParts[0])
	if err != nil {
		return time.Time{}, err
	}
	month, err := strconv.Atoi(dateParts[1])
	if err != nil {
		return time.Time{}, err
	}
	day, err := strconv.Atoi(dateParts[2])
	if err != nil {
		return time.Time{}, err
	}

	date := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)

	return date, nil
}
