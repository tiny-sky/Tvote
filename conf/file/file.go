package file

import (
	"io"
	"os"

	"github.com/tiny-sky/Tvote/conf"
	"gopkg.in/yaml.v3"
)

type File struct {
	path string
	*conf.Settings
}

func NewFile(path string) *File {
	return &File{path: path}
}

func (f *File) Load() (*conf.Settings, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	file, err := os.Open(dir + f.path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	byteAll, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(byteAll, &f.Settings)
	return f.Settings, err
}
