package cli

import (
	flag "github.com/urfave/cli/v2"
	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw4-fractal-flame/internal/application"
	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw4-fractal-flame/internal/config"
	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw4-fractal-flame/internal/infrastructure/histogram"
	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw4-fractal-flame/internal/infrastructure/renderer"
	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw4-fractal-flame/internal/infrastructure/transformations"
)

type App struct {
	engine *application.FractalFlame
}

func NewApp() *flag.App {
	histogramFactory := histogram.NewHistogramFactory()
	transformationFactory := transformations.NewTransformationFactory()
	colorMapper := renderer.NewColorMapper()
	render := renderer.NewPNGRenderer()
	symmetryApplier := renderer.NewSymmetryApplier()

	app := App{
		application.NewFractalFlame(
			histogramFactory,
			transformationFactory,
			colorMapper,
			render,
			symmetryApplier,
		),
	}

	return &flag.App{
		Name:     "fractal-flame",
		Usage:    "generate an image of a structural change based on the Chaos game",
		HideHelp: true,
		Flags:    app.flags(),
		Action:   app.run,
		OnUsageError: func(_ *flag.Context, err error, _ bool) error {
			return application.NewErrUser(err.Error())
		},
	}
}

func (a *App) flags() []flag.Flag {
	return []flag.Flag{
		&flag.StringFlag{
			Name:    "config",
			Aliases: []string{"c"},
			Usage:   "Path to JSON configuration file",
		},
		&flag.IntFlag{
			Name:    "width",
			Aliases: []string{"w"},
			Usage:   "Image width",
			Value:   1920,
		},
		&flag.IntFlag{
			Name:    "height",
			Aliases: []string{"h"},
			Usage:   "Image height",
			Value:   1080,
		},
		&flag.IntFlag{
			Name:    "iterations",
			Aliases: []string{"i"},
			Usage:   "Number of iterations",
			Value:   2500,
		},
		&flag.StringFlag{
			Name:    "output",
			Aliases: []string{"o"},
			Usage:   "Output image path",
			Value:   "result.png",
		},
		&flag.IntFlag{
			Name:    "threads",
			Aliases: []string{"t"},
			Usage:   "Number of threads",
			Value:   1,
		},
		&flag.Int64Flag{
			Name:  "seed",
			Usage: "Random generator seed value",
			Value: 5,
		},
		&flag.StringFlag{
			Name:    "affine-params",
			Aliases: []string{"ap"},
			Usage:   "Affine transformation parameters",
			Value:   "0.5,0.2,0.0,0.0,0.5,0.0/0.1,0.04,0.1,-0.04,0.85,0.01",
		},
		&flag.StringFlag{
			Name:    "functions",
			Aliases: []string{"f"},
			Usage:   "Transformation functions",
			Value:   "heart:0.3,swirl:0.5,spherical:0.8",
		},
		&flag.BoolFlag{
			Name:    "gamma-correction",
			Aliases: []string{"g"},
			Usage:   "Enable gamma correction",
		},
		&flag.Float64Flag{
			Name:  "gamma",
			Usage: "Gamma value for correction",
			Value: 2.2,
		},
		&flag.IntFlag{
			Name:    "symmetry-level",
			Aliases: []string{"s"},
			Usage:   "Symmetry level",
			Value:   1,
		},
	}
}

func (a *App) run(ctx *flag.Context) error {
	cfg, err := config.NewConfig(ctx)
	if err != nil {
		return application.NewErrUser(err.Error())
	}

	err = a.engine.Generate(cfg)
	if err != nil {
		return err
	}

	return nil
}
