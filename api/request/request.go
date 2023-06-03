package request

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/gookit/ini/v2"
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
		utils.Fail("Error status code : " + strconv.Itoa(error.Status))
		utils.Failed(error.ResponseType + " " + error.Detail)
	}

	return bodyString, nil
}

func AirflowGetRequest(endpoint string, params ...[2]string) json.RawMessage {
	profile_name, auth_method, err := config.GetActiveProfile()
	if err != nil {
		log.Fatal("Error ", err)
	}
	var header [1]string

	switch auth_method {
	case "user/password":
		profile := config.GetUserPasswordProfile(profile_name)
		header = [1]string{"Authorization: Basic " + config.BasicAuth(profile)}
	case "jwt":
		profile := config.GetJwtProfile(profile_name)
		token := config.GetToken(profile)
		header = [1]string{"Authorization: Bearer " + token}
	default:
		utils.Failed("no such possibility")
	}

	// emptiness has been checked in GetActiveProfile
	baseUrl := ini.String(profile_name + ".url")
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
	profile_name, auth_method, err := config.GetActiveProfile()
	if err != nil {
		log.Fatal("Error ", err)
	}
	var header [2]string
	switch auth_method {
	case "user/password":
		profile := config.GetUserPasswordProfile(profile_name)
		header[0] = "Authorization: Basic " + config.BasicAuth(profile)
	case "jwt":
		profile := config.GetJwtProfile(profile_name)
		token := config.GetToken(profile)
		header[0] = "Authorization: Bearer " + token
	default:
		utils.Failed("no such possibility")
	}
	header[1] = "Content-Type: application/json"

	// emptiness has been checked in GetActiveProfile
	baseUrl := ini.String(profile_name + ".url")
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
	profile_name, auth_method, err := config.GetActiveProfile()
	if err != nil {
		log.Fatal("Error ", err)
	}
	var header [2]string
	switch auth_method {
	case "user/password":
		profile := config.GetUserPasswordProfile(profile_name)
		header[0] = "Authorization: Basic " + config.BasicAuth(profile)
	case "jwt":
		profile := config.GetJwtProfile(profile_name)
		token := config.GetToken(profile)
		header[0] = "Authorization: Bearer " + token
	default:
		utils.Failed("no such possibility")
	}
	header[1] = "Content-Type: application/json"

	// emptiness has been checked in GetActiveProfile
	baseUrl := ini.String(profile_name + ".url")
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
