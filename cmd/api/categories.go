package main

import (
	"meus_gastos/internal/data"
	"meus_gastos/internal/validator"
	"net/http"
)

func (app *application) listCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name string
		data.Filters
	}

	v := validator.New()

	qs := r.URL.Query()
	input.Name = app.readString(qs, "name", "")
	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
	input.Filters.Sort = app.readString(qs, "sort", "id")
	input.Filters.SortSafelist = []string{"id", "name", "-id", "-name"}

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	userID := app.contextGetUser(r).ID
	categories, metadata, err := app.models.Categories.GetAll(input.Name, userID, input.Filters)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	user, err := app.models.Users.GetByID(userID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	for _, c := range categories {
		c.User = user
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"categories": categories, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
