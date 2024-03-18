package filmoteka

import (
	"fmt"
	"net/http"
)

func (f *Filmoteka) Api() {
	authMiddleware := f.AuthMiddleware

	http.Handle("/actors/add", authMiddleware(http.HandlerFunc(f.handleAddActor)))
	http.Handle("/actors/update", authMiddleware(http.HandlerFunc(f.handleUpdateActor)))
	http.Handle("/actors/delete", authMiddleware(http.HandlerFunc(f.handleDeleteActor)))
	http.Handle("/movies/add", authMiddleware(http.HandlerFunc(f.handleAddMovie)))
	http.Handle("/movies/update", authMiddleware(http.HandlerFunc(f.handleUpdateMovie)))
	http.Handle("/movies/update_actors", authMiddleware(http.HandlerFunc(f.handleUpdateMovieActors)))
	http.Handle("/movies/delete", authMiddleware(http.HandlerFunc(f.handleDeleteMovie)))
	http.Handle("/movies", authMiddleware(http.HandlerFunc(f.handleGetMovies)))
	http.Handle("/movies/search", authMiddleware(http.HandlerFunc(f.handleSearchMoviesByTitleOrActor)))
	http.Handle("/movies/search_by_actor", authMiddleware(http.HandlerFunc(f.handleSearchMoviesByActorName)))
	http.Handle("/actors", authMiddleware(http.HandlerFunc(f.handleGetActors)))
	http.Handle("/actors/movies", authMiddleware(http.HandlerFunc(f.handleGetActorMovies)))

	f.Logger.Info("Server start")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		f.Logger.Error(fmt.Sprintf("Server error: %v", err))
	}
}
