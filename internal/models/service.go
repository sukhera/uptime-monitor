package models

import "time"

type Service struct {
	Name           string            `bson:"name" json:"name"`
	Slug           string            `bson:"slug" json:"slug"`
	URL            string            `bson:"url" json:"url"`
	Headers        map[string]string `bson:"headers" json:"headers"`
	ExpectedStatus int               `bson:"expected_status" json:"expected_status"`
	Enabled        bool              `bson:"enabled" json:"enabled"`
}

type StatusLog struct {
	ServiceName string    `bson:"service_name" json:"service_name"`
	Status      string    `bson:"status" json:"status"`
	Latency     int64     `bson:"latency_ms" json:"latency_ms"`
	StatusCode  int       `bson:"status_code" json:"status_code"`
	Error       string    `bson:"error,omitempty" json:"error,omitempty"`
	Timestamp   time.Time `bson:"timestamp" json:"timestamp"`
}

type ServiceStatus struct {
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	Latency   int64     `json:"latency_ms"`
	UpdatedAt time.Time `json:"updated_at"`
}