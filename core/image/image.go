package image

import (
	"fmt"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"

	"github.com/disintegration/imaging"

	_ "golang.org/x/image/webp"
)

func ImageToImage(inputPath, outputPath string) error {
	img, err := imaging.Open(inputPath)
	if err != nil {
		return fmt.Errorf("file couldn't open: %v", err)
	}

	err = imaging.Save(img, outputPath)
	if err != nil {
		return fmt.Errorf("an error occured while converting the image: %v", err)
	}

	return nil
}

func ImageToPDF(imagePaths []string, outputPath string) error {
	conf := model.NewDefaultConfiguration()
	err := api.ImportImagesFile(imagePaths, outputPath, nil, conf)
	if err != nil {
		return fmt.Errorf("an error occured while creating pdf: %v", err)
	}

	return nil
}
