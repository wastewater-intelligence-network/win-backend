package core

import (
	"github.com/gin-gonic/gin"
	"github.com/wastewater-intelligence-network/win-api/db"
)

var (
	SAMPLE_COLLECTION_DB        = "collection_points"
	SAMPLE_COLLECTION_RECORD_DB = "sample_collection"
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
	app.gin.Use(AuthMiddleware())
	app.setRoutes()
	return app, nil
}

func (win WinApp) Run() {
	win.gin.Run("127.0.0.1:8080")
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
		}
	}
}
