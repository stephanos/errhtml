package errhtml

import (
	"bytes"
	. "github.com/101loops/bdd"
)

var _ = Describe("Error", func() {

	Context("create", func() {

		It("from error", func() {
			err, _ := NewError(simpleErr()).(*errContext)

			Check(err, NotNil)
			Check(err.MetaError, IsEmpty)

			Check(err.Title, Equals, "Error")
			Check(err.Message, Equals, "open nonsense: no such file or directory")

			Check(err.SourceTrace, HasLen, 0)

			Check(err.Error(), Equals, "Error: open nonsense: no such file or directory")
		})

		It("from template error", func() {
			err, _ := NewError(templateErr(), NewFileSource("fixture_tmpl.html", 5)).(*errContext)

			Check(err, NotNil)
			Check(err.MetaError, IsEmpty)

			Check(err.Title, Equals, "Error")
			Check(err.Message, Equals, `template: fixture_tmpl.html:5: function "undefined_action" not defined`)

			Check(err.SourceTrace, HasLen, 0)
			Check(err.SourceContext, HasLen, 7)

			Check(err.Error(), Equals, `Error (in fixture_tmpl.html:5): template: fixture_tmpl.html:5: function "undefined_action" not defined`)
		})

		It("from panic", func() {
			err, _ := panicErr().(*errContext)

			Check(err, NotNil)
			Check(err.MetaError, IsEmpty)

			Check(err.Title, Equals, "Panic")
			Check(err.Message, Equals, "runtime error: integer divide by zero")

			topFrame := err.Source
			Check(topFrame.Line, Equals, 15)
			Check(topFrame.Text, Equals, "div(1, 0)")
			Check(topFrame.File, HasSuffix, "/errhtml/fixture_test.go")

			stackTrace := err.SourceTrace
			Check(len(stackTrace), IsGreaterThan, 5)
			Check(stackTrace[0], Equals, *topFrame)

			lastFrame := stackTrace[len(stackTrace)-1]
			Check(lastFrame.File, HasSuffix, "runtime/proc.c")

			Check(err.Error(), Equals, `Panic (in <GOPATH>/github.com/101loops/errhtml/fixture_test.go:15): runtime error: integer divide by zero`)
		})

		It("from other Error", func() {
			err1, _ := NewError(simpleErr()).(*errContext)
			err2, _ := NewError(err1).(*errContext)

			Check(err1.Message, Equals, err2.Message)
		})

		It("from string", func() {
			err, _ := NewError("runtime error").(*errContext)

			Check(err, NotNil)
			Check(err.Title, Equals, "Error")
			Check(err.Message, Equals, "runtime error")
		})

		It("from template error with missing file", func() {
			err, _ := NewError(templateErr(), NewFileSource("nonsense.html", 10)).(*errContext)

			Check(err, NotNil)
			Check(err.MetaError, Contains, `unable to load error source "nonsense.html"`)
		})

		It("from nil", func() {
			Check(NewError(nil), IsNil)
		})
	})

	Context("render", func() {

		render := func(err interface{}) string {
			var buf bytes.Buffer
			Render(err, &buf)
			return buf.String()
		}

		It("an error", func() {
			render(simpleErr())
		})

		It("a template error", func() {
			err, _ := NewError(templateErr(), NewFileSource("fixture_tmpl.html", 5)).(*errContext)
			render(err)
		})

		It("a panic", func() {
			render(panicErr())
		})

		It("a string", func() {
			render("runtime error")
		})
	})
})
