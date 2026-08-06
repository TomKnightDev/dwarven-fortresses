package worldgen

import (
	"math/rand"

	"github.com/ojrac/opensimplex-go"
)

type WorldGen struct {
	Octaves     int
	Lacunarity  float64
	Scale       float64
	Persistance float64

	sampler opensimplex.Noise
}

func New() WorldGen {
	return WorldGen{
		sampler: opensimplex.New(rand.Int63()),
	}
}

func GenerateWorld() {

}

func (w *WorldGen) GenerateXYTile(x, y int) float64 {
	amplitude := 1.0
	frequency := 1.0
	noiseHeight := 1.0

	for i := 0; i < w.Octaves; i++ {
		sampleX := float64(x) / w.Scale * frequency
		sampleY := float64(y) / w.Scale * frequency

		perlinValue := w.sampler.Eval2(sampleX, sampleY)
		noiseHeight += perlinValue * amplitude

		amplitude *= w.Persistance
		frequency *= w.Lacunarity
	}

	return noiseHeight
}

func (w *WorldGen) GenerateXYZTile(x, y, z int) float64 {
	amplitude := 1.0
	frequency := 1.0
	noiseHeight := 1.0

	for i := 0; i < w.Octaves; i++ {
		sampleX := float64(x) / w.Scale * frequency
		sampleY := float64(y) / w.Scale * frequency
		sampleZ := float64(z) / w.Scale * frequency

		perlinValue := w.sampler.Eval3(sampleX, sampleY, sampleZ)
		noiseHeight += perlinValue * amplitude

		amplitude *= w.Persistance
		frequency *= w.Lacunarity
	}

	return noiseHeight
}

func (w *WorldGen) GenerateXYZTileWithFactors(x, y, z, octaves int, scale, persistance, lacunarity float64) float64 {
	amplitude := 1.0
	frequency := 1.0
	noiseHeight := 1.0

	for i := 0; i < octaves; i++ {
		sampleX := float64(x) / scale * frequency
		sampleY := float64(y) / scale * frequency
		sampleZ := float64(z) / scale * frequency

		perlinValue := w.sampler.Eval3(sampleX, sampleY, sampleZ)
		noiseHeight += perlinValue * amplitude

		amplitude *= persistance
		frequency *= lacunarity
	}

	return noiseHeight
}
