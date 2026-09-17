package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestOracleCompatibilityAdapters(t *testing.T) {
	ddl := []byte(`
CREATE TABLE PARENT (
  ID INTEGER CONSTRAINT NN_PARENT_ID NOT NULL,
  CODE VARCHAR2(20) DEFAULT 'x' CONSTRAINT NN_PARENT_CODE NOT NULL,
  CONSTRAINT PK_PARENT PRIMARY KEY (ID),
  CONSTRAINT UQ_PARENT_CODE UNIQUE (CODE),
  CONSTRAINT CK_PARENT CHECK (ID > 0)
) ORGANIZATION INDEX;
CREATE TABLE CHILD (
  ID INTEGER CONSTRAINT NN_CHILD_ID NOT NULL,
  PARENT_ID INTEGER,
  CONSTRAINT PK_CHILD PRIMARY KEY (ID)
);
ALTER TABLE CHILD ADD CONSTRAINT FK_CHILD_PARENT
  FOREIGN KEY (PARENT_ID) REFERENCES PARENT (ID) ON DELETE CASCADE;
CREATE UNIQUE INDEX IDX_PARENT_CODE_CI ON PARENT (NLSSORT("CODE", 'nls_sort=''BINARY_CI'''));
CREATE SEQUENCE SEQ_CHILD;
CREATE ROLE READER;
GRANT SELECT ON CHILD TO READER;
`)
	schema, err := parseSchema("fixture.sql", ddl, "oracle")
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if len(schema.Tables) != 2 || columnCount(schema) != 4 || len(schema.ForeignKeys) != 1 {
		t.Fatalf("unexpected relational counts: tables=%d columns=%d fks=%d", len(schema.Tables), columnCount(schema), len(schema.ForeignKeys))
	}
	if len(schema.Auxiliary) != 3 {
		t.Fatalf("auxiliary count = %d, want 3", len(schema.Auxiliary))
	}
	parent := schema.Tables[1]
	if parent.Name != "PARENT" {
		parent = schema.Tables[0]
	}
	if parent.Organization != "INDEX" {
		t.Fatalf("organization = %q", parent.Organization)
	}
	if parent.Columns[1].Default == "" || parent.Columns[1].Constraints[1].Name != "NN_PARENT_CODE" {
		t.Fatalf("named constraint/default not preserved: %#v", parent.Columns[1])
	}
	if got := parent.Indexes[0].Parts[0].Expression; got == "" {
		t.Fatal("function index expression was not preserved")
	}
	if schema.ForeignKeys[0].ReferencedConstraint != "PK_PARENT" || schema.ForeignKeys[0].OnDelete != "CASCADE" {
		t.Fatalf("foreign key semantics not resolved: %#v", schema.ForeignKeys[0])
	}
}

func TestMySQLHashComments(t *testing.T) {
	ddl := []byte("# heading with ; punctuation\nCREATE TABLE `T` (`ID` INTEGER NOT NULL, CONSTRAINT `PK_T` PRIMARY KEY (`ID`)) ENGINE = InnoDB;\n")
	schema, err := parseSchema("fixture.sql", ddl, "mysql")
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if len(schema.Tables) != 1 || len(schema.Tables[0].Columns) != 1 {
		t.Fatalf("hash comment corrupted statement splitting: %#v", schema.Tables)
	}
}

func TestInvalidReferenceIsFatal(t *testing.T) {
	ddl := []byte(`CREATE TABLE T (ID INTEGER, CONSTRAINT PK_T PRIMARY KEY (ID));
ALTER TABLE T ADD CONSTRAINT FK_BAD FOREIGN KEY (MISSING) REFERENCES T (ID);`)
	if _, err := parseSchema("bad.sql", ddl, "oracle"); err == nil {
		t.Fatal("invalid source column was accepted")
	}
}

func TestRealSchemasAndDeterminism(t *testing.T) {
	tests := []struct {
		file, dialect        string
		tables, columns, fks int
	}{
		{"../../../static/schema/DDL/create-oracle-schema.sql", "oracle", 28, 146, 31},
		{"../../../static/schema/DDL/create-mysql-schema.sql", "mysql", 32, 158, 33},
	}
	for _, test := range tests {
		t.Run(test.dialect, func(t *testing.T) {
			input, err := os.ReadFile(test.file)
			if err != nil {
				t.Fatal(err)
			}
			first, err := parseSchema(filepath.Base(test.file), input, test.dialect)
			if err != nil {
				t.Fatal(err)
			}
			if len(first.Tables) != test.tables || columnCount(first) != test.columns || len(first.ForeignKeys) != test.fks {
				t.Fatalf("counts: tables=%d columns=%d fks=%d", len(first.Tables), columnCount(first), len(first.ForeignKeys))
			}
			second, err := parseSchema(filepath.Base(test.file), input, test.dialect)
			if err != nil {
				t.Fatal(err)
			}
			one, _ := renderJSON(first)
			two, _ := renderJSON(second)
			if !bytes.Equal(one, two) {
				t.Fatal("JSON output is nondeterministic")
			}
		})
	}
}

func TestAtomicWriterReplacesCompleteFiles(t *testing.T) {
	dir := t.TempDir()
	one, two := filepath.Join(dir, "schema.json"), filepath.Join(dir, "schema.md")
	if err := os.WriteFile(one, []byte("old"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := writeOutputsAtomically(map[string][]byte{one: []byte("new-json"), two: []byte("new-md")}); err != nil {
		t.Fatal(err)
	}
	got, _ := os.ReadFile(one)
	if string(got) != "new-json" {
		t.Fatalf("first output = %q", got)
	}
	got, _ = os.ReadFile(two)
	if string(got) != "new-md" {
		t.Fatalf("second output = %q", got)
	}
}

func TestMarkdownAtlasIsCompactAndColumnPrecise(t *testing.T) {
	path := "../../../static/schema/DDL/create-oracle-schema.sql"
	input, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	schema, err := parseSchema(filepath.Base(path), input, "oracle")
	if err != nil {
		t.Fatal(err)
	}
	doc := renderMarkdown(schema)
	if got := strings.Count(doc, "```mermaid"); got != 7 {
		t.Fatalf("Mermaid diagram count = %d, want 7", got)
	}
	if strings.Count(doc, "<details") != strings.Count(doc, "</details>") {
		t.Fatal("unbalanced details sections")
	}
	for _, fk := range schema.ForeignKeys {
		for i, source := range fk.SourceColumns {
			edge := fmt.Sprintf("%s -->|\"%s ·", mermaidID(fk.SourceTable+"__"+source), fk.Constraint)
			target := mermaidID(fk.TargetTable + "__" + fk.TargetColumns[i])
			if !strings.Contains(doc, edge) || !strings.Contains(doc, "| "+target) {
				t.Fatalf("missing exact column edge for %s", fk.Constraint)
			}
		}
	}
	for _, heading := range []string{"## How to read this atlas", "## Architecture", "## Visualizations", "## Foreign-key registry", "## Table dictionary", "## Other database objects"} {
		if !strings.Contains(doc, heading) {
			t.Fatalf("missing section %q", heading)
		}
	}
}

func TestAllVisualizationModesAndNativeSVG(t *testing.T) {
	path := "../../../static/schema/DDL/create-oracle-schema.sql"
	input, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	schema, err := parseSchema(filepath.Base(path), input, "oracle")
	if err != nil {
		t.Fatal(err)
	}
	views, err := parseVisualizations("all")
	if err != nil {
		t.Fatal(err)
	}
	if len(views) != 5 {
		t.Fatalf("all selected %d modes, want 5", len(views))
	}
	svgs, err := renderNativeSVGs(schema, "compact")
	if err != nil {
		t.Fatal(err)
	}
	if len(svgs) != 5 {
		t.Fatalf("native SVG count = %d, want 5", len(svgs))
	}
	secondSVGs, err := renderNativeSVGs(schema, "compact")
	if err != nil {
		t.Fatal(err)
	}
	refs := map[string]string{}
	maxWidths := map[string]float64{"core": 1400, "class": 1000, "parent": 700, "config": 1000, "ops": 700}
	for key, data := range svgs {
		if !bytes.Equal(data, secondSVGs[key]) {
			t.Fatalf("%s SVG is nondeterministic", key)
		}
		refs[key] = "dbs-oracle." + key + ".svg"
		var document struct {
			XMLName xml.Name
			Width   float64 `xml:"width,attr"`
		}
		if err := xml.Unmarshal(data, &document); err != nil {
			t.Fatalf("%s SVG is invalid XML: %v", key, err)
		}
		if document.XMLName.Local != "svg" || !bytes.Contains(data, []byte(`marker-end="url(#fk-arrow)"`)) && len(foreignKeysForGroup(schema.ForeignKeys, key)) != 0 {
			t.Fatalf("%s SVG is missing its root or FK arrows", key)
		}
		if document.Width > maxWidths[key] {
			t.Fatalf("%s compact SVG width = %.0f, want at most %.0f", key, document.Width, maxWidths[key])
		}
		for _, fk := range foreignKeysForGroup(schema.ForeignKeys, key) {
			if !bytes.Contains(data, []byte(">"+fk.Constraint+"</tspan>")) || !bytes.Contains(data, []byte(">"+deleteRule(fk.OnDelete)+"</tspan>")) {
				t.Fatalf("%s SVG lost the label for %s", key, fk.Constraint)
			}
		}
	}
	doc := renderMarkdownWithOptions(schema, RenderOptions{Visualizations: views, Layout: "compact", SVGReferences: refs})
	for _, marker := range []string{"Full flowchart", "ER view", "Schema-wide key map", "Domain map:", "Go-native SVG:"} {
		if !strings.Contains(doc, marker) {
			t.Fatalf("all-mode document is missing %q", marker)
		}
	}
	if got := strings.Count(doc, "```mermaid"); got != 10 {
		t.Fatalf("all-mode Mermaid diagram count = %d, want 10", got)
	}
	if !strings.Contains(doc, `"nodeSpacing": 8`) {
		t.Fatal("compact spacing directive is missing")
	}
}

func TestVisualizationSelection(t *testing.T) {
	for _, name := range []string{"full", "er", "keys", "domains", "svg"} {
		views, err := parseVisualizations(name)
		if err != nil || len(views) != 1 || !views[name] {
			t.Fatalf("selection %q failed: views=%v err=%v", name, views, err)
		}
	}
	if _, err := parseVisualizations("unknown"); err == nil {
		t.Fatal("unknown visualization was accepted")
	}
}
