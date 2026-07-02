package document

import (
	"os"
	"testing"
)

// -----------------------------------------
// getLibreOfficePath - priority order
// -----------------------------------------

func TestGetLibreOfficePath_CustomFlag(t *testing.T) {
	// set the package variable directly, which has the same effect as the flag
	original := customLibreOfficePath
	customLibreOfficePath = "/custom/soffice"
	defer func() { customLibreOfficePath = original }()

	// env var is also set; --libreoffice-path flag must take priority
	t.Setenv("LIBREOFFICE_PATH", "/env/soffice")

	path, err := getLibreOfficePath()
	if err != nil {
		t.Fatalf("getLibreOfficePath() unexpected error: %v", err)
	}
	if path != "/custom/soffice" {
		t.Errorf("getLibreOfficePath() = %q, expected %q (flag must take priority)", path, "/custom/soffice")
	}
}

func TestGetLibreOfficePath_EnvVar(t *testing.T) {
	// flag is clear; only env var is set
	original := customLibreOfficePath
	customLibreOfficePath = ""
	defer func() { customLibreOfficePath = original }()

	fakeLibreOfficePath := "/env/path/to/soffice"
	t.Setenv("LIBREOFFICE_PATH", fakeLibreOfficePath)

	path, err := getLibreOfficePath()
	if err != nil {
		t.Fatalf("getLibreOfficePath() unexpected error: %v", err)
	}
	if path != fakeLibreOfficePath {
		t.Errorf("getLibreOfficePath() = %q, expected %q", path, fakeLibreOfficePath)
	}
}

func TestGetLibreOfficePath_NotFound(t *testing.T) {
	// no source is set; expect an error if LibreOffice is not installed
	original := customLibreOfficePath
	customLibreOfficePath = ""
	defer func() { customLibreOfficePath = original }()

	os.Unsetenv("LIBREOFFICE_PATH")

	path, err := getLibreOfficePath()
	if err != nil {
		// expected: LibreOffice is not installed
		t.Logf("getLibreOfficePath() returned expected error: %v", err)
		return
	}
	// LibreOffice is present on this system
	t.Logf("LibreOffice found on system, skipping not-found check: %s", path)
}

// -----------------------------------------
// PdfFromOffice - error path without LibreOffice
// -----------------------------------------

func TestPdfFromOffice_LibreOfficeNotFound(t *testing.T) {
	// point LibreOffice to a path that does not exist
	original := customLibreOfficePath
	customLibreOfficePath = "/nonexistent/path/to/soffice"
	defer func() { customLibreOfficePath = original }()

	// exec.Command will fail because the executable does not exist,
	// and the error is now propagated correctly via return err.
	err := PdfFromOffice("input.docx")
	if err == nil {
		t.Error("expected error for invalid LibreOffice path, got nil")
	}
}

func TestDocFromPdf_LibreOfficeNotFound(t *testing.T) {
	original := customLibreOfficePath
	customLibreOfficePath = "/nonexistent/path/to/soffice"
	defer func() { customLibreOfficePath = original }()

	err := DocFromPdf("input.pdf")
	if err == nil {
		t.Error("expected error for invalid LibreOffice path, got nil")
	}
}
