package embed

import (
	"os"
	"path"
	"strconv"
	"testing"

	"github.com/anthonyheidenreich/gadget/generator"
	"github.com/stretchr/testify/assert"
)

func TestTemplateInitialize(t *testing.T) {
	assert := assert.New(t)

	module := NewTemplateEmbedder("").(*templateEmbedder)
	assert.Equal(defaultPackageName, module.Context.PackageName)
	expected := generator.String(10)
	module = NewTemplateEmbedder(expected).(*templateEmbedder)
	assert.Equal(expected, module.Context.PackageName)
}

func TestTemplateIncludeFile(t *testing.T) {
	assert := assert.New(t)

	module := NewTemplateEmbedder("bar").(*templateEmbedder)
	fileName := "test.tpl"
	expected := "Do Stuff"
	contents := []byte(expected)

	err := module.EmbedFile(fileName, contents)
	assert.Equal(templateContext{Name: fileName, Data: strconv.Quote(expected)}, module.Context.Templates[0])
	assert.NoError(err)
}

func TestTemplateFinalize(t *testing.T) {
	assert := assert.New(t)

	module := NewTemplateEmbedder("bar")
	fileName := "test.tpl"
	contents := []byte("Do Stuff")
	err := module.EmbedFile(fileName, contents)
	assert.NoError(err)

	outputFile := path.Join(os.TempDir(), "template.finalize.test")
	f, err := os.Create(outputFile)
	assert.NoError(err)

	err = module.Finalize(f)
	assert.NoError(err)
	expected := []byte("package bar")
	actual := make([]byte, len(expected))

	f.Seek(0, 0)
	assert.NoError(err)
	_, err = f.Read(actual)
	f.Close()
	assert.NoError(err)

	assert.Equal(expected, actual)
}
