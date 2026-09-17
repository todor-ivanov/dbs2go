package main

type SourceInfo struct {
	File    string `json:"file"`
	Dialect string `json:"dialect"`
	SHA256  string `json:"sha256"`
}

type ParserInfo struct {
	Name     string   `json:"name"`
	Version  string   `json:"version"`
	Adapters []string `json:"compatibility_adapters,omitempty"`
}

type Schema struct {
	FormatVersion int               `json:"format_version"`
	Source        SourceInfo        `json:"source"`
	Parser        ParserInfo        `json:"parser"`
	Tables        []Table           `json:"tables"`
	ForeignKeys   []ForeignKey      `json:"foreign_keys"`
	Auxiliary     []AuxiliaryObject `json:"auxiliary_objects,omitempty"`
}

type Table struct {
	Name         string       `json:"name"`
	Columns      []Column     `json:"columns"`
	Constraints  []Constraint `json:"constraints,omitempty"`
	Indexes      []Index      `json:"indexes,omitempty"`
	Organization string       `json:"organization,omitempty"`
	SourceLine   int          `json:"source_line"`
}

type Column struct {
	Name        string             `json:"name"`
	Type        string             `json:"type"`
	Nullable    bool               `json:"nullable"`
	Default     string             `json:"default,omitempty"`
	Constraints []ColumnConstraint `json:"constraints,omitempty"`
}

type ColumnConstraint struct {
	Name string `json:"name,omitempty"`
	Type string `json:"type"`
}

type Constraint struct {
	Name                 string   `json:"name,omitempty"`
	Type                 string   `json:"type"`
	Columns              []string `json:"columns,omitempty"`
	Expression           string   `json:"expression,omitempty"`
	ReferencedTable      string   `json:"referenced_table,omitempty"`
	ReferencedColumns    []string `json:"referenced_columns,omitempty"`
	ReferencedConstraint string   `json:"referenced_constraint,omitempty"`
	OnDelete             string   `json:"on_delete,omitempty"`
	OnUpdate             string   `json:"on_update,omitempty"`
	SourceLine           int      `json:"source_line"`
}

type Index struct {
	Name       string      `json:"name"`
	Unique     bool        `json:"unique"`
	Method     string      `json:"method,omitempty"`
	Parts      []IndexPart `json:"parts"`
	Where      string      `json:"where,omitempty"`
	SourceLine int         `json:"source_line"`
}

type IndexPart struct {
	Column     string `json:"column,omitempty"`
	Expression string `json:"expression,omitempty"`
	Direction  string `json:"direction,omitempty"`
	Collation  string `json:"collation,omitempty"`
}

type ForeignKey struct {
	Constraint           string   `json:"constraint"`
	SourceTable          string   `json:"source_table"`
	SourceColumns        []string `json:"source_columns"`
	TargetTable          string   `json:"target_table"`
	TargetColumns        []string `json:"target_columns"`
	ReferencedConstraint string   `json:"referenced_constraint"`
	OnDelete             string   `json:"on_delete,omitempty"`
	OnUpdate             string   `json:"on_update,omitempty"`
	SourceLine           int      `json:"source_line"`
}

type AuxiliaryObject struct {
	Kind       string `json:"kind"`
	Name       string `json:"name,omitempty"`
	SQL        string `json:"sql"`
	SourceLine int    `json:"source_line"`
}

type statement struct {
	SQL  string
	Line int
}
