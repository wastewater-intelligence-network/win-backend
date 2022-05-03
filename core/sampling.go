package core

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shaj13/go-guardian/v2/auth"
	"github.com/wastewater-intelligence-network/win-api/model"
	"github.com/wastewater-intelligence-network/win-api/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (win WinApp) GetNearbyPoints(gin *gin.Context, request model.SamplingRequest) ([]model.CollectionPoint, error) {
	distance := 0.1
	c := gin.Request.Context()
	var query bson.M

	if reflect.DeepEqual(request.Location, model.Location{}) {
		u, ok := gin.Get("user")
		if !ok {
			return nil, fmt.Errorf("User data not found")
		}
		user := u.(auth.Info) //.GetExtensions().Get("locationId")
		locationId := user.GetExtensions().Get("locationId")

		if locationId == "" {
			return nil, fmt.Errorf("No location ID  available for this user")
		}

		query = bson.M{
			"locationId": locationId,
		}
	} else {
		query = bson.M{
			"location": bson.M{
				"$geoWithin": bson.M{
					"$center": []interface{}{
						[]float64{request.Location.Coordinates[0], request.Location.Coordinates[1]},
						distance,
					},
				},
			},
		}
	}

	cur, err := win.conn.Find(SAMPLE_COLLECTION_DB, query)
	if err != nil {
		return nil, err
	}

	var collectionPoints []model.CollectionPoint
	if err = cur.All(c, &collectionPoints); err != nil {
		return nil, err
	}

	return collectionPoints, nil
}

func (win WinApp) InsertSampleCollectionRecord(c context.Context, request model.SamplingRequest, point model.CollectionPoint) (*model.Sample, error) {
	// Validation to check if sample was collected earlier on this location
	// Validation if the container id is used to collect a sample earlier
	filter := bson.M{
		"$and": []bson.M{
			{
				"sampleTakenOn": bson.M{
					"$gt": utils.GetDayTime(0, 0, 0, 0, ""),
					"$lt": utils.GetDayTime(23, 59, 59, 0, ""),
				},
			},
			{
				"$or": []bson.M{
					{
						"sampleCollectionLocation.pointId": point.PointId,
					},
					{
						"containerId": request.ContainerId,
					},
				},
			},
		},
	}
	cursor, err := win.conn.Find(SAMPLE_COLLECTION_RECORD_DB, filter)
	if err != nil {
		return nil, err
	}

	var res []model.Sample
	err = cursor.All(c, &res)
	if err != nil {
		return nil, err
	}

	if len(res) == 0 {
		sample := model.Sample{
			SampleTakenOn:            time.Now(),
			ContainerId:              request.ContainerId,
			SampleCollectionLocation: point,
			Status:                   model.SampleStatusCollected,
			StatusLogList: []model.StatusLog{
				{
					Timestamp: time.Now(),
					Status:    model.SampleStatusCollected,
				},
			},
			AdditionalData: request.AdditionalData,
		}
		res, err := win.conn.Insert(SAMPLE_COLLECTION_RECORD_DB, sample)
		if err != nil {
			return nil, err
		}
		sId, ok := res.InsertedID.(primitive.ObjectID)
		if ok {
			sample.SampleId = sId
		}
		return &sample, nil
	} else {
		return nil, fmt.Errorf("Record with container id or collection point exist")
	}
}

func (win WinApp) getSamplesBetweenTime(c context.Context, start, end time.Time) ([]model.Sample, error) {
	filter := bson.M{
		"sampleTakenOn": bson.M{
			"$gt": start,
			"$lt": end,
		},
	}
	cursor, err := win.conn.Find(SAMPLE_COLLECTION_RECORD_DB, filter)
	if err != nil {
		return nil, err
	}

	res := make([]model.Sample, 0)
	err = cursor.All(c, &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (win WinApp) isValidateSampleStatusPatch(statusPatch string, sample model.Sample) bool {
	var idxPatch int = 0
	var idxStatus int = 0
	for i, s := range model.StatusSequence {
		if string(s) == statusPatch {
			idxPatch = i
		}
		if s == sample.Status {
			idxStatus = i
		}
	}

	if idxPatch < 1 || idxPatch <= idxStatus {
		return false
	} else {
		return true
	}

	return false
}

func (win WinApp) getPreviousStatusList(statusPatch string, len int) []model.SampleStatus {
	for i, s := range model.StatusSequence {
		if string(s) == statusPatch {
			return model.StatusSequence[len:i]
		}
	}
	return []model.SampleStatus{}
}
