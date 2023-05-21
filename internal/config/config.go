package config

import (
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	"github.com/gookit/ini/v2"
	"github.com/mitchellh/go-homedir"
	"github.com/pjgaetan/airflow-cli/internal/flag"
	"golang.org/x/exp/slices"
)

var config = ""

type ProfileUserPassword struct {
	url      string
	user     string
	password string
}
type ProfileJwt struct {
	url     string
	token   string
	isShell bool
}

type Options struct {
	// default section name. default "__default"
	DefSection string
}

func init() {
	if config == "" {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Fprintf(os.Stderr, fmt.Sprintf("\u001B[0;31m✗\u001B[0m %s\n", err))
			os.Exit(1)
		}
		config = home + "/.config/.airflow/.config"
	}
	ini.New()
	ini.WithOptions(ini.ParseEnv, ini.ParseVar, func(opts *ini.Options) {
		opts.DefSection = "__default"
	})

	err := ini.LoadExists(config)
	if err != nil {
		panic(err)
	}
}

func GetProfiles() map[string]ini.Section {
	return ini.Default().Data()
}

func GetJwtProfile(profile_name string) ProfileJwt {
	p := ini.StringMap(profile_name)
	isShell, err := strconv.ParseBool(p["isShell"])
	if err != nil {
		fmt.Fprintf(os.Stderr, fmt.Sprintf("\u001B[0;31m✗\u001B[0m %s\n", "Not a valable isShell bool for "+profile_name))
		os.Exit(1)
	}
	profile := ProfileJwt{
		url:     p["url"],
		token:   p["token"],
		isShell: isShell,
	}
	return profile
}

func GetUserPasswordProfile(profile_name string) ProfileUserPassword {
	p := ini.StringMap(profile_name)
	profile := ProfileUserPassword{
		url:      p["url"],
		user:     p["user"],
		password: p["password"],
	}
	return profile
}

func GetActiveProfile() (string, string, error) {
	profile_name := flag.Flag
	if !slices.Contains(ini.SectionKeys(true), profile_name) {
		fmt.Fprintf(os.Stderr, fmt.Sprintf("\u001B[0;31m✗\u001B[0m %s\n", "no such a profile "+profile_name))
		os.Exit(1)
	}

	if ini.String(profile_name+".url") == "" {
		fmt.Fprintf(os.Stderr, fmt.Sprintf("\u001B[0;31m✗\u001B[0m %s\n", "No url defined for this profile."))
		os.Exit(1)
	}

	if ini.String(profile_name+".isShell") != "" && ini.String(profile_name+".token") != "" {
		return profile_name, "jwt", nil
	}

	if ini.String(profile_name+".user") != "" && ini.String(profile_name+".password") != "" {
		return profile_name, "user/password", nil
	}
	return profile_name, "", errors.New("The profile does not have any auth methode with all field filled.")
}

func create(p string) (*os.File, error) {
	if err := os.MkdirAll(filepath.Dir(p), 0770); err != nil {
		return nil, err
	}
	return os.Create(p)
}

func GetToken(p ProfileJwt) string {
	if !p.isShell {
		return p.token
	}
	out, err := exec.Command("bash", "-c", p.token).Output()
	if err != nil {
		log.Fatal(err)
	}
	return string(out[:])
}

func BasicAuth(p ProfileUserPassword) string {
	auth := p.user + ":" + p.password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func WriteConfig() {
	// no config yet
	if _, err := os.Stat(config); errors.Is(err, os.ErrNotExist) {
		_, err = create(config)
		if err != nil {
			fmt.Fprintf(os.Stderr, fmt.Sprintf("\u001B[0;31m✗\u001B[0m %s\n", err))
			os.Exit(1)
		}
	}

	_, err := ini.Default().WriteToFile(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, fmt.Sprintf("\u001B[0;31m✗\u001B[0m %s\n", err))
		os.Exit(1)
	}
}
