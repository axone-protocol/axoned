package main

import (
	"encoding/json"
	"path/filepath"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestPredicateSectionName(t *testing.T) {
	Convey("Given Prolog predicate documentation", t, func() {
		Convey("built-in bootstrap predicates are grouped in stdlib", func() {
			got := predicateSectionName(prologPredicateDocumentation{
				ModulePath: "bootstrap.pl",
				IsBuiltin:  true,
			})

			So(got, ShouldEqual, stdlibSection)
		})

		Convey("built-in stdlib predicates are grouped in stdlib", func() {
			got := predicateSectionName(prologPredicateDocumentation{
				ModulePath: "stdlib.pl",
				IsBuiltin:  true,
			})

			So(got, ShouldEqual, stdlibSection)
		})

		Convey("library predicates use the module name", func() {
			got := predicateSectionName(prologPredicateDocumentation{
				ModulePath: "apply.pl",
			})

			So(got, ShouldEqual, "apply")
		})

		Convey("nested module paths use the file name", func() {
			got := predicateSectionName(prologPredicateDocumentation{
				ModulePath: "helpers/bech32.pl",
			})

			So(got, ShouldEqual, "bech32")
		})
	})
}

func TestBuildPredicateSections(t *testing.T) {
	Convey("Given predicates in multiple sections", t, func() {
		sections := buildPredicateSections([]predicateEntry{
			{name: "zeta/1", section: "bank"},
			{name: "beta/1", section: stdlibSection},
			{name: "alpha/1", section: "bank"},
			{name: "gamma/1", section: "apply"},
		})

		Convey("sections are ordered with stdlib first", func() {
			So(sections, ShouldHaveLength, 3)
			So(sections[0].name, ShouldEqual, stdlibSection)
			So(sections[1].name, ShouldEqual, "apply")
			So(sections[2].name, ShouldEqual, "bank")
		})

		Convey("predicates are sorted inside each section", func() {
			So(sections[2].predicates, ShouldHaveLength, 2)
			So(sections[2].predicates[0].name, ShouldEqual, "alpha/1")
			So(sections[2].predicates[1].name, ShouldEqual, "zeta/1")
		})
	})
}

func TestPredicateDocPath(t *testing.T) {
	Convey("Given a section and predicate file name", t, func() {
		got := predicateDocPath("apply", "foldl_4")

		So(got, ShouldEqual, filepath.Join(outputPath, "apply", "foldl_4.md"))
	})
}

func TestRenderPredicateCategory(t *testing.T) {
	Convey("Given a category label and position", t, func() {
		content, err := renderPredicateCategory("Apply library", 2)

		So(err, ShouldBeNil)

		var metadata struct {
			Label    string `json:"label"`
			Position int    `json:"position"`
		}
		err = json.Unmarshal([]byte(content), &metadata)

		So(err, ShouldBeNil)
		So(metadata.Label, ShouldEqual, "Apply library")
		So(metadata.Position, ShouldEqual, 2)
	})
}

func TestPredicateSectionLabel(t *testing.T) {
	Convey("Given a predicate section", t, func() {
		Convey("stdlib uses a reader-facing label", func() {
			So(predicateSectionLabel(stdlibSection), ShouldEqual, "Standard library")
		})

		Convey("library sections are humanized", func() {
			So(predicateSectionLabel("apply"), ShouldEqual, "Apply library")
			So(predicateSectionLabel("bech32"), ShouldEqual, "Bech32 library")
		})
	})
}

func TestWriteToFileCreatesParentDirectories(t *testing.T) {
	Convey("Given a target file in a missing directory tree", t, func() {
		target := filepath.Join(t.TempDir(), "docs", "predicate", "apply", "foldl_4.md")
		err := writeToFile(target, "content")

		So(err, ShouldBeNil)
	})
}
