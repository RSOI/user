// QUESTION SERVICE

package main

import (
	"fmt"
	"os"

	"github.com/RSOI/user/controller"
	"github.com/RSOI/user/database"
	"github.com/RSOI/user/utils"
	"github.com/valyala/fasthttp"
)

// PORT application port
const PORT = 8082

func main() {
	if len(os.Args) > 1 {
		utils.DEBUG = os.Args[1] == "debug"
	}
	utils.LOG("Launched in debug mode...")
	utils.LOG(fmt.Sprintf("User service is starting on localhost: %d", PORT))

	controller.Init(database.Connect())
	fasthttp.ListenAndServe(fmt.Sprintf(":%d", PORT), initRoutes().Handler)
}
