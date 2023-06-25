package request

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"

	"github.com/pjgaetan/airflow-cli/internal/config"
	"github.com/pjgaetan/airflow-cli/pkg/utils"
)

type errorResponse struct {
	Detail       string `json:"detail"`
	Status       int    `json:"status"`
	Title        string `json:"title"`
	ResponseType string `json:"type"`
}

func MakeRequest(payload, url, method string, header []string) (string, error) {
	var reader io.Reader
	if payload != "" {
		reader = bytes.NewReader([]byte(payload))
	}

	request, err := http.NewRequest(method, url, reader)
	utils.ExitIfError(err)
	client := &http.Client{}

	for _, s := range header {
		keys := strings.Split(s, ":")
		if len(keys) < 2 {
			utils.Failed("Error: ", "not enought param in header "+s)
		}

		request.Header.Add(strings.TrimSpace(keys[0]), strings.TrimSpace(keys[1]))
	}

	resp, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer func() {
		err := resp.Body.Close()
		utils.ExitIfError(err)
	}()

	bodyBytes, _ := io.ReadAll(resp.Body)
	bodyString := string(bodyBytes)

	const SUCCESS_CODE_LOWER_BOUND, SUCCESS_CODE_UPPER_BOUND = 200, 300
	statusOK := resp.StatusCode >= SUCCESS_CODE_LOWER_BOUND && resp.StatusCode < SUCCESS_CODE_UPPER_BOUND
	if !statusOK {
		var error errorResponse
		if err := json.Unmarshal([]byte(bodyString), &error); err != nil {
			panic(err)
		}
		utils.Fail("Error status code : " + strconv.Itoa(resp.StatusCode))

		// for airflow behind jwt, the response can have any schema
		if reflect.DeepEqual(error, errorResponse{}) {
			utils.Failed(bodyString)
		}
		utils.Failed(error.ResponseType + " " + error.Detail)
	}

	return bodyString, nil
}

func AirflowGetRequest(endpoint string, params ...[2]string) json.RawMessage {
	header := [1]string{config.AuthorizationHeader}

	// emptiness has been checked in GetActiveProfile
	baseUrl := config.Url
	if !strings.HasSuffix(baseUrl, "/") {
		baseUrl += "/"
	}
	baseUrl += endpoint

	// construct url
	queryParams := url.Values{}

	u, _ := url.ParseRequestURI(baseUrl)
	for _, param := range params {
		queryParams.Add(param[0], param[1])
	}
	u.RawQuery = queryParams.Encode()

	response, err := MakeRequest(
		"",
		u.String(),
		"GET",
		header[:],
	)
	if err != nil {
		log.Fatal("Error ", err)
	}
	var encapsulation map[string]json.RawMessage
	if err := json.Unmarshal([]byte(response), &encapsulation); err != nil {
		panic(err)
	}
	keys := maps.Keys(encapsulation)
	if slices.Contains(keys, "response") {
		return encapsulation["response"]
	}
	var r json.RawMessage
	if err := json.Unmarshal([]byte(response), &r); err != nil {
		panic(err)
	}
	return r
}

func AirflowPostRequest(endpoint string, payload string) json.RawMessage {
	var header [2]string
	header[0] = config.AuthorizationHeader
	header[1] = "Content-Type: application/json"

	// emptiness has been checked in GetActiveProfile
	baseUrl := config.Url
	if !strings.HasSuffix(baseUrl, "/") {
		baseUrl += "/"
	}
	baseUrl += endpoint

	response, err := MakeRequest(
		payload,
		baseUrl,
		"POST",
		header[:],
	)
	if err != nil {
		log.Fatal("Error ", err)
	}
	var encapsulation map[string]json.RawMessage
	if err := json.Unmarshal([]byte(response), &encapsulation); err != nil {
		fmt.Println(response)
		panic(err)
	}
	keys := maps.Keys(encapsulation)
	if slices.Contains(keys, "response") {
		return encapsulation["response"]
	}
	var r json.RawMessage
	if err := json.Unmarshal([]byte(response), &r); err != nil {
		panic(err)
	}
	return r
}

func AirflowPatchRequest(endpoint string, payload string, params ...[2]string) json.RawMessage {
	var header [2]string
	header[0] = config.AuthorizationHeader
	header[1] = "Content-Type: application/json"

	// emptiness has been checked in GetActiveProfile
	baseUrl := config.Url
	if !strings.HasSuffix(baseUrl, "/") {
		baseUrl += "/"
	}
	baseUrl += endpoint

	// construct url
	queryParams := url.Values{}

	u, _ := url.ParseRequestURI(baseUrl)
	for _, param := range params {
		queryParams.Add(param[0], param[1])
	}
	u.RawQuery = queryParams.Encode()

	response, err := MakeRequest(
		payload,
		u.String(),
		"PATCH",
		header[:],
	)
	if err != nil {
		log.Fatal("Error ", err)
	}
	var encapsulation map[string]json.RawMessage
	if err := json.Unmarshal([]byte(response), &encapsulation); err != nil {
		fmt.Println(response)
		panic(err)
	}
	keys := maps.Keys(encapsulation)
	if slices.Contains(keys, "response") {
		return encapsulation["response"]
	}
	var r json.RawMessage
	if err := json.Unmarshal([]byte(response), &r); err != nil {
		panic(err)
	}
	return r
}
