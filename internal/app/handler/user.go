package handler

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"lr4/internal/app/ds"
	"lr4/internal/app/role"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type LoginRequest struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	ExpiresIn   int64  `json:"expires_in"`
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

// RegisterUser регистрирует нового пользователя
// @Summary Register user
// @Description Register new user account
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body ds.ChangeUserDTO true "User registration data"
// @Success 201 {object} map[string]interface{} "User registered successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /auth/register [post]
func (h *Handler) RegisterUser(ctx *gin.Context) {
	var input ds.ChangeUserDTO
	if err := ctx.BindJSON(&input); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	if input.Login == "" || input.Password == "" {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("логин и пароль обязательны"))
		return
	}

	hashedPassword := generateHashString(input.Password)

	user := &ds.User{
		Login:    input.Login,
		Password: hashedPassword,
		Role:     role.Buyer,
	}

	err := h.Repository.CreateUser(user)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Пользователь успешно зарегистрирован",
	})
}

// LoginUser аутентифицирует пользователя
// @Summary User login
// @Description Authenticate user and return JWT token
// @Tags Auth
// @Accept json
// @Produce json
// @Param login body LoginRequest true "Login credentials"
// @Success 200 {object} LoginResponse "Login successful"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /auth/login [post]
func (h *Handler) LoginUser(ctx *gin.Context) {
	var req LoginRequest
	if err := ctx.BindJSON(&req); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	if req.Login == "" || req.Password == "" {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("логин и пароль обязательны"))
		return
	}

	hashedPassword := generateHashString(req.Password)

	user, err := h.Repository.GetUserByLogin(req.Login)
	if err != nil {
		h.errorHandler(ctx, http.StatusUnauthorized, fmt.Errorf("неверный логин или пароль"))
		return
	}

	if user.Password != hashedPassword {
		h.errorHandler(ctx, http.StatusUnauthorized, fmt.Errorf("неверный логин или пароль"))
		return
	}

	// Используем методы-геттеры для получения JWT параметров
	jwtSecret := h.AuthMiddleware.GetJWTSecret()
	jwtExpiresIn := h.AuthMiddleware.GetJWTExpiresIn()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &ds.JWTClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(jwtExpiresIn).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "lr2-admin",
			Subject:   fmt.Sprintf("%d", user.UserID),
		},
		UserUUID: user.UUID,
		Role:     user.Role,
	})

	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, LoginResponse{
		ExpiresIn:   int64(jwtExpiresIn.Seconds()),
		AccessToken: tokenString,
		TokenType:   "Bearer",
	})
}

// GetProfile получает профиль пользователя
// @Summary Get user profile
// @Description Get current user profile information
// @Tags User
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]interface{} "Success response with user profile"
// @Failure 404 {object} map[string]interface{} "User not found"
// @Router /user/profile [get]
func (h *Handler) GetProfile(ctx *gin.Context) {
	userID := h.getCurrentUserIDFromContext(ctx)

	user, err := h.Repository.GetUserByID(userID)
	if err != nil {
		h.errorHandler(ctx, http.StatusNotFound, err)
		return
	}

	userDTO := ds.UserDTO{
		UserID:      user.UserID,
		Login:       user.Login,
		IsModerator: user.Role == role.Manager || user.Role == role.Admin,
	}

	h.successResponse(ctx, userDTO)
}

// UpdateProfile обновляет профиль пользователя
// @Summary Update user profile
// @Description Update current user profile information
// @Tags User
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param user body ds.ChangeUserDTO true "User update data"
// @Success 200 {object} map[string]interface{} "Profile updated successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /user/profile [put]
func (h *Handler) UpdateProfile(ctx *gin.Context) {
	userID := h.getCurrentUserIDFromContext(ctx)

	var userUpdates ds.ChangeUserDTO
	if err := ctx.BindJSON(&userUpdates); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	if userUpdates.Password != "" {
		userUpdates.Password = generateHashString(userUpdates.Password)
	}

	updatedUser, err := h.Repository.UpdateUser(userID, userUpdates.Login, userUpdates.Password)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	userDTO := ds.UserDTO{
		UserID:      updatedUser.UserID,
		Login:       updatedUser.Login,
		IsModerator: updatedUser.Role == role.Manager || updatedUser.Role == role.Admin,
	}

	h.successResponse(ctx, userDTO)
}

// LogoutUser выходит из системы
// @Summary User logout
// @Description Logout user and invalidate token
// @Tags Auth
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "Logout successful"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /user/logout [post]
func (h *Handler) LogoutUser(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("отсутствует токен авторизации"))
		return
	}

	tokenString := authHeader[len("Bearer "):]

	err := h.RedisClient.WriteJWTToBlacklist(ctx.Request.Context(), tokenString, 24*time.Hour)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Успешный выход из системы",
	})
}

func generateHashString(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}
