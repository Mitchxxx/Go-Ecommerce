package server

import (
	"github/Mitchxxx/Go-Ecommerce/internal/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (s *Server) createOrder(c *gin.Context) {
	userID := c.GetUint("user_id")

	order, err := s.orderService.CreateOrder(userID)
	if err != nil {
		utils.BadRequestResponse(c, "Failed to create order", err)
		return
	}

	utils.CreatedResponse(c, "Order created successfully", order)
}

func (s *Server) getOrders(c *gin.Context) {
	userID := c.GetUint("user_id")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	orders, meta, err := s.orderService.GetOrders(userID, page, limit)
	if err != nil {
		utils.BadRequestResponse(c, "Failed to fetch order ", err)
		return
	}

	utils.PaginatedSuccessResponse(c, "", orders, *meta)
}

func (s *Server) getOrder(c *gin.Context) {
	userID := c.GetUint("user_id")

	id, err := strconv.ParseUint(c.Param("id"), 10, 2)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid order ID", err)
		return
	}

	order, err := s.orderService.GetOrder(userID, uint(id))
	if err != nil {
		utils.NotFoundResponse(c, "Order not found")
		return
	}
	utils.SuccessResponse(c, "Order retrieved successfully", order)

}
