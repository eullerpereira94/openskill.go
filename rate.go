package openskill

import (
	"math"
	"sort"

	"github.com/samber/lo"
)

// Rate rates a group of teams with the provided optional parameters for classification
func Rate(teams []Team, options Options) []Team {
	var model Model
	var processedTeams = make([]Team, len(teams))

	if options.Model != nil {
		model = *options.Model
	} else {
		model = PlackettLuce
	}

	if options.tau != nil {
		tauSquared := math.Pow(*options.tau, 2)

		processedTeams = lo.Map(teams, func(item Team, index int) Team {
			return lo.Map([]*Rating(item), func(item *Rating, index int) *Rating {
				item.SkillUncertaintyDegree = math.Sqrt(math.Pow(item.SkillUncertaintyDegree, 2) + tauSquared)
				return item
			})
		})
	} else {
		copy(processedTeams, teams)
	}

	var rank []int64

	if len(options.Rankings) > 0 {
		rank = lo.Map(options.Rankings, func(item int64, index int) int64 {
			return item
		})
	} else if len(options.Scores) > 0 {
		rank = lo.Map(options.Scores, func(item int64, index int) int64 {
			return -item
		})
	} else {
		rank = lo.RangeFrom[int64](1, len(teams))
	}

	orderedTeams, tenet := unwind(rank, processedTeams)

	_newRanks := make([]int64, len(rank))
	copy(_newRanks, rank)

	sort.Slice(_newRanks, func(i, j int) bool {
		return _newRanks[i] < _newRanks[j]
	})

	options.Rankings = _newRanks

	newRatings := model(orderedTeams, &options)

	reorderedTeams, _ := unwind(tenet, newRatings)

	if options.tau != nil && options.PreventUncertaintyIncrease != nil && *options.PreventUncertaintyIncrease {
		reorderedTeams = lo.Map(reorderedTeams, func(item Team, index int) Team {
			return lo.Map([]*Rating(item), func(localItem *Rating, localIndex int) *Rating {
				localItem.SkillUncertaintyDegree = math.Min(localItem.SkillUncertaintyDegree, teams[index][localIndex].SkillUncertaintyDegree)
				return localItem
			})
		})
	}

	return reorderedTeams
}
