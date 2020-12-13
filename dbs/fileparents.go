package dbs

import (
	"errors"
	"fmt"
	"net/http"
)

// FileParents API
func (API) FileParents(params Record, w http.ResponseWriter) (int64, error) {
	// variables we'll use in where clause
	var args []interface{}
	where := "WHERE "

	// parse dataset argument
	fileparent := getValues(params, "logical_file_name")
	if len(fileparent) > 1 {
		msg := "The fileparent API does not support list of fileparent"
		return 0, errors.New(msg)
	} else if len(fileparent) == 1 {
		op, val := opVal(fileparent[0])
		cond := fmt.Sprintf(" F.LOGICAL_FILE_NAME %s %s", op, placeholder("logical_file_name"))
		where += addCond(where, cond)
		args = append(args, val)
	}
	// get SQL statement from static area
	stm := getSQL("fileparent")
	// use generic query API to fetch the results from DB
	return executeAll(w, stm+where, args...)
}

// InsertFileParents DBS API
func (API) InsertFileParents(values Record) error {
	return InsertData("insert_file_parents", values)
}