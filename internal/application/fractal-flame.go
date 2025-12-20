package application

import (
	"log/slog"
	"math/rand"
	"sync"

	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw4-fractal-flame/internal/config"
	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw4-fractal-flame/internal/domain"
	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw4-fractal-flame/internal/infrastructure/utils"
)

type FractalFlame struct {
	histogramFactory      domain.HistogramFactory
	transformationFactory domain.TransformationFactory
	colorMapper           domain.ColorMapper
	renderer              domain.Renderer
	symmetryApplier       domain.SymmetryApplier
}

func NewFractalFlame(
	histogramFactory domain.HistogramFactory,
	transformationFactory domain.TransformationFactory,
	colorMapper domain.ColorMapper,
	renderer domain.Renderer,
	symmetryApplier domain.SymmetryApplier,
) *FractalFlame {
	return &FractalFlame{
		histogramFactory:      histogramFactory,
		transformationFactory: transformationFactory,
		colorMapper:           colorMapper,
		renderer:              renderer,
		symmetryApplier:       symmetryApplier,
	}
}

func (f *FractalFlame) Generate(cfg *config.Config) error {
	colorPalette := f.colorMapper.GeneratePalette(len(cfg.AffineParams), cfg.Seed)
	for i := range cfg.AffineParams {
		cfg.AffineParams[i].Color = colorPalette[i]
	}

	transformations := make([]domain.Transformation, len(cfg.Functions))
	for i, function := range cfg.Functions {
		transformation, err := f.transformationFactory.Create(function.Name, function.Weight)
		if err != nil {
			return NewErrUser(err.Error())
		}
		transformations[i] = transformation
	}

	slog.Info("start of generation")
	histogram, err := f.generateMultiThread(cfg, transformations)
	if err != nil {
		return NewErrGenerating(err.Error())
	}

	if cfg.GammaCorrection {
		if err = histogram.ApplyCorrection(cfg.Gamma); err != nil {
			return NewErrGenerating(err.Error())
		}
	}

	img := f.renderer.Render(histogram)
	err = f.renderer.Save(img, cfg.OutputPath)
	if err != nil {
		return NewErrFile(err.Error())
	}

	slog.Info("finish of generation")

	return nil
}

func (f *FractalFlame) generateMultiThread(
	cfg *config.Config,
	transformations []domain.Transformation,
) (domain.Histogram, error) {
	var wg sync.WaitGroup
	var mx sync.Mutex

	threadHists := make([]domain.Histogram, cfg.Threads)
	for i := range threadHists {
		threadHists[i] = f.histogramFactory.Create(cfg.Size.Width, cfg.Size.Height)
	}

	iterationsPerThread := cfg.IterationCount / cfg.Threads
	remainder := cfg.IterationCount % cfg.Threads
	completedIterations := 0
	lastLogPercent := -1

	for t := 0; t < cfg.Threads; t++ {
		wg.Add(1)
		threadIterations := iterationsPerThread
		if t < remainder {
			threadIterations++
		}
		go func(threadID int) {
			defer wg.Done()
			r := rand.New(rand.NewSource(cfg.Seed + int64(threadID)))

			x, y := r.Float64(), r.Float64()
			totalWeight := utils.CalculateTotalWeight(transformations)

			for i := -domain.SkipSteps; i < threadIterations; i++ {
				if i%domain.LogFrequency == 0 && i > 0 {
					mx.Lock()
					completedIterations += domain.LogFrequency
					progress := int(float64(completedIterations) / float64(cfg.IterationCount) * 100)

					if progress != lastLogPercent {
						slog.Info("generating image", "progress", progress)
						lastLogPercent = progress
					}
					mx.Unlock()
				}

				affineIdx := r.Intn(len(cfg.AffineParams))
				affineParam := cfg.AffineParams[affineIdx]

				xNew := affineParam.A*x + affineParam.B*y + affineParam.C
				yNew := affineParam.D*x + affineParam.E*y + affineParam.F

				transformation := utils.SelectTransformation(r, transformations, totalWeight)

				xNew, yNew = transformation.Apply(xNew, yNew)

				if i >= 0 {
					points := f.symmetryApplier.Apply(xNew, yNew, cfg.SymmetryLevel)
					for _, p := range points {
						threadHists[threadID].Add(p.X, p.Y, affineParam.Color)
					}
				}

				x, y = xNew, yNew
			}
		}(t)
	}

	wg.Wait()

	histogram := threadHists[0]
	for i := 1; i < cfg.Threads; i++ {
		err := histogram.Merge(threadHists[i])
		if err != nil {
			return nil, err
		}
	}

	return histogram, nil
}
