package kepler

import "sync"

func merge(cnannels ...<-chan Message) <-chan Message {

	out := make(chan Message)
	var wg sync.WaitGroup
	wg.Add(len(cnannels))
	for _, c := range cnannels {
		go func(c <-chan Message) {
			for v := range c {
				out <- v
			}
			wg.Done()
		}(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
