package main

import (
	"errors"
	"fmt"
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

	user := app.contextGetUser(r)
	categories, metadata, err := app.models.Categories.GetAll(input.Name, user.ID, input.Filters)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	categoriesDTO := []*data.CategoryDTO{}
	for _, c := range categories {
		c.User = user
		categoriesDTO = append(categoriesDTO, c.ToDTO())
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"categories": categoriesDTO, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) createCategoryHandler(w http.ResponseWriter, r *http.Request) {
	var dto data.CategoryDTO
	err := app.readJSON(w, r, &dto)

	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := app.contextGetUser(r)
	dto.User = user.ToDTO()
	category := dto.ToModel()

	v := validator.New()

	if data.ValidateCategory(v, category); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Categories.Insert(category)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/categories/%d", category.ID))

	category.User = user

	err = app.writeJSON(w, http.StatusCreated, envelope{"category": category.ToDTO()}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) showCategoryHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil || id < 1 {
		app.notFoundResponse(w, r)
		return
	}
	user := app.contextGetUser(r)
	category, err := app.models.Categories.GetByID(id, user.ID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	category.User = user

	err = app.writeJSON(w, http.StatusOK, envelope{"category": category.ToDTO()}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateCategoryHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil || id < 1 {
		app.notFoundResponse(w, r)
		return
	}

	var dto data.CategoryDTO
	err = app.readJSON(w, r, &dto)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := app.contextGetUser(r)
	v := validator.New()
	dto.User = user.ToDTO()
	if data.ValidateCategory(v, dto.ToModel()); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	category, err := app.models.Categories.GetByID(id, user.ID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	category = dto.ToDTOUpdateCategory(category)

	err = app.models.Categories.Update(category, user.ID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	category.User = user

	err = app.writeJSON(w, http.StatusOK, envelope{"category": category.ToDTO()}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteCategoryHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil || id < 1 {
		app.notFoundResponse(w, r)
		return
	}

	user := app.contextGetUser(r)
	err = app.models.Categories.Delete(id, user.ID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "category successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
