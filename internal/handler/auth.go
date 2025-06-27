package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go-recipes-server/internal/dto"
	"go-recipes-server/internal/middleware"
	"go-recipes-server/internal/model"
	"go-recipes-server/internal/util"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type AuthHandler struct {
	DB *gorm.DB
}

func NewAuthHandler(db *gorm.DB, r *gin.RouterGroup) *AuthHandler {
	h := &AuthHandler{DB: db}
	h.RegisterRoutes(r)
	return h
}

func (h *AuthHandler) RegisterRoutes(r *gin.RouterGroup) {
	{
		r.GET("/me", middleware.JWTMiddleware(), h.GetProfile)
		r.POST("/login", h.Login)
		r.POST("/register", h.Register)
		r.POST("/refresh-token", h.RefreshToken)
		r.POST("/logout", middleware.JWTMiddleware(), h.Logout)
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var input dto.LoginInput

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Error:   err.Error(),
			Message: "invalid request",
		})
		return
	}

	accessToken, refreshTokenRaw, err := h.Authenticate(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Error:   err.Error(),
			Message: "invalid credentials",
		})
		return
	}

	SetAuthCookies(c, accessToken, refreshTokenRaw)

	c.JSON(http.StatusOK, dto.Response{
		Message: "success",
	})
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	refreshCookie, err := c.Request.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Error:   "Invalid cookie",
			Message: "Invalid cookie, please login first",
		})
		return
	}

	refreshToken := refreshCookie.Value

	claims, err := util.VerifyRefreshToken(refreshToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Error:   err.Error(),
			Message: "invalid refresh token",
		})
		return
	}

	accessToken, refreshToken, err := util.GenerateTokens(claims.Subject, claims.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token generation failed"})
		return
	}

	SetAuthCookies(c, accessToken, refreshToken)

	c.JSON(http.StatusOK, dto.Response{
		Message: "success",
	})

}

func (h *AuthHandler) Register(c *gin.Context) {
	var input dto.RegisterInput

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Error:   err.Error(),
			Message: "invalid request",
		})
		return
	}

	u := model.User{
		Email:    input.Email,
		Password: input.Password,
	}

	if err := h.DB.Create(&u).Error; err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Error:   err.Error(),
			Message: "error creating user",
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Data:    u,
		Message: "user created",
	})
}

func (h *AuthHandler) Authenticate(input dto.LoginInput) (accessToken string, refreshToken string, err error) {
	u := model.User{}

	err = h.DB.Where("email = ?", input.Email).Take(&u).Error

	if err != nil {
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(input.Password))

	if err != nil && errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return
	}

	accessToken, refreshToken, err = util.GenerateTokens(u.ID, u.Email)

	if err != nil {
		return
	}

	return
}

func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID, ok := c.Get("userID")
	u := model.User{}
	if !ok {
		c.JSON(http.StatusUnauthorized, dto.Response{
			Error:   "unauthorized",
			Message: "unauthorized to get profile, please login first",
		})
		return
	}
	h.DB.Where("id = ?", userID).First(&u)
	c.JSON(http.StatusOK, dto.Response{
		Data:    u,
		Message: "success",
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	SetCookie(c, "access_token", "", -1)
	SetCookie(c, "refresh_token", "", -1)
}

func SetAuthCookies(c *gin.Context, accessToken string, refreshToken string) {
	accessCookieMaxAge := util.AccessExpiryTime + time.Minute
	refreshCookieMaxAge := util.RefreshExpiryTime + time.Minute
	SetCookie(c, "access_token", accessToken, int(accessCookieMaxAge.Seconds()))
	SetCookie(c, "refresh_token", refreshToken, int(refreshCookieMaxAge.Seconds()))
}

func SetCookie(c *gin.Context, key string, token string, maxAge int) {
	isProd := gin.Mode() == gin.ReleaseMode
	c.SetCookie(key, token, maxAge, "/", "", isProd, true)
}
