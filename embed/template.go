package embed

import (
	"io"
	"strconv"
	"text/template"

	"github.com/anthonyheidenreich/gadget/log"
)

const defaultPackageName = "templates"

var templatesTemplate = template.Must(template.New("template").Parse(`package {{ .PackageName }}

// THIS IS A GENERATED FILE. DO NOT MODIFY

import "text/template"

// GetTemplates returns a template that has the all the other templates parsed into it accessible via their filename.
func GetTemplates() *template.Template {
    master := template.New("Template")
    {{ range $index, $template := .Templates }}
    // {{ $template.Name }}
    template.Must(master.New("{{ $template.Name }}").Parse(string(` + "{{ $template.Data }}" + `)))
    {{ end }}
    return master
}
`))

type templateContext struct {
	Name string
	Data string
}

type context struct {
	PackageName string
	Templates   []templateContext
}

// NewTemplateEmbedder for including templates in a go binary.
func NewTemplateEmbedder(packageName string) Embedder {
	if packageName == "" {
		log.Infof("Using default package name: %s", defaultPackageName)
		packageName = defaultPackageName
	}
	return &templateEmbedder{Context: &context{PackageName: packageName, Templates: []templateContext{}}}
}

// templateEmbedder embeds templates into Go for compilation
type templateEmbedder struct {
	Context *context
}

func (module *templateEmbedder) EmbedFile(fileName string, contents []byte) error {
	module.Context.Templates = append(module.Context.Templates, templateContext{
		Name: fileName,
		Data: strconv.Quote(string(contents)),
	})
	return nil
}

func (module *templateEmbedder) Finalize(fileDescriptor io.Writer) error {
	err := templatesTemplate.Execute(fileDescriptor, module.Context)
	if nil == err {
		log.Infof("templates included: %d", len(module.Context.Templates))
	}
	return err
}
