package fileutils

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// File holds the filename and the full path of a file
type File struct {
	Name     string
	FullPath string
}

// Folder holds all information of the folder, including the subfolder, files, parentfolder its name and path
type Folder struct {
	Name       string
	FullPath   string
	Parent     *Folder
	SubFolders []*Folder
	Files      []*File
}

// FileFilter is a function which can be used to filter various files in the ReadFolderContent function
type FileFilter func(os.FileInfo) bool

// ReadFolderContent goes through all content of the given path and returns the top folder with all its content
// It can be passed a FileFilter function where files can be excluded from the result
func ReadFolderContent(startPath string, fileFilter FileFilter) (*Folder, error) {
	folderName := filepath.Base(startPath)
	rootFolder := newFolder(folderName, startPath, nil)
	if err := rootFolder.traverseContent(fileFilter); err != nil {
		return nil, err
	}
	return rootFolder, nil
}

// WriteStructure writes the folder tree to the given io.Writer
// the indent should start with 0
func (f *Folder) WriteStructure(writer io.Writer, indent int) {
	fmt.Fprintf(writer, "%s- %s\n", strings.Repeat("\t", indent), f.Name)
	for _, file := range f.Files {
		fmt.Fprintf(writer, "%s\t* %s\n", strings.Repeat("\t", indent), file.Name)
	}
	for _, folder := range f.SubFolders {
		folder.WriteStructure(writer, indent+1)
	}
}

func (f *Folder) traverseContent(fileFilter FileFilter) error {
	cont, err := ioutil.ReadDir(f.FullPath)
	if err != nil {
		return err
	}
	for _, fi := range cont {
		fullPath := f.FullPath + string(filepath.Separator) + fi.Name()
		switch mode := fi.Mode(); {
		case mode.IsDir():
			folder := newFolder(fi.Name(), fullPath, f)
			f.SubFolders = append(f.SubFolders, folder)
			if err := folder.traverseContent(fileFilter); err != nil {
				return err
			}
		case mode.IsRegular():
			if fileFilter == nil || fileFilter(fi) {
				file := newFile(fi.Name(), fullPath)
				f.Files = append(f.Files, file)
			}
		}
	}
	return nil
}

func newFolder(name, fullPath string, parentFolder *Folder) *Folder {
	return &Folder{Name: name, FullPath: fullPath, Parent: parentFolder}
}

func newFile(name, fullPath string) *File {
	return &File{Name: name, FullPath: fullPath}
}
