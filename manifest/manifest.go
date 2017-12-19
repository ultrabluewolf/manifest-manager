package manifest

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/ultrabluewolf/manifest-manager/files"

	mmlogger "github.com/ultrabluewolf/manifest-manager/logger"
)

var logger = mmlogger.New()

type Manifest struct {
	Filename string
	Files    map[string]bool
}

func New(filename string) *Manifest {
	return &Manifest{
		Filename: filename,
		Files:    map[string]bool{},
	}
}

// parse manifest and generate struct
func ParseManifestFile(filename string) (*Manifest, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	manifest := &Manifest{
		Filename: filename,
		Files:    map[string]bool{},
	}

	// read manifest file lines
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			continue
		}
		manifest.Files[line] = true
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return manifest, nil
}

func (manifest *Manifest) FileList() []string {
	files := []string{}
	for file, isActive := range manifest.Files {
		if !isActive {
			continue
		}
		files = append(files, file)
	}
	sort.Strings(files)

	return files
}

// remove all stale files from the manifest
func (manifest *Manifest) Prune() error {
	fileList := manifest.FileList()
	logger.Info(fmt.Sprintf("manifest size %d", len(fileList)))

	count := 0
	for _, fileItem := range fileList {
		if !files.Exists(fileItem) || files.IsDir(fileItem) {
			manifest.Files[fileItem] = false
			count += 1
		}
	}
	logger.Info(fmt.Sprintf("removed %d stale files", count))

	return nil
}

func (manifest *Manifest) Save() error {
	// create all parent folders in case they don't exist already
	dir, _ := filepath.Split(manifest.Filename)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	// initialize file buffer
	file, err := os.Create(manifest.Filename)
	if err != nil {
		return err
	}
	defer file.Close()

	buf := bufio.NewWriter(file)
	fileList := manifest.FileList()
	for _, fileItem := range fileList {
		_, err := buf.WriteString(fmt.Sprintf("%s\n", fileItem))
		if err != nil {
			return err
		}
	}
	buf.Flush()

	logger.Info("manifest saved -", manifest.Filename)
	return nil
}

// add all files included in given glob pattern
func (manifest *Manifest) Add(pattern string) error {
	filepaths, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}
	logger.Info(fmt.Sprintf("found %d files", len(filepaths)))

	count := 0
	for _, fileItem := range filepaths {
		if isActive, ok := manifest.Files[fileItem]; !ok || !isActive {
			count += 1
		}
		manifest.Files[fileItem] = true
	}
	logger.Info(fmt.Sprintf("added %d files", count))

	if err = manifest.Prune(); err != nil {
		return err
	}

	return nil
}

// remove all files included in given glob pattern
func (manifest *Manifest) Remove(pattern string) error {
	filepaths, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}
	logger.Info(fmt.Sprintf("found %d files", len(filepaths)))

	count := 0
	for _, fileItem := range filepaths {
		if isActive, ok := manifest.Files[fileItem]; ok && isActive {
			count += 1
		}
		manifest.Files[fileItem] = false
	}
	logger.Info(fmt.Sprintf("removed %d files", count))

	if err = manifest.Prune(); err != nil {
		return err
	}

	return nil
}
