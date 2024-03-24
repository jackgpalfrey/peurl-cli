package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/jackgpalfrey/peurl-cli/config"
)

func GetDomain() string {
	homeDir := GetHomeDir()
	domainURL, err := os.ReadFile(homeDir + "/.config/peurl/domain")
	if err != nil || string(domainURL) == "" {
		ExitWithError("No domain set. Use 'peurl domain <url>' to set", err)
	}

	return string(domainURL)
}

func GetURL(path string) string {
	return GetDomain() + path
}

func SetDomain(url string) {
	if url[:7] != "http://" && url[:8] != "https://" {
		url = "https://" + url
	}

	if url[len(url)-1] == '/' {
		url = url[:len(url)-1]
	}

	SaveConfig("domain", url)
	fmt.Println("Domain set to", url)
}

func EnsureURLStartsWithProto(url string, defaultProto string) string {
	for idx, char := range url {
		if char == ':' {
			if len(url) > idx+3 && url[idx+1] == '/' && url[idx+2] == '/' {
				return url
			} else {
				break
			}
		}
	}

	return defaultProto + "://" + url
}

func SaveConfig(key string, value string) {
	configDirPath := GetHomeDir() + "/.config/peurl"

	err := os.Mkdir(configDirPath, 0o700)
	if err != nil && !os.IsExist(err) {
		ExitWithError("Failed to create config directory", err)
	}

	configPath := configDirPath + "/" + key
	err = os.WriteFile(configPath, []byte(value), 0o700)
	if err != nil {
		ExitWithError("Failed to save "+key+" in config", err)
	}
}

func GetConfig(key string) string {
	path := GetHomeDir() + "/.config/peurl/" + key

	data, err := os.ReadFile(path)
	if err != nil {
		ExitWithError("'"+key+"' not configured", err)
	}
	return string(data)
}

func GetHomeDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		ExitWithError("Failed to get user home directory", err)
	}
	return homeDir
}

func AddAuthCookie(req *http.Request) {
	req.AddCookie(&http.Cookie{Name: "jwt", Value: GetConfig("jwt")})
}

func SendAuthedRequest(method string, url string, body io.Reader) *http.Response {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		ExitWithError("Failed to build request", err)
	}
	AddAuthCookie(req)
	req.Header.Set("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		ExitWithError("Failed to send request", err)
	}

	return res
}

func BuildJSONBody(data map[string]string) *bytes.Buffer {
	requestBodyJSON, err := json.Marshal(data)
	if err != nil {
		ExitWithError("Failed to build request body", err)
	}
	requestBody := bytes.NewBuffer(requestBodyJSON)
	return requestBody
}

func OpenURL(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		ExitWithError("Failed to open URL", err)
	}
}

func CheckVersionCompat() {
	res := SendAuthedRequest("GET", GetURL("/version"), nil)
	if res.StatusCode != 200 {
		ExitWithError("Failed to check server version", nil)
	}
	body, _ := io.ReadAll(res.Body)
	serverVersion := string(body)

	versionNum := strings.Split(serverVersion, "v")[1]
	if versionNum[:3] != config.VERSION[:3] {
		ExitWithError("Incompatible server version: "+serverVersion, nil)
	}
}
