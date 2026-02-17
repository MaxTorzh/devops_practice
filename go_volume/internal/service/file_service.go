package service

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"go_volume/internal/models"
)

type FileService struct {
	baseDir string
}

func NewFileService(baseDir string) *FileService {
	os.MkdirAll(baseDir, 0755)
	return &FileService{baseDir: baseDir}
}

func (s *FileService) validateFilename(filename string) error {
	if filename == "" {
		return fmt.Errorf("filename cannot be empty")
	}
	if strings.Contains(filename, "..") || strings.Contains(filename, "/") || strings.Contains(filename, "\\") {
		return fmt.Errorf("invalid filename: %s", filename)
	}
	return nil
}

func (s *FileService) safePath(filename string) (string, error) {
	if err := s.validateFilename(filename); err != nil {
		return "", err
	}
	return filepath.Join(s.baseDir, filename), nil
}

func (s *FileService) fileExists(path string) (os.FileInfo, error) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("file not found")
		}
		return nil, fmt.Errorf("failed to access file: %w", err)
	}
	return info, nil
}

func (s *FileService) ensureIsFile(info os.FileInfo, path string) error {
	if info.IsDir() {
		return fmt.Errorf("path is a directory: %s", path)
	}
	return nil
}

func (s *FileService) ListFiles() ([]models.FileInfo, error) {
	entries, err := os.ReadDir(s.baseDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	var files []models.FileInfo
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}
		files = append(files, models.FromFileInfo(entry.Name(), info))
	}
	return files, nil
}

func (s *FileService) GetFile(filename string) ([]byte, error) {
	path, err := s.safePath(filename)
	if err != nil {
		return nil, err
	}

	info, err := s.fileExists(path)
	if err != nil {
		return nil, err
	}

	if err := s.ensureIsFile(info, path); err != nil {
		return nil, err
	}

	return os.ReadFile(path)
}

func (s *FileService) CreateFile(filename string, content []byte) error {
	path, err := s.safePath(filename)
	if err != nil {
		return err
	}

	if _, err := os.Stat(path); err == nil {
		return fmt.Errorf("file already exists: %s", filename)
	}

	return os.WriteFile(path, content, 0644)
}

func (s *FileService) UpdateFile(filename string, content []byte) error {
	path, err := s.safePath(filename)
	if err != nil {
		return err
	}

	info, err := s.fileExists(path)
	if err != nil {
		return err
	}

	if err := s.ensureIsFile(info, path); err != nil {
		return err
	}

	return os.WriteFile(path, content, 0644)
}

func (s *FileService) DeleteFile(filename string) error {
	path, err := s.safePath(filename)
	if err != nil {
		return err
	}

	info, err := s.fileExists(path)
	if err != nil {
		return err
	}

	if err := s.ensureIsFile(info, path); err != nil {
		return err
	}

	return os.Remove(path)
}

func (s *FileService) GetFileInfo(filename string) (*models.FileInfo, error) {
	path, err := s.safePath(filename)
	if err != nil {
		return nil, err
	}

	info, err := s.fileExists(path)
	if err != nil {
		return nil, err
	}

	fileInfo := models.FromFileInfo(filename, info)
	return &fileInfo, nil
}
