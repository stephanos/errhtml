package errhtml

import (
	"fmt"
	"io"
	"runtime/debug"
	"strings"
)

var (
	contextPadding = 5
)

type Error struct {
	Title, Message string // Description of the error, as presented to the user.
	Source         *source
	SourceTrace    []source
	SourceContext  []source
	MetaError      string // Error that occurred producing the error page.
}

// Find the deepest stack from in user code and provide a code listing of
// that, on the line that eventually triggered the panic.  Returns nil if no
// relevant stack frame can be found.
func NewError(err interface{}, sources ...Source) error {
	if err == nil {
		return nil
	}

	if err, ok := err.(*Error); ok {
		return err
	}

	var hasTrace bool
	stackTrace := getStackTrace()
	if stackTrace != nil && len(stackTrace) > 0 {
		hasTrace = true
	}

	message := "Unspecified error"
	if err != nil {
		message = fmt.Sprint(err)
	}

	title := "Error"
	if hasTrace {
		title = "Panic"
	}

	var topFrame *source
	if hasTrace {
		topFrame = &stackTrace[0]
	}

	var context []source
	var metaError string
	if len(sources) > 0 {
		src := sources[0].source

		if !hasTrace {
			topFrame = src
		}

		if ctx, err := sourceContext(src); err == nil {
			context = ctx
		} else {
			metaError = err.Error()
		}
	}

	return &Error{
		Title:         title,
		Source:        topFrame,
		Message:       message,
		SourceTrace:   stackTrace,
		SourceContext: context,
		MetaError:     metaError,
	}
}

// Construct a plaintext version of the error, taking account that fields are optionally set.
// Returns e.g. Compilation Error (in views/header.html:51): expected right delim in end; got "}"
func (e *Error) Error() string {
	loc := ""
	source := e.Source
	if source.File != "" {
		line := ""
		if source.Line != 0 {
			line = fmt.Sprintf(":%d", source.Line)
		}
		loc = fmt.Sprintf("(in %s%s)", source.AbbreviatedFilePath(), line)
	}
	header := loc
	if e.Title != "" {
		if loc != "" {
			header = fmt.Sprintf("%s %s: ", e.Title, loc)
		} else {
			header = fmt.Sprintf("%s: ", e.Title)
		}
	}
	return fmt.Sprintf("%s%s", header, e.Message)
}

// Returns a snippet of the source around where the error occurred.
func sourceContext(src *source) ([]source, error) {
	if src == nil || src.File == "" {
		return nil, nil
	}

	lines, err := src.Lines()
	if err != nil {
		return nil, err
	}
	if lines == nil {
		return nil, fmt.Errorf("unable to load error source %q", src.File)
	}

	start := (src.Line - 1) - contextPadding
	if start < 0 {
		start = 0
	}

	end := (src.Line - 1) + contextPadding
	if end > len(lines) {
		end = len(lines)
	}

	context := make([]source, end-start)
	for i, line := range lines[start:end] {
		fileLine := start + i + 1
		context[i] = source{
			Text:      line,
			File:      src.File,
			Line:      fileLine,
			Highlight: fileLine == src.Line,
		}
	}

	return context, nil
}

// Since this is supposed to be used in development only,
// instead of returning an error it panics.
func (e *Error) Render(w io.Writer) {
	err := errTemplate.Execute(w, e)
	if err != nil {
		panic(err)
	}
}

func getStackTrace() (stackTrace []source) {
	var includeElem bool
	var traceElem *source

	fullStackTrace := strings.Split(string(debug.Stack()), "\n")
	for i, elem := range fullStackTrace {
		elem = strings.TrimSpace(elem)
		//println(elemStr)

		if i%2 == 0 {
			file, line := fileContextFromStackElement(elem)
			if file != "" {
				traceElem = NewFileSource(file, line).source
			} else {
				traceElem = nil
			}
			continue
		}

		if traceElem == nil {
			continue
		}

		if includeElem {
			traceElem.Text, _ = codeContextFromStackElement(elem)
			stackTrace = append(stackTrace, *traceElem)
			continue
		}

		if strings.HasPrefix(elem, "sigpanic: ") {
			includeElem = true
		}
	}
	return
}

func fileContextFromStackElement(elem string) (file string, line int) {
	colonIndex := strings.LastIndex(elem, ":")
	if colonIndex != -1 {
		file = elem[:colonIndex]
		fmt.Sscan(elem[colonIndex+1:], &line)
	}
	return
}

func codeContextFromStackElement(elem string) (text, fn string) {
	colonIndex := strings.Index(elem, ":")
	if colonIndex != -1 {
		text = elem[colonIndex+2:]
		fn = elem[colonIndex:]
	}
	return
}
