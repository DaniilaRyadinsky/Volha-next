package service

import (
	"context"
	"errors"
	"path/filepath"
	"strings"
)

func (s *ProductsService) IsFileReferenced(ctx context.Context, filename string) (bool, error) {
	const op = "service.IsFileReferenced"
	if err := validateFilename(filename); err != nil {
		return false, wrapServiceCheck(op, err)
	}
	repo, err := s.fileReferenceRepo(ctx)
	if err != nil {
		return false, wrapServiceCheck(op, err)
	}

	return repo.IsFileReferenced(ctx, filename)
}

func validateFilename(filename string) error {
	if strings.TrimSpace(filename) == "" {
		return errors.New("filename is empty")
	}
	if strings.ContainsAny(filename, `/\`) || filepath.Base(filename) != filename {
		return errors.New("filename must be a base name")
	}
	return nil
}
