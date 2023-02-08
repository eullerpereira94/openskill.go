package openskill

import (
	"math"
	"sort"

	"github.com/samber/lo"
	"golang.org/x/exp/constraints"
)

type teamRating struct {
	TeamMu      float64
	TeamSigmaSq float64
	Team        *Team
	Rank        int64
}

func rankings(teams []Team, ranks []int64) []int64 {
	teamScores := lo.Map(teams, func(item Team, index int) int64 {
		if index < len(ranks) {
			return ranks[index]
		}
		return int64(index)
	})

	outrank := make([]int64, len(teams))

	var s int64 = 0

	for i := 0; i < len(teamScores); i++ {
		if (i > 0) && teamScores[i-1] < teamScores[i] {
			s = int64(i)
		}
		outrank[i] = s
	}

	return outrank
}

func teamRatings(options *Options) func(game []Team) []teamRating {
	return func(game []Team) []teamRating {
		rank := rankings(game, options.Rankings)

		return lo.Map(game, func(item Team, index int) teamRating {
			mu := lo.Sum(lo.Map([]*Rating(item), func(item *Rating, index int) float64 {
				return item.AveragePlayerSkill
			}))
			sigma := lo.Sum(lo.Map([]*Rating(item), func(item *Rating, index int) float64 {
				return math.Pow(item.SkillUncertaintyDegree, 2)
			}))

			return teamRating{
				Team:        &item,
				TeamMu:      mu,
				TeamSigmaSq: sigma,
				Rank:        rank[index],
			}
		})
	}
}

func utilC(options *Options) func(teamRatings []teamRating) float64 {
	betasq := betaSq(options)

	return func(teamRatings []teamRating) float64 {
		return math.Sqrt(
			lo.Sum(lo.Map(teamRatings, func(item teamRating, index int) float64 {
				return item.TeamSigmaSq + betasq
			})),
		)
	}
}

func utilSumQ(teamRatings []teamRating, c float64) []float64 {
	return lo.Map(teamRatings, func(item teamRating, index int) float64 {
		filteredRatings := lo.Filter(teamRatings, func(localItem teamRating, index int) bool {
			return localItem.Rank >= item.Rank
		})
		mappedFilteredRatings := lo.Map(filteredRatings, func(localItem teamRating, index int) float64 {
			return math.Exp(localItem.TeamMu / c)
		})

		return lo.Sum(mappedFilteredRatings)
	})
}

func utilA(teamRatings []teamRating) []int64 {
	return lo.Map(teamRatings, func(item teamRating, index int) int64 {
		filteredRatings := lo.Filter(teamRatings, func(localItem teamRating, index int) bool {
			return item.Rank == localItem.Rank
		})

		return int64(len(filteredRatings))
	})
}

func gamma(options *Options) Gamma {
	if options.GammaFunction != nil {
		return *options.GammaFunction
	}

	return func(c float64, k int64, mu, sigmaSq float64, team *Team, qRank int64) float64 {
		return math.Sqrt(sigmaSq) / c
	}
}

func unwind[T constraints.Integer, R any](order []T, collection []R) (sortedCollection []R, stochasticTenet []T) {
	if len(collection) <= 0 {
		sortedCollection = make([]R, 0)
		stochasticTenet = make([]T, 0)
		return
	}

	zipped := []struct {
		x T
		y T
		z R
	}{}

	for i, v := range collection {
		zipped = append(zipped, struct {
			x T
			y T
			z R
		}{x: order[i], y: T(i), z: v})
	}

	sort.Slice(zipped, func(i, j int) bool {
		return zipped[i].x < zipped[j].x
	})

	for _, v := range zipped {
		sortedCollection = append(sortedCollection, v.z)
		stochasticTenet = append(stochasticTenet, v.y)
	}

	return
}
