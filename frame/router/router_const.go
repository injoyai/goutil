package router

const (
	Logo = "\n" +
		"		 __   __   __         __    ______   ____    ____   \n" +
		"		|  | |  \\ |  |       |  |  /  __  \\  \\   \\  /   /   \n" +
		"		|  | |   \\|  |       |  | |  |  |  |  \\   \\/   /    \n" +
		"		|  | |       | .--.  |  | |  |  |  |   \\_    _/     \n" +
		"		|  | |  |\\   | |  `--'  | |  `--'  |     |  |       \n" +
		"		|__| |__| \\__|  \\______/   \\______/      |__|       \n" +
		" "

	MarkExit = "EXIT"
)

var (
	CORS = map[string]string{
		"Access-Control-Allow-Methods":     "POST,GET,PUT,PATCH,OPTIONS,DELETE",
		"Access-Control-Allow-Headers":     "x-requested-with,content-type,Authorization",
		"Access-Control-Allow-Origin":      "*",
		"Access-Control-Max-Age":           "3600",
		"Access-Control-Allow-Credentials": "true",
	}
)
