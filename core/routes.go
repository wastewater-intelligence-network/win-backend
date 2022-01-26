package core

import (
	"github.com/gin-gonic/gin"
)

type RouteObj struct {
	RouteType    string
	RelativePath string
	Handler      gin.HandlerFunc
}

func (win WinApp) getRoutes() []RouteObj {

	var routeList = []RouteObj{
		{
			"GET",
			"/login",
			win.handleCreateToken,
		}, {
			"GET",
			"/getSchedule",
			win.handleGetSchedule,
		}, {
			"POST",
			"/setSchedule",
			win.handleSetSchedule,
		}, {
			"POST",
			"/samplingRequest",
			win.handleSamplingRequest,
		}, {
			"PATCH",
			"/samplingStatus",
			win.handleSamplingStatus,
		}, {
			"POST",
			"/setCollectionPoints",
			win.handleSetCollectionPoints,
		}, {
			"GET",
			"/getCollectionPoints",
			win.handleGetCollectionPoints,
		},
	}
	return routeList
}
