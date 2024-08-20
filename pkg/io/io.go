package io

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, status int, data interface{}) {
	js, err := json.Marshal(data)
	if err != nil {
		SendError(w, "Error sending resposne", http.StatusInternalServerError)
		log.Println("WriteJSON error: ", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(js)
	if err != nil {
		SendError(w, "Internal Error", http.StatusInternalServerError)
		log.Println("WriteJSON error: ", err)
		return
	}
}

func ReadJSON(r *http.Request, dest interface{}) error {
	content, err := io.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		log.Println("ReadJSON error: ", err)
		return err
	}
	err = json.Unmarshal(content, dest)
	if err != nil {
		log.Println("ReadJSON error: ", err)
		return err
	}
	return nil
}

func SendError(w http.ResponseWriter, msg string, status int) {
	data := map[string]string{"status": msg}
	js, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "Internal Error", http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(js)
	if err != nil {
		http.Error(w, "Internal Error", http.StatusInternalServerError)
	}
}
