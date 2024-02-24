package main

import (
	"embed"
	"fmt"
	"go/build"
	"os"
	"path"
	"slices"
	"strings"

	"github.com/Masterminds/sprig/v3"
	"github.com/huandu/xstrings"
	"github.com/princjef/gomarkdoc"
	"github.com/princjef/gomarkdoc/lang"
	"github.com/princjef/gomarkdoc/logger"
	"github.com/samber/lo"
)

//go:embed templates/*.gotxt
var f embed.FS

const (
	predicatePath = "x/logic/predicate"
	outputPath    = "docs/predicate"
)

func generatePredicateDocumentation() error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	buildPkg, err := build.ImportDir(path.Join(wd, predicatePath), build.ImportComment)
	if err != nil {
		return err
	}

	log := logger.New(logger.DebugLevel)
	pkg, err := lang.NewPackageFromBuild(log, buildPkg)
	if err != nil {
		return err
	}

	funcs := lo.Filter(pkg.Funcs(), func(item *lang.Func, _ int) bool { return isFuncAPredicate(item) })
	slices.SortFunc(funcs, func(a *lang.Func, j *lang.Func) int {
		return strings.Compare(a.Name(), j.Name())
	})

	for idx, f := range funcs {
		// globalCtx used to keep track of contexts between templates.
		// (yes it's a hack).
		globalCtx := make(map[string]interface{})
		globalCtx["frontmatter"] = map[string]interface{}{
			"sidebar_position": idx + 1,
		}

		out, err := createRenderer(globalCtx)
		if err != nil {
			return err
		}

		name := strings.Replace(functorName(f), "/", "_", 1)
		content, err := out.Func(f)
		if err != nil {
			return err
		}
		err = writeToFile(fmt.Sprintf("%s/%s.md", outputPath, name), content)
		if err != nil {
			return err
		}
	}
	return nil
}

func createRenderer(ctx map[string]interface{}) (*gomarkdoc.Renderer, error) {
	templateFunctionOpts := make([]gomarkdoc.RendererOption, 0)

	templateFunctionOpts = append(
		templateFunctionOpts,
		gomarkdoc.WithTemplateOverride("text", readTemplateMust("text.gotxt")),
		gomarkdoc.WithTemplateOverride("doc", readTemplateMust("doc.gotxt")),
		gomarkdoc.WithTemplateOverride("list", readTemplateMust("list.gotxt")),
		gomarkdoc.WithTemplateOverride("import", ""),
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
			return ctx
		}),
		gomarkdoc.WithTemplateFunc("functorName", functorName),
		gomarkdoc.WithTemplateFunc("bquote", bquote),
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

func isFuncAPredicate(f *lang.Func) bool {
	signature, err := f.Signature()
	if err != nil {
		return false
	}
	return f.Receiver() == "" && strings.HasSuffix(signature, "*engine.Promise")
}

func functorName(f *lang.Func) string {
	name := xstrings.ToSnakeCase(f.Name())
	signature, _ := f.Signature()
	arity := strings.Count(signature, ",") - 2
	return fmt.Sprintf("%s/%d", name, arity)
}

func bquote(str ...interface{}) string {
	out := make([]string, 0, len(str))
	for _, s := range str {
		if s != nil {
			out = append(out, fmt.Sprintf("`%v`", s))
		}
	}
	return strings.Join(out, " ")
}
