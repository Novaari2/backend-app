package users_test

import (
	"auth-app/internal/users"
	"auth-app/internal/users/mocks"
	"auth-app/internal/utils"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegisterUserHandler_Success(t *testing.T) {
	mockService := new(mocks.MockService)

	reqBody := users.RegisterUserRequest{
		Nik:  "1234567890123456",
		Role: "admin",
	}

	respBody := &users.RegisterUserResponse{
		Nik:      reqBody.Nik,
		Role:     reqBody.Role,
		Password: "mockpass",
	}

	mockService.On("RegisterUser", context.Background(), &reqBody).Return(respBody, nil)

	handler := users.NewHTTPHandler(mockService)

	reqBytes, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(reqBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.RegisterUserHandler(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	expectedResponse := utils.ApiResponse(
		http.StatusOK,
		"User registered successfully",
		map[string]interface{}{
			"nik":      reqBody.Nik,
			"role":     reqBody.Role,
			"password": "mockpass",
		})

	expectedResponseJSON, _ := json.Marshal(expectedResponse)

	assert.JSONEq(t, string(expectedResponseJSON), w.Body.String())

	mockService.AssertExpectations(t)
}
