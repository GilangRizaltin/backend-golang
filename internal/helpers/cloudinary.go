package helpers

import (
	"fmt"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
)

type CloudinaryUploader struct {
	*cloudinary.Cloudinary
}

func InitCloudinary() (*CloudinaryUploader, error) {
	cldString := fmt.Sprintf("cloudinary://%s:%s@%s", os.Getenv("CLOUDINARY_KEY"), os.Getenv("CLOUDINARY_SECRET"), os.Getenv("CLOUDINARY_NAME"))
	cld, err := cloudinary.NewFromURL(cldString)
	if err != nil {
		return nil, err
	}
	cld.Config.URL.Secure = true
	return &CloudinaryUploader{cld}, nil
}

func (c *CloudinaryUploader) Uploader(ctx *gin.Context, file interface{}, publicId, folder string) (*uploader.UploadResult, error) {
	if folder == "" {
		folder = "Coffee Shop"
	}
	response, err := c.Upload.Upload(ctx, file, uploader.UploadParams{
		PublicID: publicId,
		Folder:   folder,
	})
	if err != nil {
		return nil, err
	}
	return response, nil
}
