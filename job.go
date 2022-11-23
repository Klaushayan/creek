package brook

import "time"

type Job struct {
	JobChan chan func()
	PingChan chan string
	Quit    chan struct{}
	ticker  *time.Ticker
	job     func()
	jobArg  func(string)
}

func NewJob(job func()) *Job {
	return &Job{
		Quit: make(chan struct{}),
		job:  job,
	}
}

func NewJobWithArgument(job func(string)) *Job {
	return &Job{
		Quit: make(chan struct{}),
		jobArg:  job,
	}
}

/* Uses a ticker as the job channel based on the interval */
func (j *Job) StartWithTicker(interval time.Duration) {
	j.ticker = time.NewTicker(interval)
	go func() {
		for {
			select {
			case <-j.ticker.C:
				j.job()
			case <-j.Quit:
				// Not thread safe
				j.ticker.Stop()
				return
			}
		}
	}()
}

/* Uses a func channel as the job channel and does not use the job function given in the constructor*/
func (j *Job) StartWithChannel() {
	go func() {
		for {
			select {
			case job := <-j.JobChan:
				job()
			case <-j.Quit:
				return
			}
		}
	}()
}

/* Uses a string channel to get strings and executes the job function given in the constructor with the string value in the channel*/
func (j *Job) StartWithArgument() {
	go func() {
		for {
			select {
			case s := <-j.PingChan:
				j.jobArg(s)
			case <-j.Quit:
				return
			}
		}
	}()
}

/* Uses a string channel to get pinged and executes the job function given in the constructor*/
func (j *Job) Start() {
	go func() {
		for {
			select {
			case <-j.PingChan:
				j.job()
			case <-j.Quit:
				return
			}
		}
	}()
}

/* This stops all jobs */
func (j *Job) Stop() {
	close(j.Quit)
}
