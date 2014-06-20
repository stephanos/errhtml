package errhtml

import (
	. "github.com/101loops/bdd"
	"path/filepath"
)

var _ = Describe("Source", func() {

	var cwd string
	var absFile = func(path string) string { return filepath.Join(cwd, path) }

	BeforeSuite(func() {
		cwd = "/Users/gopher/workspace/src/host/vendor/library"
	})

	Context("create", func() {
		It("from file", func() {
			src := NewFileSource("directory/file.go", 42).source

			Check(src.Line, Equals, 42)
			Check(src.File, Equals, "directory/file.go")

			Check(src.FileName(), Equals, "file.go")
			Check(src.AbbreviatedFilePath(), Equals, "directory/file.go")
			Check(src.AbbreviatedFilePathDirectories(), Equals, []string{"directory"})
		})

		It("from text", func() {
			src := NewTextSource("err").source

			Check(src.Text, Equals, "err")
		})
	})

	Context("return file location when", func() {
		It("in GOROOT", func() {
			loc := fileLocationFromPath(absFile("/src/pkg/runtime/panic.c"), envConf{GoRoot: cwd})

			Check(loc, Equals, &fileLocation{"GOROOT", "/runtime", "panic.c"})
			Check(loc.String(), Equals, "<GOROOT>/runtime/panic.c")
		})

		It("in GOPATH", func() {
			loc := fileLocationFromPath(absFile("file.go"), envConf{GoPaths: []string{"nonsense", "/Users/gopher/workspace"}})

			Check(loc, Equals, &fileLocation{"GOPATH", "/host/vendor/library", "file.go"})
			Check(loc.String(), Equals, "<GOPATH>/host/vendor/library/file.go")
		})

		It("in working directory", func() {
			loc := fileLocationFromPath(absFile("file.go"), envConf{CurrentDir: cwd})

			Check(loc, Equals, &fileLocation{"CWD", "/", "file.go"})
			Check(loc.String(), Equals, "<CWD>/file.go")
		})

		It("standalone", func() {
			loc := fileLocationFromPath(absFile("file.go"), envConf{})

			Check(loc, Equals, &fileLocation{"", cwd, "file.go"})
			Check(loc.String(), Equals, cwd+"/file.go")
		})
	})

})
