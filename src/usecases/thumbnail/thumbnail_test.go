package thumbnail_usecase_test

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/shirou/gopsutil/process"
	"golang.org/x/sync/errgroup"

	thumbnail_usecase "github.com/coderchirag/pdf-thumbnail-generator/usecases/thumbnail"
)

type monitoring struct {
	done        chan bool
	cpuPercent  float64
	memoryBytes uint64
}

func NewMonitoring() *monitoring {
	return &monitoring{
		done: make(chan bool),
	}
}

func (m *monitoring) StartMonitoring(b *testing.B) {
	b.Helper()
	pid := os.Getpid()
	if pid < 0 || pid > int(^uint32(0)>>1) {
		b.Fatal("PID out of int32 range")
	}
	proc, err := process.NewProcess(int32(pid))
	if err != nil {
		b.Fatal(err)
	}

	// Start monitoring
	go func() {
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-m.done:
				return
			case <-ticker.C:
				if cpuPercent, err := proc.CPUPercent(); err == nil {
					if cpuPercent > m.cpuPercent {
						m.cpuPercent = cpuPercent
					}
				}
				if memInfo, err := proc.MemoryInfo(); err == nil {
					if memInfo.RSS > m.memoryBytes {
						m.memoryBytes = memInfo.RSS
					}
				}
			}
		}
	}()
}

func (m *monitoring) StopMonitoring() {
	m.done <- true
}

func (m *monitoring) ReportMetrics(b *testing.B) {
	b.Helper()
	b.ReportMetric(m.cpuPercent*10, "PeakCPU(millicores)")
	b.ReportMetric(float64(m.memoryBytes)/1024/1024, "PeakMemory(MB)")
	b.ReportMetric(float64(b.Elapsed().Seconds())/float64(b.N), "s/op")
}

func BenchmarkGenerateThumbnailSequentially(b *testing.B) {
	testPdfUrl := "https://drive.usercontent.google.com/download?id=1Qspoh1gKWl4KS9MPj0LCQE3hKdKZInmb&export=download&authuser=0"
	ctx := context.Background()

	monitoring := NewMonitoring()
	monitoring.StartMonitoring(b)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StartTimer()
		thumbnailPath, err := thumbnail_usecase.GenerateThumbnail(ctx, testPdfUrl)
		if err != nil {
			b.Fatal(err)
		}
		b.StopTimer()

		if _, err := os.Stat(thumbnailPath); os.IsNotExist(err) {
			b.Fatal("thumbnail was not created")
		}
		_ = os.Remove(thumbnailPath)
	}

	monitoring.StopMonitoring()
	monitoring.ReportMetrics(b)
}

func BenchmarkGenerateThumbnailConcurrently(b *testing.B) {
	testPdfUrl := "https://drive.usercontent.google.com/download?id=1Qspoh1gKWl4KS9MPj0LCQE3hKdKZInmb&export=download&authuser=0"
	ctx := context.Background()

	monitoring := NewMonitoring()
	monitoring.StartMonitoring(b)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g := errgroup.Group{}
		b.StartTimer()
		for j := 0; j < 10; j++ {
			g.Go(func() error {
				thumbnailPath, err := thumbnail_usecase.GenerateThumbnail(
					ctx,
					testPdfUrl,
				)
				if err != nil {
					return err
				}
				if _, err := os.Stat(thumbnailPath); os.IsNotExist(err) {
					return errors.New("thumbnail was not created")
				}
				_ = os.Remove(thumbnailPath)
				return nil
			})
		}
		if err := g.Wait(); err != nil {
			b.Fatal(err)
		}
		b.StopTimer()
	}

	monitoring.StopMonitoring()
	monitoring.ReportMetrics(b)
	b.ReportMetric(
		(float64(float64(b.Elapsed().Seconds())/float64(b.N)) / float64(10)),
		"s/op(avg)",
	)
}
