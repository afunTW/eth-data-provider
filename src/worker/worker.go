package worker

import "context"

type BlockWorker interface {
	Start(ctx context.Context)
	Stop()
}
