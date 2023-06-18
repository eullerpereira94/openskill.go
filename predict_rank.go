package openskill

import (
	"math"
	"sort"

	"github.com/ernestosuarez/itertools"
	"gonum.org/v1/gonum/floats"
)

// RankDataMin calculates the ranks for a given slice of float64 values.
// It assigns ranks based on the values in ascending order, where the lowest value gets rank 1.
// If there are ties (equal values), the tied elements receive the same rank, with the next rank skipped.
//
// Parameters:
//
//	data: A slice of float64 values for which ranks need to be calculated.
//
// Returns:
//
//	A new slice of float64 values representing the ranks corresponding to the input data.
func RankDataMin(data []float64) []float64 {
	type IndexedData struct {
		Value float64
		Index int
	}
	indexedData := make([]IndexedData, len(data))
	ranks := make([]float64, len(data))
	for i, v := range data {
		indexedData[i] = IndexedData{Value: v, Index: i}
	}
	sort.Slice(indexedData, func(i, j int) bool {
		return indexedData[i].Value < indexedData[j].Value
	})
	var rank float64 = 1
	for i := 0; i < len(data); i++ {
		ranks[indexedData[i].Index] = rank
		if i < len(data)-1 && indexedData[i].Value != indexedData[i+1].Value {
			rank++
		}
	}
	return ranks
}

// PredictRank calculates and predicts the ranks of teams based on pairwise probabilities.
//
// Parameters:
//
//	teams: A slice of Team structs representing the teams.
//	options: A pointer to an Options struct that holds the configuration options.
//
// Returns:
//
//	predictions: A 2D slice of float64 containing the predicted ranks and corresponding probabilities.
//	             Each inner slice has two elements: rank and probability.
//	             The outer slice represents the predictions for each team.
//	             The length of predictions slice is equal to the number of teams.
func PredictRank(teams []Team, options *Options) (predictions [][]float64) {
	if len(teams) < 2 {
		return [][]float64{{1, 1}}
	}
	n := float64(len(teams))
	var totalPlayerCount float64
	var teamIDs []int
	teamMap := make(map[int]Team)
	winProbs := make(map[int]float64)
	for i, t := range teams {
		totalPlayerCount += float64(len(t))
		teamIDs = append(teamIDs, i)
		teamMap[i] = t
		winProbs[i] = 0
	}
	denom := (n * (n - 1)) / 2
	drawProbability := 1 / n
	drawMargin := math.Sqrt(totalPlayerCount) * beta(options) * ppf((1+drawProbability)/2)
	betaSq := betaSq(options)

	for matchup := range itertools.PermutationsInt(teamIDs, 2) {
		currentRatings := teamRatings(options)([]Team{teamMap[matchup[0]], teamMap[matchup[1]]})

		muA := currentRatings[0].TeamMu
		sigmaA := currentRatings[0].TeamSigmaSq
		muB := currentRatings[1].TeamMu
		sigmaB := currentRatings[1].TeamSigmaSq

		sigmaBar := math.Sqrt(n*betaSq + math.Pow(sigmaA, 2) + math.Pow(sigmaB, 2))
		winProb := cdf((muA - muB - drawMargin) / sigmaBar)
		winProbs[matchup[0]] += winProb
	}

	var rankedProbability []float64
	for _, teamProb := range winProbs {
		rankedProbability = append(rankedProbability, math.Abs(teamProb/denom))
	}

	ranks := RankDataMin(rankedProbability)
	maxOrdinal := floats.Max(ranks)
	for i, rank := range ranks {
		r := math.Abs(rank-maxOrdinal) + 1
		predictions = append(predictions, []float64{r, rankedProbability[i]})
	}
	return predictions
}
