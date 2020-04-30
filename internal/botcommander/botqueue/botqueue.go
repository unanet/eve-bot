package botqueue

import (
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"syscall"
	"time"

	"gitlab.unanet.io/devops/eve-bot/internal/botmetrics"
	"gitlab.unanet.io/devops/eve/pkg/log"
	"go.uber.org/zap"
)

// EVEBOT_QUEUE_WORKERS
type Config struct {
	QueueWorkers int `split_words:"true" required:"false"`
}

// WorkRequest is an single unit of work that needs to be processed by a worker
type WorkRequest struct {
	Channel                              string
	User                                 string
	CommandName                          string
	ReceivedTS, InProcessTS, CompletedTS time.Time
	Delay                                time.Duration
}

// Worker processes WorkRequests
type Worker struct {
	ID          int
	Work        chan WorkRequest
	WorkerQueue chan chan WorkRequest
	QuitChan    chan bool
	sigChannel  chan os.Signal
}

// We throw workers into the worker queue (worker pool)
var WorkerQueue chan chan WorkRequest

// A buffered channel that we can send work requests on
var WorkQueue = make(chan WorkRequest, 100)

// NewWorker instantiates a new Worker to join the pool
func NewWorker(id int, workerQueue chan chan WorkRequest) Worker {
	// Create, and return the worker.
	worker := Worker{
		ID:          id,
		Work:        make(chan WorkRequest),
		WorkerQueue: workerQueue,
		QuitChan:    make(chan bool),
		sigChannel:  make(chan os.Signal, 1024),
	}
	return worker
}

func (w *Worker) tallyWIPMetrics(workItem *WorkRequest) {
	workItem.InProcessTS = time.Now()
	workerID := strconv.Itoa(w.ID)
	queueWaitTimeMS := float64(time.Since(workItem.ReceivedTS).Nanoseconds()) / 1000000.0
	log.Logger.Debug("work request in process", zap.String("worker_id", workerID), zap.Float64("queue_time_ms", queueWaitTimeMS))
	botmetrics.StatBotWorkerQueueSaturationGauge.WithLabelValues(workItem.CommandName).Dec()
	botmetrics.StatBotWorkReqInProcessCount.WithLabelValues(workItem.CommandName, workerID).Inc()
	botmetrics.StatBotWorkerQueueDurationGauge.WithLabelValues(workItem.CommandName, workerID).Set(queueWaitTimeMS)
	botmetrics.StatBotWorkerQueueDurationHistogram.WithLabelValues(workItem.CommandName, workerID).Observe(queueWaitTimeMS)
	botmetrics.StatBotWorkerWIPSaturationGauge.WithLabelValues(workItem.CommandName, workerID).Inc()
}

func (w *Worker) tallyCompletedMetrics(workItem *WorkRequest) {
	workItem.CompletedTS = time.Now()
	workerID := strconv.Itoa(w.ID)
	inProcessWorkTimeMS := float64(time.Since(workItem.InProcessTS).Nanoseconds()) / 1000000.0
	log.Logger.Debug("work request complete", zap.String("worker_id", workerID), zap.Float64("wip_time_ms", inProcessWorkTimeMS))
	botmetrics.StatBotWorkerWIPSaturationGauge.WithLabelValues(workItem.CommandName, workerID).Dec()
	botmetrics.StatBotWorkReqCompletedCount.WithLabelValues(workItem.CommandName, workerID).Inc()
	botmetrics.StatBotWorkerInProcessDurationGauge.WithLabelValues(workItem.CommandName, workerID).Set(inProcessWorkTimeMS)
	botmetrics.StatBotWorkerInProcessDurationHistogram.WithLabelValues(workItem.CommandName, workerID).Observe(inProcessWorkTimeMS)
}

// This function "starts" the worker by starting a goroutine, that is
// an infinite "for-select" loop.
func (w *Worker) Start() {
	go func() {
		for {
			// Add ourselves into the worker queue.
			w.WorkerQueue <- w.Work
			select {
			case work := <-w.Work:
				// work request received - the worker starts processing the work request
				w.tallyWIPMetrics(&work)
				// <<<<<<<<==========..........simulated work here.......====================>>>>
				time.Sleep(work.Delay)
				// work request complete - the worker has finished processing the work request goes back into the worker queue (pool)
				w.tallyCompletedMetrics(&work)
			case <-w.QuitChan:
				// we have been asked to stop working...
				log.Logger.Info("received worker stop request", zap.Int("worker_id", w.ID))
				close(w.QuitChan)
				return
			}
		}
	}()
}

func (w *Worker) sigHandler() {
	for {
		sig := <-w.sigChannel
		switch sig {
		case syscall.SIGHUP:
			log.Logger.Warn("worker SIGHUP caught, nothing supports this currently")
		case os.Interrupt, syscall.SIGTERM, syscall.SIGINT:
			log.Logger.Info("caught worker shutdown signal", zap.String("signal", sig.String()))
			w.Stop()
		}
	}
}

// Stop tells the worker to stop listening for work requests.
// Note that the worker will only stop *after* it has finished its work.
func (w *Worker) Stop() {
	go func() {
		w.QuitChan <- true
	}()
}

func StartDispatcher(nworkers int) {
	// Fire up the worker queue with nworkers size
	// use NumCPU() unless supplied by caller
	if nworkers <= 0 {
		nworkers = runtime.NumCPU()
	}

	log.Logger.Info("Starting Queue Dispatcher", zap.Int("workers", nworkers))
	// Initialize the channel we are going to put the workers' work channels into.
	WorkerQueue = make(chan chan WorkRequest, nworkers)

	// Now, create all of our workers.
	for i := 0; i < nworkers; i++ {
		log.Logger.Debug("starting queue worker", zap.Int("id", i+1))
		worker := NewWorker(i+1, WorkerQueue)
		worker.Start()
		signal.Notify(worker.sigChannel, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
		go worker.sigHandler()
	}

	go func() {
		for {
			select {
			case work := <-WorkQueue:
				work.ReceivedTS = time.Now()
				log.Logger.Debug("work request received")
				botmetrics.StatBotWorkReqReceivedCount.WithLabelValues(work.CommandName).Inc()
				botmetrics.StatBotWorkerQueueSaturationGauge.WithLabelValues(work.CommandName).Inc()
				go func() {
					worker := <-WorkerQueue
					worker <- work
				}()
			}
		}
	}()
}
