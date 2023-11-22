package trace

/**
* created by mengqi on 2023/11/13
 */

const TraceName = "chaos"

type Options struct {
	Name     string  `json:"name"`
	Endpoint string  `json:"endpoint"`
	Sampler  float64 `json:"sampler"`
	Batcher  string  `json:"batcher"`
}
