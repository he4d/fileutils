package fileutils

import (
	"os"
	"path/filepath"
	"testing"
)

func FilterMp3(fileInfo os.FileInfo) bool {
	if filepath.Ext(fileInfo.Name()) == ".mp3" {
		return true
	}
	return false
}

func TestTraverseDirectoriesWithoutFilter(t *testing.T) {
	rootFolder, err := ReadFolderContent("/home/he4d/Musik", nil)
	if err != nil {
		t.Error(err)
	}
	rootFolder.WriteStructure(os.Stdout, 0)
}

func TestTraverseDirectoriesWithFilter(t *testing.T) {
	rootFolder, err := ReadFolderContent("/home/he4d/Musik", FilterMp3)
	if err != nil {
		t.Error(err)
	}
	rootFolder.WriteStructure(os.Stdout, 0)
}

func BenchmarkTraverseDirectoriesWithoutFilter(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_, err := ReadFolderContent("/home/he4d/Musik", nil)
		if err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkTraverseDirectoriesWithFilter(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_, err := ReadFolderContent("/home/he4d/Musik", FilterMp3)
		if err != nil {
			b.Error(err)
		}
	}
}
