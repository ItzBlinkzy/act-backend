package controller

import (
	"log"
	"net/http"

	"github.com/itzblinkzy/act-backend/model"
	"github.com/itzblinkzy/act-backend/repository"
	"github.com/labstack/echo/v4"
)

func ListReviews(c echo.Context) error {
	reviews, err := repository.GetAllReviews()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, reviews)
}

func ListReviewsById(c echo.Context) error {
	userID := c.Param("userId")
	reviews, err := repository.GetReviewsByUserID(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, reviews)
}

func PostReview(c echo.Context) error {

	var review model.Review
	if err := c.Bind(&review); err != nil {
		log.Printf("Error binding review data: %v", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid input provided"})
	}

	if review.Description == "" {
		log.Println("Error: Description cannot be empty")
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Description cannot be empty"})
	}

	if review.Stars < 1 {
		log.Println("Error: Stars must be at least 1")
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Stars must be at least 1"})
	}

	log.Println("Attempting to create review in database")
	err := repository.CreateReview(&review)
	if err != nil {
		log.Printf("Error creating review: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to create review, please try again"})
	}

	log.Println("Review created successfully")
	return c.JSON(http.StatusCreated, review)
}

func DeleteReview(c echo.Context) error {
	reviewID := c.Param("reviewId")
	deleted, err := repository.DeleteReview(reviewID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	if !deleted {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "No review found with the provided ID or it has already been deleted"})
	}
	return c.JSON(http.StatusOK, "Review deleted successfully")
}
