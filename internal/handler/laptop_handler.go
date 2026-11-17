package handler

import (
	"errors"
	"net/http"
	"strconv"

	"biblioteca-api/internal/models"
	"biblioteca-api/internal/service"

	"github.com/gin-gonic/gin"
)

type LaptopHandler struct {
	laptopService *service.LaptopService
}

func NewLaptopHandler(laptopService *service.LaptopService) *LaptopHandler {
	return &LaptopHandler{laptopService: laptopService}
}

// Create registra una laptop (solo administrador).
func (h *LaptopHandler) Create(c *gin.Context) {
	var req models.LaptopRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "datos de laptop inválidos"})
		return
	}

	laptop, err := h.laptopService.Create(req)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "no se pudo registrar la laptop; verifica que el número de serie no esté repetido"})
		return
	}

	c.JSON(http.StatusCreated, laptop)
}

// List devuelve todas las laptops.
func (h *LaptopHandler) List(c *gin.Context) {
	laptops, err := h.laptopService.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudieron consultar las laptops"})
		return
	}

	c.JSON(http.StatusOK, laptops)
}

// ListAvailable devuelve laptops disponibles.
func (h *LaptopHandler) ListAvailable(c *gin.Context) {
	laptops, err := h.laptopService.ListAvailable()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudieron consultar las laptops disponibles"})
		return
	}

	c.JSON(http.StatusOK, laptops)
}

// Borrow registra un préstamo.
func (h *LaptopHandler) Borrow(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "usuario no autenticado"})
		return
	}

	id, ok := userID.(int)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "identificador de usuario inválido"})
		return
	}

	var req models.LaptopLoanRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.LaptopID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "debes proporcionar un laptop_id válido"})
		return
	}

	loan, err := h.laptopService.Borrow(id, req)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, loan)
}

// Return devuelve una laptop.
func (h *LaptopHandler) Return(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "usuario no autenticado"})
		return
	}

	id, ok := userID.(int)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "identificador de usuario inválido"})
		return
	}

	loanID, err := strconv.Atoi(c.Param("id"))
	if err != nil || loanID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "identificador de préstamo inválido"})
		return
	}

	err = h.laptopService.Return(id, loanID)
	if err != nil {
		status := http.StatusConflict
		if errors.Is(err, service.ErrLaptopLoanNotFound) {
			status = http.StatusNotFound
		}
		if errors.Is(err, service.ErrLaptopLoanForbidden) {
			status = http.StatusForbidden
		}

		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "laptop devuelta correctamente"})
}

// MyLoans devuelve los préstamos del usuario autenticado.
func (h *LaptopHandler) MyLoans(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "usuario no autenticado"})
		return
	}

	id, ok := userID.(int)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "identificador de usuario inválido"})
		return
	}

	loans, err := h.laptopService.ListByUser(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudieron consultar los préstamos"})
		return
	}

	c.JSON(http.StatusOK, loans)
}
