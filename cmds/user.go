package cmds

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/jackgpalfrey/peurl-cli/utils"
)

func User(args []string) {
	subsubcommand := args[0]

	switch subsubcommand {
	case "list", "ls":
		UserList()
	case "inspect":
		UserInspect(args[1:])
	case "create":
		UserCreate(args[1:])
	case "delete", "rm":
		UserDelete(args[1:])
	case "passwd":
		UserPasswd(args[1:])
	default:
		fmt.Printf("Invalid sub command '%s'\n", "user "+subsubcommand)
	}
}

type UserDetails struct {
	Username string `json:"username"`
}

func UserList() {
	url := utils.GetURL("/users")

	res := utils.SendAuthedRequest("GET", url, nil)
	switch res.StatusCode {
	case 200:
	case 401:
		utils.ExitWithError("Unauthorized", nil)
	case 403:
		utils.ExitWithError("Forbidden", nil)
	default:
		utils.ExitWithError("Error code "+fmt.Sprint(res.StatusCode), nil)
	}

	body, _ := io.ReadAll(res.Body)
	users := []UserDetails{}
	json.Unmarshal(body, &users)

	for _, user := range users {
		fmt.Println(user.Username)
	}
}

func UserInspect(args []string) {
	if len(args) != 1 {
		utils.UsageError("peurl user inspect <username>")
	}

	url := utils.GetURL("/users/" + args[0])

	res := utils.SendAuthedRequest("GET", url, nil)
	switch res.StatusCode {
	case 200:
	case 401:
		utils.ExitWithError("Unauthorized", nil)
	case 403:
		utils.ExitWithError("Forbidden", nil)
	default:
		utils.ExitWithError("Error code "+fmt.Sprint(res.StatusCode), nil)
	}

	body, _ := io.ReadAll(res.Body)
	user := UserDetails{}
	json.Unmarshal(body, &user)
	fmt.Println(user.Username)
}

func UserCreate(args []string) {
	if len(args) != 2 {
		utils.UsageError("peurl user create <username> <password>")
	}

	url := utils.GetURL("/users")
	requestBody := utils.BuildJSONBody(map[string]string{
		"username": args[0],
		"password": args[1],
	})
	res := utils.SendAuthedRequest("POST", url, requestBody)
	switch res.StatusCode {
	case 201:
	case 401:
		utils.ExitWithError("Unauthorized", nil)
	case 403:
		utils.ExitWithError("Forbidden", nil)
	case 409:
		utils.ExitWithError("User already exists", nil)
	default:
		utils.ExitWithError("Error code "+fmt.Sprint(res.StatusCode), nil)
	}

	fmt.Println(args[0], "created")
}

func UserDelete(args []string) {
	if len(args) != 1 {
		utils.UsageError("peurl user delete <username>")
	}

	url := utils.GetURL("/users/" + args[0])

	res := utils.SendAuthedRequest("DELETE", url, nil)
	switch res.StatusCode {
	case 200:
	case 401:
		utils.ExitWithError("Unauthorized", nil)
	case 403:
		utils.ExitWithError("Forbidden", nil)
	default:
		utils.ExitWithError("Error code "+fmt.Sprint(res.StatusCode), nil)
	}

	fmt.Println(args[0], "deleted")
}

func UserPasswd(args []string) {
	if len(args) != 2 {
		utils.UsageError("peurl user passwd <username> <password>")
	}

	url := utils.GetURL("/users/" + args[0])
	requestBody := utils.BuildJSONBody(map[string]string{
		"password": args[1],
	})

	res := utils.SendAuthedRequest("PATCH", url, requestBody)
	switch res.StatusCode {
	case 200:
	case 401:
		utils.ExitWithError("Unauthorized", nil)
	case 403:
		utils.ExitWithError("Forbidden", nil)
	case 409:
		utils.ExitWithError("User already exists", nil)
	default:
		utils.ExitWithError("Error code "+fmt.Sprint(res.StatusCode), nil)
	}

	fmt.Println(args[0], "updated")
}
