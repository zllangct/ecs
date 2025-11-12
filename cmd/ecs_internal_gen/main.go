package main

import (
	"bytes"
	"embed"
	"flag"
	"fmt"
	"go/format"
	"io"
	"os"
	"path/filepath"
	"text/template"
)

//go:embed *_template.*
var templateFiles embed.FS

func main() {
	flag.Usage = func() {
		fmt.Println("Usage: internal gen <command> [<args>]")
		fmt.Println("Commands:")
		fmt.Println("  FixedString")
		os.Exit(1)
	}
	flag.Parse()

	var fn func() error
	switch flag.Arg(0) {
	case "FixedString":
		fn = func() error {
			genFixedString()
			return nil
		}
	case "FixedCompound":
		fn = func() error {
			genFixedCompound()
			return nil
		}
	}

	if fn == nil {
		flag.Usage()
		os.Exit(1)
	}

	if err := fn(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}

func genFixedString() {
	flags := flag.NewFlagSet("FixedString", flag.ExitOnError)

	var output string
	flags.StringVar(&output, "o", ".", "Output directory path.")
	var packageName string
	flags.StringVar(&packageName, "p", "ecs", "package name.")
	if err := flags.Parse(flag.Args()[1:]); err != nil {
		return
	}

	t := template.New("")
	t, err := t.ParseFS(templateFiles, "fixed_string_template.*")
	if err != nil {
		panic(err)
	}
	type Data struct {
		MaxSize     []byte
		PackageName string
	}

	data := Data{
		PackageName: packageName,
		MaxSize:     make([]byte, 0, 1024),
	}
	for _ = range 1025 {
		data.MaxSize = append(data.MaxSize, byte(0))
	}
	var buffer bytes.Buffer
	if err := t.ExecuteTemplate(&buffer, "fixed_string", data); err != nil {
		panic(err)
	}

	outputFile, err := os.Create(filepath.Join(output, "fixed_string_generated.go"))
	if err != nil {
		panic(err)
	}

	err = finish(outputFile, &buffer)
	if err != nil {
		panic(err)
	}

	if err := outputFile.Close(); err != nil {
		panic(err)
	}
}

func genFixedCompound() {
	flags := flag.NewFlagSet("FixedCompound", flag.ExitOnError)

	var output string
	flags.StringVar(&output, "o", ".", "Output directory path.")
	var packageName string
	flags.StringVar(&packageName, "p", "ecs", "package name.")
	if err := flags.Parse(flag.Args()[1:]); err != nil {
		return
	}

	t := template.New("")
	t, err := t.ParseFS(templateFiles, "fixed_compound_template.*")
	if err != nil {
		panic(err)
	}
	type Data struct {
		MaxSize     []byte
		PackageName string
	}

	data := Data{
		PackageName: packageName,
		MaxSize:     make([]byte, 0, 1024),
	}
	for _ = range 1025 {
		data.MaxSize = append(data.MaxSize, byte(0))
	}
	var buffer bytes.Buffer
	if err := t.ExecuteTemplate(&buffer, "fixed_compound", data); err != nil {
		panic(err)
	}

	outputFile, err := os.Create(filepath.Join(output, "fixed_compound_generated.go"))
	if err != nil {
		panic(err)
	}

	err = finish(outputFile, &buffer)
	if err != nil {
		panic(err)
	}

	if err := outputFile.Close(); err != nil {
		panic(err)
	}
}

func finish(output io.Writer, buffer io.Reader) error {
	buf, err := io.ReadAll(buffer)
	if err != nil {
		return err
	}
	b, err := format.Source(buf)
	if err != nil {
		_, err = output.Write(buf)
		return err
	}
	_, err = output.Write(b)
	return err
}
