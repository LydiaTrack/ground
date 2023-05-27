package test_support

import (
	"context"
	"lydia-track-base/internal/mongodb"
)

func TestWithMongo() {
	context := context.Background()
	container, err := mongodb.StartContainer(context)
	if err != nil {
		panic(err)
	}

	defer container.Terminate(context)
}
