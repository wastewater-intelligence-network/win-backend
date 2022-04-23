package core

import (
	"github.com/gin-gonic/gin"
	"github.com/wastewater-intelligence-network/win-api/db"
	"github.com/wastewater-intelligence-network/win-api/utils"
)

var (
	SAMPLE_COLLECTION_DB        = "collection_points"
	SAMPLE_COLLECTION_RECORD_DB = "sample_collection"
	WIN_COLLECTION_USERS        = "users"
	SURVEY_SAMPLING_SITE        = "survey_sampling_site"
	COLLECTION_SCHEDULES        = "schedules"
)

type WinApp struct {
	gin  *gin.Engine
	conn *db.DBConnection
}

func NewWinApp() (*WinApp, error) {
	conn, err := db.NewDBConnection()
	if err != nil {
		return nil, err
	}

	app := &WinApp{
		gin:  gin.Default(),
		conn: conn,
	}
	policy := app.getPolicy()
	app.gin.Use(AuthMiddleware(policy, conn))
	app.gin.Use(PolicyMiddleware(policy))
	app.setRoutes()
	return app, nil
}

func (win WinApp) Run() {
	win.gin.Run("0.0.0.0:8002")
}

func (win WinApp) getPolicy() *utils.Policy {
	var policy = utils.NewPolicy()
	for _, r := range win.getRoutes() {
		policy.AddRule(r.RelativePath, r.PolicyRoles)
	}
	return policy
}

func (win WinApp) setRoutes() {
	for _, r := range win.getRoutes() {
		switch r.RouteType {
		case "GET":
			win.gin.GET(r.RelativePath, r.Handler)
			break
		case "POST":
			win.gin.POST(r.RelativePath, r.Handler)
			break
		case "PATCH":
			win.gin.PATCH(r.RelativePath, r.Handler)
			break
		}
	}
}
