package dice

import (
	"fmt"
	"github.com/kyokomi/emoji"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	Skilled Kind = iota
	Unskilled
	Lucky
	Unlucky
)

const (
	MaxDice = 3
)

type (
	Kind int

	Die struct {
		Sides           int // The number of sides this die has.
		SuccessfulSides int // The number of sides on this die which counts as aa success.
		Kind            Kind
	}

	Dice struct {
		Count int
		Die   Die
	}

	Pool struct {
		Skilled   []Die
		Unskilled []Die
		Lucky     []Die
		Unlucky   []Die
	}

	Result struct {
		Value     int
		IsSuccess bool
		Die       Die
	}

	Results struct {
		Skill []Result
		Luck  []Result
	}
)

var (
	regex        *regexp.Regexp
	GoodDie      = Die{Sides: 6, SuccessfulSides: 2} // Represents Skilled and Lucky.
	PoorDie      = Die{Sides: 6, SuccessfulSides: 1} // Represents Unskilled and Unlucky.
	SkilledDie   = Die{Sides: 6, SuccessfulSides: 2, Kind: Skilled}
	UnskilledDie = Die{Sides: 6, SuccessfulSides: 1, Kind: Unskilled}
	LuckyDie     = Die{Sides: 6, SuccessfulSides: 2, Kind: Lucky}
	UnluckyDie   = Die{Sides: 6, SuccessfulSides: 1, Kind: Unlucky}
	RandSource   = time.Now().UnixNano()
	roller       *rand.Rand
)

func init() {
	regex = regexp.MustCompile(`((?P<skilledcount>\d)s)((?P<luckycount>\d)l)`)
	source := rand.NewSource(RandSource)
	roller = rand.New(source)
}

func Parse(diceString string) (Pool, error) {
	matches := regex.FindStringSubmatch(diceString)
	pool := Pool{Skilled: []Die{}, Unskilled: []Die{}, Lucky: []Die{}, Unlucky: []Die{}}
	result := make(map[string]int)

	for i, name := range regex.SubexpNames() {
		if i != 0 && name != "" {
			if name == "skilledcount" || name == "luckycount" {
				count, err := strconv.Atoi(matches[i])
				if err != nil {
					return pool, fmt.Errorf("cannot get a number from \"%s\"", matches[i])
				}

				result[name] = count
			}
		}
	}

	err := addSkilledDice(result["skilledcount"], &pool)
	if err != nil {
		return pool, err
	}

	err = addLuckyDice(result["luckycount"], &pool)
	if err != nil {
		return pool, err
	}

	return pool, nil
}

func addSkilledDice(skilledCount int, pool *Pool) error {
	if skilledCount > MaxDice {
		return fmt.Errorf("the number of skilled dice, %d, exceeds the maximum of %d", skilledCount, MaxDice)
	}

	err := addDice(skilledCount, pool, Skilled)
	if err != nil {
		return err
	}

	if len(pool.Skilled)+len(pool.Unskilled) != MaxDice {
		return fmt.Errorf("somehow the number of skilled dice, %d, and the number of unskilled dice, %d, is not equal to %d", len(pool.Skilled), len(pool.Unskilled), MaxDice)
	}

	return nil
}

func addLuckyDice(luckyCount int, pool *Pool) error {
	if luckyCount > MaxDice {
		return fmt.Errorf("the number of lucky dice, %d, exceeds the maximum of %d", luckyCount, MaxDice)
	}

	err := addDice(luckyCount, pool, Lucky)
	if err != nil {
		return err
	}

	if len(pool.Lucky)+len(pool.Unlucky) != MaxDice {
		return fmt.Errorf("somehow the number of lucky dice, %d, and the number of unlucky dice, %d, is not equal to %d", len(pool.Lucky), len(pool.Unlucky), MaxDice)
	}

	return nil
}

func addDice(count int, pool *Pool, kind Kind) error {
	remaining := MaxDice - count

	if kind == Skilled {
		for i := 0; i < count; i++ {
			pool.Skilled = append(pool.Skilled, SkilledDie)
		}

		for i := 0; i < remaining; i++ {
			pool.Unskilled = append(pool.Unskilled, UnskilledDie)
		}
	} else if kind == Lucky {
		for i := 0; i < count; i++ {
			pool.Lucky = append(pool.Lucky, LuckyDie)
		}

		for i := 0; i < remaining; i++ {
			pool.Unlucky = append(pool.Unlucky, UnluckyDie)
		}
	}

	return nil
}

func Roll(pool Pool) Results {
	var results Results

	for _, die := range pool.Skilled {
		results.Skill = append(results.Skill, rollDie(die))
	}

	for _, die := range pool.Unskilled {
		results.Skill = append(results.Skill, rollDie(die))
	}

	for _, die := range pool.Lucky {
		results.Luck = append(results.Luck, rollDie(die))
	}

	for _, die := range pool.Unlucky {
		results.Luck = append(results.Luck, rollDie(die))
	}

	return results
}

func rollDie(d Die) Result {
	var result Result

	roll := roller.Intn(d.Sides) + 1
	result.Value = roll
	result.Die = d

	if roll > d.Sides-d.SuccessfulSides {
		result.IsSuccess = true
	}

	return result
}

func (d *Die) String() string {
	switch d.Kind {
	case Skilled:
		return emoji.Sprint(":green_square:")
	case Unskilled:
		return emoji.Sprint(":red_square:")
	case Lucky:
		return emoji.Sprint(":white_large_square:")
	case Unlucky:
		return emoji.Sprint(":black_large_square:")
	}

	return ""
}

func (r *Results) String() string {
	var rollsSkill, rollsLuck string
	var successesSkill, successesLuck int

	for _, result := range r.Skill {
		if result.IsSuccess {
			successesSkill += 1
		}

		rollsSkill += strconv.Itoa(result.Value) + " " + result.Die.String() + "  "
	}

	for _, result := range r.Luck {
		if result.IsSuccess {
			successesLuck += 1
		}

		rollsLuck += strconv.Itoa(result.Value) + " " + result.Die.String() + "  "
	}

	return fmt.Sprintf("**Skill Successes:** **%d**\n%s\n**Luck Successes:** **%d**\n%s", successesSkill, strings.Trim(rollsSkill, " "), successesLuck, strings.Trim(rollsLuck, " "))
}
