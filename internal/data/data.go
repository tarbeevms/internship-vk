package data

import "github.com/tarantool/go-tarantool"

type DataRepository struct {
	tConn *tarantool.Connection
}

func NewDataRepo(conn *tarantool.Connection) *DataRepository {
	return &DataRepository{
		tConn: conn,
	}
}

type WriteDataReq struct {
	Data map[string]interface{} `json:"data"`
}

type WriteDataResp struct {
	Status string `json:"status"`
}

type ReadDataReq struct {
	Keys []string `json:"keys"`
}

type ReadDataResp struct {
	Data map[string]interface{} `json:"data"`
}
