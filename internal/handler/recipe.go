package handler

import (
	"github.com/gin-gonic/gin"
	"go-recipes-server/internal/dto"
	"go-recipes-server/internal/model"
	"gorm.io/gorm"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
)

type RecipeHandler struct {
	DB *gorm.DB
}

func NewRecipeHandler(db *gorm.DB, r *gin.RouterGroup) *RecipeHandler {
	h := &RecipeHandler{DB: db}
	h.RegisterRoutes(r)
	return h
}

func (h *RecipeHandler) RegisterRoutes(r *gin.RouterGroup) {
	{
		r.GET("/recipes", h.GetAll)
		r.POST("/recipes", h.Create)
	}
}

func (h *RecipeHandler) GetAll(c *gin.Context) {
	var recipes []model.Recipe
	var userID, ok = c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, dto.Response{
			Error:   "unauthorized",
			Message: "unauthorized to get recipes, please login first",
		})
		return
	}
	h.DB.Where("user_id = ?", userID).Find(&recipes)
	c.JSON(http.StatusOK, dto.Response{Data: recipes, Error: "", Message: "success"})
}

func (h *RecipeHandler) Create(c *gin.Context) {
	var req dto.CreateRecipeRequest
	var userID, ok = c.Get("userID")

	if !ok {
		c.JSON(http.StatusUnauthorized, dto.Response{
			Error:   "unauthorized",
			Message: "unauthorized to create recipe, please login first",
		})
		return
	}

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(
			http.StatusBadRequest,
			dto.Response{
				Error:   err.Error(),
				Message: "invalid request",
			},
		)
		return
	}

	var existingRecipe model.Recipe
	if err := h.DB.Where("name = ?", req.Name).First(&existingRecipe).Error; err == nil {
		c.JSON(
			http.StatusConflict,
			dto.Response{
				Error:   "recipe already exists",
				Message: "recipe already exists",
			},
		)
		return
	}
	pagenumber, err := strconv.Atoi(req.Pagenumber)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			dto.Response{
				Error:   err.Error(),
				Message: "invalid request",
			},
		)
	}

	filename := strings.ReplaceAll(strings.ToLower(req.Name), " ", "-")
	recipe := model.Recipe{
		Name:       req.Name,
		Cookbook:   req.Cookbook,
		Pagenumber: pagenumber,
		UserID:     userID.(string),
	}

	if filename != "" {
		recipe.Image = SaveImage(c, filename)
	}

	if err := h.DB.Create(&recipe).Error; err != nil {
		c.JSON(
			http.StatusInternalServerError,
			dto.Response{
				Error:   err.Error(),
				Message: "internal server error",
			},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		dto.Response{
			Data:    recipe,
			Error:   "",
			Message: "success",
		},
	)
}

func SaveImage(c *gin.Context, filename string) string {
	file, _ := c.FormFile("image")
	if file == nil {
		return ""
	}
	path := "public/images/" + filename + filepath.Ext(file.Filename)
	err := c.SaveUploadedFile(file, path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return ""
	}
	return path
}
