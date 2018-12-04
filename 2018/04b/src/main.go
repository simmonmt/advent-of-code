package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"time"

	"intmath"
)

var (
	recordPattern = regexp.MustCompile(`^\[(....-..-.. ..:..)\] (.*)$`)
	guardPattern  = regexp.MustCompile(`^Guard #([0-9]+) begins shift`)
)

type Action int

const (
	ACTION_UNKNOWN Action = iota
	ACTION_BEGIN
	ACTION_ASLEEP
	ACTION_WAKEUP
)

func (a Action) String() string {
	switch a {
	case ACTION_UNKNOWN:
		return "unknown"
	case ACTION_BEGIN:
		return "begin shift"
	case ACTION_ASLEEP:
		return "falls asleep"
	case ACTION_WAKEUP:
		return "wakes up"
	default:
		panic(fmt.Sprintf("bad action %v", int(a)))
	}
}

type Event struct {
	Tm       time.Time
	Type     Action
	GuardNum int
}

func (e *Event) String() string {
	return fmt.Sprintf("%s %s %d", e.Tm, e.Type, e.GuardNum)
}

type ByTime []*Event

func (a ByTime) Len() int           { return len(a) }
func (a ByTime) Swap(i, j int)      { *a[i], *a[j] = *a[j], *a[i] }
func (a ByTime) Less(i, j int) bool { return a[i].Tm.Before(a[j].Tm) }

func ReadEvents() ([]*Event, error) {
	events := []*Event{}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()

		recordParts := recordPattern.FindStringSubmatch(line)
		if recordParts == nil {
			return nil, fmt.Errorf("failed to parse record %v", line)
		}

		tm, err := time.Parse("2006-01-02 15:04", recordParts[1])
		if err != nil {
			return nil, fmt.Errorf("failed to parse time from %v: %v", line, err)
		}

		msg := recordParts[2]
		eventType := ACTION_UNKNOWN
		guardNum := -1
		switch msg {
		case "falls asleep":
			eventType = ACTION_ASLEEP
			break
		case "wakes up":
			eventType = ACTION_WAKEUP
			break
		default:
			parts := guardPattern.FindStringSubmatch(msg)
			if parts == nil {
				return nil, fmt.Errorf("failed to parse guard msg in %v", line)
			}
			eventType = ACTION_BEGIN
			guardNum = intmath.AtoiOrDie(parts[1])
			break
		}

		events = append(events, &Event{
			Tm:       tm,
			Type:     eventType,
			GuardNum: guardNum,
		})
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read failed: %v", err)
	}

	return events, nil
}

func main() {
	events, err := ReadEvents()
	if err != nil {
		log.Fatalf("failed to read events: %v")
	}

	sort.Sort(ByTime(events))

	sleeps := map[int]int{}
	guards := map[int]bool{}
	eventsByGuard := map[int][]*Event{}
	curGuard := -1
	var asleepAt time.Time

	for _, event := range events {
		//fmt.Println(event)
		switch event.Type {
		case ACTION_BEGIN:
			curGuard = event.GuardNum
			if curGuard < 0 {
				panic("bad curGuard")
			}
			guards[curGuard] = true
			asleepAt = time.Time{}
			break
		case ACTION_ASLEEP:
			if !asleepAt.IsZero() {
				panic("asleep")
			}
			asleepAt = event.Tm

			if eventsByGuard[curGuard] == nil {
				eventsByGuard[curGuard] = []*Event{event}
			} else {
				eventsByGuard[curGuard] = append(eventsByGuard[curGuard], event)
			}
			break
		case ACTION_WAKEUP:
			if asleepAt.IsZero() {
				panic("not asleep")
			}
			if curGuard == -1 {
				panic("no guard")
			}

			sleptFor := event.Tm.Sub(asleepAt).Minutes()
			//fmt.Printf("guard %v slept for %v\n", curGuard, sleptFor)
			sleeps[curGuard] += int(sleptFor)

			if eventsByGuard[curGuard] == nil {
				eventsByGuard[curGuard] = []*Event{event}
			} else {
				eventsByGuard[curGuard] = append(eventsByGuard[curGuard], event)
			}

			asleepAt = time.Time{}
			break
		default:
			panic("bad action")
		}
	}

	maxAllGuard := -1
	maxAllMinNum := -1
	maxAllMin := -1

	for guardNum, _ := range guards {
		minutes := map[int]int{}
		for _, event := range eventsByGuard[guardNum] {
			switch event.Type {
			case ACTION_ASLEEP:
				asleepAt = event.Tm
				break
			case ACTION_WAKEUP:
				for cur := asleepAt; cur != event.Tm; cur = cur.Add(time.Minute) {
					minutes[cur.Minute()]++
				}
			}
		}

		maxMin := -1
		maxMinNum := -1
		for min, num := range minutes {
			if maxMin == -1 || num > maxMinNum {
				maxMin = min
				maxMinNum = num
			}
		}

		if maxMin == -1 {
			continue
		}

		if maxAllGuard == -1 || maxMinNum > maxAllMinNum {
			maxAllGuard = guardNum
			maxAllMinNum = maxMinNum
			maxAllMin = maxMin
		}
	}

	fmt.Printf("guard %v min %v num %v ret %v\n", maxAllGuard, maxAllMin, maxAllMinNum, maxAllGuard*maxAllMin)

}
