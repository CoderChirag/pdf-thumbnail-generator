package thumbnail_controller

import (
	"fmt"
	"net/http"
)

func RegisterRoutes() {
	fmt.Println("Registering route: /thumbnail")
	http.HandleFunc("/thumbnail", generateThumbnail)
}
