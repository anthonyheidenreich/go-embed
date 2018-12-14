package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/anthonyheidenreich/gadget/log"
	"github.com/anthonyheidenreich/go-embed/embed"
)

const (
	suffix = ".tmpl"

	// EmbedderTypeTemplate is the name of the template embedder module
	EmbedderTypeTemplate = "template"
)

var (
	packageName = flag.String("p", "", "the package name to use for the created go file")
)

// Usage is a replacement usage function for the flags package.
func Usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "\tinclude [flags] [module] [directory] [output]\n")
	fmt.Fprintf(os.Stderr, "Arguments:\n")
	fmt.Fprintf(os.Stderr, "\t [module] REQUIRED: The module to use when including files (template)\n")
	fmt.Fprintf(os.Stderr, "\t [directory] REQUIRED: The directory to scan for include files\n")
	fmt.Fprintf(os.Stderr, "\t [output] REQUIRED: The path to the generated go file to output to, existing file will "+
		"be overwritten\n")
	fmt.Println("Flags:")
	flag.PrintDefaults()
}

const (
	// ArgumentCount is the number of expected arguments
	ArgumentCount = 3
)

func main() {
	log.NewGlobal("go-embed", log.FunctionFromEnv())
	flag.Usage = Usage
	flag.Parse()
	if len(flag.Args()) != ArgumentCount {
		log.Infof("invalid argument count %d", len(flag.Args()))
		Usage()
		os.Exit(1)
	}
	moduleName := flag.Arg(0)
	directory := flag.Arg(1)
	// ensure a trailing slash
	if directory[len(directory)-1] != '/' {
		directory = directory + "/"
	}
	outputFile := flag.Arg(2)
	var module embed.Embedder
	switch moduleName {
	case "template":
		module = embed.NewTemplateEmbedder(*packageName)
	default:
		log.Fatalf("'%s' is not a valid module name. Must be 'template'", moduleName)
		os.Exit(2)
	}
	fs, _ := ioutil.ReadDir(directory)
	log.Infof("looking for includes matching '%s * %s'", directory, suffix)
	for _, f := range fs {
		if strings.HasSuffix(f.Name(), suffix) {
			filepath := directory + f.Name()
			bytes, err := ioutil.ReadFile(filepath)
			if nil != err {
				log.Fatalf("failed to read from file %s\n%#v", filepath, err)
				os.Exit(3)
			}
			log.Infof("including %s", filepath)
			module.EmbedFile(f.Name(), bytes)
		}
	}

	fd, err := os.OpenFile(outputFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.FileMode(0666))
	if nil != err {
		log.Fatalf("could not open output file %#v", err)
		os.Exit(4)
	}
	if err = module.Finalize(fd); err != nil {
		log.Fatalf("module finalize failed: %#v", err)
		os.Exit(5)
	}
}
