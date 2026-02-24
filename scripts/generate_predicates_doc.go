package main

import (
	"embed"
	"errors"
	"fmt"
	"go/build"
	"io/fs"
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
	predicatesPath          = "x/logic/predicate"
	prologPredicatesPath    = "x/logic/lib"
	bootstrapPredicatesPath = "x/logic/interpreter/bootstrap"
	featuresPath            = "x/logic"
	outputPath              = "docs/predicate"
)

var (
	featureRegEx            = regexp.MustCompile(`^.+\.feature$`)
	prologPredicateHeadReg  = regexp.MustCompile(`^([a-z][A-Za-z0-9_]*)(?:\((.*)\))?\s*(?::-|\.)\s*$`)
	prologPredicateSigRegEx = regexp.MustCompile(`^([a-z][A-Za-z0-9_]*)(?:\((.*)\))?\s+is\s+[A-Za-z_]+\s*\.?$`)
)

type prologPredicateDocumentation struct {
	Predicate   string
	Signature   string
	Description string
	ModulePath  string
	IsBuiltin   bool
}

func generatePredicateDocumentation() error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	funcs, err := loadGoPredicates(wd)
	if err != nil {
		return err
	}

	features, err := loadFeatures(featuresPath)
	if err != nil {
		return err
	}

	written, err := generateGoPredicateDocs(funcs, features)
	if err != nil {
		return err
	}

	return generatePrologPredicateDocs(wd, len(funcs), features, written)
}

func generateGoPredicateDocs(
	funcs []*lang.Func, features map[string]*messages.Feature,
) (map[string]struct{}, error) {
	written := make(map[string]struct{}, len(funcs))

	for idx, f := range funcs {
		name := predicateFileName(functorName(f))

		// globalCtx used to keep track of contexts between templates.
		// (yes it's a hack).
		globalCtx := make(map[string]any)
		globalCtx["frontmatter"] = map[string]any{
			"sidebar_position": idx + 1,
		}
		globalCtx["feature"] = features[name]

		out, err := createRenderer(globalCtx)
		if err != nil {
			return nil, err
		}

		content, err := out.Func(f)
		if err != nil {
			return nil, err
		}
		err = writeToFile(fmt.Sprintf("%s/%s.md", outputPath, name), content)
		if err != nil {
			return nil, err
		}
		written[name] = struct{}{}
	}

	return written, nil
}

func generatePrologPredicateDocs(
	wd string, goPredicateCount int, features map[string]*messages.Feature, written map[string]struct{},
) error {
	builtinPredicates, err := loadPrologPredicates(path.Join(wd, bootstrapPredicatesPath), true)
	if err != nil {
		return err
	}
	libPredicates, err := loadPrologPredicates(path.Join(wd, prologPredicatesPath), false)
	if err != nil {
		return err
	}

	builtinPredicates = append(builtinPredicates, libPredicates...)
	slices.SortFunc(builtinPredicates, func(a, b prologPredicateDocumentation) int {
		return strings.Compare(a.Predicate, b.Predicate)
	})

	for idx, predicate := range builtinPredicates {
		name := predicateFileName(predicate.Predicate)
		if _, exists := written[name]; exists {
			return fmt.Errorf("predicate %s is documented in both Go and Prolog sources", predicate.Predicate)
		}

		content := renderPrologPredicateMarkdown(predicate, goPredicateCount+idx+1, features[name])
		if err := writeToFile(fmt.Sprintf("%s/%s.md", outputPath, name), content); err != nil {
			return err
		}
		written[name] = struct{}{}
	}

	return nil
}

func loadGoPredicates(wd string) ([]*lang.Func, error) {
	buildPkg, err := build.ImportDir(path.Join(wd, predicatesPath), build.ImportComment)
	if err != nil {
		return nil, err
	}

	log := logger.New(logger.DebugLevel)
	pkg, err := lang.NewPackageFromBuild(log, buildPkg)
	if err != nil {
		return nil, err
	}

	funcs := lo.Filter(pkg.Funcs(), func(item *lang.Func, _ int) bool { return isFuncAPredicate(item) })
	slices.SortFunc(funcs, func(a *lang.Func, j *lang.Func) int {
		return strings.Compare(a.Name(), j.Name())
	})

	return funcs, nil
}

func loadPrologPredicates(root string, isBuiltin bool) ([]prologPredicateDocumentation, error) {
	docs := make(map[string]prologPredicateDocumentation)

	err := filepath.Walk(root, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || filepath.Ext(info.Name()) != ".pl" {
			return nil
		}

		modulePath, err := filepath.Rel(root, filePath)
		if err != nil {
			return err
		}
		modulePath = filepath.ToSlash(modulePath)
		modulePath = strings.TrimPrefix(modulePath, "./")

		fileDocs, err := parsePrologPredicateDocs(filePath, modulePath)
		if err != nil {
			return err
		}
		for _, predicateDoc := range fileDocs {
			if _, exists := docs[predicateDoc.Predicate]; exists {
				return fmt.Errorf("duplicate Prolog documentation for predicate %s", predicateDoc.Predicate)
			}
			predicateDoc.IsBuiltin = isBuiltin
			docs[predicateDoc.Predicate] = predicateDoc
		}

		return nil
	})
	if errors.Is(err, fs.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return lo.Values(docs), nil
}

//nolint:funlen
func parsePrologPredicateDocs(filePath, modulePath string) ([]prologPredicateDocumentation, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(content), "\n")
	docs := make([]prologPredicateDocumentation, 0)
	seen := make(map[string]struct{})

	for i := 0; i < len(lines); i++ {
		if !isPrologDocCommentLine(lines[i]) {
			continue
		}

		commentBlock := make([]string, 0)
		commentBlock = append(commentBlock, prologCommentLine(lines[i]))
		i++

		for i < len(lines) && isPrologCommentLine(lines[i]) {
			commentBlock = append(commentBlock, prologCommentLine(lines[i]))
			i++
		}

		nextLine := i
		for nextLine < len(lines) && strings.TrimSpace(lines[nextLine]) == "" {
			nextLine++
		}
		if nextLine >= len(lines) {
			break
		}

		name, arity, ok := parsePrologPredicateHead(lines[nextLine])
		if !ok {
			i = nextLine - 1
			continue
		}

		signature := prologPredicateSignature(commentBlock, name, arity)
		description := prologPredicateDescription(commentBlock, signature)
		if signature == "" && description == "" {
			i = nextLine
			continue
		}

		predicateName := fmt.Sprintf("%s/%d", name, arity)
		if _, exists := seen[predicateName]; exists {
			i = nextLine
			continue
		}
		seen[predicateName] = struct{}{}

		docs = append(docs, prologPredicateDocumentation{
			Predicate:   predicateName,
			Signature:   signature,
			Description: description,
			ModulePath:  modulePath,
		})

		i = nextLine
	}

	return docs, nil
}

func isPrologCommentLine(line string) bool {
	return strings.HasPrefix(strings.TrimSpace(line), "%")
}

func isPrologDocCommentLine(line string) bool {
	return strings.HasPrefix(strings.TrimSpace(line), "%!")
}

func prologCommentLine(line string) string {
	trimmed := strings.TrimSpace(line)
	// Remove %! for doc comments, or just % for continuation lines
	if strings.HasPrefix(trimmed, "%!") {
		trimmed = strings.TrimPrefix(trimmed, "%!")
	} else {
		trimmed = strings.TrimPrefix(trimmed, "%")
	}
	return strings.TrimPrefix(trimmed, " ")
}

func parsePrologPredicateHead(line string) (string, int, bool) {
	matches := prologPredicateHeadReg.FindStringSubmatch(strings.TrimSpace(line))
	if len(matches) == 0 {
		return "", 0, false
	}

	arguments := strings.TrimSpace(matches[2])
	arity := 0
	if arguments != "" {
		arity = prologArity(arguments)
	}

	return matches[1], arity, true
}

func parsePrologPredicateSignature(line string) (string, int, bool) {
	matches := prologPredicateSigRegEx.FindStringSubmatch(strings.TrimSpace(line))
	if len(matches) == 0 {
		return "", 0, false
	}

	arguments := strings.TrimSpace(matches[2])
	arity := 0
	if arguments != "" {
		arity = prologArity(arguments)
	}

	return matches[1], arity, true
}

func prologPredicateSignature(commentBlock []string, name string, arity int) string {
	for _, line := range commentBlock {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}

		sigName, sigArity, ok := parsePrologPredicateSignature(trimmed)
		if !ok || sigName != name || sigArity != arity {
			continue
		}

		return strings.TrimSuffix(trimmed, ".")
	}

	return ""
}

func prologPredicateDescription(commentBlock []string, signature string) string {
	descriptionLines := make([]string, 0, len(commentBlock))
	signature = strings.TrimSuffix(strings.TrimSpace(signature), ".")

	for _, line := range commentBlock {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			descriptionLines = append(descriptionLines, "")
			continue
		}

		candidate := strings.TrimSuffix(trimmed, ".")
		if signature != "" && candidate == signature {
			continue
		}

		descriptionLines = append(descriptionLines, strings.TrimRight(line, "\r"))
	}

	descriptionLines = trimEmptyLines(descriptionLines)
	return strings.Join(descriptionLines, "\n")
}

func trimEmptyLines(lines []string) []string {
	start := 0
	for start < len(lines) && strings.TrimSpace(lines[start]) == "" {
		start++
	}

	end := len(lines) - 1
	for end >= start && strings.TrimSpace(lines[end]) == "" {
		end--
	}

	if start > end {
		return nil
	}

	return lines[start : end+1]
}

func prologArity(arguments string) int {
	arguments = strings.TrimSpace(arguments)
	if arguments == "" {
		return 0
	}

	state := prologArityState{arity: 1}
	for _, char := range arguments {
		state.consume(char)
	}

	return state.arity
}

type prologArityState struct {
	arity         int
	roundDepth    int
	squareDepth   int
	curlyDepth    int
	inSingleQuote bool
	inDoubleQuote bool
}

func (s *prologArityState) consume(char rune) {
	if s.consumeQuoted(char) {
		return
	}
	if s.consumeQuoteStart(char) {
		return
	}
	if char == ',' && s.roundDepth == 0 && s.squareDepth == 0 && s.curlyDepth == 0 {
		s.arity++
		return
	}

	s.consumeDepth(char)
}

func (s *prologArityState) consumeQuoted(char rune) bool {
	if s.inSingleQuote {
		if char == '\'' {
			s.inSingleQuote = false
		}
		return true
	}
	if s.inDoubleQuote {
		if char == '"' {
			s.inDoubleQuote = false
		}
		return true
	}

	return false
}

func (s *prologArityState) consumeQuoteStart(char rune) bool {
	if char == '\'' {
		s.inSingleQuote = true
		return true
	}
	if char == '"' {
		s.inDoubleQuote = true
		return true
	}

	return false
}

func (s *prologArityState) consumeDepth(char rune) {
	switch char {
	case '(':
		s.roundDepth++
	case ')':
		if s.roundDepth > 0 {
			s.roundDepth--
		}
	case '[':
		s.squareDepth++
	case ']':
		if s.squareDepth > 0 {
			s.squareDepth--
		}
	case '{':
		s.curlyDepth++
	case '}':
		if s.curlyDepth > 0 {
			s.curlyDepth--
		}
	}
}

func predicateFileName(predicate string) string {
	return strings.Replace(predicate, "/", "_", 1)
}

func renderPrologPredicateMarkdown(
	predicateDoc prologPredicateDocumentation, sidebarPosition int, feature *messages.Feature,
) string {
	var out strings.Builder
	out.WriteString(fmt.Sprintf("---\nsidebar_position: %d\n---\n", sidebarPosition))
	out.WriteString("[//]: # (This file is auto-generated. Please do not modify it yourself.)\n\n")
	out.WriteString(fmt.Sprintf("# %s\n\n", predicateDoc.Predicate))

	out.WriteString("## Module\n\n")
	if predicateDoc.IsBuiltin {
		out.WriteString("Built-in predicate.\n")
	} else {
		out.WriteString(fmt.Sprintf("This predicate is provided by `%s`.\n\n", predicateDoc.ModulePath))
		out.WriteString("Load this module before using the predicate:\n\n")
		out.WriteString("```prolog\n")
		out.WriteString(fmt.Sprintf(":- consult('/v1/lib/%s').\n", predicateDoc.ModulePath))
		out.WriteString("```\n")
	}

	out.WriteString("## Description\n\n")
	description := strings.TrimSpace(predicateDoc.Description)
	if description == "" {
		description = fmt.Sprintf("`%s` is a predicate.", predicateDoc.Predicate)
	}
	out.WriteString(description + "\n")

	if predicateDoc.Signature != "" {
		out.WriteString("\n## Signature\n\n")
		out.WriteString("```text\n")
		out.WriteString(strings.TrimSuffix(predicateDoc.Signature, "."))
		out.WriteString("\n```\n")
	}

	examples := renderFeatureExamples(feature)
	if examples != "" {
		out.WriteString("\n")
		out.WriteString(examples)
		out.WriteString("\n")
	}

	return out.String()
}

func renderFeatureExamples(feature *messages.Feature) string {
	if feature == nil {
		return ""
	}

	scenarioDocs := make([]string, 0)
	for _, child := range feature.Children {
		scenario := child.Scenario
		if scenario == nil || !tagged("@great_for_documentation", scenario) {
			continue
		}

		var out strings.Builder
		out.WriteString(fmt.Sprintf("### %s\n\n", scenario.Name))

		description := strings.TrimSpace(dedent.String(scenario.Description))
		if description != "" {
			out.WriteString(description + "\n\n")
		}

		out.WriteString("Here are the steps of the scenario:\n\n")
		for _, step := range scenario.Steps {
			out.WriteString(fmt.Sprintf("- **%s** %s\n", strings.TrimSpace(step.Keyword), step.Text))

			if step.DocString != nil {
				mediaType := step.DocString.MediaType
				if mediaType == "" {
					mediaType = "text"
				}
				out.WriteString(fmt.Sprintf("\n``` %s\n%s\n```\n", mediaType, step.DocString.Content))
			}

			if step.DataTable != nil {
				out.WriteString("\n| key | value |\n| --- | ----- |\n")
				for _, row := range step.DataTable.Rows {
					if len(row.Cells) < 2 {
						continue
					}
					out.WriteString(fmt.Sprintf("| %s | %s |\n", row.Cells[0].Value, row.Cells[1].Value))
				}
			}
		}

		scenarioDocs = append(scenarioDocs, strings.TrimRight(out.String(), "\n"))
	}

	if len(scenarioDocs) == 0 {
		return ""
	}

	return "## Examples\n\n" + strings.Join(scenarioDocs, "\n\n")
}

func createRenderer(ctx map[string]any) (*gomarkdoc.Renderer, error) {
	sprigFuncs := sprig.GenericFuncMap()
	templateFunctionOpts := make([]gomarkdoc.RendererOption, 0, 13+len(sprigFuncs))

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

	for k, v := range sprigFuncs {
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
		gomarkdoc.WithTemplateFunc("globalCtx", func() map[string]any {
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

func bquote(str ...any) string {
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
