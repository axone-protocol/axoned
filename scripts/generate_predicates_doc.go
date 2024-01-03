package main

import (
	"embed"
	"fmt"
	"go/build"
	"os"
	"path"
	"strings"

	"github.com/Masterminds/sprig/v3"
	"github.com/princjef/gomarkdoc"
	"github.com/princjef/gomarkdoc/lang"
	"github.com/princjef/gomarkdoc/logger"
)

//go:embed templates/*.gotxt
var f embed.FS

// globalCtx used to keep track of contexts between templates.
// (yes it's a hack).
var globalCtx = make(map[string]interface{})

func GeneratePredicateDocumentation() error {
	out, err := createRenderer()
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

	return writeToFile("docs/predicate/predicates.md", content)
}

func createRenderer() (*gomarkdoc.Renderer, error) {
	templateFunctionOpts := make([]gomarkdoc.RendererOption, 0)

	templateFunctionOpts = append(
		templateFunctionOpts,
		gomarkdoc.WithTemplateOverride("text", readTemplateMust("text.gotxt")),
		gomarkdoc.WithTemplateOverride("doc", readTemplateMust("doc.gotxt")),
		gomarkdoc.WithTemplateOverride("list", readTemplateMust("list.gotxt")),
		gomarkdoc.WithTemplateOverride("import", ""),
		gomarkdoc.WithTemplateOverride("package", readTemplateMust("package.gotxt")),
		gomarkdoc.WithTemplateOverride("file", ""),
		gomarkdoc.WithTemplateOverride("func", readTemplateMust("func.gotxt")),
		gomarkdoc.WithTemplateOverride("index", ""),
	)

	for k, v := range sprig.GenericFuncMap() {
		templateFunctionOpts = append(
			templateFunctionOpts,
			gomarkdoc.WithTemplateFunc(k, v),
		)
	}

	templateFunctionOpts = append(
		templateFunctionOpts,
		gomarkdoc.WithTemplateFunc("countSubstr", func(substr string, s string) int {
			return strings.Count(s, substr)
		}),
		gomarkdoc.WithTemplateFunc("globalCtx", func() map[string]interface{} {
			return globalCtx
		}),
	)
	return gomarkdoc.NewRenderer(templateFunctionOpts...)
}

func readTemplateMust(templateName string) string {
	template, err := f.ReadFile("templates/" + templateName)
	if err != nil {
		panic(fmt.Errorf("failed to read file %s: %w", templateName, err))
	}

	return string(template)
}
