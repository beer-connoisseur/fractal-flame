package config

import (
	defaultflag "flag"
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	flag "github.com/urfave/cli/v2"
)

func createMockContext(flags map[string]string) *flag.Context {
	flagSet := defaultflag.NewFlagSet("test", defaultflag.ContinueOnError)
	flagSet.Int("width", 800, "")
	flagSet.Int("height", 600, "")
	flagSet.Int("iterations", 10000, "")
	flagSet.String("output", "output.png", "")
	flagSet.Int("threads", 1, "")
	flagSet.Int64("seed", 12345, "")
	flagSet.Bool("gamma-correction", false, "")
	flagSet.Float64("gamma", 2.2, "")
	flagSet.Int("symmetry-level", 1, "")
	flagSet.String("affine-params", "", "")
	flagSet.String("functions", "", "")
	flagSet.String("config", "", "")

	for name, value := range flags {
		err := flagSet.Set(name, value)
		if err != nil {
			slog.Warn(err.Error())
		}
	}

	app := &flag.App{
		Name: "test",
		Flags: []flag.Flag{
			&flag.IntFlag{Name: "width"},
			&flag.IntFlag{Name: "height"},
			&flag.IntFlag{Name: "iterations"},
			&flag.StringFlag{Name: "output"},
			&flag.IntFlag{Name: "threads"},
			&flag.Int64Flag{Name: "seed"},
			&flag.BoolFlag{Name: "gamma-correction"},
			&flag.Float64Flag{Name: "gamma"},
			&flag.IntFlag{Name: "symmetry-level"},
			&flag.StringFlag{Name: "affine-params"},
			&flag.StringFlag{Name: "functions"},
			&flag.StringFlag{Name: "config"},
		},
	}

	ctx := flag.NewContext(app, flagSet, nil)
	return ctx
}

func TestNewConfig(t *testing.T) {
	tests := []struct {
		name    string
		flags   map[string]string
		want    *Config
		wantErr bool
	}{
		{
			name: "valid minimal config",
			flags: map[string]string{
				"width":          "800",
				"height":         "600",
				"iterations":     "10000",
				"threads":        "4",
				"symmetry-level": "2",
				"affine-params":  "0.5,0,0,0,0.5,0/0.5,0,0.5,0,0.5,0",
				"functions":      "linear:1.0,sinusoidal:0.5",
			},
			want: &Config{
				Size:           Size{Width: 800, Height: 600},
				IterationCount: 10000,
				OutputPath:     "output.png",
				Threads:        4,
				Seed:           12345,
				Gamma:          2.2,
				SymmetryLevel:  2,
				AffineParams: []AffineParam{
					{A: 0.5, B: 0, C: 0, D: 0, E: 0.5, F: 0},
					{A: 0.5, B: 0, C: 0.5, D: 0, E: 0.5, F: 0},
				},
				Functions: []FunctionConfig{
					{Name: "linear", Weight: 1.0},
					{Name: "sinusoidal", Weight: 0.5},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid width too low",
			flags: map[string]string{
				"width":  "0",
				"height": "600",
			},
			wantErr: true,
		},
		{
			name: "invalid width too high",
			flags: map[string]string{
				"width":  "8000",
				"height": "600",
			},
			wantErr: true,
		},
		{
			name: "invalid height too low",
			flags: map[string]string{
				"width":  "800",
				"height": "0",
			},
			wantErr: true,
		},
		{
			name: "invalid iterations",
			flags: map[string]string{
				"iterations": "0",
			},
			wantErr: true,
		},
		{
			name: "invalid threads",
			flags: map[string]string{
				"threads": "0",
			},
			wantErr: true,
		},
		{
			name: "invalid symmetry level",
			flags: map[string]string{
				"symmetry-level": "0",
			},
			wantErr: true,
		},
		{
			name: "gamma zero",
			flags: map[string]string{
				"gamma": "0",
			},
			wantErr: true,
		},
		{
			name: "with gamma correction enabled",
			flags: map[string]string{
				"width":            "800",
				"height":           "600",
				"iterations":       "10000",
				"gamma-correction": "true",
				"gamma":            "1.8",
				"affine-params":    "0.5,0,0,0,0.5,0",
				"functions":        "linear:1.0",
			},
			want: &Config{
				Size:            Size{Width: 800, Height: 600},
				IterationCount:  10000,
				OutputPath:      "output.png",
				Threads:         1,
				Seed:            12345,
				GammaCorrection: true,
				Gamma:           1.8,
				SymmetryLevel:   1,
				AffineParams: []AffineParam{
					{A: 0.5, B: 0, C: 0, D: 0, E: 0.5, F: 0},
				},
				Functions: []FunctionConfig{
					{Name: "linear", Weight: 1.0},
				},
			},
			wantErr: false,
		},
		{
			name: "with custom seed",
			flags: map[string]string{
				"width":         "800",
				"height":        "600",
				"iterations":    "10000",
				"seed":          "99999",
				"affine-params": "0.5,0,0,0,0.5,0",
				"functions":     "linear:1.0",
			},
			want: &Config{
				Size:           Size{Width: 800, Height: 600},
				IterationCount: 10000,
				OutputPath:     "output.png",
				Threads:        1,
				Seed:           99999,
				Gamma:          2.2,
				SymmetryLevel:  1,
				AffineParams: []AffineParam{
					{A: 0.5, B: 0, C: 0, D: 0, E: 0.5, F: 0},
				},
				Functions: []FunctionConfig{
					{Name: "linear", Weight: 1.0},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := createMockContext(tt.flags)
			got, err := NewConfig(ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.want != nil {
				require.Equal(t, tt.want.Size, got.Size)
				require.Equal(t, tt.want.IterationCount, got.IterationCount)
				require.Equal(t, tt.want.OutputPath, got.OutputPath)
				require.Equal(t, tt.want.Threads, got.Threads)
				require.Equal(t, tt.want.Seed, got.Seed)
				require.Equal(t, tt.want.GammaCorrection, got.GammaCorrection)
				require.Equal(t, tt.want.Gamma, got.Gamma)
				require.Equal(t, tt.want.SymmetryLevel, got.SymmetryLevel)
				require.Equal(t, len(tt.want.AffineParams), len(got.AffineParams))
				require.Equal(t, len(tt.want.Functions), len(got.Functions))

				for i, param := range tt.want.AffineParams {
					require.InDelta(t, param.A, got.AffineParams[i].A, 0.001)
					require.InDelta(t, param.B, got.AffineParams[i].B, 0.001)
					require.InDelta(t, param.C, got.AffineParams[i].C, 0.001)
					require.InDelta(t, param.D, got.AffineParams[i].D, 0.001)
					require.InDelta(t, param.E, got.AffineParams[i].E, 0.001)
					require.InDelta(t, param.F, got.AffineParams[i].F, 0.001)
				}

				for i, fn := range tt.want.Functions {
					require.Equal(t, fn.Name, got.Functions[i].Name)
					require.InDelta(t, fn.Weight, got.Functions[i].Weight, 0.001)
				}
			}
		})
	}
}

func Test_loadFromJSON(t *testing.T) {
	tests := []struct {
		name       string
		setup      func() string
		wantConfig *Config
		wantErr    bool
	}{
		{
			name: "valid json config",
			setup: func() string {
				tmpFile, err := os.CreateTemp("", "config-*.json")
				require.NoError(t, err)
				_, err = tmpFile.WriteString(`{
					"size": {
            			"width": 1920,
            			"height": 1080
        			},
					"iteration_count": 50000,
					"output_path": "test.png",
					"threads": 8,
					"seed": 42,
					"gamma_correction": true,
					"gamma": 2.0,
					"symmetry_level": 4,
					"affine_params": [
						{"a": 0.7, "b": 0.1, "c": 0.0, "d": -0.1, "e": 0.7, "f": 0.0}
					],
					"functions": [
						{"name": "linear", "weight": 1.0},
						{"name": "spherical", "weight": 0.5}
					]
				}`)
				require.NoError(t, err)
				err = tmpFile.Close()
				if err != nil {
					slog.Warn(err.Error())
				}
				return tmpFile.Name()
			},
			wantConfig: &Config{
				Size:            Size{Width: 1920, Height: 1080},
				IterationCount:  50000,
				OutputPath:      "test.png",
				Threads:         8,
				Seed:            42,
				GammaCorrection: true,
				Gamma:           2.0,
				SymmetryLevel:   4,
				AffineParams: []AffineParam{
					{A: 0.7, B: 0.1, C: 0.0, D: -0.1, E: 0.7, F: 0.0},
				},
				Functions: []FunctionConfig{
					{Name: "linear", Weight: 1.0},
					{Name: "spherical", Weight: 0.5},
				},
			},
			wantErr: false,
		},
		{
			name: "file not found",
			setup: func() string {
				return "non-existent-file.json"
			},
			wantConfig: nil,
			wantErr:    true,
		},
		{
			name: "invalid json",
			setup: func() string {
				tmpFile, err := os.CreateTemp("", "config-*.json")
				require.NoError(t, err)
				_, err = tmpFile.WriteString(`{
					"size": {
            			"width": wtf,
            			"height": 1080
        			},
				}`)
				require.NoError(t, err)
				err = tmpFile.Close()
				if err != nil {
					slog.Warn(err.Error())
				}
				return tmpFile.Name()
			},
			wantConfig: nil,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filename := tt.setup()
			defer func(name string) {
				err := os.Remove(name)
				if err != nil {
					slog.Warn(err.Error())
				}
			}(filename)

			config := &Config{}
			err := loadFromJSON(filename, config)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.wantConfig.Size, config.Size)
			require.Equal(t, tt.wantConfig.IterationCount, config.IterationCount)
			require.Equal(t, tt.wantConfig.OutputPath, config.OutputPath)
			require.Equal(t, tt.wantConfig.Threads, config.Threads)
			require.Equal(t, tt.wantConfig.Seed, config.Seed)
			require.Equal(t, tt.wantConfig.GammaCorrection, config.GammaCorrection)
			require.Equal(t, tt.wantConfig.Gamma, config.Gamma)
			require.Equal(t, tt.wantConfig.SymmetryLevel, config.SymmetryLevel)

			if tt.wantConfig.AffineParams != nil {
				require.Equal(t, len(tt.wantConfig.AffineParams), len(config.AffineParams))
				for i := range tt.wantConfig.AffineParams {
					require.InDelta(t, tt.wantConfig.AffineParams[i].A, config.AffineParams[i].A, 0.001)
					require.InDelta(t, tt.wantConfig.AffineParams[i].B, config.AffineParams[i].B, 0.001)
					require.InDelta(t, tt.wantConfig.AffineParams[i].C, config.AffineParams[i].C, 0.001)
					require.InDelta(t, tt.wantConfig.AffineParams[i].D, config.AffineParams[i].D, 0.001)
					require.InDelta(t, tt.wantConfig.AffineParams[i].E, config.AffineParams[i].E, 0.001)
					require.InDelta(t, tt.wantConfig.AffineParams[i].F, config.AffineParams[i].F, 0.001)
				}
			}

			if tt.wantConfig.Functions != nil {
				require.Equal(t, len(tt.wantConfig.Functions), len(config.Functions))
				for i := range tt.wantConfig.Functions {
					require.Equal(t, tt.wantConfig.Functions[i].Name, config.Functions[i].Name)
					require.InDelta(t, tt.wantConfig.Functions[i].Weight, config.Functions[i].Weight, 0.001)
				}
			}
		})
	}
}

func Test_parseAffineParam(t *testing.T) {
	tests := []struct {
		name    string
		coeffs  []string
		want    AffineParam
		wantErr bool
	}{
		{
			name:   "valid coefficients",
			coeffs: []string{"0.5", "0.1", "0.0", "-0.2", "0.6", "0.1"},
			want:   AffineParam{A: 0.5, B: 0.1, C: 0.0, D: -0.2, E: 0.6, F: 0.1},
		},
		{
			name:   "valid with spaces",
			coeffs: []string{" 0.5 ", " 0.1 ", " 0.0 ", " -0.2 ", " 0.6 ", " 0.1 "},
			want:   AffineParam{A: 0.5, B: 0.1, C: 0.0, D: -0.2, E: 0.6, F: 0.1},
		},
		{
			name:    "invalid number",
			coeffs:  []string{"0.5", "invalid", "0.0", "-0.2", "0.6", "0.1"},
			wantErr: true,
		},
		{
			name:   "negative coefficients",
			coeffs: []string{"-1.5", "2.0", "-3.0", "4.5", "-5.0", "6.5"},
			want:   AffineParam{A: -1.5, B: 2.0, C: -3.0, D: 4.5, E: -5.0, F: 6.5},
		},
		{
			name:   "scientific notation",
			coeffs: []string{"1.5e-1", "2.0E2", "3.0e+0", "4.5E-2", "5.0", "6.5"},
			want:   AffineParam{A: 0.15, B: 200.0, C: 3.0, D: 0.045, E: 5.0, F: 6.5},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseAffineParam(tt.coeffs)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseAffineParam() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.InDelta(t, tt.want.A, got.A, 0.0001)
			require.InDelta(t, tt.want.B, got.B, 0.0001)
			require.InDelta(t, tt.want.C, got.C, 0.0001)
			require.InDelta(t, tt.want.D, got.D, 0.0001)
			require.InDelta(t, tt.want.E, got.E, 0.0001)
			require.InDelta(t, tt.want.F, got.F, 0.0001)
		})
	}
}

func Test_parseAffineParams(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []AffineParam
		wantErr bool
	}{
		{
			name:  "single affine param",
			input: "0.5,0.1,0.0,-0.2,0.6,0.1",
			want: []AffineParam{
				{A: 0.5, B: 0.1, C: 0.0, D: -0.2, E: 0.6, F: 0.1},
			},
		},
		{
			name:  "multiple affine params",
			input: "0.5,0,0,0,0.5,0/0.5,0,0.5,0,0.5,0/0.5,0,0.25,0,0.5,0.5",
			want: []AffineParam{
				{A: 0.5, B: 0, C: 0, D: 0, E: 0.5, F: 0},
				{A: 0.5, B: 0, C: 0.5, D: 0, E: 0.5, F: 0},
				{A: 0.5, B: 0, C: 0.25, D: 0, E: 0.5, F: 0.5},
			},
		},
		{
			name:  "multiple with spaces",
			input: " 0.5 , 0 , 0 , 0 , 0.5 , 0 / 0.5 , 0 , 0.5 , 0 , 0.5 , 0 ",
			want: []AffineParam{
				{A: 0.5, B: 0, C: 0, D: 0, E: 0.5, F: 0},
				{A: 0.5, B: 0, C: 0.5, D: 0, E: 0.5, F: 0},
			},
		},
		{
			name:    "wrong separator",
			input:   "0.5,0,0,0,0.5,0;0.5,0,0.5,0,0.5,0",
			wantErr: true,
		},
		{
			name:    "invalid coefficients",
			input:   "0.5,invalid,0,0,0.5,0",
			wantErr: true,
		},
		{
			name:    "too few coefficients in one param",
			input:   "0.5,0,0,0,0.5/0.5,0,0.5,0,0.5,0",
			wantErr: true,
		},
		{
			name:    "empty string",
			input:   "",
			wantErr: true,
		},
		{
			name:    "just slashes",
			input:   "//",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseAffineParams(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseAffineParams() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.Equal(t, len(tt.want), len(got))
			for i := range tt.want {
				require.InDelta(t, tt.want[i].A, got[i].A, 0.0001)
				require.InDelta(t, tt.want[i].B, got[i].B, 0.0001)
				require.InDelta(t, tt.want[i].C, got[i].C, 0.0001)
				require.InDelta(t, tt.want[i].D, got[i].D, 0.0001)
				require.InDelta(t, tt.want[i].E, got[i].E, 0.0001)
				require.InDelta(t, tt.want[i].F, got[i].F, 0.0001)
			}
		})
	}
}

func Test_parseFunctions(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []FunctionConfig
		wantErr bool
	}{
		{
			name:  "single function",
			input: "linear:1.0",
			want: []FunctionConfig{
				{Name: "linear", Weight: 1.0},
			},
		},
		{
			name:  "multiple functions",
			input: "linear:1.0,sinusoidal:0.5,spherical:0.3",
			want: []FunctionConfig{
				{Name: "linear", Weight: 1.0},
				{Name: "sinusoidal", Weight: 0.5},
				{Name: "spherical", Weight: 0.3},
			},
		},
		{
			name:  "with spaces",
			input: " linear : 1.0 , sinusoidal : 0.5 ",
			want: []FunctionConfig{
				{Name: "linear", Weight: 1.0},
				{Name: "sinusoidal", Weight: 0.5},
			},
		},
		{
			name:    "invalid format - missing colon",
			input:   "linear1.0,sinusoidal0.5",
			wantErr: true,
		},
		{
			name:    "invalid weight",
			input:   "linear:invalid,sinusoidal:0.5",
			wantErr: true,
		},
		{
			name:    "empty string",
			input:   "",
			wantErr: true,
		},
		{
			name:    "just commas",
			input:   ",,",
			wantErr: true,
		},
		{
			name:  "zero weight",
			input: "linear:0.0,sinusoidal:1.0",
			want: []FunctionConfig{
				{Name: "linear", Weight: 0.0},
				{Name: "sinusoidal", Weight: 1.0},
			},
		},
		{
			name:  "negative weight",
			input: "linear:-1.0,sinusoidal:2.0",
			want: []FunctionConfig{
				{Name: "linear", Weight: -1.0},
				{Name: "sinusoidal", Weight: 2.0},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseFunctions(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseFunctions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.Equal(t, len(tt.want), len(got))
			for i := range tt.want {
				require.Equal(t, tt.want[i].Name, got[i].Name)
				require.InDelta(t, tt.want[i].Weight, got[i].Weight, 0.0001)
			}
		})
	}
}
