package repository

import (
	"time"

	"github.com/itzblinkzy/act-backend/database"
	"github.com/itzblinkzy/act-backend/model"
)

func GetAllReviews() ([]model.Review, error) {
	var reviews []model.Review
	query := `SELECT * FROM reviews WHERE deleted_at IS NULL`
	err := database.GetDB().Select(&reviews, query)
	return reviews, err
}

func GetReviewsByUserID(userID string) ([]model.Review, error) {
	var reviews []model.Review
	query := `SELECT * FROM reviews WHERE user_id = $1 AND deleted_at IS NULL`
	err := database.GetDB().Select(&reviews, query, userID)
	return reviews, err
}

func CreateReview(review *model.Review) error {
	query := `INSERT INTO reviews (description, stars, user_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)`
	now := time.Now()
	_, err := database.GetDB().Exec(query, review.Description, review.Stars, review.UserID, now, now)
	return err
}

func DeleteReview(reviewID string) (bool, error) {
	query := `UPDATE reviews SET deleted_at = $1 WHERE id = $2 AND deleted_at IS NULL`
	now := time.Now()
	result, err := database.GetDB().Exec(query, now, reviewID)
	if err != nil {
		return false, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	return rowsAffected > 0, nil
}
