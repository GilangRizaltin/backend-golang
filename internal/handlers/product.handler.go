package handlers

import (
	"Backend_Golang/internal/helpers"
	"Backend_Golang/internal/models"
	"Backend_Golang/internal/repositories"
	"fmt"
	"log"
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
		ctx.JSON(http.StatusBadRequest, helpers.NewResponse("Error binding query get product", nil, nil))
		log.Println(err.Error())
		return
	}
	if query.ProductName != "" {
		isValid := helpers.ValidateInput(query.ProductName)
		if !isValid {
			ctx.JSON(http.StatusBadRequest, helpers.NewResponse("Wrong input for product name", nil, nil))
			return
		}
	}
	if query.MaximumPrice != 0 && query.MinimumPrice != 0 {
		if query.MaximumPrice < query.MinimumPrice {
			ctx.JSON(http.StatusBadRequest, helpers.NewResponse("Maximum price must greater than minimum price", nil, nil))
			return
		}
	}
	if _, err := govalidator.ValidateStruct(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, helpers.NewResponse("Wrong input after validation", nil, nil))
		log.Println(err.Error())
		return
	}
	result, err := h.RepositoryGet(&query)
	data, _ := h.RepositoryCountProduct(&query)
	if err != nil {
		if strings.Contains(err.Error(), "trailing junk after numeric literal") {
			ctx.JSON(http.StatusBadRequest, helpers.NewResponse("trailing junk after numeric literal", nil, nil))
			return
		}
		ctx.JSON(http.StatusInternalServerError, helpers.NewResponse("Internal Server Error", nil, nil))
		log.Println(err.Error())
		return
	}
	if len(result) == 0 {
		ctx.JSON(http.StatusNotFound, helpers.NewResponse("Data not found", nil, nil))
		return
	}
	meta := helpers.GetPagination(ctx, data, query.Page)
	ctx.JSON(http.StatusOK, helpers.NewResponse("Successfully Get Product", result, &meta))
}

func (h *HandlerProduct) GetProductDetail(ctx *gin.Context) {
	ID, _ := strconv.Atoi(ctx.Param("id"))
	result, err := h.RepositoryGetDetail(ID)
	if len(result) == 0 {
		ctx.JSON(http.StatusNotFound, helpers.NewResponse("Data not found", nil, nil))
		return
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helpers.NewResponse("Internal Server Error", nil, nil))
		log.Println(err.Error())
		return
	}
	ctx.JSON(http.StatusOK, helpers.NewResponse("Successfully Get Product", result, nil))
}

func (h *HandlerProduct) CreateProduct(ctx *gin.Context) {
	var newProduct models.ProductModel
	if err := ctx.ShouldBind(&newProduct); err != nil {
		ctx.JSON(http.StatusBadRequest, helpers.NewResponse("Error binding body request", nil, nil))
		log.Println(err.Error())
		return
	}
	if _, err := govalidator.ValidateStruct(&newProduct); err != nil {
		ctx.JSON(http.StatusBadRequest, helpers.NewResponse("Wrong inpu after validation", nil, nil))
		log.Println(err.Error())
		return
	}
	//cloud upload
	cld, errCloud := helpers.InitCloudinary()
	if errCloud != nil {
		ctx.JSON(http.StatusInternalServerError, helpers.NewResponse("Failed during initialization", nil, nil))
		log.Println(errCloud.Error())
		return
	}
	formFiles, _ := ctx.MultipartForm()
	var dataUrls []string
	if formFiles != nil {
		files := formFiles.File["Product_photo"]
		for idx, formFile := range files {
			file, err := formFile.Open()
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, helpers.NewResponse("Error during open form file", nil, nil))
				log.Println(err.Error())
				return
			}
			defer file.Close()
			publicID := fmt.Sprintf("%s %s_%s-%d", "Products", newProduct.Product_name, "Product_photo", idx+1)
			folder := ""
			res, err := cld.Uploader(ctx, file, publicID, folder)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, helpers.NewResponse("Error uploading image", nil, nil))
				log.Println(err.Error())
				return
			}
			dataUrls = append(dataUrls, res.SecureURL)
		}
	}
	if err := h.RepositoryCreateProduct(&newProduct, dataUrls); err != nil {
		if strings.Contains(err.Error(), "unique_product_name") {
			ctx.JSON(http.StatusBadRequest, helpers.NewResponse("Product name already used", nil, nil))
			return
		}
		ctx.JSON(http.StatusInternalServerError, helpers.NewResponse("Internal Server Error", nil, nil))
		log.Println(err.Error())
		return
	}
	ctx.JSON(http.StatusCreated, helpers.NewResponse("Successfully create product", newProduct.Product_name, nil))
}

func (h *HandlerProduct) UpdateProduct(ctx *gin.Context) {
	var updateProduct models.UpdateProduct
	ID, _ := strconv.Atoi(ctx.Param("id"))
	if err := ctx.ShouldBind(&updateProduct); err != nil {
		ctx.JSON(http.StatusBadRequest, helpers.NewResponse("Error binding product body for update", nil, nil))
		return
	}
	if _, err := govalidator.ValidateStruct(&updateProduct); err != nil {
		ctx.JSON(http.StatusBadRequest, helpers.NewResponse("Wrong input after validation", nil, nil))
		log.Println(err.Error())
		return
	}
	//cloud upload
	cld, errCloud := helpers.InitCloudinary()
	if errCloud != nil {
		ctx.JSON(http.StatusInternalServerError, helpers.NewResponse("Error during initialization", nil, nil))
		log.Println(errCloud.Error())
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
				ctx.JSON(http.StatusInternalServerError, helpers.NewResponse("Error during open file", nil, nil))
				log.Println(err.Error())
				return
			}
			defer file.Close()
			publicID := fmt.Sprintf("%s %s_%s-%d", "Products", ctx.Param("id"), "Product_photo", dataIndexPhoto[idx])
			folder := ""
			res, err := cld.Uploader(ctx, file, publicID, folder)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, helpers.NewResponse("Error during upload file", nil, nil))
				log.Println(err.Error())
				return
			}
			dataUrls = append(dataUrls, res.SecureURL)
		}
	}
	result, err := h.RepositoryUpdateProduct(ID, &updateProduct, dataUrls)
	if err != nil {
		if strings.Contains(err.Error(), "unique_product_name") {
			ctx.JSON(http.StatusBadRequest, helpers.NewResponse("Product name already used", nil, nil))
			return
		}
		ctx.JSON(http.StatusInternalServerError, helpers.NewResponse("Internal Server Error", nil, nil))
		log.Println(err.Error())
		return
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, helpers.NewResponse("Product not found to update", nil, nil))
		return
	}
	ctx.JSON(http.StatusCreated, helpers.NewResponse("Successfully update product", updateProduct, nil))
}

func (h *HandlerProduct) DeleteProduct(ctx *gin.Context) {
	ID, _ := strconv.Atoi(ctx.Param("id"))
	result, err := h.RepositoryDeleteProduct(ID)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, helpers.NewResponse("Internal Erver Error", nil, nil))
		return
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helpers.NewResponse("Internal Server Error", nil, nil))
		log.Print(err.Error())
		return
	}
	if rowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, helpers.NewResponse("Product that will deleted not found", nil, nil))
		return
	}
	ctx.JSON(http.StatusOK, helpers.NewResponse("Successfully delete product", ID, nil))
}

func (h *HandlerProduct) GetStatisticProduct(ctx *gin.Context) {
	dateStart := ctx.Query("date-start")
	dateEnd := ctx.Query("date-end")
	result, err := h.RepositoryStatisticProduct(dateStart, dateEnd, "")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helpers.NewResponse("Internal Server Error", nil, nil))
		log.Println(err.Error())
		return
	}
	if len(result) < 1 {
		ctx.JSON(http.StatusNotFound, helpers.NewResponse("Data not found", nil, nil))
		return
	}
	ctx.JSON(http.StatusOK, helpers.NewResponse("Successfully get product statistic", result, nil))
}

func (h *HandlerProduct) GetFavouriteProduct(ctx *gin.Context) {
	dateStart := ctx.Query("date-start")
	dateEnd := ctx.Query("date-end")
	data, err := h.RepositoryStatisticProduct(dateStart, dateEnd, "favourite")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helpers.NewResponse("Internal Server Error Statistic Product", nil, nil))
		log.Println(err.Error())
		return
	}
	result, err := h.RepositoryFavouriteProduct(data)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helpers.NewResponse("Internal Server Error", nil, nil))
		log.Println(err.Error())
		return
	}
	ctx.JSON(http.StatusOK, helpers.NewResponse("Successfully get popular products", result, nil))
}
