package core

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shaj13/go-guardian/v2/auth"
	"github.com/wastewater-intelligence-network/win-api/model"
	"github.com/wastewater-intelligence-network/win-api/utils"
	"go.mongodb.org/mongo-driver/bson"
)

func (win WinApp) handleNewUser(c *gin.Context) {
	decoder := json.NewDecoder(c.Request.Body)

	var user model.User
	err := decoder.Decode(&user)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{
				"status":  505,
				"message": err.Error(),
			},
		)
		return
	}

	passwordHash := sha1.Sum([]byte(user.Password))
	user.Hash = hex.EncodeToString(passwordHash[:])
	user.Password = ""

	_, err = win.conn.Insert(WIN_COLLECTION_USERS, user)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{
				"status":  506,
				"message": err.Error(),
			},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"status":  200,
			"message": "User Created",
		},
	)
}

func (win WinApp) handleCreateToken(c *gin.Context) {
	token := uuid.New().String()

	u, ok := c.Get("user")
	if !ok {
		c.AbortWithError(http.StatusBadRequest, errors.New("User not parsed"))
	}

	user := u.(auth.Info)
	auth.Append(tokenStrategy, token, user)
	c.Header("Authorization", "Bearer "+token)
	c.JSON(http.StatusOK, gin.H{
		"token":                token,
		"roles":                user.GetGroups(),
		"tokenExpiryTimestamp": time.Now().Add(TokenExpiryHours).Unix(),
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
	_, err := win.conn.Insert(COLLECTION_SCHEDULES, schedule)
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, schedule)
}

func (win WinApp) handleGetSchedule(c *gin.Context) {
	var schedules []model.CollectionPointSchedule
	schedules = []model.CollectionPointSchedule{}

	cursor, err := win.conn.Find(COLLECTION_SCHEDULES, gin.H{})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err = cursor.All(c.Request.Context(), &schedules); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, schedules)
}

func (win WinApp) handleSamplingRequest(c *gin.Context) {
	decoder := json.NewDecoder(c.Request.Body)

	var sampleCollRequest model.SamplingRequest
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

	points, err := win.GetNearbyPoints(c, sampleCollRequest)
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
		sample, err := win.InsertSampleCollectionRecord(c.Request.Context(), sampleCollRequest, points[0])
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
			"status":        200,
			"message":       "New sample added to the database",
			"sampleDetails": sample,
		})
	} else {
		if sampleCollRequest.PointId != "" {
			for _, p := range points {
				if p.PointId == sampleCollRequest.PointId {
					sample, err := win.InsertSampleCollectionRecord(c.Request.Context(), sampleCollRequest, p)
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
						"status":        200,
						"message":       "New sample added to the database",
						"sampleDetails": sample,
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

func (win WinApp) handleSamplingStatusPatch(c *gin.Context) {
	decoder := json.NewDecoder(c.Request.Body)

	var samplingStatusRequest model.SamplingStatusRequest
	if err := decoder.Decode(&samplingStatusRequest); err != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{
				"status":  500,
				"message": "Request body is not correct",
			},
		)
		return
	}

	result := win.conn.FindOne(SAMPLE_COLLECTION_RECORD_DB, bson.M{
		"containerId": samplingStatusRequest.ContainerId,
	})

	var sample model.Sample
	err := result.Decode(&sample)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusOK,
			gin.H{
				"status":  500,
				"message": "Cannot parse the sampleId. Err: " + err.Error(),
			},
		)
		return
	}

	if win.isValidateSampleStatusPatch(samplingStatusRequest.StatusPatch, sample) {
		sample.Status = model.SampleStatus(samplingStatusRequest.StatusPatch)

		prevStatusList := win.getPreviousStatusList(samplingStatusRequest.StatusPatch, len(sample.StatusLogList))

		fmt.Println(prevStatusList)

		for _, prevStatus := range prevStatusList {
			sample.StatusLogList = append(sample.StatusLogList, model.StatusLog{
				Timestamp: time.Now(),
				Status:    prevStatus,
				Error:     "step_skipped",
				Message:   "This sample is not marked for this step",
			})
		}

		sample.StatusLogList = append(sample.StatusLogList, model.StatusLog{
			Timestamp: time.Now(),
			Status:    sample.Status,
		})
		_, err := win.conn.UpdateOne(SAMPLE_COLLECTION_RECORD_DB, bson.M{
			"containerId": samplingStatusRequest.ContainerId,
		}, bson.M{
			"$set": sample,
		})
		if err != nil {
			c.AbortWithStatusJSON(
				http.StatusOK,
				gin.H{
					"status":  504,
					"message": "Cannot update the sample status: Err: " + err.Error(),
				},
			)
			return
		}
		c.JSON(
			http.StatusOK,
			gin.H{
				"status":  200,
				"message": "Status changed to " + samplingStatusRequest.StatusPatch,
				"sample":  sample,
			},
		)
	} else {
		c.AbortWithStatusJSON(
			http.StatusOK,
			gin.H{
				"status":  504,
				"message": "Cannot update the sample status: Err: New status is not compatible",
			},
		)
	}
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
		_, err = win.conn.Insert(SAMPLE_COLLECTION_DB, collectionPoint)
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
		return
	}

	if err = cursor.All(c.Request.Context(), &collectionPoints); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	fmt.Print(collectionPoints[0])

	c.JSON(http.StatusOK, collectionPoints)
}

func (win WinApp) handleGetSamplesCollectedOn(c *gin.Context) {
	date := c.Query("date")
	fmt.Println(date)
	start := utils.GetDayTime(0, 0, 0, 0, date)
	end := utils.GetDayTime(23, 59, 59, 0, date)

	samples, err := win.getSamplesBetweenTime(c.Request.Context(), start, end)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{
				"status":  500,
				"message": "Could not get samples data. Err: " + err.Error(),
			},
		)
		return
	}

	c.JSON(http.StatusOK, samples)
}

func (win WinApp) handleSetPointForSurvey(c *gin.Context) {
	var collectionPoint model.CollectionPoint

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = json.Unmarshal(body, &collectionPoint)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	_, err = win.conn.Insert(SURVEY_SAMPLING_SITE, collectionPoint)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, collectionPoint)
}

func (win WinApp) handleGetAllPointsSurvey(c *gin.Context) {
	var collectionPoints []model.CollectionPoint

	cursor, err := win.conn.Find(SURVEY_SAMPLING_SITE, gin.H{})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err = cursor.All(c.Request.Context(), &collectionPoints); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, collectionPoints)
}

func (win WinApp) handleUpgradeSurveyPointToCollectionPoint(c *gin.Context) {
	decoder := json.NewDecoder(c.Request.Body)

	var collectionPoint model.CollectionPoint
	if err := decoder.Decode(&collectionPoint); err != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{
				"status":  500,
				"message": "Could not parse the request body",
			},
		)
		return
	}

	res := win.conn.FindOne(SURVEY_SAMPLING_SITE, gin.H{
		"pointId": collectionPoint.PointId,
	})
	err := res.Decode(&collectionPoint)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	_, err = win.conn.Insert(SAMPLE_COLLECTION_DB, collectionPoint)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, collectionPoint)
}

// get schedule
// start sample collection (container_id, location)
// start transportation
// receive sample
