// The dump_leaderboard command pretty-prints a single day's worth of a Advent
// of Code private leaderboard, making it easy to see the completion order for
// each star.
//
// To use, download the JSON version of the private leaderboard. Pass the path
// to that file using --path. Specify the day of interest using --day.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"sort"
	"strconv"
	"time"
)

var (
	path    = flag.String("path", "", "leaderboard path")
	dayFlag = flag.Int("day", 0, "day num")
	sortFlag = flag.String("sort", "default", "member sort")

	allowedSorts = "names, stars"
)

// The structs used to decode the JSON leaderboard

type LeaderboardJSON struct {
	Members map[int]MemberJSON
}

type MemberJSON struct {
	Name               string
	Stars              int
	CompletionDayLevel map[string]map[string]StarJSON `json:"completion_day_level"`
}

type StarJSON struct {
	GetStarTs string `json:"get_star_ts"`
}

// The friendlier native struct used to represent the leaderboard
type Member struct {
	Name  string
	Stars int
	Ranks [25][2]string

	// The first map is keyed by day number, the second by star number. The
	// timestamp is the completion time.
	Completions map[int]map[int]time.Time
}

type ByName []Member

func (a ByName) Len() int           { return len(a) }
func (a ByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByName) Less(i, j int) bool { return a[i].Name < a[j].Name }

type ByStars []Member

func (a ByStars) Len() int           { return len(a) }
func (a ByStars) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByStars) Less(i, j int) bool { return a[i].Stars > a[j].Stars }

// Contains the completion time for a single star for a single user. This
// container exists largely to enable sorting.
type Result struct {
	Name string
	Ts   time.Time
}

type ByTimestamp []Result

func (a ByTimestamp) Len() int           { return len(a) }
func (a ByTimestamp) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByTimestamp) Less(i, j int) bool { return a[i].Ts.Before(a[j].Ts) }

func AtoiOrDie(s string) int {
	val, err := strconv.Atoi(s)
	if err != nil {
		panic(fmt.Sprintf("failed to parse %v: %v", s, err))
	}
	return val
}

func main() {
	flag.Parse()

	if *path == "" {
		log.Fatalf("--path is required")
	}

	contents, err := ioutil.ReadFile(*path)
	if err != nil {
		log.Fatal(err)
	}

	var board LeaderboardJSON
	if err := json.Unmarshal([]byte(contents), &board); err != nil {
		log.Fatal(err)
	}

	// Transform the raw JSON-based structs into a friendlier version.
	members := []Member{}
	for _, jsonMember := range board.Members {
		member := Member{
			Name:  jsonMember.Name,
			Stars: jsonMember.Stars,
			Ranks: [25][2]string{},
		}

		completions := map[int]map[int]time.Time{}

		for dayNumStr, day := range jsonMember.CompletionDayLevel {
			dayNum := AtoiOrDie(dayNumStr)
			completions[dayNum] = map[int]time.Time{}
			for starNumStr, jsonStar := range day {
				starNum := AtoiOrDie(starNumStr)
				starTs := time.Unix(int64(AtoiOrDie(jsonStar.GetStarTs)), 0)
				completions[dayNum][starNum] = starTs
			}
		}

		member.Completions = completions

		members = append(members, member)
	}

	if *sortFlag == "default" {
		sort.Sort(ByName(members))
	} else if *sortFlag == "stars" {
		sort.Sort(ByStars(members))
	} else if *sortFlag == "names" {
		sort.Sort(ByName(members))
	} else {
		log.Fatal(fmt.Sprintf("invalid sort: %s. Allowed: %s", *sortFlag, allowedSorts))
	}

	// Build the completions for the given day.
	var results [25]map[int][]Result

  // Initialize results map for each day.
	for day := range results {
		results[day] = map[int][]Result{
			1: []Result{},
			2: []Result{},
		}
	}

	for _, member := range members {
		for day, completion := range member.Completions {
			results[day][1] = append(results[day][1], Result{member.Name, completion[1]})
			results[day][2] = append(results[day][2], Result{member.Name, completion[2]})
		}
	}

	if *dayFlag == 0 {
		fmt.Printf("\nUse --day flag for day ranks with times\n\n")
		dailyRanks(results, members)
		return
	}

	// Dump the results.
	for starNum := 1; starNum <= 2; starNum++ {
		starResults := results[*dayFlag][starNum]
		sort.Sort(ByTimestamp(starResults))

		fmt.Printf("== star %d\n", starNum)
		for _, r := range starResults {
			fmt.Printf("%-20s %v\n", r.Name, r.Ts)
		}
	}
}

func dailyRanks(results [25]map[int][]Result, members []Member) {
	// Sort the results for each day.
	for day := range results {
		sort.Sort(ByTimestamp(results[day][1]))
		sort.Sort(ByTimestamp(results[day][2]))
	}

	// Assign rank to members for each day.
	for day := range results {
		for ndx, member := range members {
			members[ndx].Ranks[day][0] = getRank(results[day][1], member.Name)
			members[ndx].Ranks[day][1] = getRank(results[day][2], member.Name)
		}
	}

	var dayNums = ""
	for day := range [25]int{} {
		dayNums += strconv.Itoa((day + 1) % 10) + "_ "
	}
	fmt.Printf("%-20s %s\n", "== Day: ", dayNums)

	for _, member := range members {
		ranks := ""
		for day := range [25]int{} {
			ranks += member.Ranks[day][0] + member.Ranks[day][1] + " "
		}
		fmt.Printf("%-20s %s\n",
			map[bool]string{true: member.Name, false: "Anonymous"} [ len(member.Name) > 0 ],
			ranks)
	}

	if *sortFlag == "default" {
		fmt.Printf("== Add --sort flag for day ranks with times (flags: %s)\n", allowedSorts)
	} else {
		fmt.Printf("== Other sort flags: %s\n", allowedSorts) // Could filter out current flag.
	}
}

func getRank(sortedDayResults []Result, name string) string {
	rank := "."
	for ndx, result := range sortedDayResults {
		if result.Name == name {
			// Avoid weird map ternary, since not inline.
			if ndx > (10 - 1 - 1) {
				rank = ">"
			} else {
				rank = strconv.Itoa(ndx + 1)
			}
		}
	}
	return rank
}
