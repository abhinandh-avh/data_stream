package route

import (
	"datastream/api"
	"net/http"
)

func SetupRoutes() {
	http.HandleFunc("/", api.HomePage)

}
