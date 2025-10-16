package templates

import (
	"bytes"
	"math/rand"
	"text/template"
	"time"
)

const (
	randNumFunc = "randNum"
)

// TemplateBuilder builder go templates.
type TemplateBuilder struct {
	funcMap template.FuncMap
}

// NewTemplateBuilder create a new TemplateBuilder.
func NewTemplateBuilder() *TemplateBuilder {
	tb := &TemplateBuilder{}
	funcMap := template.FuncMap{
		randNumFunc: tb.randNum,
	}
	tb.funcMap = funcMap
	return tb
}

// Process processing message and return string.
func (b *TemplateBuilder) Process(str string) (string, error) {
	tmpl := template.New("").Funcs(b.funcMap)

	tmpl, err := tmpl.Parse(str)
	if err != nil {
		return "", err
	}

	var buffer bytes.Buffer
	err = tmpl.Execute(&buffer, nil)
	if err != nil {
		return "", err
	}

	return buffer.String(), nil
}

// randNum return random number with specified range.
func (b *TemplateBuilder) randNum(min, max int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return min + r.Intn(max-min+1)
}
