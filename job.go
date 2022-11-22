package brook

import "time"

type Job struct {
	ticker *time.Ticker
	quit   chan struct{}
	job    func()
}

func NewJob(interval time.Duration, job func()) *Job {
	return &Job{
		ticker: time.NewTicker(interval),
		quit:   make(chan struct{}),
		job:    job,
	}
}

func (j *Job) Start() {
	go func() {
		for {
			select {
			case <-j.ticker.C:
				j.job()
			case <-j.quit:
				j.ticker.Stop()
				return
			}
		}
	}()
}

func (j *Job) Stop() {
	close(j.quit)
}