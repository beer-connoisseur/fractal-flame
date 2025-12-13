package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	flag "github.com/urfave/cli/v2"
	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw4-fractal-flame/internal/domain"
)

type Size struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

type AffineParam struct {
	A     float64 `json:"a"`
	B     float64 `json:"b"`
	C     float64 `json:"c"`
	D     float64 `json:"d"`
	E     float64 `json:"e"`
	F     float64 `json:"f"`
	Color domain.Color
}

type FunctionConfig struct {
	Name   string  `json:"name"`
	Weight float64 `json:"weight"`
}

type Config struct {
	Size            Size             `json:"size"`
	IterationCount  int              `json:"iteration_count"`
	OutputPath      string           `json:"output_path"`
	Threads         int              `json:"threads"`
	Seed            int64            `json:"seed"`
	GammaCorrection bool             `json:"gamma_correction"`
	Gamma           float64          `json:"gamma"`
	SymmetryLevel   int              `json:"symmetry_level"`
	Functions       []FunctionConfig `json:"functions"`
	AffineParams    []AffineParam    `json:"affine_params"`
}

func NewConfig(ctx *flag.Context) (*Config, error) {
	cfg := &Config{}

	// set default values
	cfg.Size.Width = ctx.Int("width")
	cfg.Size.Height = ctx.Int("height")
	cfg.IterationCount = ctx.Int("iterations")
	cfg.OutputPath = ctx.String("output")
	cfg.Threads = ctx.Int("threads")
	cfg.Seed = ctx.Int64("seed")
	cfg.GammaCorrection = ctx.Bool("gamma-correction")
	cfg.Gamma = ctx.Float64("gamma")
	cfg.SymmetryLevel = ctx.Int("symmetry-level")

	affineParams := ctx.String("affine-params")
	params, err := parseAffineParams(affineParams)
	if err != nil {
		return nil, err
	}
	cfg.AffineParams = params

	functions := ctx.String("functions")
	funcs, err := parseFunctions(functions)
	if err != nil {
		return nil, err
	}
	cfg.Functions = funcs

	if configFile := ctx.String("config"); configFile != "" {
		err = loadFromJSON(configFile, cfg)
		if err != nil {
			return nil, fmt.Errorf("error loading JSON: %w", err)
		}
	}

	if ctx.IsSet("width") {
		cfg.Size.Width = ctx.Int("width")
	}
	if cfg.Size.Width <= 0 || cfg.Size.Width > 7680 {
		return nil, fmt.Errorf("width must be between 1 and 7680, got %d", cfg.Size.Width)
	}

	if ctx.IsSet("height") {
		cfg.Size.Height = ctx.Int("height")
	}
	if cfg.Size.Height <= 0 || cfg.Size.Height > 4320 {
		return nil, fmt.Errorf("height must be between 1 and 4320, got %d", cfg.Size.Height)
	}

	if ctx.IsSet("iterations") {
		cfg.IterationCount = ctx.Int("iterations")
	}
	if cfg.IterationCount <= 0 {
		return nil, fmt.Errorf("iterations must be > 0, got %d", cfg.IterationCount)
	}

	if ctx.IsSet("output") {
		cfg.OutputPath = ctx.String("output")
	}

	if ctx.IsSet("threads") {
		cfg.Threads = ctx.Int("threads")
	}
	if cfg.Threads <= 0 {
		return nil, fmt.Errorf("threads must be > 0, got %d", cfg.Threads)
	}

	if ctx.IsSet("seed") {
		cfg.Seed = ctx.Int64("seed")
	}

	if ctx.IsSet("gamma-correction") {
		cfg.GammaCorrection = ctx.Bool("gamma-correction")
	}

	if ctx.IsSet("gamma") {
		cfg.Gamma = ctx.Float64("gamma")
	}
	if cfg.Gamma == 0 {
		return nil, errors.New("gamma shouldn't be zero")
	}

	if ctx.IsSet("symmetry-level") {
		cfg.SymmetryLevel = ctx.Int("symmetry-level")
	}
	if cfg.SymmetryLevel <= 0 {
		return nil, fmt.Errorf("symmetry-level must be > 0, got %d", cfg.SymmetryLevel)
	}

	if ctx.IsSet("affine-params") {
		affineParams = ctx.String("affine-params")
		params, err = parseAffineParams(affineParams)
		if err != nil {
			return nil, err
		}
		cfg.AffineParams = params
	}

	if ctx.IsSet("functions") {
		functions = ctx.String("functions")
		funcs, err = parseFunctions(functions)
		if err != nil {
			return nil, err
		}
		cfg.Functions = funcs
	}

	return cfg, nil
}

func loadFromJSON(filename string, config *Config) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(data, config); err != nil {
		return err
	}

	return nil
}

func parseAffineParams(input string) ([]AffineParam, error) {
	parts := strings.Split(input, "/")
	params := make([]AffineParam, 0, len(parts))

	for _, part := range parts {
		coeffs := strings.Split(part, ",")
		if len(coeffs) != 6 {
			return nil, fmt.Errorf("incorrect format of affine params: %s", part)
		}

		p, err := parseAffineParam(coeffs)
		if err != nil {
			return nil, fmt.Errorf("error while parsing affine params: %w", err)
		}

		params = append(params, p)
	}

	return params, nil
}

func parseAffineParam(coeffs []string) (AffineParam, error) {
	var p AffineParam
	var err error

	parseFloat := func(str string) (float64, error) {
		return strconv.ParseFloat(strings.TrimSpace(str), 64)
	}

	if p.A, err = parseFloat(coeffs[0]); err != nil {
		return p, err
	}
	if p.B, err = parseFloat(coeffs[1]); err != nil {
		return p, err
	}
	if p.C, err = parseFloat(coeffs[2]); err != nil {
		return p, err
	}
	if p.D, err = parseFloat(coeffs[3]); err != nil {
		return p, err
	}
	if p.E, err = parseFloat(coeffs[4]); err != nil {
		return p, err
	}
	if p.F, err = parseFloat(coeffs[5]); err != nil {
		return p, err
	}

	return p, nil
}

func parseFunctions(input string) ([]FunctionConfig, error) {
	parts := strings.Split(input, ",")
	functions := make([]FunctionConfig, 0, len(parts))

	for _, part := range parts {
		subParts := strings.Split(part, ":")
		if len(subParts) != 2 {
			return nil, fmt.Errorf("incorrect format of functions: %s", part)
		}

		name := strings.TrimSpace(subParts[0])
		weight, err := strconv.ParseFloat(strings.TrimSpace(subParts[1]), 64)
		if err != nil {
			return nil, fmt.Errorf("error while parsing function weight: %w", err)
		}

		functions = append(functions, FunctionConfig{
			Name:   name,
			Weight: weight,
		})
	}

	return functions, nil
}
