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

	It("return file name", func() {
		src := source{File: absFile("file.go")}
		Check(src.FileName(), Equals, "file.go")
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
