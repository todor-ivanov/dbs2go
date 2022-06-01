package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/dmwm/dbs2go/dbs"
	"github.com/dmwm/dbs2go/web"
)

// this file contains logic for fileArray API
//
// the HTTP POST request body is defined by multiple structs defined in this file
// each of those structs are named based on the POST fields
//
// the basic HTTP response body is defined by fileResponse struct in test/int_files.go
// the detailed HTTP response body is defined by fileDetailResponse struct in test/int_files.go
// the HTTP response body for run_num param is defined by fileRunResponse struct in test/int_files.go
// the HTTP handlers and endpoints are defined in the EndpointTestCase struct defined in test/integration_cases.go

// fileArray request with dataset parameter
type fileArrayDatasetRequest struct {
	Dataset string `json:"dataset"`
}

// fileArray request with dataset and validFileOnly parameter
type fileArrayDatasetValidFileRequest struct {
	Dataset       string `json:"dataset"`
	ValidFileOnly string `json:"validFileOnly"`
}

// fileArray request with dataset, validFileOnly, detail and sumOverLumi parameter
type fileArrayDatasetValidFileDetailSumOverLumiRequest struct {
	Dataset       string `json:"dataset"`
	ValidFileOnly string `json:"validFileOnly"`
	Detail        string `json:"detail"`
	SumOverLumi   string `json:"sumOverLumi"`
}

// fileArray request with block_name parameter
type fileArrayBlockNameRequest struct {
	BlockName string `json:"block_name"`
}

// fileArray request with block_name and detail parameter
type fileArrayBlockNameDetailRequest struct {
	BlockName string `json:"block_name"`
	Detail    string `json:"detail"`
}

// fileArray request with block_name and validFileOnly parameter
type fileArrayBlockNameValidFileRequest struct {
	BlockName     string `json:"block_name"`
	ValidFileOnly string `json:"validFileOnly"`
}

// fileArray request with block_name and validFileOnly parameter
type fileArrayBlockNameDetailValidFileRequest struct {
	BlockName     string `json:"block_name"`
	Detail        string `json:"detail"`
	ValidFileOnly string `json:"validFileOnly"`
}

// fileArray request with block_name, run_num, and lumi_list parameter
type fileArrayBlockNameRunNumLumiListRequest struct {
	BlockName string `json:"block_name"`
	RunNum    string `json:"run_num"`
	LumiList  string `json:"lumi_list"`
}

// fileArray request with block_name, run_num, lumi_list, and detail parameter
type fileArrayBlockNameRunNumLumiListDetailRequest struct {
	BlockName string `json:"block_name"`
	RunNum    string `json:"run_num"`
	LumiList  string `json:"lumi_list"`
	Detail    string `json:"detail"`
}

// fileArray request with block_name, run_num, lumi_list, and validFileOnly parameter
type fileArrayBlockNameRunNumLumiListValidFileRequest struct {
	BlockName     string `json:"block_name"`
	RunNum        string `json:"run_num"`
	LumiList      string `json:"lumi_list"`
	ValidFileOnly string `json:"validFileOnly"`
}

// fileArray request with block_name, run_num, lumi_list, detail, and validFileOnly parameter
type fileArrayBlockNameRunNumLumiListValidFileDetailRequest struct {
	BlockName     string `json:"block_name"`
	RunNum        string `json:"run_num"`
	LumiList      string `json:"lumi_list"`
	ValidFileOnly string `json:"validFileOnly"`
	Detail        string `json:"detail"`
}

// fileArray request with block_name, run_num, lumi_list and sumOverLumi parameter
type fileArrayBlockNameRunNumLumiListSumOverLumiRequest struct {
	BlockName   string `json:"block_name"`
	RunNum      string `json:"run_num"`
	LumiList    string `json:"lumi_list"`
	SumOverLumi string `json:"sumOverLumi"`
}

// fileArray request with block_name, run_num, lumi_list and sumOverLumi parameter
type fileArrayBlockNameRunNumLumiListSumOverLumiDetailRequest struct {
	BlockName   string `json:"block_name"`
	RunNum      string `json:"run_num"`
	LumiList    string `json:"lumi_list"`
	SumOverLumi string `json:"sumOverLumi"`
	Detail      string `json:"detail"`
}

// fileArray request with block_name, run_num, sumOverLumi, detail parameter
type fileArrayBlockNameRunNumSumOverLumiDetailRequest struct {
	BlockName   string `json:"block_name"`
	RunNum      string `json:"run_num"`
	SumOverLumi string `json:"sumOverLumi"`
	Detail      string `json:"detail"`
}

// fileArray request with block_name, run_num, sumOverLumi parameter
type fileArrayBlockNameRunNumSumOverLumiRequest struct {
	BlockName   string `json:"block_name"`
	RunNum      string `json:"run_num"`
	SumOverLumi string `json:"sumOverLumi"`
}

// logical_file_name structs

// fileArray request with logical_file_name
type fileArrayLFNRequest struct {
	LogicalFileName string `json:"logical_file_name"`
}

// fileArray request with logical_file_name, validFileOnly
type fileArrayLFNValidFileRequest struct {
	LogicalFileName string `json:"logical_file_name"`
	ValidFileOnly   string `json:"validFileOnly"`
}

// fileArray request with logical_file_name, run_num, lumi_list
type fileArrayLFNRunNumLumiListRequest struct {
	LogicalFileName string `json:"logical_file_name"`
	RunNum          string `json:"run_num"`
	LumiList        string `json:"lumi_list"`
}

// fileArray request with logical_file_name, run_num, lumi_list, detail
type fileArrayLFNRunNumLumiListDetailRequest struct {
	LogicalFileName string `json:"logical_file_name"`
	RunNum          string `json:"run_num"`
	LumiList        string `json:"lumi_list"`
	Detail          string `json:"detail"`
}

// fileArray request with logical_file_name, run_num, lumi_list, validFileOnly
type fileArrayLFNRunNumLumiListValidFileRequest struct {
	LogicalFileName string `json:"logical_file_name"`
	RunNum          string `json:"run_num"`
	LumiList        string `json:"lumi_list"`
	ValidFileOnly   string `json:"validFileOnly"`
}

// fileArray request with logical_file_name, run_num, lumi_list, detail
type fileArrayLFNRunNumLumiListValidFileDetailRequest struct {
	LogicalFileName string `json:"logical_file_name"`
	RunNum          string `json:"run_num"`
	LumiList        string `json:"lumi_list"`
	ValidFileOnly   string `json:"validFileOnly"`
	Detail          string `json:"detail"`
}

// fileArray request with logical_file_name, run_num, lumi_list, sumOverLumi
type fileArrayLFNRunNumLumiListSumOverLumiRequest struct {
	LogicalFileName string `json:"logical_file_name"`
	RunNum          string `json:"run_num"`
	LumiList        string `json:"lumi_list"`
	SumOverLumi     string `json:"sumOverLumi"`
}

// fileArray request with logical_file_name, run_num, sumOverLumi, detail
type fileArrayLFNRunNumSumOverLumiDetailRequest struct {
	LogicalFileName string `json:"logical_file_name"`
	RunNum          string `json:"run_num"`
	SumOverLumi     string `json:"sumOverLumi"`
	Detail          string `json:"detail"`
}

// fileArray request with dataset, release_version
type fileArrayDatasetReleaseRequest struct {
	Dataset        string `json:"dataset"`
	ReleaseVersion string `json:"release_version"`
}

// fileArray request with dataset, release_version, validFileOnly
type fileArrayDatasetReleaseValidFileRequest struct {
	Dataset        string `json:"dataset"`
	ReleaseVersion string `json:"release_version"`
	ValidFileOnly  string `json:"validFileOnly"`
}

// fileArray request with dataset, output_module_config fields
type fileArrayDatasetOutputModRequest struct {
	Dataset           string `json:"dataset"`
	ReleaseVersion    string `json:"release_version"`
	PsetHash          string `json:"pset_hash"`
	AppName           string `json:"app_name"`
	OutputModuleLabel string `json:"output_module_label"`
}

// fileArray request with lfn, output_module_config fields
type fileArrayLFNOutputModRequest struct {
	LogicalFileName   string `json:"logical_file_name"`
	ReleaseVersion    string `json:"release_version"`
	PsetHash          string `json:"pset_hash"`
	AppName           string `json:"app_name"`
	OutputModuleLabel string `json:"output_module_label"`
}

// fileArray request with lfn, output_module_config fields, validFileOnly
type fileArrayLFNOutputModValidFileRequest struct {
	LogicalFileName   string `json:"logical_file_name"`
	ReleaseVersion    string `json:"release_version"`
	PsetHash          string `json:"pset_hash"`
	AppName           string `json:"app_name"`
	OutputModuleLabel string `json:"output_module_label"`
	ValidFileOnly     string `json:"validFileOnly"`
}

// test fileArray
func getFileArrayTestTable(t *testing.T) []EndpointTestCase {
	fileLumiList := []dbs.FileLumi{
		{LumiSectionNumber: 27414, RunNumber: 97},
		{LumiSectionNumber: 26422, RunNumber: 98},
		{LumiSectionNumber: 29838, RunNumber: 99},
	}

	var lfns []Response
	var lfnsRun97 []Response
	var detailResp []Response
	var detailRunResp []Response
	var detailRunSumLumiResp []Response
	for i := 1; i <= 10; i++ {
		lfn := fmt.Sprintf("/store/mc/Fall08/BBJets250to500-madgraph/GEN-SIM-RAW/IDEAL_/%v/%v.root", TestData.UID, i)
		lfns = append(lfns, fileResponse{LOGICAL_FILE_NAME: lfn})
		fileParentLFN := fmt.Sprintf("/store/mc/Fall08/BBJets250to500-madgraph/GEN-SIM-RAW/IDEAL_/p%v/%v.root", TestData.UID, i)
		fileParentList := []dbs.FileParentLFNRecord{
			{
				FILE_PARENT_LFN: fileParentLFN,
			},
		}
		fileRecord := createFileRecord(i, TestData.Dataset, TestData.Block, fileLumiList, lfn, fileParentList)
		r := createFileDetailedResponse(i+10, 1, 1, fileRecord)
		fileRunResp := fileRunResponse{
			LOGICAL_FILE_NAME: lfn,
			RUN_NUM:           97,
		}
		lfnsRun97 = append(lfnsRun97, fileRunResp)

		if i == 1 {
			r.LAST_MODIFIED_BY = "DBS-workflow"
			r.IS_FILE_VALID = 0
		}
		detailResp = append(detailResp, r)
		var detailRun fileDetailRunEventResponse
		d, err := json.Marshal(r)
		if err != nil {
			t.Fatal(err.Error())
		}
		err = json.Unmarshal(d, &detailRun)
		if err != nil {
			t.Fatal(err.Error())
		}
		detailRun.RUN_NUM = 97
		detailRun.EventCount = 1619
		detailRunResp = append(detailRunResp, detailRun)

		var detailRunResp fileDetailRunResponse
		d, err = json.Marshal(detailRun)
		if err != nil {
			t.Fatal(err.Error())
		}
		err = json.Unmarshal(d, &detailRunResp)
		if err != nil {
			t.Fatal(err.Error())
		}
		if i == 1 {
			for j := 97; j < 100; j++ {
				for i := 0; i < 3; i++ {
					detailRunResp.RUN_NUM = int64(j)
					detailRunSumLumiResp = append(detailRunSumLumiResp, detailRunResp)
				}
			}
		}
	}

	dbsError := dbs.DBSError{
		Reason:   dbs.InvalidParamErr.Error(),
		Code:     dbs.ParametersErrorCode,
		Message:  "cannot supply more than one list (lfn, run_num or lumi) at one query",
		Function: "dbs.files.Files",
	}
	hrec := createHTTPError("POST", "/dbs/fileArray")
	errorResp := createServerErrorResponse(hrec, &dbsError)

	dbsError2 := dbs.DBSError{
		Reason:   dbs.InvalidParamErr.Error(),
		Code:     dbs.ParametersErrorCode,
		Message:  "When sumOverLumi=1, no run_num list is allowed",
		Function: "dbs.files.Files",
	}
	errorResp2 := createServerErrorResponse(hrec, &dbsError2)

	var largeFileResp []Response
	err := readJsonFile(t, "./data/integration/files_response_data.json", &largeFileResp)
	if err != nil {
		t.Fatal(err.Error())
	}

	return []EndpointTestCase{
		{
			description:     "Test fileArray API with dataset parameter",
			defaultHandler:  web.FileArrayHandler,
			defaultEndpoint: "/dbs/fileArray",
			testCases: []testCase{
				{
					description: "Test POST with datasets", // DBSClientReader.test03200
					method:      "POST",
					serverType:  "DBSReader",
					input: fileArrayDatasetRequest{
						Dataset: TestData.Dataset,
					},
					output:   lfns,
					respCode: http.StatusOK,
				},
				{
					description: "Test POST with datasets, validFileOnly true", // DBSClientReader.test03200a
					method:      "POST",
					serverType:  "DBSReader",
					input: fileArrayDatasetValidFileRequest{
						Dataset:       TestData.Dataset,
						ValidFileOnly: "1",
					},
					output:   lfns[1:],
					respCode: http.StatusOK,
				},
				{
					description: "Test POST with datasets, validFileOnly false", // DBSClientReader.test03200b
					method:      "POST",
					serverType:  "DBSReader",
					input: fileArrayDatasetValidFileRequest{
						Dataset:       TestData.Dataset,
						ValidFileOnly: "0",
					},
					output:   lfns,
					respCode: http.StatusOK,
				},
				{
					description: "Test POST with datasets, validFileOnly true, detail, sumOverLumi", // DBSClientReader.test03200c
					method:      "POST",
					serverType:  "DBSReader",
					input: fileArrayDatasetValidFileDetailSumOverLumiRequest{
						Dataset:       TestData.Dataset,
						ValidFileOnly: "1",
						Detail:        "1",
						SumOverLumi:   "1",
					},
					output:   detailResp[1:],
					respCode: http.StatusOK,
				},
				{
					description: "Test POST with datasets, validFileOnly false, detail, sumOverLumi", // DBSClientReader.test03200d
					method:      "POST",
					serverType:  "DBSReader",
					input: fileArrayDatasetValidFileDetailSumOverLumiRequest{
						Dataset:       TestData.Dataset,
						ValidFileOnly: "0",
						Detail:        "1",
						SumOverLumi:   "1",
					},
					output:   detailResp,
					respCode: http.StatusOK,
				},
			},
		},
		{
			description:     "Test fileArray API with block_name parameter",
			defaultHandler:  web.FileArrayHandler,
			defaultEndpoint: "/dbs/fileArray",
			testCases: []testCase{
				{
					description: "Test POST with block_name", // DBSClientReader.test03300a
					method:      "POST",
					serverType:  "DBSReader",
					input: fileArrayBlockNameRequest{
						BlockName: TestData.Block,
					},
					output:   lfns,
					respCode: http.StatusOK,
				},
				{
					description: "Test POST with block_name and detail", // DBSClientReader.test03300b
					method:      "POST",
					serverType:  "DBSReader",
					input: fileArrayBlockNameDetailRequest{
						BlockName: TestData.Block,
						Detail:    "1",
					},
					output:   detailResp,
					respCode: http.StatusOK,
				},
				{
					description: "Test POST with block_name and validFileOnly true", // DBSClientReader.test03300c
					method:      "POST",
					serverType:  "DBSReader",
					input: fileArrayBlockNameValidFileRequest{
						BlockName:     TestData.Block,
						ValidFileOnly: "1",
					},
					output:   lfns[1:],
					respCode: http.StatusOK,
				},
				{
					description: "Test POST with block_name, detail, validFileOnly true", // DBSClientReader.test03300d
					method:      "POST",
					serverType:  "DBSReader",
					input: fileArrayBlockNameDetailValidFileRequest{
						BlockName:     TestData.Block,
						Detail:        "1",
						ValidFileOnly: "1",
					},
					output:   detailResp[1:],
					respCode: http.StatusOK,
				},
				{
					description: "Test POST with block_name and validFileOnly 0", // DBSClientReader.test03300e
					method:      "POST",
					serverType:  "DBSReader",
					input: fileArrayBlockNameValidFileRequest{
						BlockName:     TestData.Block,
						ValidFileOnly: "0",
					},
					output:   lfns,
					respCode: http.StatusOK,
				},
				{
					description: "Test POST with block_name, detail, validFileOnly 0", // DBSClientReader.test03300f
					method:      "POST",
					serverType:  "DBSReader",
					input: fileArrayBlockNameDetailValidFileRequest{
						BlockName:     TestData.Block,
						Detail:        "1",
						ValidFileOnly: "0",
					},
					output:   detailResp,
					respCode: http.StatusOK,
				},
				{
					description: "Test POST with block_name, run_num, lumi_list", // DBSClientReader.test03300g
					method:      "POST",
					serverType:  "DBSReader",
					input: fileArrayBlockNameRunNumLumiListRequest{
						BlockName: TestData.Block,
						RunNum:    "97",
						LumiList:  "[27414,26422,29838]",
					},
					output:   lfnsRun97,
					respCode: http.StatusOK,
				},
				{
					description: "Test POST with block_name, run_num, nested lumi_list", // DBSClientReader.test03300h
					method:      "POST",
					serverType:  "DBSReader",
					input: fileArrayBlockNameRunNumLumiListRequest{
						BlockName: TestData.Block,
						RunNum:    "97",
						LumiList:  "[[27414 27418] [26422 26426] [29838 29842]]",
					},
					output:   lfnsRun97,
					respCode: http.StatusOK,
				},
				{
					description: "Test POST with block_name, run_num, lumi_list, detail", // DBSClientReader.test03300i
					method:      "POST",
					serverType:  "DBSReader",
					input: fileArrayBlockNameRunNumLumiListDetailRequest{
						BlockName: TestData.Block,
						RunNum:    "97",
						LumiList:  "[27414,26422,29838]",
						Detail:    "1",
					},
					output:   detailRunResp,
					respCode: http.StatusOK,
				},
				{
					description: "Test POST with block_name, run_num, nested lumi_list, detail", // DBSClientReader.test03300j
					method:      "POST",
					serverType:  "DBSReader",
					input: fileArrayBlockNameRunNumLumiListDetailRequest{
						BlockName: TestData.Block,
						RunNum:    "97",
						LumiList:  "[[27414 27418] [26422 26426] [29838 29842]]",
						Detail:    "1",
					},
					output:   detailRunResp,
					respCode: http.StatusOK,
				},
				{
					description: "Test POST with block_name, run_num, lumi_list, validFileOnly 1", // DBSClientReader.test03300k
					method:      "POST",
					serverType:  "DBSReader",
					input: fileArrayBlockNameRunNumLumiListValidFileRequest{
						BlockName:     TestData.Block,
						RunNum:        "97",
						LumiList:      "[27414,26422,29838]",
						ValidFileOnly: "1",
					},
					output:   lfnsRun97[1:],
					respCode: http.StatusOK,
				},
				{
					description: "Test POST with block_name, run_num, nested lumi_list, validFileOnly 1", // DBSClientReader.test03300l
					method:      "POST",
					serverType:  "DBSReader",
					input: fileArrayBlockNameRunNumLumiListValidFileRequest{
						BlockName:     TestData.Block,
						RunNum:        "97",
						LumiList:      "[[27414 27418] [26422 26426] [29838 29842]]",
						ValidFileOnly: "1",
					},
					output:   lfnsRun97[1:],
					respCode: http.StatusOK,
				},
				{
					description: "Test POST with block_name, run_num, lumi_list, detail, validFileOnly 1", // DBSClientReader.test03300m
					method:      "POST",
					serverType:  "DBSReader",
					input: fileArrayBlockNameRunNumLumiListValidFileDetailRequest{
						BlockName:     TestData.Block,
						RunNum:        "97",
						LumiList:      "[27414,26422,29838]",
						ValidFileOnly: "1",
						Detail:        "1",
					},
					output:   detailRunResp[1:],
					respCode: http.StatusOK,
				},
				{
					description: "Test POST with block_name, run_num, nested lumi_list, detail, validFileOnly 1", // DBSClientReader.test03300n
					method:      "POST",
					serverType:  "DBSReader",
					input: fileArrayBlockNameRunNumLumiListValidFileDetailRequest{
						BlockName:     TestData.Block,
						RunNum:        "97",
						LumiList:      "[[27414 27418] [26422 26426] [29838 29842]]",
						ValidFileOnly: "1",
						Detail:        "1",
					},
					output:   detailRunResp[1:],
					respCode: http.StatusOK,
				},
				{
					description: "Test POST with block_name, list run_num, lumi_list", // DBSClientReader.test03300o
					method:      "POST",
					serverType:  "DBSReader",
					input: fileArrayBlockNameRunNumLumiListRequest{
						BlockName: TestData.Block,
						RunNum:    "[97]",
						LumiList:  "[27414,26422,29838]",
					},
					output:   lfnsRun97,
					respCode: http.StatusOK,
				},
				{
					description: "Test POST with block_name, list run_num, lumi_list, sumOverLumi", // DBSClientReader.test03300p
					method:      "POST",
					serverType:  "DBSReader",
					input: fileArrayBlockNameRunNumLumiListSumOverLumiRequest{
						BlockName:   TestData.Block,
						RunNum:      "[97]",
						LumiList:    "[27414,26422,29838]",
						SumOverLumi: "1",
					},
					output:   []Response{errorResp},
					respCode: http.StatusBadRequest,
				},
				{
					description: "Test POST with block_name, list run_num, lumi_list, sumOverLumi, detail", // DBSClientReader.test03300q
					method:      "POST",
					serverType:  "DBSReader",
					input: fileArrayBlockNameRunNumLumiListSumOverLumiDetailRequest{
						BlockName:   TestData.Block,
						RunNum:      "[97]",
						LumiList:    "[27414,26422,29838]",
						SumOverLumi: "1",
						Detail:      "1",
					},
					output:   []Response{errorResp},
					respCode: http.StatusBadRequest,
				},
				{
					description: "Test POST with block_name, list run_num, sumOverLumi, detail", // DBSClientReader.test03300r
					method:      "POST",
					serverType:  "DBSReader",
					input: fileArrayBlockNameRunNumSumOverLumiDetailRequest{
						BlockName:   TestData.Block,
						RunNum:      "[97]",
						SumOverLumi: "1",
						Detail:      "1",
					},
					output:   []Response{errorResp2},
					respCode: http.StatusBadRequest,
				},
				{
					description: "Test POST with block_name, list run_num, sumOverLumi", // DBSClientReader.test03300s
					method:      "POST",
					serverType:  "DBSReader",
					input: fileArrayBlockNameRunNumSumOverLumiRequest{
						BlockName:   TestData.Block,
						RunNum:      "[97]",
						SumOverLumi: "1",
					},
					output:   []Response{errorResp2},
					respCode: http.StatusBadRequest,
				},
				{
					description: "Test POST with block_name, range run_num, sumOverLumi, detail", // DBSClientReader.test03300t
					method:      "POST",
					serverType:  "DBSReader",
					input: fileArrayBlockNameRunNumSumOverLumiDetailRequest{
						BlockName:   TestData.Block,
						RunNum:      "97-99",
						SumOverLumi: "1",
						Detail:      "1",
					},
					output:   largeFileResp,
					respCode: http.StatusOK,
				},
			},
		},
		{
			description:     "Test fileArray API with logical_file_name parameter",
			defaultHandler:  web.FileArrayHandler,
			defaultEndpoint: "/dbs/fileArray",
			testCases: []testCase{
				{
					description: "Test POST with logical_file_name", // DBSClientReader_t.test03400a
					method:      "POST",
					serverType:  "DBSReader",
					input: fileArrayLFNRequest{
						LogicalFileName: TestData.Files[0],
					},
					output:   lfns[:1],
					respCode: http.StatusOK,
				},
				{
					description: "Test POST with logical_file_name, validFileOnly 1", // DBSClientReader_t.test03400b
					method:      "POST",
					serverType:  "DBSReader",
					input: fileArrayLFNValidFileRequest{
						LogicalFileName: TestData.Files[0],
						ValidFileOnly:   "1",
					},
					output:   []Response{},
					respCode: http.StatusOK,
				},
				{
					description: "Test POST with logical_file_name, validFileOnly 1", // DBSClientReader_t.test03400b2
					method:      "POST",
					serverType:  "DBSReader",
					input: fileArrayLFNValidFileRequest{
						LogicalFileName: TestData.Files[1],
						ValidFileOnly:   "1",
					},
					output:   lfns[1:2],
					respCode: http.StatusOK,
				},
				{
					description: "Test POST with logical_file_name, validFileOnly 0", // DBSClientReader_t.test03400c
					method:      "POST",
					serverType:  "DBSReader",
					input: fileArrayLFNValidFileRequest{
						LogicalFileName: TestData.Files[0],
						ValidFileOnly:   "0",
					},
					output:   lfns[:1],
					respCode: http.StatusOK,
				},
				{
					description: "Test POST with logical_file_name, run_num, lumi_list", // DBSClientReader_t.test03400d
					method:      "POST",
					serverType:  "DBSReader",
					input: fileArrayLFNRunNumLumiListRequest{
						LogicalFileName: TestData.Files[0],
						RunNum:          "97",
						LumiList:        "[27414,26422,29838]",
					},
					output:   lfnsRun97[:1],
					respCode: http.StatusOK,
				},
				{
					description: "Test POST with logical_file_name, run_num, nested lumi_list", // DBSClientReader_t.test03400e
					method:      "POST",
					serverType:  "DBSReader",
					input: fileArrayLFNRunNumLumiListRequest{
						LogicalFileName: TestData.Files[0],
						RunNum:          "97",
						LumiList:        "[[27414 27418] [26422 26426] [29838 29842]]",
					},
					output:   lfnsRun97[:1],
					respCode: http.StatusOK,
				},
				{
					description: "Test POST with logical_file_name, run_num, lumi_list, detail", // DBSClientReader_t.test03400f
					method:      "POST",
					serverType:  "DBSReader",
					input: fileArrayLFNRunNumLumiListDetailRequest{
						LogicalFileName: TestData.Files[0],
						RunNum:          "97",
						LumiList:        "[27414,26422,29838]",
						Detail:          "1",
					},
					output:   detailRunResp[:1],
					respCode: http.StatusOK,
				},
				{
					description: "Test POST with logical_file_name, run_num, nested lumi_list, detail", // DBSClientReader_t.test03400g
					method:      "POST",
					serverType:  "DBSReader",
					input: fileArrayLFNRunNumLumiListDetailRequest{
						LogicalFileName: TestData.Files[0],
						RunNum:          "97",
						LumiList:        "[[27414 27418] [26422 26426] [29838 29842]]",
						Detail:          "1",
					},
					output:   detailRunResp[:1],
					respCode: http.StatusOK,
				},
				{
					description: "Test POST with logical_file_name, run_num, lumi_list, validFileOnly", // DBSClientReader_t.test03400h
					method:      "POST",
					serverType:  "DBSReader",
					input: fileArrayLFNRunNumLumiListValidFileRequest{
						LogicalFileName: TestData.Files[1],
						RunNum:          "97",
						LumiList:        "[27414,26422,29838]",
						ValidFileOnly:   "1",
					},
					output:   lfnsRun97[1:2],
					respCode: http.StatusOK,
				},
				{
					description: "Test POST with logical_file_name, run_num, nested lumi_list, validFileOnly", // DBSClientReader_t.test03400i
					method:      "POST",
					serverType:  "DBSReader",
					input: fileArrayLFNRunNumLumiListValidFileRequest{
						LogicalFileName: TestData.Files[1],
						RunNum:          "97",
						LumiList:        "[[27414 27418] [26422 26426] [29838 29842]]",
						ValidFileOnly:   "1",
					},
					output:   lfnsRun97[1:2],
					respCode: http.StatusOK,
				},
				{
					description: "Test POST with logical_file_name, run_num, lumi_list, validFileOnly, detail", // DBSClientReader_t.test03400j
					method:      "POST",
					serverType:  "DBSReader",
					input: fileArrayLFNRunNumLumiListValidFileDetailRequest{
						LogicalFileName: TestData.Files[1],
						RunNum:          "97",
						LumiList:        "[27414,26422,29838]",
						ValidFileOnly:   "1",
						Detail:          "1",
					},
					output:   detailRunResp[1:2],
					respCode: http.StatusOK,
				},
				{
					description: "Test POST with logical_file_name, run_num, nested lumi_list, validFileOnly, detail", // DBSClientReader_t.test03400k
					method:      "POST",
					serverType:  "DBSReader",
					input: fileArrayLFNRunNumLumiListValidFileDetailRequest{
						LogicalFileName: TestData.Files[1],
						RunNum:          "97",
						LumiList:        "[[27414 27418] [26422 26426] [29838 29842]]",
						ValidFileOnly:   "1",
						Detail:          "1",
					},
					output:   detailRunResp[1:2],
					respCode: http.StatusOK,
				},
				{
					description: "Test POST with logical_file_name, list run_num, lumi_list", // DBSClientReader_t.test03400l
					method:      "POST",
					serverType:  "DBSReader",
					input: fileArrayLFNRunNumLumiListRequest{
						LogicalFileName: TestData.Files[0],
						RunNum:          "[97]",
						LumiList:        "[27414,26422,29838]",
					},
					output:   lfnsRun97[:1],
					respCode: http.StatusOK,
				},
				{
					description: "Test POST with logical_file_name, run_num, lumi_list", // DBSClientReader_t.test03400m
					method:      "POST",
					serverType:  "DBSReader",
					input: fileArrayLFNRunNumLumiListRequest{
						LogicalFileName: TestData.Files[0],
						RunNum:          "97",
						LumiList:        "[27414,26422,29838]",
					},
					output:   lfnsRun97[:1],
					respCode: http.StatusOK,
				},
				{
					description: "Test POST with logical_file_name, list run_num, sumOverLumi, detail", // DBSClientReader_t.test03400n
					method:      "POST",
					serverType:  "DBSReader",
					input: fileArrayLFNRunNumSumOverLumiDetailRequest{
						LogicalFileName: TestData.Files[0],
						RunNum:          "[97]",
						SumOverLumi:     "1",
						Detail:          "1",
					},
					output:   []Response{errorResp2},
					respCode: http.StatusBadRequest,
				},
				{
					description: "Test POST with logical_file_name, ranged run_num, sumOverLumi, detail", // DBSClientReader_t.test03400o
					method:      "POST",
					serverType:  "DBSReader",
					input: fileArrayLFNRunNumSumOverLumiDetailRequest{
						LogicalFileName: TestData.Files[0],
						RunNum:          "97-99",
						SumOverLumi:     "1",
						Detail:          "1",
					},
					output:   detailRunSumLumiResp,
					respCode: http.StatusOK,
				},
			},
		},
		{
			description:     "Test fileArray with dataset and output_module_config parameters",
			defaultHandler:  web.FileArrayHandler,
			defaultEndpoint: "/dbs/fileArray",
			testCases: []testCase{
				{
					description: "Test POST with dataset and release_version", // DBSClientReader_t.test03500a
					method:      "POST",
					serverType:  "DBSReader",
					input: fileArrayDatasetReleaseRequest{
						Dataset:        TestData.Dataset,
						ReleaseVersion: TestData.ReleaseVersion,
					},
					output:   lfns,
					respCode: http.StatusOK,
				},
				{
					description: "Test POST with dataset, release_version, validFileOnly 1", // DBSClientReader_t.test03500b
					method:      "POST",
					serverType:  "DBSReader",
					input: fileArrayDatasetReleaseValidFileRequest{
						Dataset:        TestData.Dataset,
						ReleaseVersion: TestData.ReleaseVersion,
						ValidFileOnly:  "1",
					},
					output:   lfns[1:],
					respCode: http.StatusOK,
				},
				{
					description: "Test POST with dataset, release_version, pset_hash, app_name, output_module_label", // DBSClientReader_t.test03600
					method:      "POST",
					serverType:  "DBSReader",
					input: fileArrayDatasetOutputModRequest{
						Dataset:           TestData.Dataset,
						ReleaseVersion:    TestData.ReleaseVersion,
						PsetHash:          TestData.PsetHash,
						AppName:           TestData.AppName,
						OutputModuleLabel: TestData.OutputModuleLabel,
					},
					output:   lfns,
					respCode: http.StatusOK,
				},
			},
		},
		{
			description:     "Test fileArray with logical_file_name and output_module_config parameters",
			defaultHandler:  web.FileArrayHandler,
			defaultEndpoint: "/dbs/fileArray",
			testCases: []testCase{
				{
					description: "Test POST with dataset and release_version", // DBSClientReader_t.test03700a
					method:      "POST",
					serverType:  "DBSReader",
					input: fileArrayLFNOutputModRequest{
						LogicalFileName:   TestData.Files[0],
						ReleaseVersion:    TestData.ReleaseVersion,
						PsetHash:          TestData.PsetHash,
						AppName:           TestData.AppName,
						OutputModuleLabel: TestData.OutputModuleLabel,
					},
					output:   lfns[:1],
					respCode: http.StatusOK,
				},
				{
					description: "Test POST with dataset, release_version, pset_hash, app_name, output_module_label", // DBSClientReader_t.test03700b
					method:      "POST",
					serverType:  "DBSReader",
					input: fileArrayLFNOutputModValidFileRequest{
						LogicalFileName:   TestData.Files[1],
						ReleaseVersion:    TestData.ReleaseVersion,
						PsetHash:          TestData.PsetHash,
						AppName:           TestData.AppName,
						OutputModuleLabel: TestData.OutputModuleLabel,
						ValidFileOnly:     "1",
					},
					output:   lfns[1:2],
					respCode: http.StatusOK,
				},
			},
		},
		{
			description:     "Test fileArray with non-existing fields",
			defaultHandler:  web.FileArrayHandler,
			defaultEndpoint: "/dbs/fileArray",
			testCases: []testCase{
				{
					description: "Test POST with non-existing dataset", // DBSClientReader_t.test03800
					method:      "POST",
					serverType:  "DBSReader",
					input: fileArrayDatasetRequest{
						Dataset: "/does/not/EXIST",
					},
					output:   []Response{},
					respCode: http.StatusOK,
				},
				{
					description: "Test POST with non-existing block_name", // DBSClientReader_t.test03900
					method:      "POST",
					serverType:  "DBSReader",
					input: fileArrayBlockNameRequest{
						BlockName: "/does/not/EXIST#123",
					},
					output:   []Response{},
					respCode: http.StatusOK,
				},
				{
					description: "Test POST with non-existing logical_file_name", // DBSClientReader_t.test0400
					method:      "POST",
					serverType:  "DBSReader",
					input: fileArrayLFNRequest{
						LogicalFileName: "/store/mc/does/not/EXIST/NotReally/0815/doesnotexist.root",
					},
					output:   []Response{},
					respCode: http.StatusOK,
				},
			},
		},
	}
}
