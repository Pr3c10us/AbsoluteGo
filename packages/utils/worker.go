package utils

import "sync"

func RunWorkerPool[T any](items []T, maxWorkers int, fn func(T) error) error {
	if len(items) == 0 {
		return nil
	}
	if maxWorkers <= 0 {
		maxWorkers = 1
	}
	if maxWorkers > len(items) {
		maxWorkers = len(items)
	}

	var (
		wg       sync.WaitGroup
		once     sync.Once
		firstErr error
		work     = make(chan T, len(items))
	)

	for _, item := range items {
		work <- item
	}
	close(work)

	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for item := range work {
				if err := fn(item); err != nil {
					once.Do(func() { firstErr = err })
					return
				}
			}
		}()
	}

	wg.Wait()
	return firstErr
}
