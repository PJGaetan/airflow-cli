package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"

	"github.com/pjgaetan/airflow-cli/commands"
)

func TestDag(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Log(r.URL.Path)
		// if r.URL.Path != "/dags" {
		// 	t.Errorf("Expected to request '/dags', got: %s", r.URL.Path)
		// }
		// if r.Header.Get("Accept") != "application/json" {
		// 	t.Errorf("Expected Accept: application/json header, got: %s", r.Header.Get("Accept"))
		// }
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(loadFixture(t, "input/dags.json")))
		if err != nil {
			t.Error("Unabale to write fixture as string input")
		}
	}))
	defer server.Close()

	file, e := os.CreateTemp("~/.config/.airflow/", ".config")
	defer file.Close()
	if e != nil {
		t.Error(e)
	}

	tests := []struct {
		name    string
		args    []string
		fixture string
	}{
		// {"dag", []string{"dag"}, "dag.golden"},
		{"dag list", []string{"dag", "list"}, "dag-list.golden"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bufOut := new(bytes.Buffer)
			bufErr := new(bytes.Buffer)

			cmd := commands.NewRootCmd()
			cmd.SetArgs(tt.args)
			cmd.SetOut(bufOut)
			cmd.SetErr(bufErr)
			cmd.Execute()

			expected := loadFixture(t, tt.fixture)

			// e is nil if nothing in it
			outErr := bufErr.String()
			out := bufOut.String()
			// fmt.Println(outErr)
			if !reflect.DeepEqual(outErr, expected) {
				// t.Fatalf("actual = %s, expected = %s", out, expected)
			}
			if !reflect.DeepEqual(out, expected) {
				// t.Fatalf("actual = %s, expected = %s", out, expected)
			}

			// output, err := runBinary(tt.args)
			// // list.NewList()
			// if err != nil {
			// 	t.Log(string(debug.Stack()))
			// 	t.Fatal(err)
			// }
			//
			// if *update {
			// 	writeFixture(t, tt.fixture, output)
			// }
			//
			// actual := string(output)
			//
			// expected := loadFixture(t, tt.fixture)
			//
			// if !reflect.DeepEqual(actual, expected) {
			// 	t.Fatalf("actual = %s, expected = %s", actual, expected)
			// }
		})
	}
}
