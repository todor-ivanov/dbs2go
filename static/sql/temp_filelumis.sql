CREATE PRIVATE TEMPORARY TABLE {{.Owner}}.TEMP_FILE_LUMIS
AS SELECT * FROM {{.Owner}}.FILE_LUMIS
ON COMMIT DROP DEFINITION