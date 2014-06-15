package errhtml

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var (
	env envConf
)

func init() {
	env = initEnv()
}

func initEnv() envConf {
	cwd, err := os.Getwd()
	if err != nil {
		panic(fmt.Errorf("unable to determine current working directory: %v", err))
	}

	return envConf{
		CurrentDir: cwd,
		GoRoot:     runtime.GOROOT(),
		GoPaths:    strings.Split(os.Getenv("GOPATH"), ":"),
	}
}

type envConf struct {
	GoRoot     string
	GoPaths    []string
	CurrentDir string
}

type Source struct {
	source *source
}

type source struct {
	File      string
	Text      string
	Line      int
	Highlight bool
}

type fileLocation struct {
	Root, Package, File string
}

func NewFileSource(filePath string, line int) Source {
	return Source{&source{File: filePath, Line: line}}
}

func NewTextSource(text string) Source {
	return Source{&source{Text: text}}
}

func (loc *fileLocation) String() string {
	root := loc.Root
	if root != "" {
		root = "<" + root + ">"
	}
	return filepath.Join(root, loc.Package, loc.File)
}

func (s *source) FileName() string {
	return s.FileLocation().File
}

func (s *source) FileLocation() *fileLocation {
	return fileLocationFromPath(s.File, env)
}

func (s *source) AbbreviatedFilePath() string {
	return s.FileLocation().String()
}

func (s *source) Lines() ([]string, error) {
	filePath := s.File

	paths := []string{filePath}
	if filepath.IsAbs(filePath) {
		paths = append(paths, filePath)
	} else {
		paths = append(paths, filepath.Join(env.CurrentDir, filePath))
		for _, goPath := range env.GoPaths {
			paths = append(paths, filepath.Join(goPath, filePath))
		}
	}

	for _, path := range paths {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			continue
		}

		bytes, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, err
		}

		return strings.Split(string(bytes), "\n"), nil
	}

	return nil, nil
}

func (s *source) AbbreviatedFilePathDirectories() []string {
	dirs := strings.Split(s.AbbreviatedFilePath(), "/")
	if len(dirs) > 0 {
		dirs = dirs[:len(dirs)-1]
	}
	return dirs
}

func fileLocationFromPath(path string, env envConf) *fileLocation {
	var root, prefix string

	if env.GoRoot != "" {
		goRootPrefix := filepath.Join(env.GoRoot, "src", "pkg")
		if strings.HasPrefix(path, goRootPrefix) {
			root = "GOROOT"
			prefix = goRootPrefix
		}
	}

	if root == "" && len(env.GoPaths) > 0 {
		for _, goPath := range env.GoPaths {
			goPathPrefix := filepath.Join(goPath, "src")
			if strings.HasPrefix(path, goPathPrefix) {
				root = "GOPATH"
				prefix = goPathPrefix
				break
			}
		}
	}

	cwd := env.CurrentDir
	if root == "" && len(cwd) > 0 {
		if strings.HasPrefix(path, cwd) {
			root = "CWD"
			prefix = cwd
		}
	}

	file := filepath.Base(path)
	pckg := filepath.Dir(strings.Replace(path, prefix, "", 1))
	return &fileLocation{root, pckg, file}
}
