package filerepository

import (
	"fmt"
	"log"
	"os"
	"time"
)

type FileRepository struct {
	files map[string]*os.File
}

func NewFileRepository() *FileRepository {
	return &FileRepository{
		files: make(map[string]*os.File),
	}
}

func (r *FileRepository) CreateAudioFile() string {
	filename := fmt.Sprintf("audio_%v.wav", time.Now().Unix())
	file, err := os.Create(filename)
	if err != nil {
		log.Println("File creation error:", err)
		return ""
	}

	r.files[filename] = file
	return filename
}

func (r *FileRepository) WriteAudioData(filename string, data []byte) error {
	file, ok := r.files[filename]
	if !ok {
		return fmt.Errorf("file not found: %s", filename)
	}

	_, err := file.Write(data)

	fmt.Println("asdasd")
	return err
}

func (r *FileRepository) CloseAudioFile(filename string) {
	if file, ok := r.files[filename]; ok {
		file.Close()
		delete(r.files, filename)
	}
}

func (r *FileRepository) DeleteAudioFile(filename string) {

	os.Remove(filename)

}
