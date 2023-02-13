package workerpool

import (
	"context"
	"devices-test/configs"
	"log"
	"sync"
)

func StartWorkerpool(cfg *configs.Config) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	wg := &sync.WaitGroup{}
	for i := 1; i <= cfg.WorkersNum; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			worker := NewWorker(i)
			if err := worker.Work(ctx, cfg); err != nil {
				cancel()
				return
			}
			log.Printf("Worker %d successfully passed the test", worker.ID)
		}(i)
	}
	wg.Wait()
}
