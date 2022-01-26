package core

import (
	"context"
	"fmt"
	"time"

	"github.com/wastewater-intelligence-network/win-api/model"
	"github.com/wastewater-intelligence-network/win-api/utils"
	"go.mongodb.org/mongo-driver/bson"
)

func (win WinApp) GetNearbyPoints(c context.Context, request model.SampleCollectionRequest) ([]model.CollectionPoint, error) {
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

func (win WinApp) InsertSampleCollectionRecord(c context.Context, request model.SampleCollectionRequest, point model.CollectionPoint) error {
	// Validation to check if sample was collected earlier on this location
	// Validation if the container id is used to collect a sample earlier
	filter := bson.M{
		"$and": []bson.M{
			{
				"sampleTakenOn": bson.M{
					"$gt": utils.GetDayStartTime(),
					"$lt": utils.GetDayEndTime(),
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

	var res []model.SampleCollection
	err = cursor.All(c, &res)
	if err != nil {
		return err
	}

	if len(res) == 0 {
		win.conn.Insert(SAMPLE_COLLECTION_RECORD_DB, model.SampleCollection{
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
