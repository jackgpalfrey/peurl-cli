package cmds

import (
	"fmt"
	"io"
	"net/http"

	"github.com/jackgpalfrey/peurl-cli/utils"
)

func Login(args []string) {
	if len(args) != 2 {
		utils.UsageError("peurl login <username> <password>")
	}

	url := utils.GetURL("/login")

	username, password := args[0], args[1]
	requestBody := utils.BuildJSONBody(map[string]string{
		"username": username,
		"password": password,
	})

	resp, err := http.Post(url, "application/json", requestBody)
	if err != nil {
		utils.ExitWithError("Failed to send login request", err)
	}
	switch resp.StatusCode {
	case 200:
	case 401:
		utils.ExitWithError("Incorect credentials", nil)
	case 403:
		utils.ExitWithError("Forbidden", nil)
	default:
		utils.ExitWithError("Error code "+fmt.Sprint(resp.StatusCode), nil)
	}

	jwt := resp.Cookies()[0].Value
	utils.SaveConfig("jwt", jwt)
	fmt.Println("Logged in as", username)
}

func Logout(args []string) {
	utils.SaveConfig("jwt", "")
	fmt.Println("Logged out")
}

func Whoami() {
	url := utils.GetURL("/whoami")
	res := utils.SendAuthedRequest("GET", url, nil)
	body, _ := io.ReadAll(res.Body)
	fmt.Println(string(body))
}
