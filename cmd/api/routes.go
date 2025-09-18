package main

import (
	"expvar"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)

	router.HandlerFunc(http.MethodGet, "/v1/categories", app.requireActivatedUser(app.listCategoriesHandler))
	router.HandlerFunc(http.MethodPost, "/v1/categories", app.requireActivatedUser(app.createCategoryHandler))
	router.HandlerFunc(http.MethodGet, "/v1/categories/:id", app.requireActivatedUser(app.showCategoryHandler))
	router.HandlerFunc(http.MethodPut, "/v1/categories/:id", app.requireActivatedUser(app.updateCategoryHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/categories/:id", app.requireActivatedUser(app.deleteCategoryHandler))

	router.HandlerFunc(http.MethodGet, "/v1/transactions", app.requireActivatedUser(app.listTransactionsHandler))
	router.HandlerFunc(http.MethodPost, "/v1/transactions", app.requireActivatedUser(app.createTransactionHandler))
	router.HandlerFunc(http.MethodGet, "/v1/transactions/find/:id", app.requireActivatedUser(app.showTransactionHandler))
	router.HandlerFunc(http.MethodPut, "/v1/transactions/update/:id", app.requireActivatedUser(app.updateTransactionHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/transactions/delete/:id", app.requireActivatedUser(app.deleteTransactionHandler))
	router.HandlerFunc(http.MethodGet, "/v1/transactions/category/:id", app.requireActivatedUser(app.listTransactionsByCategoryIDHandler))

	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)

	router.Handler(http.MethodGet, "/debug/vars", expvar.Handler())

	return app.metrics(app.recoverPanic(app.enableCORS(app.rateLimit(app.authenticate(router)))))
}
