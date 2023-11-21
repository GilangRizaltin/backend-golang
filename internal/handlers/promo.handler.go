package handlers

import (
	"Backend_Golang/internal/models"
	"Backend_Golang/internal/repositories"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type HandlerPromo struct {
	*repositories.PromoRepository
}

func InitializePromoHandler(r *repositories.PromoRepository) *HandlerPromo {
	return &HandlerPromo{r}
}

func (h *HandlerPromo) GetPromo(ctx *gin.Context) {
	Promo_code := ctx.Query("promo-code")
	Time_end := ctx.Query("time-end")
	Sort := ctx.Query("sort")
	page, _ := strconv.Atoi(ctx.Query("page"))
	if page == 0 {
		page = 1
	}
	conditions := []string{
		Promo_code,
		Time_end,
		Sort,
	}
	result, err := h.RepositoryGetPromo(conditions, page)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	if len(result) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "Data Not Found",
			"result":  result,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Get all promo success",
		"data":    result,
		"page":    page,
	})
}

func (h *HandlerPromo) CreatePromo(ctx *gin.Context) {
	var newPromo models.PromoModel
	if err := ctx.ShouldBind(&newPromo); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := h.RepositoryCreatePromo(&newPromo)
	if err != nil {
		log.Fatalln(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{
		"message":    "Product created successfully",
		"Promo_Code": newPromo.Promo_code,
		"Ended at":   newPromo.Ended_at,
	})
}

// func (h *HandlerPromo) UpdatePromo(ctx *gin.Context) {
// 	var updatePromo models.PromoModel
// 	ID, _ := strconv.Atoi(ctx.Param("id"))
// 	if err := ctx.ShouldBind(&updatePromo); err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}
// 	err := h.RepositoryUpdatePromo(ID, &updatePromo)
// 	if err != nil {
// 		log.Fatalln(err)
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	ctx.JSON(http.StatusCreated, gin.H{
// 		"message": "Promo successfully updated",
// 	})
// }

// func (h *HandlerPromo) DeletePromo(ctx *gin.Context) {
// 	ID, _ := strconv.Atoi(ctx.Param("id"))
// 	result, err := h.RepositoryDeletePromo(ID)
// 	if err != nil {
// 		log.Print(err)
// 		ctx.JSON(http.StatusInternalServerError, err)
// 		return
// 	}
// 	rowsAffected, err := result.RowsAffected()
// 	if err != nil {
// 		log.Print(err)
// 		ctx.JSON(http.StatusInternalServerError, err)
// 		return
// 	}
// 	if rowsAffected == 0 {
// 		ctx.JSON(http.StatusNotFound, gin.H{
// 			"message": "Product not found",
// 		})
// 		return
// 	}
// 	ctx.JSON(http.StatusOK, gin.H{
// 		"message": "Product successfully deleted",
// 	})
// }
