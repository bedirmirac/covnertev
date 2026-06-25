package document

import (
	"os"
	"testing"
)

// ─────────────────────────────────────────
// getLibreOfficePath — öncelik sırası
// ─────────────────────────────────────────

func TestGetLibreOfficePath_CustomFlag(t *testing.T) {
	// Paket değişkenini doğrudan set et (flag ile aynı etki)
	original := customLibreOfficePath
	customLibreOfficePath = "/custom/soffice"
	defer func() { customLibreOfficePath = original }()

	// Env var da set edilmiş olsa bile --libreoffice-path flag öncelikli
	t.Setenv("LIBREOFFICE_PATH", "/env/soffice")

	path, err := getLibreOfficePath()
	if err != nil {
		t.Fatalf("getLibreOfficePath() beklenmedik hata: %v", err)
	}
	if path != "/custom/soffice" {
		t.Errorf("getLibreOfficePath() = %q, beklenen %q (flag öncelikli olmalı)", path, "/custom/soffice")
	}
}

func TestGetLibreOfficePath_EnvVar(t *testing.T) {
	// Flag temiz, sadece env var
	original := customLibreOfficePath
	customLibreOfficePath = ""
	defer func() { customLibreOfficePath = original }()

	fakeLibreOfficePath := "/env/path/to/soffice"
	t.Setenv("LIBREOFFICE_PATH", fakeLibreOfficePath)

	path, err := getLibreOfficePath()
	if err != nil {
		t.Fatalf("getLibreOfficePath() beklenmedik hata: %v", err)
	}
	if path != fakeLibreOfficePath {
		t.Errorf("getLibreOfficePath() = %q, beklenen %q", path, fakeLibreOfficePath)
	}
}

func TestGetLibreOfficePath_NotFound(t *testing.T) {
	// Hiçbir kaynak set edilmemiş; LibreOffice kurulu değilse hata döner
	original := customLibreOfficePath
	customLibreOfficePath = ""
	defer func() { customLibreOfficePath = original }()

	os.Unsetenv("LIBREOFFICE_PATH")

	path, err := getLibreOfficePath()
	if err != nil {
		// Beklenen durum: kurulu değil
		t.Logf("getLibreOfficePath() hata döndü (bekleniyor): %v", err)
		return
	}
	// LibreOffice sistemde kurulu — test ortamına göre geçer
	t.Logf("LibreOffice sistemde mevcut, atlıyor: %s", path)
}

// ─────────────────────────────────────────
// PdfFromOffice — LibreOffice gerekmeden hata yolu
// ─────────────────────────────────────────

func TestPdfFromOffice_LibreOfficeNotFound(t *testing.T) {
	// LibreOffice'i kasıtlı olarak geçersiz bir yola yönlendir
	original := customLibreOfficePath
	customLibreOfficePath = "/nonexistent/path/to/soffice"
	defer func() { customLibreOfficePath = original }()

	// Geçersiz yol verildiğinde exec.Command hata verecek
	err := PdfFromOffice("input.docx")
	// NOT: mevcut kod log.Fatal çağırıyor — bu testin geçmesi için
	// fonksiyonun error döndürmesi gerekir (Bug #2'nin düzeltilmiş hali).
	// Şu anki haliyle log.Fatal çağrısı process'i öldürür ve test paniklenir.
	// Bu test Bug #2 düzeltilince düzgün çalışacak.
	if err == nil {
		t.Error("geçersiz LibreOffice yolu için hata bekleniyor, nil döndü")
	}
}

func TestDocFromPdf_LibreOfficeNotFound(t *testing.T) {
	original := customLibreOfficePath
	customLibreOfficePath = "/nonexistent/path/to/soffice"
	defer func() { customLibreOfficePath = original }()

	err := DocFromPdf("input.pdf")
	// Aynı şekilde Bug #2 düzeltilince bu test düzgün çalışacak
	if err == nil {
		t.Error("geçersiz LibreOffice yolu için hata bekleniyor, nil döndü")
	}
}
