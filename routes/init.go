package routes

import (
	"github.com/NickLovera/KzStats/mgr"
	rest "github.com/NickLovera/rest-utils-go"
	"github.com/gobuffalo/packr/v2"
	"github.com/gorilla/mux"
	"net/http"
)

type KzStatsServer struct {
	requestService mgr.IRequestService
}

func NewKzStatsServer(requestService mgr.IRequestService) *KzStatsServer {
	return &KzStatsServer{
		requestService: requestService,
	}
}

// InitServer - initialize web server
func InitServer(requestService mgr.IRequestService) {
	server := NewKzStatsServer(requestService)
	server.defineRoutes()
}

func (kzs *KzStatsServer) defineRoutes() {
	apiPrefix := "/KzStats/v1"

	ar := mux.NewRouter().PathPrefix(apiPrefix).Subrouter()

	kzs.addRoutes(ar)

	r := mux.NewRouter()

	swaggerUiBox := packr.New("SwaggerUiBox", "../resources/swagger-ui")
	swaggerSpecBox := packr.New("SwaggerSpecBox", "../resources/api/server")
	r.PathPrefix("/swagger-ui").Handler(http.StripPrefix("/swagger-ui", http.FileServer(swaggerUiBox)))
	r.PathPrefix("/swagger").Handler(http.StripPrefix("/swagger", http.FileServer(swaggerSpecBox)))

	rest.InitServer().Add(r, "", rest.RouteInfo{SubRoutes: ar, IsSecure: false, ApiPrefix: apiPrefix}).Metrics().Start()
}
