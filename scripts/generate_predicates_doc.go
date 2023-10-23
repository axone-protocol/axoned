package main

import (
	"bufio"
	"embed"
	"fmt"
	"go/build"
	"os"
	"path"
	"strings"

	"github.com/Masterminds/sprig/v3"
	"github.com/princjef/gomarkdoc"
	"github.com/princjef/gomarkdoc/format"
	"github.com/princjef/gomarkdoc/lang"
	"github.com/princjef/gomarkdoc/logger"
)

//go:embed templates/*.gotxt
var f embed.FS

func GeneratePredicateDocumentation() error {
	// Create a renderer to output data
	out, err := gomarkdoc.NewRenderer(
		gomarkdoc.WithTemplateOverride("text", readTemplateMust("text.gotxt")),
		gomarkdoc.WithTemplateOverride("doc", readTemplateMust("doc.gotxt")),
		gomarkdoc.WithTemplateOverride("list", readTemplateMust("list.gotxt")),
		gomarkdoc.WithTemplateOverride("import", ""),
		gomarkdoc.WithTemplateOverride("package", readTemplateMust("package.gotxt")),
		gomarkdoc.WithTemplateOverride("file", ""),
		gomarkdoc.WithTemplateOverride("func", readTemplateMust("func.gotxt")),
		gomarkdoc.WithTemplateOverride("index", ""),
		gomarkdoc.WithTemplateFunc("snakecase", sprig.TxtFuncMap()["snakecase"]),
		gomarkdoc.WithTemplateFunc("hasSuffix", sprig.TxtFuncMap()["hasSuffix"]),
		gomarkdoc.WithTemplateFunc("countSubstr", func(substr string, s string) int {
			return strings.Count(s, substr)
		}),
		gomarkdoc.WithTemplateFunc("dict", sprig.TxtFuncMap()["dict"]),
		gomarkdoc.WithFormat(&format.GitHubFlavoredMarkdown{}),
	)
	if err != nil {
		return err
	}

	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	buildPkg, err := build.ImportDir(path.Join(wd, "x", "logic", "predicate"), build.ImportComment)
	if err != nil {
		return err
	}

	log := logger.New(logger.DebugLevel)
	pkg, err := lang.NewPackageFromBuild(log, buildPkg)
	if err != nil {
		return err
	}
	content, err := out.Package(pkg)
	if err != nil {
		return err
	}

	if err := writeToFile("docs/predicate/predicates.md", content); err != nil {
		return err
	}

	return nil
}

func readTemplateMust(templateName string) string {
	template, err := f.ReadFile("templates/" + templateName)

	if err != nil {
		panic(fmt.Errorf("failed to read file %s: %s", templateName, err))
	}

	return string(template)
}

func writeToFile(filePath string, content string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", filePath, err)
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	if _, err = w.WriteString(content); err != nil {
		return fmt.Errorf("failed to write to file %s: %w", filePath, err)
	}

	if err = w.Flush(); err != nil {
		return fmt.Errorf("failed to flush writer: %w", err)
	}

	return nil
}
