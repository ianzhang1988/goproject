package lambda_service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	inter "goproject/prototype/lambda/server/internal"
)

func ProbeJob(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("read err: %s", err.Error()), 400)
		return
	}

	job := inter.Job{}
	err = json.Unmarshal(data, &job)
	if err != nil {
		http.Error(w, fmt.Sprintf("parse body err: %s", err.Error()), 400)
		return
	}

	// send j to process and get return status

	w.Write([]byte("OK"))
}
