package helpers

import (
	"fmt"
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

type Meta struct {
	Page      int    `json:"page,omitempty"`
	NoPage    string `json:"no_page,omitempty"`
	NextPage  string `json:"next,omitempty"`
	PrevPage  string `json:"prev,omitempty"`
	TotalData int    `json:"total_data,omitempty"`
	TotalPage int    `json:"total_page,omitempty"`
}

func GetPagination(ctx *gin.Context, totalData []int, page int) Meta {
	var nextPage, prevPage string
	url := fmt.Sprintf("%s%s", ctx.Request.Host, ctx.Request.URL.RequestURI())
	pages := 1
	if page != 0 {
		pages = page
	}
	nextPage = url[:len(url)-1] + strconv.Itoa(page+1)
	prevPage = url[:len(url)-1] + strconv.Itoa(page-1)
	lastPage := int(math.Ceil(float64(totalData[0]) / 6))
	if page == 0 {
		nextPage = fmt.Sprintf("%s&page=%d", url, pages+1)
		prevPage = "null"
		if pages == lastPage {
			nextPage = "null"
		}
	}
	if page == lastPage {
		nextPage = "null"
	}
	if page == 1 {
		prevPage = "null"
	}
	return Meta{
		Page:      pages,
		NoPage:    url[:len(url)-1],
		NextPage:  nextPage,
		PrevPage:  prevPage,
		TotalPage: lastPage,
		TotalData: totalData[0],
	}
}

func NewResponse(message string, data interface{}, meta *Meta) Response {
	return Response{
		Message: message,
		Data:    data,
		Meta:    meta,
	}
}
