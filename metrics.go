package main

import (
	"context"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"time"
)

const (
	serverURL      = "http://34.86.236.100/"
	influxDBToken  = "I_UycfPULIG3VFr6eT-b0EzSIESMVb6rxZlS3n49zwHAcmpjPXQPS4u0eaZNY69hsWIVErE--T3lodcHQyx5rA=="
	influxDBOrg    = "API-Observability"
	influxDBBucket = "gRPC-Metrics"
)

func MetricsInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()
	resp, err := handler(ctx, req)
	duration := time.Since(start)
	// Measure request and response size (assuming they can be converted to string)
	reqSize := proto.Size(req.(proto.Message))
	respSize := 0
	if resp != nil {
		respSize = proto.Size(resp.(proto.Message))
	}

	// Get method name
	methodName := info.FullMethod

	// Record metrics to InfluxDB (or print to console, log, etc.)
	writeMetrics(duration, err, reqSize, respSize, methodName)

	return resp, err
}

func writeMetrics(duration time.Duration, err error, reqSize, respSize int, methodName string) {
	// Create a new InfluxDB client
	client := influxdb2.NewClient(serverURL, influxDBToken)
	defer client.Close()

	// Create a write API (this can be reused)
	writeAPI := client.WriteAPI(influxDBOrg, influxDBBucket)

	// Create a point to write (measurement name is "gRPCMetrics")
	point := influxdb2.NewPoint(
		"gRPCMetrics",
		map[string]string{"unit": "seconds", "method": methodName},
		map[string]interface{}{
			"duration":     duration.Seconds(),
			"error":        err != nil,
			"requestSize":  reqSize,
			"responseSize": respSize,
		},
		time.Now(),
	)

	// Write the point
	writeAPI.WritePoint(point)

	// Ensure data is written
	writeAPI.Flush()
}
