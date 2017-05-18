/**
 * go-mapbox Maps Module Cache
 * Provides a simple file system cache to avoid hitting the maps API quite so regularly
 *
 * https://github.com/ryankurte/go-mapbox
 * Copyright 2017 Ryan Kurte
 */

package maps

import (
	"fmt"
	"image"
	"os"
	"strings"
)

// FileCache is a simple file-based caching implementation for map tiles
// This does not implement any mechanisms for deletion / removal, and as such is not suitable for production use
type FileCache struct {
	basePath string
}

// NewFileCache creates a new file cache instance
func NewFileCache(basePath string) (*FileCache, error) {
	fc := &FileCache{basePath}

	err := os.Mkdir(basePath, 0777)
	if err != nil && !os.IsExist(err) {
		return nil, err
	}

	return fc, nil
}

func (fc *FileCache) getName(mapID MapID, x, y, level uint64, format MapFormat, highDPI bool) string {
	dpiString := ""
	if highDPI {
		dpiString = "@2x"
	}
	return fmt.Sprintf("%s-%d-%d-%d%s.%s", mapID, x, y, level, dpiString, format)
}

// Save saves an image to the file cache
func (fc *FileCache) Save(mapID MapID, x, y, level uint64, format MapFormat, highDPI bool, img image.Image) error {
	name := fc.getName(mapID, x, y, level, format, highDPI)
	path := fmt.Sprintf("%s/%s", fc.basePath, name)

	// Check if file already exists
	if _, err := os.Stat(path); err == nil {
		return nil
	}

	// TODO: cannot currently save pngraw
	if strings.Contains(string(format), "png") && strings.Contains(string(format), "pngraw") {
		return SaveImagePNG(img, path)
	}

	if strings.Contains(string(format), "jpg") || strings.Contains(string(format), "jpeg") {
		return SaveImageJPG(img, path)
	}

	return fmt.Errorf("Unrecognized file type (%s)", format)
}

// Fetch fetches an image from the file cache if possible
func (fc *FileCache) Fetch(mapID MapID, x, y, level uint64, format MapFormat, highDPI bool) (image.Image, *image.Config, error) {
	name := fc.getName(mapID, x, y, level, format, highDPI)
	path := fmt.Sprintf("%s/%s", fc.basePath, name)

	if format == MapFormatPngRaw {
		return nil, nil, nil
	}

	if _, err := os.Stat(path); err != nil {
		return nil, nil, nil
	}

	img, cfg, err := LoadImage(path)

	return img, cfg, err
}
