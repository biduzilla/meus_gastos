package main

import (
	"errors"
	"meus_gastos/internal/data"
	"meus_gastos/internal/validator"
	"net/http"
)

func (app *application) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Cod   int    `json:"cod"`
		Email string `json:"email"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()

	if data.ValidateEmail(v, input.Email); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	user, err := app.models.Users.GetByCodAndEmail(input.Cod, input.Email)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			v.AddError("code", "invalid validation code or email")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	user.Activated = true
	user.Cod = 0

	err = app.models.Users.Update(user)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var userDTO data.UserSaveDTO
	err := app.readJSON(w, r, &userDTO)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user, err := userDTO.ToModel()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	codActivation := app.generateRandomCod()
	user.Cod = codActivation

	v := validator.New()

	if data.ValidateUser(v, user); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Users.Insert(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateEmail):
			v.AddError("email", "a user with this email address already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}

		return
	}

	/*
		MANDA COD VERIFICAÇÃO EMAIL
			app.background(func() {
				data := map[string]interface{}{
					"activationToken": codActivation,
					"userID":          user.ID,
				}

				err = app.mailer.Send(user.Email, "user_welcome.tmpl", data)
				if err != nil {
					app.logger.PrintError(err, nil)
				}
			})
	*/

	err = app.writeJSON(w, http.StatusCreated, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
