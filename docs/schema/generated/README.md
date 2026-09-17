# Generated schema reference

[`../README.md`](../README.md) is the GitHub-native DBS schema atlas generated
from the checked-in Oracle DDL. `dbs-oracle.json` is its validated
machine-readable model. The `dbs-oracle.*.svg` files are deterministic,
Go-native domain diagrams embedded by the atlas; they do not require Graphviz.

Regenerate from `docs/schema/tool`:

```sh
go run . \
  -ddl ../../../static/schema/DDL/create-oracle-schema.sql \
  -dialect oracle \
  -out ../README \
  -svg-out ../generated/dbs-oracle \
  -emit markdown \
  -visualizations all \
  -layout compact

go run . \
  -ddl ../../../static/schema/DDL/create-oracle-schema.sql \
  -dialect oracle \
  -out ../generated/dbs-oracle \
  -emit json
```
