package core

import (
	"context"
	"fmt"
	"time"

	"github.com/wastewater-intelligence-network/win-api/model"
	"github.com/wastewater-intelligence-network/win-api/utils"
	"go.mongodb.org/mongo-driver/bson"
)

func (win WinApp) GetNearbyPoints(c context.Context, request model.SamplingRequest) ([]model.CollectionPoint, error) {
	cur, err := win.conn.Find(SAMPLE_COLLECTION_DB, bson.M{
		"location": bson.M{
			"$geoWithin": bson.M{
				"$center": []interface{}{
					[]float64{request.Location.Coordinates[0], request.Location.Coordinates[1]},
					0.01,
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}

	var collectionPoints []model.CollectionPoint
	if err = cur.All(c, &collectionPoints); err != nil {
		return nil, err
	}

	return collectionPoints, nil
}

func (win WinApp) InsertSampleCollectionRecord(c context.Context, request model.SamplingRequest, point model.CollectionPoint) error {
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
		return err
	}

	var res []model.Sample
	err = cursor.All(c, &res)
	if err != nil {
		return err
	}

	if len(res) == 0 {
		win.conn.Insert(SAMPLE_COLLECTION_RECORD_DB, model.Sample{
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
		})
	} else {
		return fmt.Errorf("Record with container id or collection point exist")
	}
	return nil
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
	var idx int = 0
	for i, s := range model.StatusSequence {
		if string(s) == statusPatch {
			idx = i
		}
	}

	if idx < 1 {
		return false
	} else if sample.Status == model.StatusSequence[idx-1] {
		return true
	}

	return false
}
