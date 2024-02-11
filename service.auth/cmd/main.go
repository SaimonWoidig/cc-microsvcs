package main

import (
	"context"

	"github.com/SaimonWoidig/cc-microsvcs/common/otel"
)

func main() {
	otel.NewResource(context.Background(), otel.ResourceConfig{
		Name:       "auth-service",
		Version:    "1.0.0",
		Addr:       "localhost",
		Port:       8080,
		InstanceID: "123123",
	})
}
