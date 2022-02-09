{{.TokenGenerator}}
{{if .Runs}}
SELECT DISTINCT
{{else}}
SELECT
{{end}}
{{if .Detail}}
        D.DATASET_ID, D.DATASET, D.PREP_ID, 
        D.XTCROSSSECTION, 
        D.CREATION_DATE, D.CREATE_BY, 
        D.LAST_MODIFICATION_DATE, D.LAST_MODIFIED_BY,
        P.PRIMARY_DS_NAME,
        PDT.PRIMARY_DS_TYPE,
        PD.PROCESSED_DS_NAME,
        DT.DATA_TIER_NAME,
        DP.DATASET_ACCESS_TYPE,
        AE.ACQUISITION_ERA_NAME,
        PE.PROCESSING_VERSION,
        PH.PHYSICS_GROUP_NAME 
{{if .ParentDataset}}
        ,PDS.DATASET PARENT_DATASET
{{end}}
{{if .Version}}
        ,OMC.OUTPUT_MODULE_LABEL
        ,OMC.GLOBAL_TAG
        ,RV.RELEASE_VERSION,
        ,PSH.PSET_HASH,
        ,AEX.APP_NAME
{{end}}
{{else}}
        D.DATASET
{{end}}
       
FROM {{.Owner}}.DATASETS D
{{if .Lfns}}
JOIN {{.Owner}}.FILES FL on FL.DATASET_ID = D.DATASET_ID
{{end}}
{{if .Runs}}
{{if .Lfns}}
JOIN {{.Owner}}.FILE_LUMIS FLLU on FLLU.FILE_ID=FL.FILE_ID
{{else}}
JOIN {{.Owner}}.FILES FL on FL.DATASET_ID = D.DATASET_ID
JOIN {{.Owner}}.FILE_LUMIS FLLU on FLLU.FILE_ID=FL.FILE_ID
{{end}}
{{end}}
JOIN {{.Owner}}.PRIMARY_DATASETS P ON P.PRIMARY_DS_ID = D.PRIMARY_DS_ID
JOIN {{.Owner}}.PRIMARY_DS_TYPES PDT ON PDT.PRIMARY_DS_TYPE_ID = P.PRIMARY_DS_TYPE_ID
JOIN {{.Owner}}.PROCESSED_DATASETS PD ON PD.PROCESSED_DS_ID = D.PROCESSED_DS_ID
JOIN {{.Owner}}.DATA_TIERS DT ON DT.DATA_TIER_ID = D.DATA_TIER_ID
JOIN {{.Owner}}.DATASET_ACCESS_TYPES DP on DP.DATASET_ACCESS_TYPE_ID= D.DATASET_ACCESS_TYPE_ID

LEFT OUTER JOIN {{.Owner}}.ACQUISITION_ERAS AE ON AE.ACQUISITION_ERA_ID = D.ACQUISITION_ERA_ID
LEFT OUTER JOIN {{.Owner}}.PROCESSING_ERAS PE ON PE.PROCESSING_ERA_ID = D.PROCESSING_ERA_ID
LEFT OUTER JOIN {{.Owner}}.PHYSICS_GROUPS PH ON PH.PHYSICS_GROUP_ID = D.PHYSICS_GROUP_ID
{{if .ParentDataset}}
LEFT OUTER JOIN {{.Owner}}.DATASET_PARENTS DSP ON DSP.THIS_DATASET_ID = D.DATASET_ID
LEFT OUTER JOIN {{.Owner}}.DATASETS PDS ON PDS.DATASET_ID = DSP.PARENT_DATASET_ID
{{end}}

{{if .Version}}
LEFT OUTER JOIN {{.Owner}}.DATASET_OUTPUT_MOD_CONFIGS DOMC ON DOMC.DATASET_ID = D.DATASET_ID
LEFT OUTER JOIN {{.Owner}}.OUTPUT_MODULE_CONFIGS OMC ON OMC.OUTPUT_MOD_CONFIG_ID = DOMC.OUTPUT_MOD_CONFIG_ID
LEFT OUTER JOIN {{.Owner}}.RELEASE_VERSIONS RV ON RV.RELEASE_VERSION_ID = OMC.RELEASE_VERSION_ID
LEFT OUTER JOIN {{.Owner}}.PARAMETER_SET_HASHES PSH ON PSH.PARAMETER_SET_HASH_ID = OMC.PARAMETER_SET_HASH_ID
LEFT OUTER JOIN {{.Owner}}.APPLICATION_EXECUTABLES AEX ON AEX.APP_EXEC_ID = OMC.APP_EXEC_ID
{{end}}
