package util

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
		fb.filter[fieldName] = v
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
	raw := strings.TrimSpace(q.Get(fieldName))
	if raw == "" {
		return fb
	}
	parts := []string{}
	for p := range strings.SplitSeq(raw, ",") {
		if t := strings.TrimSpace(p); t != "" {
			parts = append(parts, t)
		}
	}
	if len(parts) > 0 {
		fb.filter[fieldName] = bson.M{"$all": parts}
	}
	return fb
}

// WithRegexField allows mapping a query param to a field with case-insensitive regex
func (fb *FilterBuilder) WithRegexField(q url.Values, paramName, fieldName string) *FilterBuilder {
	if _, exists := fb.filter[fieldName]; exists {
		return fb
	}
	if v := strings.TrimSpace(q.Get(paramName)); v != "" {
		fb.filter[fieldName] = bson.M{"$regex": v, "$options": "i"}
	}
	return fb
}
