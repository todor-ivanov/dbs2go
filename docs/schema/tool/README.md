# DBS schema visualizer

This is the canonical, standalone replacement for the experiments in the
parent directory. It parses DBS DDL into a validated schema model and emits
deterministic JSON plus a GitHub-native Markdown/Mermaid atlas. The primary
documentation is `../README.md`; it needs no HTML, GitHub Pages, or external
renderer.

It uses GoSQLX v1.14.0 for the relational AST. Small, explicit compatibility
adapters preserve Oracle syntax GoSQLX cannot currently represent: named
column-level `NOT NULL` constraints, `ORGANIZATION INDEX`, and function-index
expressions. Comments and non-relational setup statements are handled by the
script layer. Unknown or unresolved relational statements are fatal; the tool
never presents a partial schema as complete.

## Run

From this directory:

```sh
go run . \
  -ddl ../../../static/schema/DDL/create-oracle-schema.sql \
  -dialect oracle \
  -out ../README \
  -emit markdown

go run . \
  -ddl ../../../static/schema/DDL/create-oracle-schema.sql \
  -dialect oracle \
  -out ../generated/dbs-oracle \
  -emit json
```

For MySQL, use `create-mysql-schema.sql` and `-dialect mysql`.

The Markdown contains a compact architecture view, a whole-schema overview,
five collapsible field-level relational maps, a grouped foreign-key registry,
collapsed table dictionaries, and non-table objects such as sequences, roles,
and grants. The JSON is the machine-readable source for later API-to-database
mapping work.

## Verify

```sh
go test ./...
go vet ./...
```

The integration tests parse both checked-in DBS DDL files and assert their full
table, column, and foreign-key counts. The root dbs2go Go module is not changed.
