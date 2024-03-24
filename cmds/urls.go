package cmds

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"text/tabwriter"

	"github.com/jackgpalfrey/peurl-cli/utils"
)

func UrlShorten(args []string) {
	short := ""
	long := ""

	for _, arg := range args {
		if len(arg) > 7 && arg[:7] == "-short=" {
			short = arg[7:]
		} else {
			if long != "" {
				utils.UsageError("peurl shorten [-short=<short>] <url>")
			}
			if arg[:1] == "-" {
				utils.UsageError("peurl shorten [-short=<short>] <url>")
			}

			long = arg
		}
	}

	long = utils.EnsureURLStartsWithProto(long, "http")

	url := utils.GetURL("/shorten")
	requestBody := utils.BuildJSONBody(map[string]string{
		"short": short,
		"long":  long,
	})
	res := utils.SendAuthedRequest("POST", url, requestBody)
	switch res.StatusCode {
	case 201:
	case 401:
		utils.ExitWithError("Unauthorized", nil)
	case 403:
		utils.ExitWithError("Forbidden", nil)
	case 409:
		utils.ExitWithError("Short URL '"+short+"' already exists", nil)
	default:
		utils.ExitWithError("Error code "+fmt.Sprint(res.StatusCode), nil)
	}
	body, _ := io.ReadAll(res.Body)
	urlDetails := UrlDetails{}
	json.Unmarshal(body, &urlDetails)

	fmt.Println(utils.GetURL("/" + utils.Blue + urlDetails.Short + utils.Reset))
}

type UrlDetails struct {
	Short    string `json:"short"`
	Long     string `json:"long"`
	Username string `json:"username"`
}

func GetURLDetails(short string) UrlDetails {
	url := utils.GetURL("/inspect/" + short)

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
	urlDetails := UrlDetails{}
	json.Unmarshal(body, &urlDetails)
	return urlDetails
}

func UrlInspect(args []string) {
	if len(args) != 1 {
		utils.UsageError("peurl inspect <short>")
	}

	urlDetails := GetURLDetails(args[0])

	fmt.Println("Short:", urlDetails.Short)
	fmt.Println("Long:", urlDetails.Long)
	fmt.Println("Username:", urlDetails.Username)
}

func UrlExpand(args []string) {
	if len(args) != 1 {
		utils.UsageError("peurl expand <short>")
	}
	urlDetails := GetURLDetails(args[0])
	fmt.Println(urlDetails.Long)
}

func UrlGo(args []string) {
	if len(args) != 1 {
		utils.UsageError("peurl go <short>")
	}

	utils.OpenURL(args[0])
}

func UrlDelete(args []string) {
	if len(args) != 1 {
		utils.UsageError("peurl delete <short>")
	}
	url := utils.GetURL("/inspect/" + args[0])
	res := utils.SendAuthedRequest("DELETE", url, nil)
	switch res.StatusCode {
	case 200:
	case 401:
		utils.ExitWithError("Unauthorized", nil)
	case 403:
		utils.ExitWithError("Forbidden", nil)
	case 404:
		utils.ExitWithError("Short URL '"+args[0]+"' not found", nil)
	default:
		utils.ExitWithError("Error code "+fmt.Sprint(res.StatusCode), nil)
	}

	fmt.Println("Deleted", args[0])
}

func UrlList(args []string) {
	url := utils.GetURL("/inspect")
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
	urls := []UrlDetails{}
	json.Unmarshal(body, &urls)

	w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
	for _, url := range urls {
		fmt.Fprintln(w, url.Short+"\t> "+url.Long+"\t("+url.Username+")")
	}
	w.Flush()
}
