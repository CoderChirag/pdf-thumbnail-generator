package thumbnail_usecase

import (
	"context"
	"fmt"
	"os"

	"github.com/coderchirag/pdf-thumbnail-generator/entities/file_storage"
	"github.com/coderchirag/pdf-thumbnail-generator/entities/thumbnail"
)

var thumbnailPipeline *ThumbnailPipeline

type downloadPdfInput struct {
	ctx    context.Context
	result chan *pipelineOutput
	pdfUrl string
}

type generateThumbnailInput struct {
	ctx     context.Context
	result  chan *pipelineOutput
	pdfPath string
}

type pipelineOutput struct {
	err           error
	thumbnailPath string
}

type ThumbnailPipeline struct {
	downloadPdfInput       chan *downloadPdfInput
	generateThumbnailInput chan *generateThumbnailInput
	shutdown               chan bool
	workers                int
}

func NewThumbnailPipeline(workers int) *ThumbnailPipeline {
	if thumbnailPipeline != nil {
		return thumbnailPipeline
	}
	p := &ThumbnailPipeline{
		downloadPdfInput:       make(chan *downloadPdfInput),
		generateThumbnailInput: make(chan *generateThumbnailInput),
		workers:                workers,
		shutdown:               make(chan bool),
	}
	thumbnailPipeline = p
	p.start()
	return p
}

func GetThumbnailPipeline() *ThumbnailPipeline {
	return thumbnailPipeline
}

func (p *ThumbnailPipeline) Process(ctx context.Context, pdfUrl string) (string, error) {
	resultChan := make(chan *pipelineOutput)
	defer close(resultChan)

	select {
	case p.downloadPdfInput <- &downloadPdfInput{
		ctx:    ctx,
		result: resultChan,
		pdfUrl: pdfUrl,
	}:
	case <-ctx.Done():
		return "", ctx.Err()
	}

	select {
	case output := <-resultChan:
		return output.thumbnailPath, output.err
	case <-ctx.Done():
		return "", ctx.Err()
	}
}

func (p *ThumbnailPipeline) Shutdown() {
	close(p.downloadPdfInput)
	close(p.generateThumbnailInput)
	close(p.shutdown)
	fmt.Println("Thumbnail pipeline closed")
}

func (p *ThumbnailPipeline) start() {
	for i := 0; i < p.workers; i++ {
		go p.worker(i)
	}
	fmt.Println("Thumbnail pipeline started")
}

func (p *ThumbnailPipeline) worker(id int) {
	fmt.Printf("Worker %d started\n", id)
	go p.downloadPdfWorker(id)
	go p.generateThumbnailWorker(id)
	<-p.shutdown
	fmt.Printf("Worker %d closed\n", id)
}

func (p *ThumbnailPipeline) downloadPdfWorker(id int) {
	fmt.Printf("Download PDF Worker %d started\n", id)
	for input := range p.downloadPdfInput {
		pdfPath, err := file_storage.DownloadFileToTempDir(
			input.ctx,
			input.pdfUrl,
			baseDir,
			pdfExt,
		)
		if err != nil {
			input.result <- &pipelineOutput{
				err: err,
			}
			continue
		}
		p.generateThumbnailInput <- &generateThumbnailInput{
			ctx:     input.ctx,
			result:  input.result,
			pdfPath: pdfPath,
		}
	}
	fmt.Printf("Download PDF Worker %d closed\n", id)
}

func (p *ThumbnailPipeline) generateThumbnailWorker(id int) {
	fmt.Printf("Generate Thumbnail Worker %d started\n", id)
	for input := range p.generateThumbnailInput {
		thumbnailPath, err := thumbnail.GenerateThumbnailFromPdf(
			input.ctx,
			input.pdfPath,
			quality,
		)
		if err != nil {
			input.result <- &pipelineOutput{
				err: err,
			}
			continue
		}
		input.result <- &pipelineOutput{
			thumbnailPath: thumbnailPath,
		}
	}
	fmt.Printf("Generate Thumbnail Worker %d closed\n", id)
}

func GenerateThumbnail(
	ctx context.Context,
	pdfUrl string,
) (string, error) {
	pdfPath, err := file_storage.DownloadFileToTempDir(ctx, pdfUrl, baseDir, pdfExt)
	if err != nil {
		return "", err
	}
	defer os.Remove(pdfPath)

	thumbnailPath, err := thumbnail.GenerateThumbnailFromPdf(ctx, pdfPath, quality)
	if err != nil {
		return "", err
	}

	return thumbnailPath, nil
}
