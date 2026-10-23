package handlers

import (
	"context"
	"errors"

	productsRPC "github.com/autumnterror/volha-backend/api/proto/gen"
	log "github.com/autumnterror/volha-backend/internal/platform/logger"

	"fmt"
	"io"
	"io/fs"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/autumnterror/volha-backend/internal/gateway/net/view"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

var allowedMIMEs = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
	"image/webp": ".webp",
	"image/gif":  ".gif",
}

const (
	ImagesDir      = "./images"
	PDFsDir        = "./pdfs"
	MaxUploadBytes = 10 << 20 // 10 MB
	FileCleanupTTL = 2 * time.Second
)

// UploadFile godoc
// @Summary Загрузить изображение
// @Description Сохраняет изображение в папке сервера images/. В поле "img" передается файл с расширениями jpg, png, webp, gif
// @Tags files
// @Produce json
// @Success 200 {object} view.FileUploadResponse "Цвет успешно создан"
// @Failure 400 {object} view.ErrorResponse "Неверный формат данных"
// @Failure 413 {object} view.ErrorResponse "Файл слишком большой"
// @Failure 415 {object} view.ErrorResponse "Неподдерживаемый тип файла"
// @Failure 500 {object} view.ErrorResponse "Ошибка на сервере"
// @Router /api/files [post]
func (a *Apis) UploadFile(c echo.Context) error {
	const op = "handlers.UploadFile"
	fileHeader, err := c.FormFile("file")
	if err != nil {
		fileHeader, err = c.FormFile("img")
	}
	if err != nil {
		return c.JSON(http.StatusBadRequest, view.ErrorResponse{Error: "file field 'file' or 'img' is required"})
	}

	src, err := fileHeader.Open()
	if err != nil {
		log.Error(c.Request().Context(), "open uploaded file failed", "operation", op, "error", err)
		return c.JSON(http.StatusInternalServerError, view.ErrorResponse{Error: "failed to read uploaded file"})
	}
	defer src.Close()

	buff := make([]byte, 512)
	n, err := io.ReadFull(src, buff)
	if err != nil && !errors.Is(err, io.EOF) && !errors.Is(err, io.ErrUnexpectedEOF) {
		log.Error(c.Request().Context(), "read uploaded file header failed", "operation", op, "error", err)
		return c.JSON(http.StatusInternalServerError, view.ErrorResponse{Error: "failed to read uploaded file"})
	}
	if _, err := src.Seek(0, io.SeekStart); err != nil {
		log.Error(c.Request().Context(), "seek uploaded file failed", "operation", op, "error", err)
		return c.JSON(http.StatusInternalServerError, view.ErrorResponse{Error: "cannot seek file"})
	}

	detected := http.DetectContentType(buff[:n])
	ext, ok := allowedMIMEs[detected]
	if detected == "application/pdf" {
		ext = ".pdf"
		ok = true
	}
	if !ok {
		return c.JSON(http.StatusUnsupportedMediaType, view.ErrorResponse{Error: "unsupported file type"})
	}

	filename := uuid.NewString() + ext
	dstPath, err := storedFilePath(filename)
	if err != nil {
		log.Error(c.Request().Context(), "resolve upload path failed", "operation", op, "filename", filename, "error", err)
		return c.JSON(http.StatusInternalServerError, view.ErrorResponse{Error: "failed to save file"})
	}

	if err := saveMultipartFile(src, dstPath, fileHeader); err != nil {
		code, msg := uploadSaveErrorResponse(err)
		if code == http.StatusInternalServerError {
			log.Error(c.Request().Context(), "save uploaded file failed", "operation", op, "filename", filename, "error", err)
		}
		return c.JSON(code, view.ErrorResponse{Error: msg})
	}

	return c.JSON(http.StatusOK, view.FileUploadResponse{
		Name: filename,
	})
}

var errTooLarge = errors.New("too large")

func uploadSaveErrorResponse(err error) (int, string) {
	if errors.Is(err, errTooLarge) {
		return http.StatusRequestEntityTooLarge, "file is too large"
	}
	return http.StatusInternalServerError, "failed to save file"
}

func saveMultipartFile(src multipart.File, dstPath string, fh *multipart.FileHeader) error {
	if fh.Size > MaxUploadBytes {
		return errTooLarge
	}
	if err := os.MkdirAll(filepath.Dir(dstPath), 0755); err != nil {
		return err
	}

	dst, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dst.Close()
	if _, err := io.Copy(dst, src); err != nil {
		_ = os.Remove(dstPath)
		return err
	}

	if err := os.Chmod(dstPath, 0644); err != nil {
		_ = os.Remove(dstPath)
		return err
	}

	return nil
}

// DeleteFile godoc
// @Summary Удалить файл
// @Description Удаляет файл. Требуется передовать только имя. Пример: example.gif
// @Tags files
// @Produce json
// @Param title query string true "название файла"
// @Success 200 {object} view.MessageResponse "файл успешно удалён"
// @Failure 400 {object} view.ErrorResponse "Неверное название"
// @Failure 500 {object} view.ErrorResponse "Ошибка на сервере"
// @Router /api/files [delete]
func (a *Apis) DeleteFile(c echo.Context) error {
	const op = "handlers.DeleteFile"

	filename := c.QueryParam("title")
	if filename == "" {
		return c.JSON(http.StatusBadRequest, view.ErrorResponse{Error: "empty filename"})
	}
	if err := validateImageFilename(filename); err != nil {
		return c.JSON(http.StatusBadRequest, view.ErrorResponse{Error: "invalid filename"})
	}

	ctx, cancel := requestContext(c)
	defer cancel()

	referenced, err := a.fileReferenced(ctx, filename)
	if err != nil {
		code, swg := a.mapGRPCError(ctx, op, err)
		return c.JSON(code, swg)
	}
	if referenced {
		return c.JSON(http.StatusConflict, view.ErrorResponse{Error: "file is referenced"})
	}

	if err := deleteFile(filename); err != nil {
		switch err {
		case ErrFileName:
			return c.JSON(http.StatusBadRequest, view.ErrorResponse{Error: "invalid filename"})
		case ErrFileNotFound:
			return c.JSON(http.StatusNotFound, view.ErrorResponse{Error: "file not found"})
		default:
			log.Error(ctx, "delete file failed", "operation", op, "filename", filename, "error", err)
			return c.JSON(http.StatusInternalServerError, view.ErrorResponse{Error: "check logs"})
		}
	}

	return c.JSON(http.StatusOK, view.MessageResponse{Message: "file delete successfully"})
}

var ErrFileName = errors.New("invalid filename")
var ErrFileNotFound = errors.New("file not found")
var ErrReferencedFileNotFound = errors.New("referenced file does not exist")
var ErrImageStorage = errors.New("image storage error")

func deleteFile(filename string) error {
	fullPath, err := storedFilePath(filename)
	if err != nil {
		return err
	}

	err = os.Remove(fullPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return ErrFileNotFound
		}
		return err
	}
	return nil
}

func imageFilePath(filename string) (string, error) {
	if err := validateImageFilename(filename); err != nil {
		return "", err
	}

	baseDir, err := filepath.Abs(ImagesDir)
	if err != nil {
		return "", err
	}
	baseDir = filepath.Clean(baseDir)
	fullPath, err := filepath.Abs(filepath.Join(baseDir, filename))
	if err != nil {
		return "", err
	}
	fullPath = filepath.Clean(fullPath)
	if fullPath == baseDir {
		return "", ErrFileName
	}
	rel, err := filepath.Rel(baseDir, fullPath)
	if err != nil {
		return "", err
	}
	if rel == ".." || strings.HasPrefix(rel, ".."+string(os.PathSeparator)) || filepath.IsAbs(rel) {
		return "", ErrFileName
	}

	return fullPath, nil
}

func validateImageFilename(filename string) error {
	if filename == "" || filename == "." || filename == ".." {
		return ErrFileName
	}
	if filepath.Clean(filename) != filename || filepath.IsAbs(filename) || strings.ContainsAny(filename, `/\:`) || filepath.Base(filename) != filename {
		return ErrFileName
	}
	return nil
}

func ensureImageReferences(filenames ...string) error {
	checked := make(map[string]struct{}, len(filenames))
	for _, filename := range filenames {
		if filename == "" {
			continue
		}
		if _, ok := checked[filename]; ok {
			continue
		}
		checked[filename] = struct{}{}
		if err := ensureImageFileExists(filename); err != nil {
			return err
		}
	}
	return nil
}

func ensureImageFileExists(filename string) error {
	return ensureImageFileExistsWithStat(filename, os.Stat)
}

func ensureImageFileExistsWithStat(filename string, stat func(string) (os.FileInfo, error)) error {
	path, err := imageFilePath(filename)
	if err != nil {
		return err
	}
	info, err := stat(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return fmt.Errorf("%w: %s", ErrReferencedFileNotFound, filename)
		}
		return fmt.Errorf("%w: %v", ErrImageStorage, err)
	}
	if info.IsDir() {
		return fmt.Errorf("%w: %s", ErrReferencedFileNotFound, filename)
	}
	return nil
}

func imageReferenceErrorResponse(ctx context.Context, op string, err error) (int, view.ErrorResponse) {
	switch {
	case err == nil:
		return http.StatusOK, view.ErrorResponse{}
	case errors.Is(err, ErrFileName):
		return http.StatusBadRequest, view.ErrorResponse{Error: "invalid filename"}
	case errors.Is(err, ErrReferencedFileNotFound):
		return http.StatusBadRequest, view.ErrorResponse{Error: "referenced file does not exist"}
	default:
		log.Error(ctx, "image reference validation failed", "operation", op, "error", err)
		return http.StatusInternalServerError, view.ErrorResponse{Error: "failed to validate image references"}
	}
}

func (a *Apis) fileReferenced(ctx context.Context, filename string) (bool, error) {
	res, err := a.apiProduct.API.IsFileReferenced(ctx, &productsRPC.FileReferenceRequest{Filename: filename})
	if err != nil {
		return false, err
	}

	return res.GetReferenced(), nil
}

func removedStrings(oldItems, newItems []string) []string {
	next := make(map[string]struct{}, len(newItems))
	for _, item := range newItems {
		next[item] = struct{}{}
	}
	removed := make([]string, 0)
	for _, item := range oldItems {
		if _, ok := next[item]; !ok {
			removed = append(removed, item)
		}
	}
	return removed
}

func (a *Apis) cleanupFileIfUnreferenced(ctx context.Context, op, filename string) {
	if filename == "" {
		return
	}
	referenced, err := a.fileReferenced(ctx, filename)
	if err != nil {
		log.Error(ctx, "check file reference failed", "operation", op, "filename", filename, "error", err)
		return
	}
	if referenced {
		return
	}
	if err := deleteFile(filename); err != nil && !errors.Is(err, ErrFileNotFound) {
		log.Error(ctx, "delete unreferenced file failed", "operation", op, "filename", filename, "error", err)
	}
}

func fileCleanupContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), FileCleanupTTL)
}

func storedFilePath(filename string) (string, error) {
	if strings.HasSuffix(strings.ToLower(filename), ".pdf") {
		return pdfFilePath(filename)
	}
	return imageFilePath(filename)
}
func pdfFilePath(filename string) (string, error) {
	if err := validateImageFilename(filename); err != nil {
		return "", err
	}
	if !strings.HasSuffix(filename, ".pdf") {
		return "", ErrFileName
	}
	return filepath.Abs(filepath.Join(PDFsDir, filename))
}
func normalizePDFLink(link string) (string, error) {
	name := strings.TrimPrefix(link, "/pdfs/")
	if _, err := pdfFilePath(name); err != nil {
		return "", err
	}
	return "/pdfs/" + name, nil
}
func ensurePDFReference(link string) error {
	normalized, err := normalizePDFLink(link)
	if err != nil {
		return err
	}
	path, err := pdfFilePath(strings.TrimPrefix(normalized, "/pdfs/"))
	if err != nil {
		return err
	}
	info, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		return ErrReferencedFileNotFound
	}
	if err != nil {
		return fmt.Errorf("%w: %v", ErrImageStorage, err)
	}
	if info.IsDir() {
		return ErrReferencedFileNotFound
	}
	return nil
}
