SELECT
    MIGRATION_BLOCK_ID,
    MIGRATION_BLOCK_NAME,
    MIGRATION_ORDER,
    MIGRATION_STATUS
FROM {{.Owner}}.MIGRATION_BLOCKS
WHERE MIGRATION_REQUEST_ID=:migration_request_id
AND (MIGRATION_STATUS=0 or MIGRATION_STATUS=1 or MIGRATION_STATUS=3) ORDER BY MIGRATION_ORDER DESC
