package parliament

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"time"
)

type Storage interface {
	Store(ctx context.Context, stenogram Stenogram) error
}

type FileStorage struct {
	baseLoc string
}

func NewFileStorage(loc string) *FileStorage {
	return &FileStorage{baseLoc: loc}
}

func (fs *FileStorage) Store(ctx context.Context, stenogram Stenogram) error {
	date, err := time.Parse("2006-01-02", stenogram.Date)
	if err != nil {
		return err
	}
	loc := path.Join(fs.baseLoc, fmt.Sprintf("%d/%d/%d", date.Year(), date.Month(), date.Day()))
	err = os.MkdirAll(loc, os.ModeDir)
	if err != nil {
		return err
	}

	f, err := os.Create(path.Join(loc, fmt.Sprintf("%d.json", stenogram.ID)))
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "\t")
	err = enc.Encode(stenogram)
	if err != nil {
		return err
	}

	for fileLoc, data := range stenogram.Files {
		_, name := path.Split(fileLoc)
		f, err := os.Create(path.Join(loc, name))
		if err != nil {
			return err
		}
		f.Write(data)
		f.Close()
	}

	return nil
}
