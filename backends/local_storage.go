package backends

import (
	"io"
	"os"
	"regexp"
)

type LocalStorage struct {
	path string
}

func NewLocalStorage(path string) *LocalStorage {
	return &LocalStorage{
		path: path,
	}
}

func (self *LocalStorage) WriteCloser(id string) (io.WriteCloser, error) {
	err := os.MkdirAll(self.dirPath(id), 0777)
	if err != nil {
		return nil, err
	}
	return os.Create(self.fullPath(id))
}

func (self *LocalStorage) ReadCloser(id string) (io.ReadCloser, error) {
	// get the path and filename
	return os.Open(self.fullPath(id))
}

func (self *LocalStorage) Move(from, to string) error {

	return os.Rename(self.fullPath(from), self.fullPath(to))
}

func (self *LocalStorage) Delete(id string) error {
	// get the path and filename
	return os.Remove(self.fullPath(id))
}

func (self *LocalStorage) dirPath(id string) string {
	re := regexp.MustCompile(`/\w*$`)
	return re.ReplaceAllString(self.fullPath(id), "")
}

func (self *LocalStorage) fullPath(id string) string {
	re := regexp.MustCompile("-")
	return self.path + "/" + re.ReplaceAllString(id, "/")
}
