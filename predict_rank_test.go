package openskill_test

import (
	"math"
	"reflect"
	"testing"

	"github.com/eullerpereira94/openskill"
)

func TestRankDataMin(t *testing.T) {
	// Test case 1
	data1 := []float64{4, 2, 7, 2, 9, 5, 1}
	expectedRanks1 := []float64{3, 2, 5, 2, 6, 4, 1}
	ranks1 := openskill.RankDataMin(data1)
	if !reflect.DeepEqual(ranks1, expectedRanks1) {
		t.Errorf("RankData failed for test case 1. Expected %v, but got %v", expectedRanks1, ranks1)
	}

	// Test case 2: Negative values and duplicates
	data2 := []float64{-2, 4, 7, 2, 4, -2, 9}
	expectedRanks2 := []float64{1, 3, 4, 2, 3, 1, 5}
	ranks2 := openskill.RankDataMin(data2)
	if !reflect.DeepEqual(ranks2, expectedRanks2) {
		t.Errorf("RankData failed for test case 2. Expected %v, but got %v", expectedRanks2, ranks2)
	}

	// Test case 3: Test with empty data
	data3 := []float64{}
	expectedRanks3 := []float64{}
	ranks3 := openskill.RankDataMin(data3)
	if !reflect.DeepEqual(ranks3, expectedRanks3) {
		t.Errorf("RankData failed for test case 3. Expected %v, but got %v", expectedRanks3, ranks3)
	}

	// Test case 4: Test with identical values
	data4 := []float64{5, 5, 5, 5, 5}
	expectedRanks4 := []float64{1, 1, 1, 1, 1}
	ranks4 := openskill.RankDataMin(data4)
	if !reflect.DeepEqual(ranks4, expectedRanks4) {
		t.Errorf("RankData failed for test case 4. Expected %v, but got %v", expectedRanks4, ranks4)
	}

	// Test case 5
	data5 := []float64{9, 2, 8, 6, 3, 1, 4, 7, 5, 10, 4, 5}
	expectedRanks5 := []float64{9, 2, 8, 6, 3, 1, 4, 7, 5, 10, 4, 5}
	ranks5 := openskill.RankDataMin(data5)
	if !reflect.DeepEqual(ranks5, expectedRanks5) {
		t.Errorf("RankData failed for test case 5. Expected %v, but got %v", expectedRanks5, ranks5)
	}
}

func withinTolerance(a, b, tolerance float64) bool {
	if a == b {
		return true
	}
	d := math.Abs(a - b)
	if b == 0 {
		return d < tolerance
	}
	return (d / math.Abs(b)) < tolerance
}

func TestPredictRank(t *testing.T) {
	a1 := openskill.NewRating(&openskill.NewRatingParams{32, 0.25}, nil)
	a2 := openskill.NewRating(&openskill.NewRatingParams{32, 0.25}, nil)
	a3 := openskill.NewRating(&openskill.NewRatingParams{34, 0.25}, nil)

	b1 := openskill.NewRating(&openskill.NewRatingParams{24, 0.5}, nil)
	b2 := openskill.NewRating(&openskill.NewRatingParams{22, 0.5}, nil)
	b3 := openskill.NewRating(&openskill.NewRatingParams{20, 0.5}, nil)

	team1 := openskill.NewTeam(a1, b1)
	team2 := openskill.NewTeam(a2, b2)
	team3 := openskill.NewTeam(a3, b3)

	ranks := openskill.PredictRank([]openskill.Team{team1, team2, team3}, nil)

	var totalRankProbability float64
	for _, r := range ranks {
		totalRankProbability += r[1]
	}
	drawProbability := openskill.PredictDraw([]openskill.Team{team1, team2, team3}, nil)

	p := totalRankProbability + drawProbability
	expected := 1.0
	if !withinTolerance(expected, p, 0.0000001) {
		t.Errorf("Expected %f, got %f", expected, p)
	}

	oneTeam := openskill.PredictRank([]openskill.Team{team1}, nil)
	expectedoneTeam := [][]float64{{1, 1}}
	if !reflect.DeepEqual(oneTeam, expectedoneTeam) {
		t.Errorf("RankData failed for test case 5. Expected %v, but got %v", expectedoneTeam, oneTeam)
	}
}
