package main

import (
	"embed"
	"fmt"
	"go/build"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	"github.com/Masterminds/sprig/v3"
	gherkin "github.com/cucumber/gherkin/go/v26"
	messages "github.com/cucumber/messages/go/v21"
	"github.com/huandu/xstrings"
	"github.com/muesli/reflow/dedent"
	"github.com/princjef/gomarkdoc"
	"github.com/princjef/gomarkdoc/lang"
	"github.com/princjef/gomarkdoc/logger"
	"github.com/samber/lo"
)

//go:embed templates/*.go.txt
var f embed.FS

const (
	predicatesPath = "x/logic/predicate"
	featuresPath   = "x/logic"
	outputPath     = "docs/predicate"
)

var featureRegEx = regexp.MustCompile(`^.+\.feature$`)

func generatePredicateDocumentation() error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	buildPkg, err := build.ImportDir(path.Join(wd, predicatesPath), build.ImportComment)
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

	features, err := loadFeatures(featuresPath)
	if err != nil {
		return err
	}

	for idx, f := range funcs {
		name := strings.Replace(functorName(f), "/", "_", 1)

		// globalCtx used to keep track of contexts between templates.
		// (yes it's a hack).
		globalCtx := make(map[string]interface{})
		globalCtx["frontmatter"] = map[string]interface{}{
			"sidebar_position": idx + 1,
		}
		globalCtx["feature"] = features[name]

		out, err := createRenderer(globalCtx)
		if err != nil {
			return err
		}

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
		gomarkdoc.WithTemplateOverride("text", readTemplateMust("text.go.txt")),
		gomarkdoc.WithTemplateOverride("doc", readTemplateMust("doc.go.txt")),
		gomarkdoc.WithTemplateOverride("list", readTemplateMust("list.go.txt")),
		gomarkdoc.WithTemplateOverride("import", ""),
		gomarkdoc.WithTemplateOverride("file", ""),
		gomarkdoc.WithTemplateOverride("func", readTemplateMust("func.go.txt")),
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
		gomarkdoc.WithTemplateFunc("dedent", dedent.String),
		gomarkdoc.WithTemplateFunc("tagged", tagged),
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
	signature, _ := f.Signature()
	arity := strings.Count(signature, ",") - 2

	// TODO: This is a hack to get the name of the predicate. We remove the arity from the name if present  (e.g. `Open3` becomes `open`).
	name := strings.TrimSuffix(xstrings.ToSnakeCase(f.Name()), fmt.Sprintf("%d", arity))

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

func tagged(tag string, feature *messages.Scenario) bool {
	return lo.ContainsBy(feature.Tags, func(t *messages.Tag) bool {
		return t.Name == tag
	})
}

func loadFeatures(path string) (map[string]*messages.Feature, error) {
	features := make(map[string]*messages.Feature)
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err == nil && featureRegEx.MatchString(info.Name()) {
			bs, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			r := strings.NewReader(string(bs))
			gherkinDocument, err := gherkin.ParseGherkinDocument(r, (&messages.Incrementing{}).NewId)
			if err != nil {
				return err
			}
			features[strings.TrimSuffix(info.Name(), filepath.Ext(info.Name()))] = gherkinDocument.Feature
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return features, nil
}
