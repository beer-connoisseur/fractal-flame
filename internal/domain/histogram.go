package domain

type Histogram interface {
	Add(x, y float64, color Color)
	GetColor(x, y int) Color
	ApplyCorrection(gamma float64) error
	Merge(other Histogram) error
	GetSize() (width, height int)
}

type HistogramFactory interface {
	Create(width, height int) Histogram
}
