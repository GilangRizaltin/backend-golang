package handlers

import (
	"Backend_Golang/internal/models"
	"Backend_Golang/internal/repositories"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"

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
	sortBy := ctx.Query("sort")
	page, _ := strconv.Atoi(ctx.Query("page"))
	if page == 0 {
		page = 1
	}
	if Maximum_Product_Price != "" || Minimum_Product_Price != "" {
		if Maximum_Product_Price != "" {
			_, errMax := strconv.Atoi(Maximum_Product_Price)
			if errMax != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"message": "Please input right number",
				})
				return
			}
		}
		if Minimum_Product_Price != "" {
			_, errMin := strconv.Atoi(Minimum_Product_Price)
			if errMin != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"message": "Please input right number",
				})
				return
			}
		}
		if Maximum_Product_Price != "" && Minimum_Product_Price != "" {
			priceMax, _ := strconv.Atoi(Maximum_Product_Price)
			priceMin, _ := strconv.Atoi(Minimum_Product_Price)
			if priceMax < priceMin {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"message": "Maximum price must be greater than Minimum price",
				})
				return
			}
		}
	}
	conditions := []string{
		Product_Name,
		Maximum_Product_Price,
		Minimum_Product_Price,
		Product_Category,
		sortBy,
	}
	result, err := h.RepositoryGet(conditions, page)
	data, _ := h.RepositoryCountProduct(conditions)
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
	var nextPage string
	var prevPage string
	lastPage := int(math.Ceil(float64(data[0]) / 6))
	linkPage := "localhost:6121/product?" + url
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
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if newProduct.Category == "" || newProduct.Product_name == "" || newProduct.Description == "" || newProduct.Price_default == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Please fullfill all data"})
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
	ID, _ := strconv.Atoi(ctx.Param("id"))
	if err := ctx.ShouldBind(&updateProduct); err != nil {
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
