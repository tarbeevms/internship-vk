package controllers

import (
	"myapp/internal/data"
	"myapp/pkg/io"
	"net/http"
)

func (h *HandlerLayer) WriteDataHandler(w http.ResponseWriter, r *http.Request) {
	var req data.WriteDataReq
	err := io.ReadJSON(r, &req)
	if err != nil {
		io.SendError(w, "Bad Request", http.StatusBadRequest)
		return
	}
	err = h.LogicLayer.WriteDataBatch(req.Data)
	if err != nil {
		io.WriteJSON(w, http.StatusBadRequest, data.WriteDataResp{Status: err.Error()})
		return
	}
	io.WriteJSON(w, http.StatusOK, data.WriteDataResp{Status: "Success"})
}

func (h *HandlerLayer) ReadDataHandler(w http.ResponseWriter, r *http.Request) {
	var req data.ReadDataReq
	err := io.ReadJSON(r, &req)
	if err != nil {
		io.SendError(w, "Bad Request", http.StatusBadRequest)
		return
	}
	result, err := h.LogicLayer.ReadDataBatch(req.Keys)
	if err != nil {
		io.WriteJSON(w, http.StatusInternalServerError, data.WriteDataResp{Status: err.Error()})
		return
	}
	io.WriteJSON(w, http.StatusOK, data.ReadDataResp{Data: result})
}
