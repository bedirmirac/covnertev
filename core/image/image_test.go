package image

import (
	stdimage "image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"testing"
)

// createTestPNG creates a small 10x10 PNG file inside dir for use in tests.
func createTestPNG(t *testing.T, dir, name string) string {
	t.Helper()

	path := filepath.Join(dir, name)
	f, err := os.Create(path)
	if err != nil {
		t.Fatalf("failed to create test PNG: %v", err)
	}
	defer f.Close()

	// 10x10 solid color square
	img := stdimage.NewRGBA(stdimage.Rect(0, 0, 10, 10))
	for y := range 10 {
		for x := range 10 {
			img.Set(x, y, color.RGBA{R: 200, G: 100, B: 50, A: 255})
		}
	}
	if err := png.Encode(f, img); err != nil {
		t.Fatalf("PNG encode error: %v", err)
	}
	return path
}

// -----------------------------------------
// ImageToImage
// -----------------------------------------

func TestImageToImage_PNGtoJPG(t *testing.T) {
	dir := t.TempDir()
	input := createTestPNG(t, dir, "test.png")
	output := filepath.Join(dir, "out.jpg")

	if err := ImageToImage(input, output); err != nil {
		t.Fatalf("ImageToImage() error: %v", err)
	}

	info, err := os.Stat(output)
	if os.IsNotExist(err) {
		t.Fatal("output file was not created")
	}
	if info.Size() == 0 {
		t.Error("output file is empty")
	}
}

func TestImageToImage_PNGtoPNG(t *testing.T) {
	dir := t.TempDir()
	input := createTestPNG(t, dir, "test.png")
	output := filepath.Join(dir, "out.png")

	if err := ImageToImage(input, output); err != nil {
		t.Fatalf("ImageToImage() error: %v", err)
	}

	if _, err := os.Stat(output); os.IsNotExist(err) {
		t.Fatal("output file was not created")
	}
}

// TestImageToImage_PNGtoWEBP_WriteUnsupported documents that disintegration/imaging
// can read WEBP via the golang.org/x/image/webp decoder but cannot write WEBP.
// This contradicts the README claim that WEBP is a supported target format.
func TestImageToImage_PNGtoWEBP_WriteUnsupported(t *testing.T) {
	dir := t.TempDir()
	input := createTestPNG(t, dir, "test.png")
	output := filepath.Join(dir, "out.webp")

	err := ImageToImage(input, output)
	if err == nil {
		t.Error("expected error for WEBP write (unsupported by library); the library may have been updated")
	} else {
		t.Logf("expected error (WEBP write unsupported): %v", err)
	}
}

func TestImageToImage_InvalidInput(t *testing.T) {
	dir := t.TempDir()
	output := filepath.Join(dir, "out.jpg")

	err := ImageToImage("nonexistent_file.png", output)
	if err == nil {
		t.Error("expected error for nonexistent input file, got nil")
	}
}

// -----------------------------------------
// ImageToPDF
// -----------------------------------------

func TestImageToPDF_SingleImage(t *testing.T) {
	dir := t.TempDir()
	input := createTestPNG(t, dir, "test.png")
	output := filepath.Join(dir, "out.pdf")

	if err := ImageToPDF([]string{input}, output); err != nil {
		t.Fatalf("ImageToPDF() error: %v", err)
	}

	info, err := os.Stat(output)
	if os.IsNotExist(err) {
		t.Fatal("PDF was not created")
	}
	if info.Size() == 0 {
		t.Error("PDF file is empty")
	}
}

func TestImageToPDF_MultipleImages(t *testing.T) {
	dir := t.TempDir()
	input1 := createTestPNG(t, dir, "test1.png")
	input2 := createTestPNG(t, dir, "test2.png")
	output := filepath.Join(dir, "out.pdf")

	if err := ImageToPDF([]string{input1, input2}, output); err != nil {
		t.Fatalf("ImageToPDF() error with multiple images: %v", err)
	}

	if _, err := os.Stat(output); os.IsNotExist(err) {
		t.Fatal("PDF was not created")
	}
}

func TestImageToPDF_EmptyList(t *testing.T) {
	dir := t.TempDir()
	output := filepath.Join(dir, "out.pdf")

	// an empty slice must return an error
	err := ImageToPDF([]string{}, output)
	if err == nil {
		t.Error("expected error for empty image list, got nil")
	}
}

func TestImageToPDF_InvalidInput(t *testing.T) {
	dir := t.TempDir()
	output := filepath.Join(dir, "out.pdf")

	err := ImageToPDF([]string{"nonexistent.png"}, output)
	if err == nil {
		t.Error("expected error for nonexistent input file, got nil")
	}
}
