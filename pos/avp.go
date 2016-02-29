package pos

type FeatureClass struct {
	f string
	c string
}

type AveragedPerceptron struct {
	i       int
	totals  map[FeatureClass]float64
	tstamps map[FeatureClass]int
	weights map[string]map[string]float64
}

func NewAveragedPerceptron() *AveragedPerceptron {
	return &AveragedPerceptron{
		totals:  make(map[FeatureClass]float64),
		tstamps: make(map[FeatureClass]int),
		weights: make(map[string]map[string]float64),
	}
}

func (ap *AveragedPerceptron) Predict(features map[string]int) string {

	scores := map[string]float64{}

	for feat, value := range features {
		if _, ok := ap.weights[feat]; !ok || value == 0 {
			continue
		}
		for label, weight := range ap.weights[feat] {
			scores[label] += float64(value) * weight
		}
	}
	return ap.maxScore(scores)
}

func (ap *AveragedPerceptron) Update(truth, guess string, features map[string]int) {

	ap.i += 1
	if truth == guess {
		return
	}

	var wt, wg float64

	for f := range features {
		if _, ok := ap.weights[f]; !ok {
			ap.weights[f] = make(map[string]float64)
		} else {
			wt = ap.weights[f][truth]
			wg = ap.weights[f][guess]
		}
		ap.updateFeature(truth, f, wt, 1.0)
		ap.updateFeature(guess, f, wg, -1.0)
	}
	return
}

func (ap *AveragedPerceptron) AverageWeights() {

	for feat, weights := range ap.weights {

		newWeights := map[string]float64{}

		for class, weight := range weights {

			param := FeatureClass{feat, class}
			total := ap.totals[param]
			total += float64(ap.i-ap.tstamps[param]) * weight

			averaged := Round(total/float64(ap.i), 0.5, 3)
			if averaged != 0 {
				newWeights[class] = averaged
			}
		}
		ap.weights[feat] = newWeights
	}
	return
}

func (ap *AveragedPerceptron) updateFeature(class, feature string, weight, value float64) {

	param := FeatureClass{feature, class}
	ap.totals[param] += float64(ap.i-ap.tstamps[param]) * weight
	ap.tstamps[param] = ap.i
	ap.weights[feature][class] = weight + value
}

func (ap *AveragedPerceptron) maxScore(m map[string]float64) string {
	var maxKey string
	var maxVal float64

	for k, v := range m {
		switch {
		case v > maxVal:
			maxKey = k
			maxVal = v
		case v == maxVal:
			if k < maxKey {
				maxKey = k
				maxVal = v
			}
		}
	}
	return maxKey
}
