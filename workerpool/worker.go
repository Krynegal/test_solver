package workerpool

import (
	"context"
	"devices-test/client"
	"devices-test/configs"
	"fmt"
)

type worker struct {
	ID int
}

func NewWorker(ID int) *worker {
	return &worker{
		ID: ID,
	}
}

func (w *worker) Work(ctx context.Context, cfg *configs.Config) error {
	client := client.NewClient(ctx, cfg)
	if err := client.Run(); err != nil {
		return fmt.Errorf("error: %s workerpool id: %d", err.Error(), w.ID)
	}
	return nil
}
