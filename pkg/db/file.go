package db

import (
	"fmt"
	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/properties"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type FileDB struct {
	mu      sync.RWMutex
	baseDir string
	tmpDir  string
}

func (db *FileDB) getData(ns, sparrow, fileName string) *Properties {
	db.mu.RLock()
	defer db.mu.RUnlock()
	property := newProperties(ns, sparrow, fileName)

	filePath := filepath.Join(db.baseDir, ns, sparrow, fileName)
	con := config.New("default")
	con.WithOptions(config.ParseEnv)
	con.AddDriver(properties.Driver)
	err := con.LoadFilesByFormat("properties", filePath)
	if err != nil {
		return property
	}
	property.SetAll(con.Data())
	return property
}

func (db *FileDB) updateData(data *Properties) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	filePath := filepath.Join(db.baseDir, data.ns, data.sparrow, data.fileName)
	temp, err := os.CreateTemp(db.tmpDir, filepath.Base(data.sparrow))
	if err != nil {
		return err
	}
	tmpPath := temp.Name()
	defer temp.Close()
	defer os.Remove(tmpPath)

	for key, value := range data.data {
		_, err = temp.WriteString(fmt.Sprintf("%s=%v\n", key, value))
		if err != nil {
			return fmt.Errorf("error writing to file: %v", err)
		}
	}

	if err = temp.Sync(); err != nil {
		return err
	}
	if err = os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return err
	}
	return os.Rename(tmpPath, filePath)
}

func newFileDB(config *dBConfig) (*FileDB, error) {
	path := strings.ReplaceAll(config.path, "$HOME", os.Getenv("HOME"))
	tmpDir := filepath.Join(path, ".sparrow")
	if err := os.RemoveAll(tmpDir); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		return nil, err
	}

	return &FileDB{
		baseDir: path,
		tmpDir:  tmpDir,
	}, nil
}
