package handlers

import (
	"errors"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func TestDeleteFileRejectsInvalidFilenameBeforeRPC(t *testing.T) {
	e := echo.New()

	tests := []string{
		".",
		"..",
		"../bad.png",
		"/bad.png",
		`C:\bad.png`,
		"nested/bad.png",
	}
	for _, filename := range tests {
		t.Run(filename, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodDelete, "/api/files?title="+filename, nil)
			rec := httptest.NewRecorder()

			err := (&Apis{}).DeleteFile(e.NewContext(req, rec))

			require.NoError(t, err)
			require.Equal(t, http.StatusBadRequest, rec.Code)
			require.JSONEq(t, `{"error":"invalid filename"}`, rec.Body.String())
		})
	}
}

func TestValidateImageFilenameRejectsDirectoryLikeInput(t *testing.T) {
	tests := []string{
		"",
		".",
		"..",
		"../x",
		"/x",
		`C:\x`,
		"x/y",
		`x\y`,
	}
	for _, filename := range tests {
		t.Run(filename, func(t *testing.T) {
			require.ErrorIs(t, validateImageFilename(filename), ErrFileName)
		})
	}
}

func TestUploadSaveErrorResponse(t *testing.T) {
	t.Parallel()

	code, msg := uploadSaveErrorResponse(errTooLarge)
	require.Equal(t, http.StatusRequestEntityTooLarge, code)
	require.Equal(t, "file is too large", msg)

	code, msg = uploadSaveErrorResponse(errors.New("permission denied"))
	require.Equal(t, http.StatusInternalServerError, code)
	require.Equal(t, "failed to save file", msg)
}

func TestEnsureImageFileExists(t *testing.T) {
	tmp := t.TempDir()
	t.Chdir(tmp)
	require.NoError(t, os.Mkdir(ImagesDir, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(ImagesDir, "ok.png"), []byte("ok"), 0644))
	require.NoError(t, os.Mkdir(filepath.Join(ImagesDir, "dir.png"), 0755))

	require.NoError(t, ensureImageFileExists("ok.png"))
	require.ErrorIs(t, ensureImageFileExists("missing.png"), ErrReferencedFileNotFound)
	require.ErrorIs(t, ensureImageFileExists("."), ErrFileName)
	require.ErrorIs(t, ensureImageFileExists(".."), ErrFileName)
	require.ErrorIs(t, ensureImageFileExists("dir.png"), ErrReferencedFileNotFound)

	err := ensureImageFileExistsWithStat("ok.png", func(string) (os.FileInfo, error) {
		return nil, fs.ErrPermission
	})
	require.ErrorIs(t, err, ErrImageStorage)
}

func TestCreateCategoryRejectsMissingImageBeforeRPC(t *testing.T) {
	tmp := t.TempDir()
	t.Chdir(tmp)
	require.NoError(t, os.Mkdir(ImagesDir, 0755))

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/category", strings.NewReader(`{"title":"chairs","img":"missing.png"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	err := (&Apis{}).CreateCategory(e.NewContext(req, rec))

	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, rec.Code)
	require.JSONEq(t, `{"error":"referenced file does not exist"}`, rec.Body.String())
}

func TestCreateProductRejectsOneMissingPhotoBeforeRPC(t *testing.T) {
	tmp := t.TempDir()
	t.Chdir(tmp)
	require.NoError(t, os.Mkdir(ImagesDir, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(ImagesDir, "ok.png"), []byte("ok"), 0644))

	e := echo.New()
	body := `{"title":"table","photos":["ok.png","missing.png"]}`
	req := httptest.NewRequest(http.MethodPost, "/api/product", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	err := (&Apis{}).CreateProduct(e.NewContext(req, rec))

	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, rec.Code)
	require.JSONEq(t, `{"error":"referenced file does not exist"}`, rec.Body.String())
}

func TestCreateSlideRejectsMissingSecondImageBeforeRPC(t *testing.T) {
	tmp := t.TempDir()
	t.Chdir(tmp)
	require.NoError(t, os.Mkdir(ImagesDir, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(ImagesDir, "desktop.png"), []byte("ok"), 0644))

	e := echo.New()
	body := `{"link":"/","img":"desktop.png","img762":"missing.png"}`
	req := httptest.NewRequest(http.MethodPost, "/api/slide", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	err := (&Apis{}).CreateSlide(e.NewContext(req, rec))

	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, rec.Code)
	require.JSONEq(t, `{"error":"referenced file does not exist"}`, rec.Body.String())
}
