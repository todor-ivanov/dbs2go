package dbs

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/vkuznet/dbs2go/utils"
)

// Blocks DBS API
func (API) Blocks(params Record, w http.ResponseWriter) (int64, error) {
	var args []interface{}
	var conds []string
	tmpl := make(Record)
	tmpl["Owner"] = DBOWNER
	tmpl["TokenGenerator"] = ""

	// use run_num first since it may produce TokenGenerator
	// which should contain bind variables
	runs, err := ParseRuns(getValues(params, "run_num"))
	if err != nil {
		return 0, err
	}
	if len(runs) > 0 {
		tmpl["Runs"] = true
		token, whereRuns, bindsRuns := runsClause("FLM", runs)
		tmpl["TokenGenerator"] = token
		conds = append(conds, whereRuns)
		for _, v := range bindsRuns {
			args = append(args, v)
		}
	}
	// parse arguments
	lfns := getValues(params, "logical_file_name")
	if len(lfns) == 1 {
		tmpl["Lfns"] = true
		conds, args = AddParam("logical_file_name", "FL.LOGICAL_FILE_NAME", params, conds, args)
	}

	conds, args = AddParam("block_name", "B.BLOCK_NAME", params, conds, args)
	conds, args = AddParam("dataset", "DS.DATASET", params, conds, args)
	conds, args = AddParam("origin_site_name", "B.ORIGIN_SITE_NAME", params, conds, args)
	conds, args = AddParam("cdate", "B.CREATION_DATE", params, conds, args)

	minDate := getValues(params, "min_cdate")
	maxDate := getValues(params, "max_cdate")
	if len(minDate) == 1 && len(maxDate) == 1 {
		_, minval := OperatorValue(minDate[0])
		_, maxval := OperatorValue(maxDate[0])
		if minval != "0" && maxval != "0" {
			cond := fmt.Sprintf(" B.CREATION_DATE BETWEEN %s and %s", placeholder("min_cdate"), placeholder("max_cdate"))
			conds = append(conds, cond)
			args = append(args, minval)
			args = append(args, maxval)
		} else if minval != "0" && maxval == "0" {
			cond := fmt.Sprintf(" B.CREATION_DATE > %s", placeholder("min_cdate"))
			conds = append(conds, cond)
			args = append(args, minval)
		} else if minval == "0" && maxval != "0" {
			cond := fmt.Sprintf(" B.CREATION_DATE < %s", placeholder("max_cdate"))
			conds = append(conds, cond)
			args = append(args, maxval)
		}
	}

	conds, args = AddParam("ldate", "B.LAST_MODIFICATION_DATE", params, conds, args)

	minDate = getValues(params, "min_ldate")
	maxDate = getValues(params, "max_ldate")
	if len(minDate) == 1 && len(maxDate) == 1 {
		_, minval := OperatorValue(minDate[0])
		_, maxval := OperatorValue(maxDate[0])
		if minval != "0" && maxval != "0" {
			cond := fmt.Sprintf(" B.CREATION_DATE BETWEEN %s and %s", placeholder("min_ldate"), placeholder("max_ldate"))
			conds = append(conds, cond)
			args = append(args, minval)
			args = append(args, maxval)
		} else if minval != "0" && maxval == "0" {
			cond := fmt.Sprintf(" B.CREATION_DATE > %s", placeholder("min_ldate"))
			conds = append(conds, cond)
			args = append(args, minval)
		} else if minval == "0" && maxval != "0" {
			cond := fmt.Sprintf(" B.CREATION_DATE < %s", placeholder("max_ldate"))
			conds = append(conds, cond)
			args = append(args, maxval)
		}
	}
	stm, err := LoadTemplateSQL("blocks", tmpl)
	if err != nil {
		return 0, err
	}
	stm = WhereClause(stm, conds)

	// use generic query API to fetch the results from DB
	return executeAll(w, stm, args...)
}

// Blocks
type Blocks struct {
	BLOCK_ID               int64  `json:"block_id"`
	BLOCK_NAME             string `json:"block_name"`
	DATASET_ID             int64  `json:"dataset_id"`
	OPEN_FOR_WRITING       int64  `json:"open_for_writing"`
	ORIGIN_SITE_NAME       string `json:"origin_site_name"`
	BLOCK_SIZE             int64  `json:"block_size"`
	FILE_COUNT             int64  `json:"file_count"`
	CREATION_DATE          int64  `json:"creation_date"`
	CREATE_BY              string `json:"create_by"`
	LAST_MODIFICATION_DATE int64  `json:"last_modification_date"`
	LAST_MODIFIED_BY       string `json:"last_modified_by"`
}

// Insert implementation of Blocks
func (r *Blocks) Insert(tx *sql.Tx) error {
	var tid int64
	var err error
	if r.DATASET_ID == 0 {
		if DBOWNER == "sqlite" {
			tid, err = LastInsertId(tx, "BLOCKS", "block_id")
			r.DATASET_ID = tid + 1
		} else {
			tid, err = IncrementSequence(tx, "SEQ_BK")
			r.DATASET_ID = tid
		}
		if err != nil {
			return err
		}
	}
	// get SQL statement from static area
	stm := getSQL("insert_blocks")
	if DBOWNER == "sqlite" {
		stm = getSQL("insert_blocks_sqlite")
	}
	if utils.VERBOSE > 0 {
		log.Printf("Insert Blocks\n%s\n%+v", stm, r)
	}
	// validate our record
	err = r.Validate()
	if err != nil {
		return err
	}
	_, err = tx.Exec(stm, r.BLOCK_ID, r.BLOCK_NAME, r.DATASET_ID, r.OPEN_FOR_WRITING, r.ORIGIN_SITE_NAME, r.BLOCK_SIZE, r.FILE_COUNT, r.CREATION_DATE, r.CREATE_BY, r.LAST_MODIFICATION_DATE, r.LAST_MODIFIED_BY)
	return err
}

// Validate implementation of Blocks
func (r *Blocks) Validate() error {
	if matched := blockPattern.MatchString(r.BLOCK_NAME); !matched {
		log.Println("validate Block", r)
		return errors.New("invalid pattern for block")
	}
	if matched := unixTimePattern.MatchString(fmt.Sprintf("%d", r.CREATION_DATE)); !matched {
		return errors.New("invalid pattern for createion date")
	}
	if r.CREATION_DATE == 0 {
		return errors.New("missing creation_date")
	}
	if r.CREATE_BY == "" {
		return errors.New("missing create_by")
	}
	if r.LAST_MODIFICATION_DATE == 0 {
		return errors.New("missing last_modification_date")
	}
	if r.LAST_MODIFIED_BY == "" {
		return errors.New("missing last_modified_by")
	}
	return nil
}

// Decode implementation for Blocks
func (r *Blocks) Decode(reader io.Reader) (int64, error) {
	// init record with given data record
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Println("fail to read data", err)
		return 0, err
	}
	err = json.Unmarshal(data, &r)

	//     decoder := json.NewDecoder(r)
	//     err := decoder.Decode(&rec)
	if err != nil {
		log.Println("fail to decode data", err)
		return 0, err
	}
	size := int64(len(data))
	return size, nil
}

// BlockRecord represent input recor for insert blocks API
type BlockRecord struct {
	BLOCK_NAME             string `json:"block_name"`
	OPEN_FOR_WRITING       int64  `json:"open_for_writing"`
	ORIGIN_SITE_NAME       string `json:"origin_site_name"`
	BLOCK_SIZE             int64  `json:"block_size"`
	FILE_COUNT             int64  `json:"file_count"`
	CREATION_DATE          int64  `json:"creation_date"`
	CREATE_BY              string `json:"create_by"`
	LAST_MODIFICATION_DATE int64  `json:"last_modification_date"`
	LAST_MODIFIED_BY       string `json:"last_modified_by"`
}

// InsertBlocks DBS API
func (API) InsertBlocks(r io.Reader) (int64, error) {
	// TODO: implement the following logic
	// input values: blockname
	// optional values: open_for_writing, origin_site(name), block_size, file_count, creation_date, create_by, last_modification_date, last_modified_by
	// blkinput["dataset_id"] = self.datasetid.execute(conn,  ds_name, tran)
	// blkinput["block_id"] =  self.sm.increment(conn, "SEQ_BK", tran)
	// self.blockin.execute(conn, blkinput, tran)

	// read given input
	data, err := ioutil.ReadAll(r)
	if err != nil {
		log.Println("fail to read data", err)
		return 0, err
	}
	size := int64(len(data))
	var rec BlockRecord
	err = json.Unmarshal(data, &rec)
	if err != nil {
		log.Println("fail to decode data", err)
		return 0, err
	}
	// set dependent's records
	brec := Blocks{BLOCK_NAME: rec.BLOCK_NAME, OPEN_FOR_WRITING: rec.OPEN_FOR_WRITING, ORIGIN_SITE_NAME: rec.ORIGIN_SITE_NAME, BLOCK_SIZE: rec.BLOCK_SIZE, FILE_COUNT: rec.FILE_COUNT, CREATION_DATE: rec.CREATION_DATE, CREATE_BY: rec.CREATE_BY, LAST_MODIFICATION_DATE: rec.LAST_MODIFICATION_DATE, LAST_MODIFIED_BY: rec.LAST_MODIFIED_BY}

	// start transaction
	tx, err := DB.Begin()
	if err != nil {
		msg := fmt.Sprintf("unable to get DB transaction %v", err)
		return 0, errors.New(msg)
	}
	defer tx.Rollback()

	// get all necessary IDs from different tables
	dataset := strings.Split(rec.BLOCK_NAME, "#")[0]
	dsId, err := getTxtID(tx, "DATASETS", "dataset_id", "dataset", dataset)
	if err != nil {
		log.Println("unable to find dataset_id for", dataset)
		return 0, err
	}

	// assign all Id's in dataset DB record
	brec.DATASET_ID = dsId
	err = brec.Insert(tx)
	if err != nil {
		return 0, err
	}

	// commit transaction
	err = tx.Commit()
	if err != nil {
		log.Println("faile to insert_outputconfigs_sqlite", err)
		return 0, err
	}
	return size, err
}

// InsertBulkBlocks DBS API
func (API) InsertBulkBlocks(values Record) error {
	// TODO: implement the following logic
	// /Users/vk/CMS/DMWM/GIT/DBS/Server/Python/src/dbs/business/DBSBlock.py
	/*
	   #1 insert configuration
	   configList = self.insertOutputModuleConfig(
	                   blockcontent['dataset_conf_list'], migration)
	   #2 insert dataset
	   datasetId = self.insertDataset(blockcontent, configList, migration)
	   #3 insert block & files
	   self.insertBlockFile(blockcontent, datasetId, migration)
	*/

	//     return InsertValues("insert_blocks", values)
	return nil
}
