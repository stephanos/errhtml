package errhtml

import (
	"html/template"
	"io/ioutil"
)

func panicErr() (e error) {
	defer func() {
		if err := recover(); err != nil {
			e = NewError(err)
		}
	}()

	div(1, 0)
	return e
}

func div(x, y int) int {
	return x / y
}

func simpleErr() error {
	_, err := ioutil.ReadFile("nonsense")
	return err
}

func templateErr() error {
	_, err := template.New("test").ParseFiles("fixture_tmpl.html")
	return err
}
