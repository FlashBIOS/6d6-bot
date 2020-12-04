package dice

import (
	"math/rand"
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name       string
		diceString string
		want       Pool
		wantErr    bool
	}{
		{
			"Test 3s3l",
			"3s3l",
			Pool{
				Skilled:   []Die{SkilledDie, SkilledDie, SkilledDie},
				Unskilled: []Die{},
				Lucky:     []Die{LuckyDie, LuckyDie, LuckyDie},
				Unlucky:   []Die{},
			},
			false,
		},
		{
			"Test 3s1l",
			"3s1l",
			Pool{
				Skilled:   []Die{SkilledDie, SkilledDie, SkilledDie},
				Unskilled: []Die{},
				Lucky:     []Die{LuckyDie},
				Unlucky:   []Die{UnluckyDie, UnluckyDie},
			},
			false,
		},
		{
			"Test 1s3l",
			"1s3l",
			Pool{
				Skilled:   []Die{SkilledDie},
				Unskilled: []Die{UnskilledDie, UnskilledDie},
				Lucky:     []Die{LuckyDie, LuckyDie, LuckyDie},
				Unlucky:   []Die{},
			},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.diceString)

			if (err != nil) != tt.wantErr {
				t.Errorf("Parse(%s) error = %v, wantErr: %v", tt.diceString, err, tt.wantErr)
			} else if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse(%s) = %v, want: %v", tt.diceString, got, tt.want)
			}
		})
	}
}

func TestRoll(t *testing.T) {
	tests := []struct {
		name    string
		source  int64
		pool    Pool
		want    Results
		wantErr bool
	}{
		{
			"Roll 1",
			0,
			Pool{
				Skilled:   []Die{SkilledDie},
				Unskilled: []Die{UnskilledDie, UnskilledDie},
				Lucky:     []Die{LuckyDie},
				Unlucky:   []Die{UnluckyDie, UnluckyDie},
			},
			Results{
				Skill: []Result{
					{1, false, SkilledDie},
					{1, false, UnskilledDie},
					{2, false, UnskilledDie},
				},
				Luck: []Result{
					{5, true, LuckyDie},
					{6, true, UnluckyDie},
					{5, true, UnluckyDie},
				},
			},
			false,
		},
		{
			"Roll 2",
			1,
			Pool{
				Skilled:   []Die{SkilledDie},
				Unskilled: []Die{UnskilledDie, UnskilledDie},
				Lucky:     []Die{LuckyDie},
				Unlucky:   []Die{UnluckyDie, UnluckyDie},
			},
			Results{
				Skill: []Result{
					{6, true, SkilledDie},
					{4, false, UnskilledDie},
					{6, true, UnskilledDie},
				},
				Luck: []Result{
					{6, true, LuckyDie},
					{2, false, UnluckyDie},
					{1, false, UnluckyDie},
				},
			},
			false,
		},
		{
			"Roll 3",
			2,
			Pool{
				Skilled:   []Die{SkilledDie},
				Unskilled: []Die{UnskilledDie, UnskilledDie},
				Lucky:     []Die{LuckyDie},
				Unlucky:   []Die{UnluckyDie, UnluckyDie},
			},
			Results{
				Skill: []Result{
					{5, true, SkilledDie},
					{1, false, UnskilledDie},
					{1, false, UnskilledDie},
				},
				Luck: []Result{
					{3, true, LuckyDie},
					{3, false, UnluckyDie},
					{3, false, UnluckyDie},
				},
			},
			false,
		},
		{
			"Roll 4",
			4,
			Pool{
				Skilled:   []Die{SkilledDie},
				Unskilled: []Die{UnskilledDie, UnskilledDie},
				Lucky:     []Die{LuckyDie},
				Unlucky:   []Die{UnluckyDie, UnluckyDie},
			},
			Results{
				Skill: []Result{
					{2, false, SkilledDie},
					{5, true, UnskilledDie},
					{2, false, UnskilledDie},
				},
				Luck: []Result{
					{6, true, LuckyDie},
					{4, false, UnluckyDie},
					{2, false, UnluckyDie},
				},
			},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			roller = rand.New(rand.NewSource(tt.source))
			got := Roll(tt.pool)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Roll() = %v, want: %v", got, tt.want)
			}
		})
	}
}

func TestResults_String(t *testing.T) {
	var tests = []struct {
		name    string
		results Results
		want    string
		wantErr bool
	}{
		{
			"To string conversion",
			Results{
				Skill: []Result{
					{
						Value:     5,
						IsSuccess: true,
						Die:       SkilledDie,
					},
					{
						Value:     3,
						IsSuccess: false,
						Die:       UnskilledDie,
					},
				},
				Luck: []Result{
					{
						Value:     6,
						IsSuccess: true,
						Die:       LuckyDie,
					},
					{
						Value:     4,
						IsSuccess: false,
						Die:       UnluckyDie,
					},
				},
			},
			"**Skill Successes:** **1**\n5 \U0001F7E9   3 \U0001F7E5\n**Luck Successes:** **1**\n6 ⬜   4 ⬛",
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.results.String()

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("String() = %v, want: %v", got, tt.want)
			}
		})
	}
}
