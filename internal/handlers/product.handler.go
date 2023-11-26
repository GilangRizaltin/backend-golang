package handlers

import (
	"Backend_Golang/internal/models"
	"Backend_Golang/internal/repositories"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
)

type HandlerProduct struct {
	*repositories.ProductRepository
}

func InitializeHandler(r *repositories.ProductRepository) *HandlerProduct {
	return &HandlerProduct{r}
}

func (h *HandlerProduct) GetProduct(ctx *gin.Context) {
	var query models.QueryParamsProduct
	var page int
	if query.Page == 0 {
		page = 1
	}
	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Error in binding query get user",
			"Error":   err,
		})
	}
	if query.MaximumPrice != 0 || query.MinimumPrice != 0 {
		if query.MaximumPrice < query.MinimumPrice {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "Maximum price must be greater than Minimum price",
			})
			return
		}

	}
	if _, err := govalidator.ValidateStruct(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Error in Validator",
			"Error":   err.Error(),
		})
		return
	}
	// conditions := []string{
	// 	Product_Name,
	// 	Maximum_Product_Price,
	// 	Minimum_Product_Price,
	// 	Product_Category,
	// 	sortBy,
	// }
	result, err := h.RepositoryGet(&query)
	data, _ := h.RepositoryCountProduct(&query)
	//error handling
	if err != nil {
		if strings.Contains(err.Error(), "trailing junk after numeric literal") {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "Please input right number",
			})
			return
		}
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	if len(result) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "Data Not Found",
		})
		return
	}
	url := ctx.Request.URL.RawQuery
	pages := ctx.Query("page")
	nextPage, prevPage, lastPage := pagination(url, pages, "product?", data[0], page)
	ctx.JSON(http.StatusOK, gin.H{
		"message":    "Get all products success",
		"data":       result,
		"page":       page,
		"total_data": data[0],
		"nextPage":   nextPage,
		"prevPage":   prevPage,
		"lastPage":   lastPage,
	})
}

func (h *HandlerProduct) GetProductDetail(ctx *gin.Context) {
	ID, _ := strconv.Atoi(ctx.Param("id"))
	result, err := h.RepositoryGetDetail(ID)
	if len(result) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "Data not found",
		})
		return
	}
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Get product success",
		"Product": result})
}

func (h *HandlerProduct) CreateProduct(ctx *gin.Context) {
	var newProduct models.ProductModel
	if err := ctx.ShouldBind(&newProduct); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Error in binding query product",
			"error":   err.Error(),
		})
		return
	}
	// if newProduct.Category == "" || newProduct.Product_name == "" || newProduct.Description == "" || newProduct.Price_default == 0 {
	// 	ctx.JSON(http.StatusBadRequest, gin.H{
	// 		"message": "Please fullfill all data"})
	// 	return
	// }
	if _, err := govalidator.ValidateStruct(&newProduct); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Error in Validator",
			"Error":   err.Error(),
		})
		return
	}
	if err := h.RepositoryCreateProduct(&newProduct); err != nil {
		if strings.Contains(err.Error(), "unique_product_name") {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "Product name already used",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{
		"message":      "Product created successfully",
		"Product_Name": newProduct.Product_name})
}

func (h *HandlerProduct) UpdateProduct(ctx *gin.Context) {
	var updateProduct models.UpdateProduct
	ID, _ := strconv.Atoi(ctx.Param("id"))
	if err := ctx.ShouldBind(&updateProduct); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := h.RepositoryUpdateProduct(ID, &updateProduct)
	if err != nil {
		if strings.Contains(err.Error(), "unique_product_name") {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "Product name already used",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "Product not found",
		})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Product successfully updated",
	})
}

func (h *HandlerProduct) DeleteProduct(ctx *gin.Context) {
	ID, _ := strconv.Atoi(ctx.Param("id"))
	result, err := h.RepositoryDeleteProduct(ID)
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
			"message": "Product not found",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Product successfully deleted",
	})
}

func pagination(url, pages, endpoint string, totalData, page int) (string, string, int) {
	var nextPage string
	var prevPage string
	lastPage := int(math.Ceil(float64(totalData) / 6))
	linkPage := "localhost:6121/" + endpoint + url
	nextPage = linkPage[:len(linkPage)-1] + strconv.Itoa(page+1)
	prevPage = linkPage[:len(linkPage)-1] + strconv.Itoa(page-1)
	if pages == "" {
		nextPage = linkPage + "&page=" + strconv.Itoa(page+1)
		prevPage = linkPage + "&page=" + strconv.Itoa(page-1)
	}
	if page == int(lastPage) {
		nextPage = "null"
	}
	if page == 1 {
		prevPage = "null"
	}
	return nextPage, prevPage, lastPage
}

func metaPagination(url, pages, endpoint string, totalData, page int) (string, string, int) {
	var nextPage string
	var prevPage string
	lastPage := int(math.Ceil(float64(totalData) / 6))
	linkPage := "localhost:6121/" + endpoint + url
	nextPage = linkPage[:len(linkPage)-1] + strconv.Itoa(page+1)
	prevPage = linkPage[:len(linkPage)-1] + strconv.Itoa(page-1)
	if pages == "" {
		nextPage = linkPage + "&page=" + strconv.Itoa(page+1)
		prevPage = linkPage + "&page=" + strconv.Itoa(page-1)
	}
	if page == int(lastPage) {
		nextPage = "null"
	}
	if page == 1 {
		prevPage = "null"
	}
	return nextPage, prevPage, lastPage
}
