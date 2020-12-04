package dice

import (
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
				Skilled:   []Die{GoodDie, GoodDie, GoodDie},
				Unskilled: []Die{},
				Lucky:     []Die{GoodDie, GoodDie, GoodDie},
				Unlucky:   []Die{},
			},
			false,
		},
		{
			"Test 3s1l",
			"3s1l",
			Pool{
				Skilled:   []Die{GoodDie, GoodDie, GoodDie},
				Unskilled: []Die{},
				Lucky:     []Die{GoodDie},
				Unlucky:   []Die{PoorDie, PoorDie},
			},
			false,
		},
		{
			"Test 1s3l",
			"1s3l",
			Pool{
				Skilled:   []Die{GoodDie},
				Unskilled: []Die{PoorDie, PoorDie},
				Lucky:     []Die{GoodDie, GoodDie, GoodDie},
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
			"Roll",
			1,
			Pool{
				Skilled:   []Die{GoodDie, GoodDie, GoodDie},
				Unskilled: []Die{PoorDie, PoorDie, PoorDie},
				Lucky:     []Die{GoodDie, GoodDie, GoodDie},
				Unlucky:   []Die{PoorDie, PoorDie, PoorDie},
			},
			Results{
				Skill: []Result{
					{0, false, Die{}},
					{0, false, Die{}},
					{0, false, Die{}},
					{0, false, Die{}},
					{0, false, Die{}},
					{0, false, Die{}},
				},
				Luck: []Result{
					{0, false, Die{}},
					{0, false, Die{}},
					{0, false, Die{}},
					{0, false, Die{}},
					{0, false, Die{}},
					{0, false, Die{}},
				},
			},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			RandSource = 0
			got := Roll(tt.pool)

			if reflect.TypeOf(got) != reflect.TypeOf(tt.want) {
				t.Errorf("Roll() = %v, want: %v", reflect.TypeOf(got), reflect.TypeOf(tt.want))
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
			"Skill Successes 1 (rolled 5G 3R); Luck Successes 1 (rolled 6W 4B)",
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
