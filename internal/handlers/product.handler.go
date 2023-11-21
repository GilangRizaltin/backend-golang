package handlers

import (
	"Backend_Golang/internal/models"
	"Backend_Golang/internal/repositories"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type HandlerProduct struct {
	*repositories.ProductRepository
}

func InitializeHandler(r *repositories.ProductRepository) *HandlerProduct {
	return &HandlerProduct{r}
}

func (h *HandlerProduct) GetProduct(ctx *gin.Context) {
	Product_Name := ctx.Query("search")
	Maximum_Product_Price := ctx.Query("max_price")
	Minimum_Product_Price := ctx.Query("min_price")
	Product_Category := ctx.Query("product_category")
	sortBy := ctx.Query("sortBy")
	sortOrder := ctx.Query("sortOrder")
	page, _ := strconv.Atoi(ctx.Query("page"))
	if page == 0 {
		page = 1
	}
	conditions := []string{
		Product_Name,
		Maximum_Product_Price,
		Minimum_Product_Price,
		Product_Category,
		sortBy,
		sortOrder,
	}
	result, err := h.RepositoryGet(conditions, page)
	data, _ := h.RepositoryCountProduct(conditions)
	if err != nil {
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
	// url := ctx.Request.URL.RawQuery
	// lastPage := math.Round(float64(data[0]) / 6)
	// next := page + 1
	// prev := page - 1
	// nextPage := "localhost:6121/product?" + url
	// prevPage := "localhost:6121/product" + url
	// if page == int(lastPage) {
	// 	nextPage = "null"
	// }
	// if page == 1 {
	// 	prevPage = "null"
	// }
	ctx.JSON(http.StatusOK, gin.H{
		"message":    "Get all products success",
		"data":       result,
		"page":       page,
		"total_data": data[0],
		// "url":        url,
		// "Laspage":    nextPage,
	})
}

func (h *HandlerProduct) CreateProduct(ctx *gin.Context) {
	var newProduct models.ProductModel
	if err := ctx.BindJSON(&newProduct); err != nil { //shouldBind datanya ga masuk
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.RepositoryCreateProduct(&newProduct); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{
		"message":      "Product created successfully",
		"Product_Name": newProduct.Product_name})
}

func (h *HandlerProduct) UpdateProduct(ctx *gin.Context) {
	var updateProduct models.ProductModel
	// ID := ctx.Param("id")
	ID, _ := strconv.Atoi(ctx.Param("id"))
	if err := ctx.BindJSON(&updateProduct); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.RepositoryUpdateProduct(ID, &updateProduct); err != nil {
		log.Fatalln(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
