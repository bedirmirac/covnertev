package image

import (
	stdimage "image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"testing"
)

// createTestPNG yardımcı fonksiyonu: t.TempDir() içine küçük bir PNG oluşturur.
func createTestPNG(t *testing.T, dir, name string) string {
	t.Helper()

	path := filepath.Join(dir, name)
	f, err := os.Create(path)
	if err != nil {
		t.Fatalf("test PNG oluşturulamadı: %v", err)
	}
	defer f.Close()

	// 10x10 kırmızı kare
	img := stdimage.NewRGBA(stdimage.Rect(0, 0, 10, 10))
	for y := range 10 {
		for x := range 10 {
			img.Set(x, y, color.RGBA{R: 200, G: 100, B: 50, A: 255})
		}
	}
	if err := png.Encode(f, img); err != nil {
		t.Fatalf("PNG encode hatası: %v", err)
	}
	return path
}

// ─────────────────────────────────────────
// ImageToImage
// ─────────────────────────────────────────

func TestImageToImage_PNGtoJPG(t *testing.T) {
	dir := t.TempDir()
	input := createTestPNG(t, dir, "test.png")
	output := filepath.Join(dir, "out.jpg")

	if err := ImageToImage(input, output); err != nil {
		t.Fatalf("ImageToImage() hata: %v", err)
	}

	info, err := os.Stat(output)
	if os.IsNotExist(err) {
		t.Fatal("çıktı dosyası oluşturulmadı")
	}
	if info.Size() == 0 {
		t.Error("çıktı dosyası boş")
	}
}

func TestImageToImage_PNGtoPNG(t *testing.T) {
	dir := t.TempDir()
	input := createTestPNG(t, dir, "test.png")
	output := filepath.Join(dir, "out.png")

	if err := ImageToImage(input, output); err != nil {
		t.Fatalf("ImageToImage() hata: %v", err)
	}

	if _, err := os.Stat(output); os.IsNotExist(err) {
		t.Fatal("çıktı dosyası oluşturulmadı")
	}
}

// TestImageToImage_PNGtoWEBP_WriteUnsupported — disintegration/imaging kütüphanesi
// WEBP okuyabiliyor (golang.org/x/image/webp decoder ile) ama WEBP yazamıyor.
// Bu test mevcut kısıtlamayı belgeler: WEBP hedef format olarak desteklenmiyor.
func TestImageToImage_PNGtoWEBP_WriteUnsupported(t *testing.T) {
	dir := t.TempDir()
	input := createTestPNG(t, dir, "test.png")
	output := filepath.Join(dir, "out.webp")

	err := ImageToImage(input, output)
	if err == nil {
		t.Error("WEBP yazma desteklenmiyor, hata bekleniyor — kütüphane güncellenmiş olabilir")
	} else {
		t.Logf("Beklenen hata (WEBP write desteksiz): %v", err)
	}
}

func TestImageToImage_InvalidInput(t *testing.T) {
	dir := t.TempDir()
	output := filepath.Join(dir, "out.jpg")

	err := ImageToImage("olmayan_dosya.png", output)
	if err == nil {
		t.Error("olmayan dosya için hata bekleniyor, nil döndü")
	}
}

// ─────────────────────────────────────────
// ImageToPDF
// ─────────────────────────────────────────

func TestImageToPDF_SingleImage(t *testing.T) {
	dir := t.TempDir()
	input := createTestPNG(t, dir, "test.png")
	output := filepath.Join(dir, "out.pdf")

	if err := ImageToPDF([]string{input}, output); err != nil {
		t.Fatalf("ImageToPDF() hata: %v", err)
	}

	info, err := os.Stat(output)
	if os.IsNotExist(err) {
		t.Fatal("PDF oluşturulmadı")
	}
	if info.Size() == 0 {
		t.Error("PDF dosyası boş")
	}
}

func TestImageToPDF_MultipleImages(t *testing.T) {
	dir := t.TempDir()
	input1 := createTestPNG(t, dir, "test1.png")
	input2 := createTestPNG(t, dir, "test2.png")
	output := filepath.Join(dir, "out.pdf")

	if err := ImageToPDF([]string{input1, input2}, output); err != nil {
		t.Fatalf("ImageToPDF() birden fazla görüntü için hata: %v", err)
	}

	if _, err := os.Stat(output); os.IsNotExist(err) {
		t.Fatal("PDF oluşturulmadı")
	}
}

func TestImageToPDF_EmptyList(t *testing.T) {
	dir := t.TempDir()
	output := filepath.Join(dir, "out.pdf")

	// Boş liste ile çağrılınca hata beklenir
	err := ImageToPDF([]string{}, output)
	if err == nil {
		t.Error("boş görüntü listesi için hata bekleniyor, nil döndü")
	}
}

func TestImageToPDF_InvalidInput(t *testing.T) {
	dir := t.TempDir()
	output := filepath.Join(dir, "out.pdf")

	err := ImageToPDF([]string{"olmayan.png"}, output)
	if err == nil {
		t.Error("olmayan dosya için hata bekleniyor, nil döndü")
	}
}
