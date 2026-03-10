package histogram

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw4-fractal-flame/internal/domain"
)

func TestHistFactory_Create(t *testing.T) {
	t.Run("create with factory", func(t *testing.T) {
		f := &HistFactory{}
		hist := f.Create(800, 600)

		require.NotNil(t, hist)

		width, height := hist.GetSize()
		require.Equal(t, 800, width)
		require.Equal(t, 600, height)
	})
}

func TestHistField_Add(t *testing.T) {
	tests := []struct {
		name       string
		setup      func() *HistField
		x          float64
		y          float64
		color      domain.Color
		expectedX  int
		expectedY  int
		wantAdd    bool
		wantMaxHit int
	}{
		{
			name: "add point in center",
			setup: func() *HistField {
				return NewHistogram(800, 600)
			},
			x:          0.0,
			y:          0.0,
			color:      domain.Color{R: 1.0, G: 0.5, B: 0.25},
			expectedX:  400,
			expectedY:  300,
			wantAdd:    true,
			wantMaxHit: 1,
		},
		{
			name: "add point at top left",
			setup: func() *HistField {
				return NewHistogram(800, 600)
			},
			x:          domain.XMIN,
			y:          domain.YMAX,
			color:      domain.Color{R: 0.0, G: 1.0, B: 0.0},
			wantAdd:    false,
			wantMaxHit: 0,
		},
		{
			name: "add point outside bounds",
			setup: func() *HistField {
				return NewHistogram(800, 600)
			},
			x:          domain.XMAX + 1.0,
			y:          domain.YMIN - 1.0,
			color:      domain.Color{R: 1.0, G: 1.0, B: 1.0},
			wantAdd:    false,
			wantMaxHit: 0,
		},
		{
			name: "add multiple points to same pixel",
			setup: func() *HistField {
				hist := NewHistogram(800, 600)
				hist.Add(0.0, 0.0, domain.Color{R: 1.0, G: 0.0, B: 0.0})
				return hist
			},
			x:          0.0,
			y:          0.0,
			color:      domain.Color{R: 0.0, G: 1.0, B: 0.0},
			expectedX:  400,
			expectedY:  300,
			wantAdd:    true,
			wantMaxHit: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hist := tt.setup()
			initialMaxHits := hist.maxHits

			hist.Add(tt.x, tt.y, tt.color)

			if tt.wantAdd {
				pixel := &hist.pixels[tt.expectedY][tt.expectedX]
				require.Greater(t, pixel.Counter, 0)
				require.Equal(t, tt.wantMaxHit, pixel.Counter)
				require.Equal(t, tt.wantMaxHit, hist.maxHits)

				if initialMaxHits > 0 {
					require.InDelta(t, 0.5, pixel.Color.R, 0.01)
					require.InDelta(t, 0.5, pixel.Color.G, 0.01)
					require.InDelta(t, 0.0, pixel.Color.B, 0.01)
				}
			} else {
				require.Equal(t, initialMaxHits, hist.maxHits)
			}
		})
	}
}

func TestHistField_ApplyCorrection(t *testing.T) {
	tests := []struct {
		name        string
		setup       func() *HistField
		gamma       float64
		wantErr     bool
		checkResult func(t *testing.T, hist *HistField)
	}{
		{
			name: "apply gamma correction to single pixel",
			setup: func() *HistField {
				hist := NewHistogram(800, 600)
				hist.Add(0.0, 0.0, domain.Color{R: 1.0, G: 0.5, B: 0.25})
				return hist
			},
			gamma:   2.2,
			wantErr: false,
			checkResult: func(t *testing.T, hist *HistField) {
				color := hist.GetColor(400, 300)
				require.InDelta(t, 1.0, color.R, 0.001)
				require.InDelta(t, 0.5, color.G, 0.001)
				require.InDelta(t, 0.25, color.B, 0.001)
			},
		},
		{
			name: "apply gamma correction to multiple pixels",
			setup: func() *HistField {
				hist := NewHistogram(800, 600)
				hist.Add(-1.0, -1.0, domain.Color{R: 1.0, G: 0.0, B: 0.0})
				hist.Add(-1.0, -1.0, domain.Color{R: 1.0, G: 0.0, B: 0.0})
				hist.Add(1.0, 1.0, domain.Color{R: 0.0, G: 1.0, B: 0.0})
				return hist
			},
			gamma:   1.8,
			wantErr: false,
			checkResult: func(t *testing.T, hist *HistField) {
				color1 := hist.GetColor(
					hist.width-int(((domain.XMAX+1.0)/(domain.XMAX-domain.XMIN))*float64(hist.width)),
					hist.height-int(((domain.YMAX+1.0)/(domain.YMAX-domain.YMIN))*float64(hist.height)),
				)
				color2 := hist.GetColor(
					hist.width-int(((domain.XMAX-1.0)/(domain.XMAX-domain.XMIN))*float64(hist.width)),
					hist.height-int(((domain.YMAX-1.0)/(domain.YMAX-domain.YMIN))*float64(hist.height)),
				)

				require.Greater(t, color1.R, 0.0)
				require.InDelta(t, 0.0, color2.G, 0.001)
			},
		},
		{
			name: "empty histogram",
			setup: func() *HistField {
				return NewHistogram(800, 600)
			},
			gamma:   2.2,
			wantErr: true,
			checkResult: func(t *testing.T, hist *HistField) {
				require.Equal(t, 0, hist.maxHits)
			},
		},
		{
			name: "gamma 1.0",
			setup: func() *HistField {
				hist := NewHistogram(800, 600)
				hist.Add(0.0, 0.0, domain.Color{R: 0.5, G: 0.5, B: 0.5})
				return hist
			},
			gamma:   1.0,
			wantErr: false,
			checkResult: func(t *testing.T, hist *HistField) {
				color := hist.GetColor(400, 300)
				require.InDelta(t, 0.5, color.R, 0.001)
				require.InDelta(t, 0.5, color.G, 0.001)
				require.InDelta(t, 0.5, color.B, 0.001)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hist := tt.setup()
			err := hist.ApplyCorrection(tt.gamma)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			if tt.checkResult != nil {
				tt.checkResult(t, hist)
			}
		})
	}
}

func TestHistField_GetColor(t *testing.T) {
	tests := []struct {
		name     string
		setup    func() *HistField
		x        int
		y        int
		expected domain.Color
	}{
		{
			name: "get color from empty pixel",
			setup: func() *HistField {
				return NewHistogram(800, 600)
			},
			x:        400,
			y:        300,
			expected: domain.Color{R: 0.0, G: 0.0, B: 0.0},
		},
		{
			name: "get color from populated pixel",
			setup: func() *HistField {
				hist := NewHistogram(800, 600)
				hist.Add(0.0, 0.0, domain.Color{R: 0.8, G: 0.4, B: 0.2})
				return hist
			},
			x:        400,
			y:        300,
			expected: domain.Color{R: 0.8, G: 0.4, B: 0.2},
		},
		{
			name: "get color from out of bounds",
			setup: func() *HistField {
				return NewHistogram(800, 600)
			},
			x:        -1,
			y:        -1,
			expected: domain.Color{},
		},
		{
			name: "get color from left edge",
			setup: func() *HistField {
				hist := NewHistogram(800, 600)
				hist.Add(domain.XMIN, 0.0, domain.Color{R: 1.0, G: 0.0, B: 0.0})
				return hist
			},
			x:        0,
			y:        300,
			expected: domain.Color{R: 1.0, G: 0.0, B: 0.0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hist := tt.setup()
			color := hist.GetColor(tt.x, tt.y)

			require.InDelta(t, tt.expected.R, color.R, 0.001)
			require.InDelta(t, tt.expected.G, color.G, 0.001)
			require.InDelta(t, tt.expected.B, color.B, 0.001)
		})
	}
}

func TestHistField_GetSize(t *testing.T) {
	tests := []struct {
		name   string
		width  int
		height int
	}{
		{
			name:   "standard size",
			width:  800,
			height: 600,
		},
		{
			name:   "HD size",
			width:  1920,
			height: 1080,
		},
		{
			name:   "small size",
			width:  100,
			height: 100,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hist := NewHistogram(tt.width, tt.height)
			width, height := hist.GetSize()

			require.Equal(t, tt.width, width)
			require.Equal(t, tt.height, height)
		})
	}
}

func TestHistField_Merge(t *testing.T) {
	tests := []struct {
		name        string
		setup       func() (*HistField, *HistField)
		wantErr     bool
		checkResult func(t *testing.T, merged *HistField)
	}{
		{
			name: "merge two histograms with different pixels",
			setup: func() (*HistField, *HistField) {
				hist1 := NewHistogram(800, 600)
				hist2 := NewHistogram(800, 600)

				hist1.Add(0.0, 0.0, domain.Color{R: 1.0, G: 0.0, B: 0.0})
				hist1.Add(0.0, 0.0, domain.Color{R: 1.0, G: 0.0, B: 0.0})

				hist2.Add(0.0, domain.YMIN, domain.Color{R: 0.0, G: 1.0, B: 0.0})
				hist2.Add(0.0, domain.YMIN, domain.Color{R: 0.0, G: 1.0, B: 0.0})
				hist2.Add(0.0, domain.YMIN, domain.Color{R: 0.0, G: 1.0, B: 0.0})

				return hist1, hist2
			},
			wantErr: false,
			checkResult: func(t *testing.T, merged *HistField) {
				require.Equal(t, 3, merged.maxHits)

				pixel1 := merged.pixels[300][400]
				require.Equal(t, 2, pixel1.Counter)
				require.InDelta(t, 1.0, pixel1.Color.R, 0.001)

				pixel2 := merged.pixels[0][400]
				require.Equal(t, 3, pixel2.Counter)
				require.InDelta(t, 0.0, pixel2.Color.R, 0.001)
				require.InDelta(t, 1.0, pixel2.Color.G, 0.001)
			},
		},
		{
			name: "merge histograms with overlapping pixels",
			setup: func() (*HistField, *HistField) {
				hist1 := NewHistogram(800, 600)
				hist2 := NewHistogram(800, 600)

				hist1.Add(0.0, 0.0, domain.Color{R: 1.0, G: 0.0, B: 0.0})
				hist1.Add(0.0, 0.0, domain.Color{R: 1.0, G: 0.0, B: 0.0})

				hist2.Add(0.0, 0.0, domain.Color{R: 0.0, G: 1.0, B: 0.0})
				hist2.Add(0.0, 0.0, domain.Color{R: 0.0, G: 1.0, B: 0.0})
				hist2.Add(0.0, 0.0, domain.Color{R: 0.0, G: 1.0, B: 0.0})

				return hist1, hist2
			},
			wantErr: false,
			checkResult: func(t *testing.T, merged *HistField) {
				pixel := merged.pixels[300][400]
				require.Equal(t, 5, pixel.Counter)

				expectedRed := (2.0/5.0)*1.0 + (3.0/5.0)*0.0
				expectedGreen := (2.0/5.0)*0.0 + (3.0/5.0)*1.0

				require.InDelta(t, expectedRed, pixel.Color.R, 0.001)
				require.InDelta(t, expectedGreen, pixel.Color.G, 0.001)
				require.InDelta(t, 0.0, pixel.Color.B, 0.001)
			},
		},
		{
			name: "merge with wrong histogram type",
			setup: func() (*HistField, *HistField) {
				hist1 := NewHistogram(800, 600)
				return hist1, &HistField{}
			},
			wantErr: true,
			checkResult: func(t *testing.T, merged *HistField) {
			},
		},
		{
			name: "merge histograms of different sizes",
			setup: func() (*HistField, *HistField) {
				hist1 := NewHistogram(800, 600)
				hist2 := NewHistogram(1024, 768)
				return hist1, hist2
			},
			wantErr: true,
			checkResult: func(t *testing.T, merged *HistField) {
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hist1, hist2 := tt.setup()

			if tt.name == "merge with wrong histogram type" {
				fakeHist := struct {
					domain.Histogram
				}{}
				err := hist1.Merge(fakeHist)
				require.Error(t, err)
				return
			}

			err := hist1.Merge(hist2)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			if tt.checkResult != nil {
				tt.checkResult(t, hist1)
			}
		})
	}
}

func TestNewHistogram(t *testing.T) {
	tests := []struct {
		name   string
		width  int
		height int
	}{
		{
			name:   "standard size",
			width:  800,
			height: 600,
		},
		{
			name:   "square size",
			width:  1000,
			height: 1000,
		},
		{
			name:   "small size",
			width:  10,
			height: 10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hist := NewHistogram(tt.width, tt.height)

			require.NotNil(t, hist)
			require.Equal(t, tt.width, hist.width)
			require.Equal(t, tt.height, hist.height)
			require.Equal(t, 0, hist.maxHits)

			require.Len(t, hist.pixels, tt.height)
			for y := 0; y < tt.height; y++ {
				require.Len(t, hist.pixels[y], tt.width)
				for x := 0; x < tt.width; x++ {
					require.Equal(t, 0, hist.pixels[y][x].Counter)
					require.Equal(t, 0.0, hist.pixels[y][x].Color.R)
					require.Equal(t, 0.0, hist.pixels[y][x].Color.G)
					require.Equal(t, 0.0, hist.pixels[y][x].Color.B)
				}
			}
		})
	}
}

func TestNewHistogramFactory(t *testing.T) {
	t.Run("factory creates histogram", func(t *testing.T) {
		factory := NewHistogramFactory()
		require.NotNil(t, factory)

		hist := factory.Create(800, 600)
		require.NotNil(t, hist)

		_, ok := hist.(*HistField)
		require.True(t, ok)
	})
}
