package handlers

import (
	"Backend_Golang/internal/helpers"
	"Backend_Golang/internal/models"
	"Backend_Golang/internal/repositories"
	"fmt"
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
	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, helpers.ResponseFailed("Error in binding query get user", err.Error()))
		return
	}
	var page = 1
	if query.Page != 0 {
		page = query.Page
	}
	if query.MaximumPrice != 0 || query.MinimumPrice != 0 {
		if query.MaximumPrice < query.MinimumPrice {
			ctx.JSON(http.StatusBadRequest, helpers.ResponseFailed("Maximum price must be greater than Minimum price", ""))
			return
		}
	}
	if _, err := govalidator.ValidateStruct(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, helpers.ResponseFailed("Error in validation get product", err.Error()))
		return
	}
	result, err := h.RepositoryGet(&query)
	data, _ := h.RepositoryCountProduct(&query)
	if err != nil {
		if strings.Contains(err.Error(), "trailing junk after numeric literal") {
			ctx.JSON(http.StatusBadRequest, helpers.ResponseFailed("Please input right number", ""))
			return
		}
		ctx.JSON(http.StatusInternalServerError, helpers.ResponseFailed("Error in repository get product", err.Error()))
		return
	}
	if len(result) == 0 {
		ctx.JSON(http.StatusNotFound, helpers.ResponseFailed("Data not found", ""))
		return
	}
	url := ctx.Request.URL.RawQuery
	pages := ctx.Query("page")
	meta := helpers.MetaPage(url, pages, "product?", data[0], page)
	ctx.JSON(http.StatusOK, helpers.ResponseGetData("Successfully Get Product", result, meta))
}

func (h *HandlerProduct) GetProductDetail(ctx *gin.Context) {
	ID, _ := strconv.Atoi(ctx.Param("id"))
	result, err := h.RepositoryGetDetail(ID)
	if len(result) == 0 {
		ctx.JSON(http.StatusNotFound, helpers.ResponseFailed("Data not found", ""))
		return
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helpers.ResponseFailed("Error in repository get product detail", err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, helpers.ResponseSuccess("Successfully Get Product Detail", result))
}

func (h *HandlerProduct) CreateProduct(ctx *gin.Context) {
	var newProduct models.ProductModel
	if err := ctx.ShouldBind(&newProduct); err != nil {
		ctx.JSON(http.StatusBadRequest, helpers.ResponseFailed("Error in binding body product", err.Error()))
		return
	}
	if _, err := govalidator.ValidateStruct(&newProduct); err != nil {
		ctx.JSON(http.StatusBadRequest, helpers.ResponseFailed("Error in validation body product", err.Error()))
		return
	}
	//cloud upload
	cld, errCloud := helpers.InitCloudinary()
	if errCloud != nil {
		ctx.JSON(http.StatusInternalServerError, helpers.ResponseFailed("Error initializes cloudinary", errCloud.Error()))
		return
	}
	formFiles, _ := ctx.MultipartForm()
	// if err != nil {
	// 	ctx.JSON(http.StatusInternalServerError, gin.H{
	// 		"message": "Error in multipart form",
	// 		"error":   err,
	// 	})
	// 	return
	// }
	var dataUrls []string
	if formFiles != nil {
		files := formFiles.File["Product_photo"]
		for idx, formFile := range files {
			file, err := formFile.Open()
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, helpers.ResponseFailed("Error opening file", err.Error()))
				return
			}
			defer file.Close()
			publicID := fmt.Sprintf("%s %s_%s-%d", "Products", newProduct.Product_name, "Product_photo", idx+1)
			folder := ""
			res, err := cld.Uploader(ctx, file, publicID, folder)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, helpers.ResponseFailed("Error uploading file to cloudinary", err.Error()))
				return
			}
			dataUrls = append(dataUrls, res.SecureURL)
		}
	}
	if err := h.RepositoryCreateProduct(&newProduct, dataUrls); err != nil {
		if strings.Contains(err.Error(), "unique_product_name") {
			ctx.JSON(http.StatusBadRequest, helpers.ResponseFailed("Product name already used", ""))
			return
		}
		ctx.JSON(http.StatusInternalServerError, helpers.ResponseFailed("Error in repository create product", err.Error()))
		return
	}
	ctx.JSON(http.StatusCreated, helpers.ResponseSuccess("Successfully create product", newProduct.Product_name))
}

func (h *HandlerProduct) UpdateProduct(ctx *gin.Context) {
	var updateProduct models.UpdateProduct
	ID, _ := strconv.Atoi(ctx.Param("id"))
	if err := ctx.ShouldBind(&updateProduct); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Error in binding update body",
			"error":   err.Error(),
		})
		return
	}
	if _, err := govalidator.ValidateStruct(&updateProduct); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Error in Validator",
			"Error":   err.Error(),
		})
		return
	}
	//cloud upload
	cld, errCloud := helpers.InitCloudinary()
	if errCloud != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error in initialize cloudinary",
			"error":   errCloud,
		})
		return
	}
	formFiles, _ := ctx.MultipartForm()
	// if err != nil {
	// 	ctx.JSON(http.StatusInternalServerError, gin.H{
	// 		"message": "Error in get data photo",
	// 		"error":   err.Error(),
	// 	})
	// 	return
	// }
	var dataUrls []string
	dataIndexPhoto := updateProduct.Photo_index
	if formFiles != nil {
		files := formFiles.File["Product_photo"]
		for idx, formFile := range files {
			file, err := formFile.Open()
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"message": "Error opening file",
					"error":   err.Error(),
				})
				return
			}
			defer file.Close()
			publicID := fmt.Sprintf("%s %s_%s-%d", "Products", ctx.Param("id"), "Product_photo", dataIndexPhoto[idx])
			folder := ""
			res, err := cld.Uploader(ctx, file, publicID, folder)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"message": "Error uploading file to Cloudinary",
					"error":   err.Error(),
				})
				return
			}
			dataUrls = append(dataUrls, res.SecureURL)
		}
	}
	result, err := h.RepositoryUpdateProduct(ID, &updateProduct, dataUrls)
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

func (h *HandlerProduct) GetStatistic(ctx *gin.Context) {
	dateStart := ctx.Query("date-start")
	dateEnd := ctx.Query("date-end")
	result, err := h.RepositoryStatisticProduct(dateStart, dateEnd)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error in popular products",
			"ERROR":   err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Get products selled statistic success",
		"data":    result,
	})
}

func (h *HandlerProduct) GetPopular(ctx *gin.Context) {
	dateStart := ctx.Query("date-start")
	dateEnd := ctx.Query("date-end")
	result, err := h.RepositoryPopularProduct(dateStart, dateEnd)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error in popular products",
			"ERROR":   err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Get popular products success",
		"data":    result,
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
