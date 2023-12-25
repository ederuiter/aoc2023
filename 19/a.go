package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
)

type Category string

const (
	X Category = "x"
	M          = "m"
	A          = "a"
	S          = "s"
)

const MaxRating = 4000
const MinRating = 1

type Interval struct {
	Min int
	Max int
}

type Rule struct {
	Category  Category
	Min       int
	Max       int
	Action    int
	ActionStr string
	Nested    int
}

type Workflow struct {
	Rules []Rule
}

type Part map[Category]int

type Workflows map[string]Workflow

func (r Interval) Overlaps(with Interval) bool {
	return r.Max >= with.Min && with.Max >= r.Min
}

func (r Interval) Num() int {
	return r.Max - r.Min + 1
}

func (r Interval) NonIntersecting(interval2 Interval) []Interval {
	// assert interval1 starts before or at the same time as interval2
	interval1 := r
	if interval1.Min > interval2.Min {
		fmt.Printf("%+v, %+v\n", interval1, interval2)
		panic("This should not be happening")
	}

	/*
	 * TODO: there should be a more generic way of doing this ..
	 */
	if interval1.Min == interval2.Min {
		/*
		 *  1: |====|
		 *  2: |===========|
		 */
		return []Interval{interval1, Interval{interval1.Max + 1, interval2.Max}}
	} else if interval1.Max == interval2.Max {
		/*
		 *  1: |===========|
		 *  2:        |====|
		 */
		return []Interval{Interval{interval1.Min, interval2.Min - 1}, interval2}
	} else if interval1.Max >= interval2.Min {
		/*
		 *  1: |===========|
		 *  2:             |====|
		 */

		/*
		 *  1: |===========|
		 *  2:          |====|
		 */
		return []Interval{Interval{interval1.Min, interval2.Min - 1}, {interval2.Min, interval1.Max}, {interval1.Max + 1, interval2.Max}}
	} else {
		/*
		 *  1: |===========|
		 *  2:     |====|
		 */
		return []Interval{Interval{interval1.Min, interval2.Min - 1}, {interval2.Min, interval2.Max}, {interval2.Max + 1, interval1.Max}}
	}
}

func (w Workflows) Print() {
	for name, workflow := range w {
		fmt.Printf("Workflow: %s\n", name)
		for _, rule := range workflow.Rules {
			fmt.Printf("  %+v\n", rule)
		}
	}
}

func CompareInterval(a, b Interval) int {
	res := a.Min - b.Min
	if res == 0 {
		res = a.Max - b.Max
	}
	return res
}

func mergeInterval(a, b Interval) Interval {
	return Interval{
		Min: max(a.Min, b.Min),
		Max: min(a.Max, b.Max),
	}
}

func createNonOverlappingIntervals(intervals []Interval) []Interval {
	if len(intervals) == 0 {
		return []Interval{}
	}

	slices.SortFunc(intervals, CompareInterval)
	intervals = slices.Compact(intervals)

	var currentInterval Interval
	splitIntervals := []Interval{intervals[0]}
	intervals = intervals[1:]
	for len(intervals) > 0 {
		lastInterval := splitIntervals[len(splitIntervals)-1]
		currentInterval, intervals = intervals[0], intervals[1:]
		if lastInterval.Overlaps(currentInterval) {
			nonIntersecting := lastInterval.NonIntersecting(currentInterval)
			splitIntervals[len(splitIntervals)-1] = nonIntersecting[0]
			intervals = append(nonIntersecting[1:], intervals...)
			slices.SortFunc(intervals, CompareInterval)
			intervals = slices.Compact(intervals)
		} else {
			splitIntervals = append(splitIntervals, currentInterval)
		}
	}

	return splitIntervals
}

func walkWorkflows(parentRules []Rule, workflow string, workflows Workflows) [][]Rule {
	res := [][]Rule{}
	for _, rule := range workflows[workflow].Rules {
		/*
		 * TODO:
		 *   apparently go slices are a bit more fragile that I thought .. something like this will break:
		 *   myRules := append(parentRules, rule)
		 *   ^^ this will create some sort of reference to parentRules <= need to investigate ..
		 */
		myRules := []Rule{}
		myRules = append(myRules, parentRules...)
		myRules = append(myRules, rule)
		if rule.ActionStr == "A" || rule.ActionStr == "R" {
			res = append(res, myRules)
		} else {
			r := walkWorkflows(myRules, rule.ActionStr, workflows)
			res = append(res, r...)
		}

		/*
		 * We are interested in all the rules that lead to here; that means that we need to add
		 * the inverse rules of all the previous siblings as well.
		 */
		if rule.Min > MinRating {
			parentRules = append(parentRules, Rule{
				Category:  rule.Category,
				Min:       MinRating,
				Max:       rule.Min - 1,
				Action:    0,
				ActionStr: "",
				Nested:    0,
			})
		}

		if rule.Max < MaxRating {
			parentRules = append(parentRules, Rule{
				Category:  rule.Category,
				Min:       rule.Max + 1,
				Max:       MaxRating,
				Action:    0,
				ActionStr: "",
				Nested:    0,
			})
		}
	}

	return res
}

func optimizeWorkflows(workflows Workflows) (Workflows, int) {
	res := Workflows{}
	optimized := map[string]string{}
	for name, workflow := range workflows {
		actions := ""
		for _, rule := range workflow.Rules {
			actions += rule.ActionStr
		}

		if actions == strings.Repeat(workflow.Rules[0].ActionStr, len(workflow.Rules)) {
			optimized[name] = workflow.Rules[0].ActionStr
		} else {
			res[name] = workflow
		}
	}
	for name, workflow := range res {
		for ruleIndex, rule := range workflow.Rules {
			action := rule.ActionStr
			for optimizedName, optimizedAction := range optimized {
				if action == optimizedName {
					res[name].Rules[ruleIndex].ActionStr = optimizedAction
				}
			}
		}
	}

	numOptimized := len(optimized)
	for _, workflow := range res {
		defaultRule := workflow.Rules[len(workflow.Rules)-1]
		toDelete := 0
		for i := len(workflow.Rules) - 2; i > 0; i-- {
			if workflow.Rules[i].Action == defaultRule.Action {
				toDelete++
			} else {
				break
			}
		}
		if toDelete > 0 {
			toDelete++
			numOptimized++
			newRules := []Rule{}
			newRules = append(newRules, workflow.Rules[0:len(workflow.Rules)-toDelete]...)
			newRules = append(newRules, defaultRule)
			workflow.Rules = newRules
		}
	}

	return res, len(optimized)
}

func inlineRules(workflows Workflows, workflow string) []Rule {
	res := []Rule{}
	for _, rule := range workflows[workflow].Rules {
		action := rule.ActionStr
		if action == "R" || action == "A" {
			if action == "R" {
				rule.Action = -1
			} else {
				rule.Action = 0
			}
			res = append(res, rule)
		} else {
			inlined := inlineRules(workflows, action)
			rule.Action = 1
			rule.Nested = len(inlined)
			res = append(res, rule)
			res = append(res, inlined...)
		}
	}
	return res
}

func ProcessPart(part Part, rules []Rule) bool {
	i := 0
	for i < len(rules) {
		r := rules[i]
		action := r.Nested + 1
		if part[r.Category] >= r.Min && part[r.Category] <= r.Max {
			action = r.Action
			if action <= 0 {
				return action == 0
			}
		}
		i += action
	}
	panic("the last rule should always match")
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	workflows := Workflows{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lineStr := scanner.Text()
		if lineStr == "" {
			break
		}
		line := strings.Split(scanner.Text(), "{")
		name, rulesStr := line[0], strings.Trim(line[1], "}")
		rules := strings.Split(rulesStr, ",")
		rules, defaultAction := rules[0:len(rules)-1], rules[len(rules)-1]
		workflowRules := []Rule{}
		for _, rule := range rules {
			r := strings.Split(rule, ":")
			value, _ := strconv.ParseInt(r[0][2:], 10, 64)
			minValue := MinRating
			maxValue := MaxRating
			if r[0][1] == '<' {
				maxValue = int(value) - 1
			} else {
				minValue = int(value) + 1
			}

			workflowRules = append(workflowRules, Rule{
				Category:  Category(r[0][0]),
				Min:       minValue,
				Max:       maxValue,
				ActionStr: r[1],
			})
		}
		workflowRules = append(workflowRules, Rule{
			Category:  X,
			Min:       MinRating,
			Max:       MaxRating,
			ActionStr: defaultAction,
		})
		workflows[name] = Workflow{
			Rules: workflowRules,
		}
	}

	/*
	 * TODO: not needed
	 */
	numOptimized := 1
	i := 0
	for numOptimized > 0 {
		i++
		fmt.Printf("Starting optimization round #%d\n", i)
		workflows, numOptimized = optimizeWorkflows(workflows)
		fmt.Printf(" =>%d\n", numOptimized)
	}

	workflows.Print()

	/*
	 * TODO: not needed
	 */
	inlined := inlineRules(workflows, "in")
	fmt.Printf("Inlined rules:\n")
	for _, rule := range inlined {
		fmt.Printf("  %+v\n", rule)
	}

	sum := 0
	for scanner.Scan() {
		line := strings.Split(strings.Trim(scanner.Text(), "{}"), ",")
		part := Part{X: 0, M: 0, A: 0, S: 0}
		partSum := 0
		for _, item := range line {
			value, _ := strconv.ParseInt(item[2:], 10, 64)
			category := Category(item[0])
			part[category] = int(value)
			partSum += int(value)
		}
		accepted := ProcessPart(part, inlined)
		if accepted {
			sum += partSum
			//fmt.Printf("ACCEPTED: %+v => %d\n", part, partSum)
		} else {
			//fmt.Printf("REJECTED: %+v\n", part)
		}
	}
	fmt.Printf("Total sum of accepted parts: %d\n", sum)
	sets := map[Category][]Interval{X: {{MinRating, MaxRating}}, M: {{MinRating, MaxRating}}, A: {{MinRating, MaxRating}}, S: {{MinRating, MaxRating}}}

	numAccepted := int64(0)
	numRejected := int64(0)
	parentRules := []Rule{}
	for _, ruleSet := range walkWorkflows(parentRules, "in", workflows) {

		for _, rule := range ruleSet {
			sets[rule.Category] = append(sets[rule.Category], Interval{rule.Min, rule.Max})
		}

		accepted := map[Category]Interval{X: {MinRating, MaxRating}, M: {MinRating, MaxRating}, A: {MinRating, MaxRating}, S: {MinRating, MaxRating}}
		for _, rule := range ruleSet {
			accepted[rule.Category] = mergeInterval(accepted[rule.Category], Interval{rule.Min, rule.Max})
		}
		action := ruleSet[len(ruleSet)-1].ActionStr
		if action == "A" {
			numAccepted += int64(accepted[X].Num() * accepted[M].Num() * accepted[A].Num() * accepted[S].Num())
		} else {
			numRejected += int64(accepted[X].Num() * accepted[M].Num() * accepted[A].Num() * accepted[S].Num())
		}
		fmt.Printf("%s => %+v\n", action, ruleSet)
		fmt.Printf("  => %+v\n", accepted)
	}
	fmt.Printf("Total accepted combinations: %d rejected %d\n", numAccepted, numRejected)

	/*
	 * TODO: not needed
	 */
	intervalsByCategory := map[Category][]Interval{}
	for cat, set := range sets {
		intervalsByCategory[cat] = createNonOverlappingIntervals(set)
	}

	accepted := int64(0)
	rejected := int64(0)
	part := Part{X: 0, M: 0, A: 0, S: 0}

	/*
	 * TODO: not needed
	 */
	fmt.Printf("X: %d M: %d A: %d S: %d\n", len(intervalsByCategory[X]), len(intervalsByCategory[M]), len(intervalsByCategory[A]), len(intervalsByCategory[S]))
	toCheck := len(intervalsByCategory[X]) * len(intervalsByCategory[M]) * len(intervalsByCategory[A]) * len(intervalsByCategory[S])
	checked := 0
	for _, x := range intervalsByCategory[X] {
		part[X] = x.Min
		numX := x.Max - x.Min + 1
		for _, m := range intervalsByCategory[M] {
			part[M] = m.Min
			numM := m.Max - m.Min + 1
			for _, a := range intervalsByCategory[A] {
				part[A] = a.Min
				numA := a.Max - a.Min + 1
				for _, s := range intervalsByCategory[S] {
					part[S] = s.Min
					numS := s.Max - s.Min + 1

					num := int64(numX * numM * numA * numS)
					res := ProcessPart(part, inlined)
					if res {
						//fmt.Printf("ACCEPTED: %+v => %d\n", part, num)
						accepted += num
					} else {
						//fmt.Printf("REJECTED: %+v\n", part)
						rejected += num
					}

					checked++
					if checked%10_000_000 == 0 {
						fmt.Printf("Checking %d/%d [%00d%%]\n", checked, toCheck, (checked*100)/toCheck)
					}
				}
			}
		}
	}

	fmt.Printf("Total possible combinations: %d accepted: %d rejected: %d\n", accepted+rejected, accepted, rejected)
}
