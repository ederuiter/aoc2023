package main

import (
	"bufio"
	"cmp"
	"fmt"
	"log"
	"maps"
	"os"
	"slices"
	"strconv"
	"strings"
)

type Kind int

const ( // iota is reset to 0
	KIND_UNKNOWN Kind = iota
	KIND_HIGH
	KIND_ONE_PAIR
	KIND_TWO_PAIR
	KIND_THREE
	KIND_FULLHOUSE
	KIND_FOUR
	KIND_FIVE
)

type Hand struct {
	cards string
	bid   int64
	kind  Kind
}

var strength = map[string]int{
	"A": 14,
	"K": 13,
	"Q": 12,
	"T": 10,
	"9": 9,
	"8": 8,
	"7": 7,
	"6": 6,
	"5": 5,
	"4": 4,
	"3": 3,
	"2": 2,
	"J": 1,
}

func getKind(cardMap map[uint8]int) Kind {
	/*
			5
			41
			32
		    311
			221
			2111
			11111
	*/

	var kind = KIND_UNKNOWN
	l := len(cardMap)
	if l == 1 {
		kind = KIND_FIVE
	} else if l == 2 {
		for _, num := range cardMap {
			if num == 4 || num == 1 {
				kind = KIND_FOUR
			} else if num == 3 || num == 2 {
				kind = KIND_FULLHOUSE
			} else {
				panic(num)
			}
			break
		}
	} else if l == 3 {
		for _, num := range cardMap {
			if num == 3 {
				kind = KIND_THREE
				break
			} else if num == 2 {
				kind = KIND_TWO_PAIR
				break
			}
		}
	} else if l == 4 {
		kind = KIND_ONE_PAIR
	} else if l == 5 {
		kind = KIND_HIGH
	}
	return kind
}

func getHand(cards string, bid int64) Hand {
	cardMap := map[uint8]int{}
	jokers := 0
	for i := 0; i < len(cards); i++ {
		if cards[i] == 'J' {
			jokers = jokers + 1
		} else {
			cardMap[cards[i]] = cardMap[cards[i]] + 1
		}
	}

	var kind = KIND_UNKNOWN
	if jokers == 5 {
		kind = KIND_FIVE
	} else {
		for card, _ := range cardMap {
			newCardMap := map[uint8]int{}
			maps.Copy(newCardMap, cardMap)
			newCardMap[card] = newCardMap[card] + jokers
			kind = max(kind, getKind(newCardMap))
		}
	}

	return Hand{
		cards: cards,
		bid:   bid,
		kind:  kind,
	}
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	hands := []Hand{}
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), " ")
		cards := parts[0]
		bid, _ := strconv.ParseInt(parts[1], 10, 64)
		hands = append(hands, getHand(cards, bid))
	}

	slices.SortFunc(hands, func(a Hand, b Hand) int {
		res := cmp.Compare(a.kind, b.kind)
		index := 0
		for res == 0 {
			res = cmp.Compare(strength[a.cards[index:index+1]], strength[b.cards[index:index+1]])
			index = index + 1
			if index == len(a.cards) {
				break
			}
		}
		return res
	})

	score := int64(0)
	for i, hand := range hands {
		score = score + (int64(i+1) * hand.bid)
	}
	fmt.Printf("Score: %d\n", score)
}
