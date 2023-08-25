package lambda_service

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestJobHandler(t *testing.T) {
	reader := strings.NewReader(`
	{
		"task":[
			{
				"id": "t1",
				"affinity": ["a"]
				"x": "y"
			}
		],
		"affinity":{
			"a": {
				"isp":"ct"
			}
		},
		"upload": {
			"type": "kafka",
			"topic": "test"
		},
		"lambda_behaviour":{
			"replic": 3
		}
	}
	`)

	req, err := http.NewRequest("POST", "/", reader)
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(ProbeJob)

	handler.ServeHTTP(recorder, req)

	expected := "OK"
	actual := recorder.Body.String()

	if actual != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", actual, expected)
	}
}
