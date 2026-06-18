package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/bedirmirac/convertev/core/document"
	"github.com/bedirmirac/convertev/core/image"
	"github.com/bedirmirac/convertev/core/media"
)

type Args struct {
	Mode   string
	Input  string
	Output string
}

func (a Args) Validate() {
	if a.Mode == "doc" && a.Output != "" {
		fmt.Println("Usage in 'doc' mode: convertev -mode doc -i [input]", "Do NOT use '-o'!")
		flag.Usage()
		os.Exit(1)
	}
	if a.Mode == "" || a.Input == "" {
		fmt.Println("Usage: convertev -mode [img|media] -i [input] -o [output] || convertev -mode [doc] -i [input]")
		flag.Usage()
		os.Exit(1)
	}
	if a.Mode == "media" || a.Mode == "img" && a.Output == "" {
		fmt.Println("Usage in 'media' and 'img': convertev -mode [img|media] -i [input] -o [output]", "Output is needed!")
		flag.Usage()
		os.Exit(1)
	}
}

func main() {
	mode := flag.String("mode", "", "Process type: doc, img, media")
	input := flag.String("i", "", "Input file (or folder) If it is used to convert images into PDF use img1,img2,img3 format) ")
	output := flag.String("o", "", "output file/path")
	ffmpegPath := flag.String("ffmpeg", "ffmpeg", "FFmpeg path")
	ffprobePath := flag.String("ffprobe", "ffprobe", "FFprobe path")
	flag.Parse()

	args := Args{
		Mode:   *mode,
		Input:  *input,
		Output: *output,
	}

	args.Validate()

	switch *mode {
	case "doc":
		if strings.Contains(*input, ".pdf") {
			err := document.DocFromPdf(*input)
			if err != nil {
				log.Fatalf("Error while pdf convering to doc")
			}
		} else {
			err := document.PdfFromOffice(*input)
			if err != nil {
				log.Fatalf("Error while document converting : %v", err)
			}
		}
		fmt.Println("Document converted.")

	case "img":
		if strings.Contains(*output, ".pdf") {
			imgs := strings.Split(*input, (","))
			err := image.ImageToPDF(imgs, *output)
			if err != nil {
				log.Fatalf("Error while image converting to PDF: %v", err)
			}
			fmt.Println("Image converted to PDF.")
		} else {
			err := image.ImageToImage(*input, *output)
			if err != nil {
				log.Fatalf("Error while image converting: %v", err)
			}
			fmt.Println("Image converted.")
		}

	case "media":
		duration, err := media.GetMediaDuration(*input, *ffprobePath)
		if err != nil {
			log.Fatalf("Couldn't fetch the media info: %v", err)
		}

		err = media.MediaConverter(*input, *output, *ffmpegPath, duration, func(percent float64) {
			fmt.Printf("\rProccess info: %.2f%%", percent)
		})
		if err != nil {
			log.Fatalf("\nError while media converting: %v", err)
		}
		fmt.Println("\nMedia converted.")

	default:
		fmt.Println("Invalid mod. Please select doc, img or media.")
	}
}
