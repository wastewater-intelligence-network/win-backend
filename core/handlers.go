package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shaj13/go-guardian/v2/auth"
	"github.com/wastewater-intelligence-network/win-api/model"
)

func (win WinApp) handleCreateToken(c *gin.Context) {
	token := uuid.New().String()

	user, ok := c.Get("user")
	if !ok {
		c.AbortWithError(http.StatusBadRequest, errors.New("User not parsed"))
	}

	auth.Append(tokenStrategy, token, user.(auth.Info))
	c.Header("Authorization", "Bearer "+token)
	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

func (win WinApp) handleSetSchedule(c *gin.Context) {
	var schedule model.CollectionSchedule
	c.BindJSON(&schedule)
	/* t := model.CollectionSchedule{
		Date: "29/11/2021",
		Name: "Bhesan Jahangirabad",
		Time: "06:00 AM",
		Type: "STP",
		Location: model.Location{
			Latitude:  22.3433,
			Longitude: 77.36353,
		},
	} */
	err := win.conn.Insert("test", schedule)
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, schedule)
}

func (win WinApp) handleGetSchedule(c *gin.Context) {
	user, ok := c.Get("user")
	if !ok {
		c.AbortWithError(http.StatusBadRequest, errors.New("User not parsed"))
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "schedule",
		"user":    user.(auth.Info),
		"schedule": []gin.H{
			{
				"assignedPointId":   23,
				"assignedPointName": "Bhesan Jahangirabad",
				"assignedUserId":    1,
				"type":              "STP",
				"latitude":          23.4524242,
				"longitude":         77.3534242,
				"date":              "29/11/2021",
				"time":              "06:00 AM",
			},
			{
				"name":      "Pisad",
				"type":      "STP",
				"latitude":  23.6452552,
				"longitude": 77.3645478,
				"date":      "29/11/2021",
				"time":      "08:00 AM",
			},
		},
	})
}

func (win WinApp) handleStartSampleCollection(c *gin.Context) {
	decoder := json.NewDecoder(c.Request.Body)

	var sampleCollRequest model.SampleCollectionRequest
	err := decoder.Decode(&sampleCollRequest)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{
				"status":  500,
				"message": "Request body is not correct",
			},
		)
	}

	points, err := win.GetNearbyPoints(c.Request.Context(), sampleCollRequest)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusOK,
			gin.H{
				"status":  500,
				"message": "Error while retrieving the nearby point",
			},
		)
		return
	}

	if len(points) == 0 {
		c.AbortWithStatusJSON(
			http.StatusOK,
			gin.H{
				"status":  500,
				"message": "No collection point nearby",
			},
		)
		return
	} else if len(points) == 1 {
		err = win.InsertSampleCollectionRecord(c.Request.Context(), sampleCollRequest, points[0])
		if err != nil {
			c.AbortWithStatusJSON(
				http.StatusOK,
				gin.H{
					"status":  500,
					"message": err.Error(),
				},
			)
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status":  200,
			"message": "New sample added to the database",
		})
	} else {
		if sampleCollRequest.PointId != "" {
			for _, p := range points {
				if p.PointId == sampleCollRequest.PointId {
					err = win.InsertSampleCollectionRecord(c.Request.Context(), sampleCollRequest, p)
					if err != nil {
						c.AbortWithStatusJSON(
							http.StatusOK,
							gin.H{
								"status":  500,
								"message": err.Error(),
							},
						)
						return
					}

					c.JSON(http.StatusOK, gin.H{
						"status":  200,
						"message": "New sample added to the database",
					})
					return
				}
			}
			c.JSON(http.StatusOK, gin.H{
				"status":  500,
				"message": "No Point found with pointId: " + sampleCollRequest.PointId,
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"status":  501,
				"message": "Multiple points at the given location. Pick one pointId",
				"list":    points,
			})
		}
	}
}

func (win WinApp) handleStartTransportation(c *gin.Context) {
	// Input: sample id
	// Output: status

}

func (win WinApp) handleSetCollectionPoints(c *gin.Context) {
	var collectionPoints []model.CollectionPoint

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}

	err = json.Unmarshal(body, &collectionPoints)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}

	err = win.conn.DeleteCollection(SAMPLE_COLLECTION_DB)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}

	for _, collectionPoint := range collectionPoints {
		fmt.Println(collectionPoint)
		err = win.conn.Insert(SAMPLE_COLLECTION_DB, collectionPoint)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		}
	}

	c.JSON(http.StatusOK, collectionPoints)
}

func (win WinApp) handleGetCollectionPoints(c *gin.Context) {
	var collectionPoints []model.CollectionPoint

	cursor, err := win.conn.Find(SAMPLE_COLLECTION_DB, gin.H{})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}

	if err = cursor.All(c.Request.Context(), &collectionPoints); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}

	fmt.Print(collectionPoints[0])

	c.JSON(http.StatusOK, collectionPoints)
}

// get schedule
// start sample collection (container_id, location)
// start transportation
// receive sample
