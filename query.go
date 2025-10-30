package main

import (
	"net/url"
	"strconv"
	"strings"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type FilterBuilder struct {
	filter bson.M
}

func NewFilterBuilder() *FilterBuilder {
	return &FilterBuilder{filter: bson.M{}}
}

func (fb *FilterBuilder) Build() bson.M {
	return fb.filter
}

func (fb *FilterBuilder) WithStringField(q url.Values, fieldName string) *FilterBuilder {
	if v := strings.TrimSpace(q.Get(fieldName)); v != "" {
		fb.filter[fieldName] = bson.M{"$regex": v, "$options": "i"}
	}
	return fb
}

func (fb *FilterBuilder) WithIntField(q url.Values, fieldName string) *FilterBuilder {
	if v := strings.TrimSpace(q.Get(fieldName)); v != "" {
		if intVal, err := strconv.Atoi(v); err == nil {
			fb.filter[fieldName] = intVal
		}
	}
	return fb
}

func (fb *FilterBuilder) WithArrayField(q url.Values, fieldName string) *FilterBuilder {
	parts := []string{}
	for p := range strings.SplitSeq(fieldName, ",") {
		if t := strings.TrimSpace(p); t != "" {
			parts = append(parts, t)
		}
	}
	if len(parts) > 0 {
		// match documents that contain all provided cuts
		fb.filter["cuts"] = bson.M{"$all": parts}
	}
	return fb
}
