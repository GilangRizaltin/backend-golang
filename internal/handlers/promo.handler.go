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
	var query models.QueryParamsPromo
	var page int
	if query.Page == 0 {
		page = 1
	}
	// conditions := []string{
	// 	Promo_code,
	// 	Time_end,
	// 	Sort,
	// }
	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Error in binding query get user",
			"Error":   err,
		})
	}
	result, err := h.RepositoryGetPromo(&query)
	data, _ := h.RepositoryCountPromo(&query)
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
	url := ctx.Request.URL.RawQuery
	pages := ctx.Query("page")
	nextPage, prevPage, lastPage := pagination(url, pages, "promo?", data[0], page)
	ctx.JSON(http.StatusOK, gin.H{
		"message":    "Get all promo success",
		"data":       result,
		"page":       page,
		"total_data": data[0],
		"nextPage":   nextPage,
		"prevPage":   prevPage,
		"lastPage":   lastPage,
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

func (h *HandlerPromo) UpdatePromo(ctx *gin.Context) {
	var updatePromo models.PromoModel
	ID, _ := strconv.Atoi(ctx.Param("id"))
	if err := ctx.ShouldBind(&updatePromo); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := h.RepositoryUpdatePromo(ID, &updatePromo)
	if err != nil {
		log.Fatalln(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Promo successfully updated",
	})
}

func (h *HandlerPromo) DeletePromo(ctx *gin.Context) {
	ID, _ := strconv.Atoi(ctx.Param("id"))
	result, err := h.RepositoryDeletePromo(ID)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	if rowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "Promo not found",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Promo successfully deleted",
	})
}
