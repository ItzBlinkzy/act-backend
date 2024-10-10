package v1

import (
	"github.com/itzblinkzy/act-backend/controller"
	"github.com/labstack/echo/v4"
)

func ReviewGroup(g *echo.Group) {
	g.GET("/list-reviews", controller.ListReviews)
	g.GET("/list-reviews/:userId", controller.ListReviewsById)
	g.POST("/post-review", controller.PostReview)
	g.DELETE("/delete-review/:reviewId", controller.DeleteReview)
}
