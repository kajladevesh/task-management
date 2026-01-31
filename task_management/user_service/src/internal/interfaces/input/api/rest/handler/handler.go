package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"task-management/user-service/src/internal/usecase"
	errors "task-management/user-service/src/pkg/error"
	pkg "task-management/user-service/src/pkg/jsonResponse"
)

type UserHandler struct {
	userService *usecase.UserService
}

func NewUserHandler(userService *usecase.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (u *UserHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var input usecase.RegisterInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		pkg.Error(w, http.StatusBadRequest, "failed to decode request body")
		return
	}

	//  validation empty name and password
	if strings.TrimSpace(input.UserName) == "" || strings.TrimSpace(input.Email) == "" || strings.TrimSpace(input.Password) == "" {
		pkg.Error(w, http.StatusBadRequest, "username, email and password must not be empty")
		return
	}

	if err := u.userService.RegisterUser(ctx, input); err != nil {
		pkg.Error(w, http.StatusConflict, err.Error())
		return
	}

	pkg.Created(w, nil, "User reegistered successfully")

}

//--------------------------------------------------------------------------------------------------------------

func (u *UserHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var input usecase.LoginInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		pkg.Error(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if strings.TrimSpace(input.UserName) == "" || strings.TrimSpace(input.Password) == "" {
		pkg.Error(w, http.StatusBadRequest, "username and password must not be empty")
		return
	}

	loginOutput, err := u.userService.Login(ctx, input)
	if err != nil {
		pkg.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    loginOutput.Token,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   10 * 60 * 60, //10 hr
	})

	pkg.Success(w, loginOutput, "Login successful")
}

//-------------------------------------------------------------------------------------------------------------

func (u *UserHandler) GetUserProfileHandler(w http.ResponseWriter, r *http.Request) {
	// ctx := r.Context()

	// Extract user ID from context (set during auth middleware)
	username := r.URL.Query().Get("username")
	if strings.TrimSpace(username) == "" {
		pkg.Error(w, http.StatusBadRequest, "username query parameter is required")
		return
	}

	userProfile, err := u.userService.GetUserByUsername(r.Context(), username)
	if err != nil {
		pkg.Error(w, http.StatusNotFound, "user not found")
		return
	}

	// Hide sensitive info before sending
	userProfile.Password = ""

	pkg.Success(w, userProfile, "user profile fetched successfully")
}

//--------------------------------------------------------------------------------------------------------------

func (u *UserHandler) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	var input usecase.UpdateUserInput

	// Decode JSON body
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		pkg.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	fmt.Println("userhandler request data : ", input)

	err := u.userService.UpdateUser(r.Context(), input)
	if err != nil {
		switch err {
		case errors.ErrUserNotFound:
			pkg.Error(w, http.StatusNotFound, "User not found")
		default:
			pkg.Error(w, http.StatusInternalServerError, "Failed to update user")
		}
		return
	}

	// Respond with success
	pkg.Success(w, nil, "User updated successfully")
}
