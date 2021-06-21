package main
import (
	"net/http"
	"github.com/julienschmidt/httprouter"
)
func (app *application) routes() http.Handler {

	router := httprouter.New()

	router.NotFound=http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed=http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodGet, "/v1/animes", app.listAnimesHandler)
	router.HandlerFunc(http.MethodPost, "/v1/animes", app.createAnimeHandler)
	router.HandlerFunc(http.MethodGet, "/v1/animes/:id", app.showAnimeHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/animes/:id", app.updateAnimeHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/animes/:id", app.deleteAnimeHandler)

	router.HandlerFunc(http.MethodPost,"/v1/users",app.registerUserHandler)
	return app.recoverPanic(app.rateLimit(router))

}
