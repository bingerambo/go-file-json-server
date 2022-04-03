package utils

import "testing"

func TestExists(t *testing.T) {
	dir_path := "D:/Go_projects/alert-engine/src/alert_engine/api"
	file_path := "D:/http_alert/http_client/client.go"


	if !IsDir(dir_path){
		t.Error("IsDir error")
	}

	if !IsFile(file_path){
		t.Error("IsFile error")
	}

	if !Exists(file_path){
		t.Error("Exists error")
	}

	if !Exists(dir_path){
		t.Error("Exists error")
	}
}
