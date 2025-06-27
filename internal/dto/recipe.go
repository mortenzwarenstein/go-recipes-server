package dto

import "mime/multipart"

type CreateRecipeRequest struct {
	Name       string                `form:"name" binding:"required" example:"Pasta Carbonara"`
	Cookbook   string                `form:"cookbook" binding:"required" example:"Pasta for dummies"`
	Pagenumber string                `form:"pagenumber" binding:"required" example:"1"`
	Image      *multipart.FileHeader `form:"image"`
}
