package db

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

/*
	For LLM query tools
*/

type QueryItem struct {
	Album      string       `json:"album" bson:"album" jsonschema:"description=album of query item, used to search the album"`
	Artists    []string     `json:"artists" bson:"artists" jsonschema:"description=artists of query item"`
	Title      string       `json:"title" bson:"title" jsonschema:"description=title of query item"`
	Duration   JSONDuration `json:"duration" bson:"duration" jsonschema:"description=duration to search for"`
	StartPoint *time.Time   `json:"start_point" bson:"start_point" jsonschema:"description=start point of time, used with duration"`
}
type JSONDuration time.Duration

func (d *JSONDuration) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	if s == "" {
		*d = JSONDuration(0)
		return nil
	}

	// 支持自定义单位
	unitMap := map[string]time.Duration{
		"h":  time.Hour,
		"m":  time.Minute,
		"s":  time.Second,
		"ms": time.Millisecond,
		"us": time.Microsecond,
		"µs": time.Microsecond,
		"ns": time.Nanosecond,
		"d":  24 * time.Hour,
		"w":  7 * 24 * time.Hour,
		"mo": 30 * 24 * time.Hour,
		"y":  365 * 24 * time.Hour,
	}

	// 提取数字和单位
	var numPart string
	var unitPart string
	for i, r := range s {
		if (r < '0' || r > '9') && r != '.' {
			numPart = s[:i]
			unitPart = s[i:]
			break
		}
	}
	if numPart == "" || unitPart == "" {
		return fmt.Errorf("invalid duration: %s", s)
	}

	multiplier, ok := unitMap[unitPart]
	if !ok {
		return fmt.Errorf("unknown unit %q in duration %q", unitPart, s)
	}

	value, err := strconv.ParseFloat(numPart, 64)
	if err != nil {
		return err
	}

	*d = JSONDuration(time.Duration(value * float64(multiplier)))
	return nil
}

func (d JSONDuration) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Duration(d).String())
}

func (q *QueryItem) ToFilter() bson.M {
	filter := bson.M{}

	if q.Album != "" {
		filter["album"] = q.Album
	}

	if len(q.Artists) > 0 {
		// 如果希望包含任意一个艺术家
		filter["artists"] = bson.M{"$in": q.Artists}

		// 如果希望同时包含所有指定艺术家，可以改为：
		// filter["artists"] = bson.M{"$all": q.Artists}
	}

	if q.Title != "" {
		filter["title"] = q.Title
	}
	var start, end time.Time
	now := time.Now()

	// 如果用户没指定 StartPoint，就从 “当前时间往前数 Duration”
	if q.StartPoint == nil || q.StartPoint.IsZero() {
		if q.Duration > 0 {
			start = now.Add(-time.Duration(q.Duration))
			end = now
		}
	} else {
		start = *q.StartPoint
		if q.Duration > 0 {
			end = q.StartPoint.Add(time.Duration(q.Duration))
		}
	}

	if !start.IsZero() {
		timeFilter := bson.M{"$gte": start}
		if !end.IsZero() {
			timeFilter["$lte"] = end
		}
		filter["last_modified"] = timeFilter
	}
	return filter
}

func QuerySinglesFromMongoDB(ctx context.Context, params *QueryItem) (map[string]interface{}, error) {
	filter := params.ToFilter()
	cursor, err := singlesCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var singles []AlbumSingle
	if err = cursor.All(ctx, &singles); err != nil {
		return nil, err
	}
	kv := map[string]interface{}{
		"singles": singles,
		"count":   len(singles),
	}
	return kv, nil
}
func QuerySinglesFromMongoDBFunc(ctx context.Context, params *QueryItem) (string, error) {
	kv, err := QuerySinglesFromMongoDB(ctx, params)
	if err != nil {
		return "", err
	}
	data, err := json.MarshalIndent(kv, "", "  ")
	if err != nil {
		return "Failed to MarshalIndent", err
	}
	return string(data), nil
}
