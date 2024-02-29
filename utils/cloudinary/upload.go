package cloudinary

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"my-tourist-ticket/app/configs"
	"os"
	"path/filepath"
	"strings"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type CloudinaryUploaderInterface interface {
	UploadImage(fileHeader *multipart.FileHeader) (string, error)
	Destroy(publicID string) error
}

type CloudinaryUploader struct {
}

func New() CloudinaryUploaderInterface {
	return &CloudinaryUploader{}
}

func (cu *CloudinaryUploader) UploadImage(fileHeader *multipart.FileHeader) (string, error) {
	ctx := context.Background()

	cld, err := cloudinary.NewFromURL(configs.CLD_URL)
	if err != nil {
		return "", err
	}

	file, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("error opening file: %w", err)
	}

	defer file.Close()

	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if ext != ".png" && ext != ".jpg" && ext != ".jpeg" {
		return "", fmt.Errorf("invalid file type: %w", err)
	}

	uploadParams := uploader.UploadParams{
		Folder: "BE20_MyTouristTicket",
	}

	resp, err := cld.Upload.Upload(ctx, file, uploadParams)
	if err != nil {
		return "", fmt.Errorf("error uploading to Cloudinary: %w", err)
	}

	_, err = file.Seek(0, 0)
	if err != nil {
		return "", fmt.Errorf("error seeking file: %w", err)
	}

	localFilePath := filepath.Join("utils", "images", fileHeader.Filename)
	if err := SaveImageToLocal(file, localFilePath); err != nil {
		return "", fmt.Errorf("error saving image to local: %w", err)
	}

	return resp.SecureURL, nil
}

func SaveImageToLocal(file multipart.File, destinationPath string) error {
	dst, err := os.Create(destinationPath)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		return fmt.Errorf("error copying file: %w", err)
	}

	return nil
}

func (cu *CloudinaryUploader) Destroy(publicId string) error {
	ctx := context.Background()

	cldUrl := configs.CLD_URL
	cldService, _ := cloudinary.NewFromURL(cldUrl)
	_, errDest := cldService.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID:     publicId,
		Type:         "",
		ResourceType: "",
		Invalidate:   new(bool),
	})
	if errDest != nil {
		return errDest
	}

	return nil
}

func GetPublicID(url string) string {
	splitURL := strings.Split(url, "/")
	fileNameWithExt := splitURL[len(splitURL)-1]
	fileNameWithoutExt := strings.TrimSuffix(fileNameWithExt, filepath.Ext(fileNameWithExt))
	publicID := splitURL[len(splitURL)-2] + "/" + fileNameWithoutExt
	return publicID
}
