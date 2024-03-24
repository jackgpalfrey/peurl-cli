package cmds

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/jackgpalfrey/peurl-cli/utils"
)

func Domain(args []string) {
	if len(args) == 0 {
		url := utils.GetDomain()
		utils.CheckVersionCompat()
		fmt.Println(url)
		os.Exit(0)
	}

	if len(args) != 1 {
		utils.UsageError("peurl domain <url>")
	}

	utils.SetDomain(args[0])
	utils.CheckVersionCompat()
}

func Help() {
	w := tabwriter.NewWriter(os.Stdout, 1, 1, 4, ' ', 0)
	fmt.Println("AVAILABLE SUBCOMMANDS:")
	fmt.Fprintln(w, "\thelp\tThis help page")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "CONFIGURATION")
	fmt.Fprintln(w, "\tdomain\tGet or set the url of peurl server")
	fmt.Fprintln(w, "\tlogin\tLogin to user account")
	fmt.Fprintln(w, "\tlogout\tLogout of user account")
	fmt.Fprintln(w, "\twhoami\tGet username of logged in user")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "URL MANAGEMENT")
	fmt.Fprintln(w, "\tshorten\tShorten a URL")
	fmt.Fprintln(w, "\tlist / ls\tList all shortened URLs (Requires admin)")
	fmt.Fprintln(w, "\tgo\tOpen a shortened URL in browser")
	fmt.Fprintln(w, "\tinspect\tGet details about a shortened URL")
	fmt.Fprintln(w, "\texpand\tOriginal URL of a shortened URL")
	fmt.Fprintln(w, "\tdelete / rm\tDelete a shortened URL")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "USER MANAGEMENT")
	fmt.Fprintln(w, "\tuser	fmt.Fprintl")
	fmt.Fprintln(w, "\t├─list / ls\tList all users (Requires admin)")
	fmt.Fprintln(w, "\t├─inspect\tGet details about a user (Requires admin)")
	fmt.Fprintln(w, "\t├─create\tCreate a new user (Requires admin)")
	fmt.Fprintln(w, "\t├─delete / rm\tDelete a user (Requires admin)")
	fmt.Fprintln(w, "\t└─passwd\tChange a user's password (Requires admin or same user)")
	w.Flush()
}
