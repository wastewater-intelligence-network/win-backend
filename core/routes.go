package core

import (
	"github.com/gin-gonic/gin"
	"github.com/wastewater-intelligence-network/win-api/utils"
)

type RouteObj struct {
	RouteType    string
	RelativePath string
	Handler      gin.HandlerFunc
	PolicyRoles  utils.PolicyRules
}

func (win WinApp) getRoutes() []RouteObj {

	var routeList = []RouteObj{
		{
			"POST",
			"/register",
			win.handleNewUser,
			utils.PolicyRules{"admin"},
		}, {
			"GET",
			"/login",
			win.handleCreateToken,
			utils.PolicyRules{},
		}, {
			"GET",
			"/getSchedule",
			win.handleGetSchedule,
			utils.PolicyRules{utils.PolicyOpen},
		}, {
			"POST",
			"/setSchedule",
			win.handleSetSchedule,
			utils.PolicyRules{"admin"},
		}, {
			"POST",
			"/samplingRequest",
			win.handleSamplingRequest,
			utils.PolicyRules{"collector"},
		}, {
			"PATCH",
			"/samplingStatus",
			win.handleSamplingStatusPatch,
			utils.PolicyRules{"admin", "transporter", "technician"},
		}, {
			"POST",
			"/setCollectionPoints",
			win.handleSetCollectionPoints,
			utils.PolicyRules{"admin"},
		}, {
			"GET",
			"/getCollectionPoints",
			win.handleGetCollectionPoints,
			utils.PolicyRules{utils.PolicyAllUsers},
		}, {
			"GET",
			"/getSamplesCollectedOn",
			win.handleGetSamplesCollectedOn,
			utils.PolicyRules{"admin", "collector"},
		}, {
			"POST",
			"/setPointForSurvey",
			win.handleSetPointForSurvey,
			utils.PolicyRules{"admin"},
		}, {
			"GET",
			"/getAllPointsSurvey",
			win.handleGetAllPointsSurvey,
			utils.PolicyRules{"admin"},
		},
		{
			"PATCH",
			"/upgradeSurveyPointToCollectionPoint",
			win.handleUpgradeSurveyPointToCollectionPoint,
			utils.PolicyRules{"admin"},
		},
	}
	return routeList
}
