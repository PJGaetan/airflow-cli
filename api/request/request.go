package request

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gookit/ini/v2"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"

	"github.com/pjgaetan/airflow-cli/internal/config"
	"github.com/pjgaetan/airflow-cli/pkg/model"
)

func MakeRequest(payload, url, method string, header []string) (string, error) {
	var reader io.Reader
	if payload != "" && payload != "{}" {
		reader = bytes.NewReader([]byte(payload))
	}

	request, err := http.NewRequest(method, url, reader)
	if err != nil {
		log.Fatal("Error: ", err)
	}
	client := &http.Client{}

	for _, s := range header {
		keys := strings.Split(s, ":")
		if len(keys) < 2 {
			log.Fatal("Error: ", "not enought param in header "+s)
		}

		request.Header.Add(strings.TrimSpace(keys[0]), strings.TrimSpace(keys[1]))
	}

	resp, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	bodyString := string(bodyBytes)

	return bodyString, nil
}

func AirflowGetRequest(endpoint string) model.Dags {
	profile_name, auth_method, err := config.GetActiveProfile()
	if err != nil {
		log.Fatal("Error ", err)
	}
	var header [1]string
	if auth_method == "user/password" {
		profile := config.GetUserPasswordProfile(profile_name)
		header = [1]string{"Authorization: Basic " + config.BasicAuth(profile)}
	} else if auth_method == "jwt" {
		profile := config.GetJwtProfile(profile_name)
		token := config.GetToken(profile)
		header = [1]string{"Authorization: Bearer " + token}
	} else {
		fmt.Fprintf(os.Stderr, "%s\n", "no such possibility")
		os.Exit(1)
	}

	// emptiness has been checked in GetActiveProfile
	base_url := ini.String(profile_name + ".url")
	if !strings.HasSuffix(base_url, "/") {
		base_url = base_url + "/"
	}
	url := base_url + endpoint

	response, err := MakeRequest(
		"",
		url,
		"GET",
		header[:],
	)
	if err != nil {
		log.Fatal("Error ", err)
	}
	var encapsulation map[string]interface{}
	if err := json.Unmarshal([]byte(response), &encapsulation); err != nil {
		panic(err)
	}
	keys := maps.Keys(encapsulation)
	var dag model.Dags
	if slices.Contains(keys, "response") {
		var dat model.ResponseDag

		if err := json.Unmarshal([]byte(response), &dat); err != nil {
			panic(err)
		}
		dag = dat.Response

	} else {
		if err := json.Unmarshal([]byte(response), &dag); err != nil {
			panic(err)
		}
	}
	return dag
}
