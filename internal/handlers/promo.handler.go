package handlers

import (
	"Backend_Golang/internal/helpers"
	"Backend_Golang/internal/models"
	"Backend_Golang/internal/repositories"
	"log"
	"net/http"
	"strconv"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
)

type HandlerPromo struct {
	repositories.IPromoRepository
}

func InitializePromoHandler(r repositories.IPromoRepository) *HandlerPromo {
	return &HandlerPromo{r}
}

func (h *HandlerPromo) GetPromo(ctx *gin.Context) {
	var query models.QueryParamsPromo
	// conditions := []string{
	// 	Promo_code,
	// 	Time_end,
	// 	Sort,
	// }
	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, helpers.NewResponse("Error binding query promo", nil, nil))
		log.Println(err.Error())
	}
	if _, err := govalidator.ValidateStruct(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, helpers.NewResponse("Wrong input after validation", nil, nil))
		log.Println(err.Error())
		return
	}
	result, err := h.RepositoryGetPromo(&query)
	data, _ := h.RepositoryCountPromo(&query)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, helpers.NewResponse("Internal Server Error", nil, nil))
		return
	}
	if len(result) == 0 {
		ctx.JSON(http.StatusNotFound, helpers.NewResponse("Promo not found", nil, nil))
		return
	}
	meta := helpers.GetPagination(ctx, data, query.Page, 6)
	ctx.JSON(http.StatusOK, helpers.NewResponse("Successfully get all promo", result, &meta))
}

func (h *HandlerPromo) CreatePromo(ctx *gin.Context) {
	var newPromo models.PromoModel
	if err := ctx.ShouldBind(&newPromo); err != nil {
		ctx.JSON(http.StatusBadRequest, helpers.NewResponse("Error binding body request promo", nil, nil))
		return
	}
	if _, err := govalidator.ValidateStruct(&newPromo); err != nil {
		ctx.JSON(http.StatusBadRequest, helpers.NewResponse("Wrong input after validation", nil, nil))
		log.Println(err.Error())
		return
	}
	err := h.RepositoryCreatePromo(&newPromo)
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, helpers.NewResponse("Internal Server Error", nil, nil))
		return
	}
	ctx.JSON(http.StatusCreated, helpers.NewResponse("Successfully created promo", newPromo, nil))
}

func (h *HandlerPromo) UpdatePromo(ctx *gin.Context) {
	var updatePromo models.UpdatePromoModel
	ID, _ := strconv.Atoi(ctx.Param("id"))
	if err := ctx.ShouldBind(&updatePromo); err != nil {
		ctx.JSON(http.StatusBadRequest, helpers.NewResponse("Error binding body request update promo", nil, nil))
		log.Println(err.Error())
		return
	}
	if _, err := govalidator.ValidateStruct(&updatePromo); err != nil {
		ctx.JSON(http.StatusBadRequest, helpers.NewResponse("Wrong input after validation", nil, nil))
		log.Println(err.Error())
		return
	}
	result, err := h.RepositoryUpdatePromo(ID, &updatePromo)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helpers.NewResponse("Internal Server Error", nil, nil))
		log.Println(err.Error())
		return
	}
	var dataNotFound int64 = 1
	if result < dataNotFound {
		ctx.JSON(http.StatusNotFound, helpers.NewResponse("Promo not found", nil, nil))
		return
	}
	ctx.JSON(http.StatusCreated, helpers.NewResponse("Successfully update promo", updatePromo, nil))
}

func (h *HandlerPromo) DeletePromo(ctx *gin.Context) {
	ID, _ := strconv.Atoi(ctx.Param("id"))
	result, err := h.RepositoryDeletePromo(ID)
	if err != nil {
		log.Print(err.Error())
		ctx.JSON(http.StatusInternalServerError, helpers.NewResponse("Internal Server Error", nil, nil))
		return
	}
	// rowsAffected, err := result.RowsAffected()
	var dataNotFound int64 = 1
	if result < dataNotFound {
		ctx.JSON(http.StatusNotFound, helpers.NewResponse("Promo not found", nil, nil))
		return
	}
	ctx.JSON(http.StatusOK, helpers.NewResponse("Successfully delete promo", ID, nil))
}
