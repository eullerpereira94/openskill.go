package openskill

import (
	"reflect"
	"testing"
)

func TestRate(t *testing.T) {
	type args struct {
		teams   []Team
		options Options
	}
	tests := []struct {
		name string
		args args
		want []Team
	}{
		{
			name: "rate accepts and runs a placket-luce model by default",
			args: args{
				teams: []Team{
					[]*Rating{{AveragePlayerSkill: 29.182, SkillUncertaintyDegree: 4.782}},
					[]*Rating{{AveragePlayerSkill: 27.174, SkillUncertaintyDegree: 4.922}},
					[]*Rating{{AveragePlayerSkill: 16.672, SkillUncertaintyDegree: 6.217}},
					[]*Rating{NewRating(nil, nil)},
				},
				options: Options{},
			},
			want: []Team{
				[]*Rating{{AveragePlayerSkill: 30.209971908310553, SkillUncertaintyDegree: 4.764898977359521}},
				[]*Rating{{AveragePlayerSkill: 27.64460833689499, SkillUncertaintyDegree: 4.882789305097372}},
				[]*Rating{{AveragePlayerSkill: 17.403586731283518, SkillUncertaintyDegree: 6.100723440599442}},
				[]*Rating{{AveragePlayerSkill: 19.214790707434826, SkillUncertaintyDegree: 7.8542613981643985}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Rate(tt.args.teams, tt.args.options); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Rate() = %v, want %v", got, tt.want)
			}
		})
	}
}
