package main

import (
	"errors"
	"fmt"
	"meus_gastos/internal/data"
	"meus_gastos/internal/validator"
	"net/http"
	"time"
)

func (app *application) listTransactionsByCategoryIDHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil || id < 1 {
		app.notFoundResponse(w, r)
		return
	}

	var input struct {
		Name       string
		CategoryID int
		StartDate  time.Time
		EndDate    time.Time
		data.Filters
	}

	v := validator.New()

	qs := r.URL.Query()
	input.Name = app.readString(qs, "description", "")
	input.StartDate = *app.readDate(qs, "start", "2006-01-02")
	input.EndDate = *app.readDate(qs, "end", "2006-01-02")
	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
	input.Filters.Sort = app.readString(qs, "sort", "id")
	input.Filters.SortSafelist = []string{"id", "description", "-id", "-description"}

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	user := app.contextGetUser(r)
	transactions, metadata, err := app.models.Transactions.GetAllByUserAndCategory(input.Name, user.ID, id, input.StartDate, input.EndDate, input.Filters)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	transactionsDTO := []*data.TransactionDTO{}
	for _, t := range transactions {
		err = prepareTransactionForResponse(app, t, user)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
		transactionsDTO = append(transactionsDTO, t.ToDTO())
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"transactions": transactionsDTO, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listTransactionsHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name string
		data.Filters
		CategoryType data.TypeCategoria
		StartDate    time.Time
		EndDate      time.Time
	}

	v := validator.New()

	qs := r.URL.Query()
	categoryStr := app.readString(qs, "category_type", "")
	if categoryStr != "" {
		input.CategoryType = data.TypeCategoriaFromString(categoryStr)
	}
	input.StartDate = *app.readDate(qs, "start", "2006-01-02")
	input.EndDate = *app.readDate(qs, "end", "2006-01-02")
	input.Name = app.readString(qs, "description", "")
	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
	input.Filters.Sort = app.readString(qs, "sort", "id")
	input.Filters.SortSafelist = []string{"id", "description", "-id", "-description"}

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	user := app.contextGetUser(r)
	transactions, metadata, err := app.models.Transactions.GetAllByUser(
		input.Name, user.ID,
		input.StartDate,
		input.EndDate,
		input.CategoryType,
		input.Filters,
	)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	transactionsDTO := []*data.TransactionDTO{}
	for _, t := range transactions {
		err = prepareTransactionForResponse(app, t, user)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
		transactionsDTO = append(transactionsDTO, t.ToDTO())
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"transactions": transactionsDTO, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) createTransactionHandler(w http.ResponseWriter, r *http.Request) {
	var dto data.TransactionDTO
	err := app.readJSON(w, r, &dto)

	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := app.contextGetUser(r)
	dto.User = user.ToDTO()
	transaction := dto.ToModel()

	v := validator.New()

	if data.ValidateTransaction(v, transaction); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	if transaction.Category == nil {
		app.errorResponse(w, r, http.StatusUnprocessableEntity, "must category be provided")
		return
	}

	err = app.models.Transactions.Insert(transaction)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = prepareTransactionForResponse(app, transaction, user)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/transactions/%d", transaction.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"transaction": transaction.ToDTO()}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) showTransactionHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil || id < 1 {
		app.notFoundResponse(w, r)
		return
	}

	user := app.contextGetUser(r)
	transaction, err := app.models.Transactions.GetByID(id, user.ID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = prepareTransactionForResponse(app, transaction, user)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"transaction": transaction.ToDTO()}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateTransactionHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil || id < 1 {
		app.notFoundResponse(w, r)
		return
	}

	var dto data.TransactionDTO
	err = app.readJSON(w, r, &dto)

	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := app.contextGetUser(r)

	dto.User = user.ToDTO()
	v := validator.New()
	if data.ValidateTransaction(v, dto.ToModel()); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	transaction, err := app.models.Transactions.GetByID(id, user.ID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	dto.ToDTOUpdateTransaction(transaction)

	err = app.models.Transactions.Update(transaction, user.ID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = prepareTransactionForResponse(app, transaction, user)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"transaction": transaction.ToDTO()}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteTransactionHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil || id < 1 {
		app.notFoundResponse(w, r)
		return
	}

	user := app.contextGetUser(r)
	err = app.models.Transactions.Delete(id, user.ID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "transaction successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func prepareTransactionForResponse(app *application, transaction *data.Transaction, user *data.User) error {
	transaction.User = user

	category, err := app.models.Categories.GetByID(transaction.Category.ID, user.ID)

	if err != nil {
		return err
	}

	if category != nil {
		transaction.Category = category
	}

	return nil
}
