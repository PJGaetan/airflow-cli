package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

var update = flag.Bool("update", false, "update golden files")

var binaryName = "airflow-cli-coverage"

var binaryPath = ""

func fixturePath(t *testing.T, fixture string) string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatalf("problems recovering caller information")
	}

	return filepath.Join(filepath.Dir(filename), fixture)
}

func writeFixture(t *testing.T, fixture string, content []byte) {
	err := ioutil.WriteFile(fixturePath(t, fixture), content, 0644)
	if err != nil {
		t.Fatal(err)
	}
}

func loadFixture(t *testing.T, fixture string) string {
	content, err := ioutil.ReadFile(fixturePath(t, fixture))
	if err != nil {
		t.Fatal(err)
	}

	return string(content)
}

func TestMain(m *testing.M) {
	err := os.Chdir("..")
	if err != nil {
		fmt.Printf("could not change dir: %v", err)
		os.Exit(1)
	}

	dir, err := os.Getwd()
	if err != nil {
		fmt.Printf("could not get current dir: %v", err)
	}

	binaryPath = filepath.Join(dir, binaryName)
	// composeFilePaths := []string{"testdata/docker-compose.yml"}
	// identifier := strings.ToLower(uuid.New().String())
	//
	// compose := tc.NewLocalDockerCompose(composeFilePaths, identifier)
	// execError := compose.
	// 	WithCommand([]string{"up", "-d"}).
	// 	WithEnv(map[string]string{
	// 		"key1": "value1",
	// 		"key2": "value2",
	// 	}).
	// 	Invoke()
	//
	// errExec := execError.Error
	// if errExec != nil {
	// 	fmt.Printf("Could not run compose file: %v - %v", composeFilePaths, err)
	// 	os.Exit(1)
	// }

	// compose, err := tc.NewDockerCompose("testdata/docker-compose.yml")
	// if err != nil {
	// 	fmt.Printf("could not get docker-compose file: %v", err)
	// }
	//
	// ctx, cancel := context.WithCancel(context.Background())
	//
	// compose.Up(ctx, tc.Wait(true))
	// if err != nil {
	// 	fmt.Printf("could not get docker-compose up: %v", err)
	// }
	// defer cancel()

	code := m.Run()

	// execErrorDown := compose.Down()
	// errDown := execErrorDown.Error
	// if errDown != nil {
	// 	fmt.Printf("Could not run compose file: %v - %v", composeFilePaths, err)
	// 	os.Exit(1)
	// }

	os.Exit(code)
}

func runBinary(args []string) ([]byte, error) {
	cmd := exec.Command(binaryPath, args...)
	cmd.Env = append(os.Environ(), "GOCOVERDIR=.coverdata")
	return cmd.CombinedOutput()
}
