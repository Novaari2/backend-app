package users

import (
	"auth-app/internal/api/resp"
	"auth-app/internal/utils"
	"encoding/json"
	"net/http"
)

type httpHandler struct {
	svc Service
}

func NewHTTPHandler(svc Service) *httpHandler {
	return &httpHandler{svc: svc}
}

func (h *httpHandler) RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req RegisterUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response := utils.ApiResponse(http.StatusBadRequest, "Invalid request body", nil)
		resp.WriteJSON(w, http.StatusBadRequest, response)
		return
	}

	if req.Nik == "" || req.Role == "" {
		response := utils.ApiResponse(http.StatusBadRequest, "Missing required fields", nil)
		resp.WriteJSON(w, http.StatusBadRequest, response)
		return
	}

	if len(req.Nik) != 16 {
		response := utils.ApiResponse(http.StatusBadRequest, "NIK must be 16 characters", nil)
		resp.WriteJSON(w, http.StatusBadRequest, response)
		return
	}

	user, err := h.svc.RegisterUser(ctx, &req)
	if err != nil {
		response := utils.ApiResponse(http.StatusInternalServerError, "Failed to register user", nil)
		resp.WriteJSON(w, http.StatusInternalServerError, response)
		return
	}

	response := utils.ApiResponse(http.StatusOK, "User registered successfully", user)
	resp.WriteJSON(w, http.StatusOK, response)
}

func (h *httpHandler) LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req LoginUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response := utils.ApiResponse(http.StatusBadRequest, "Invalid request body", nil)
		resp.WriteJSON(w, http.StatusBadRequest, response)
		return
	}

	if req.Nik == "" || req.Password == "" {
		response := utils.ApiResponse(http.StatusBadRequest, "Missing required fields", nil)
		resp.WriteJSON(w, http.StatusBadRequest, response)
		return
	}

	if len(req.Nik) != 16 {
		response := utils.ApiResponse(http.StatusBadRequest, "NIK must be 16 characters", nil)
		resp.WriteJSON(w, http.StatusBadRequest, response)
		return
	}

	user, err := h.svc.LoginUser(ctx, &req)
	if err != nil {
		response := utils.ApiResponse(http.StatusBadRequest, "Failed to login user, NIK or Password is incorrect", nil)
		resp.WriteJSON(w, http.StatusBadRequest, response)
		return
	}

	response := utils.ApiResponse(http.StatusOK, "User logged in successfully", user)
	resp.WriteJSON(w, http.StatusOK, response)
}

func (h *httpHandler) ValidateTokenHandler(w http.ResponseWriter, r *http.Request) {
	tokenString, ok := r.Context().Value("token").(string)
	if !ok {
		http.Error(w, "Token not found", http.StatusUnauthorized)
		return
	}

	// PrettyJSON(tokenString)
	claims, err := h.svc.ValidateToken(r.Context(), tokenString)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	response := utils.ApiResponse(http.StatusOK, "Token is valid", claims)
	resp.WriteJSON(w, http.StatusOK, response)

}
