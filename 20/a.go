package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"strings"
)

var debug = false

type ModuleType uint8

const (
	Broadcast ModuleType = iota
	Conjunction
	FlipFlop
	RX
)

var pulseStrMap = map[bool]string{true: "high", false: "low"}
var typeStrMap = map[ModuleType]string{Broadcast: "BR", Conjunction: "CJ", FlipFlop: "FF", RX: "RX"}

type Module struct {
	Index                 uint8
	Station               string
	Type                  ModuleType
	Destinations          []string
	Sources               []string
	numHighPulsesReceived int
	numLowPulsesReceived  int
	SourceMask            uint64
}

type Pulse struct {
	High      bool
	FromIndex uint8
	From      string
	To        []string
}

type BroadCastSystem struct {
	modules           map[string]*Module
	queue             []Pulse
	state             []*uint64
	lastPulse         map[uint8]bool
	numHighPulsesSent int
	numLowPulsesSent  int
	numButtonPresses  int
}

func (m *BroadCastSystem) AddModule(module *Module) {
	m.modules[module.Station] = &Module{
		Station:      module.Station,
		Type:         module.Type,
		Destinations: module.Destinations,
	}
}

func (m *BroadCastSystem) ProcessButtonPress() {
	//fmt.Printf("%+v\n", m.state)
	m.numButtonPresses++
	current := []*Pulse{{false, 0, "button", []string{"broadcaster"}}}
	for len(current) > 0 {
		next := []*Pulse{}
		for _, pulse := range current {
			for _, to := range pulse.To {
				if debug {
					fmt.Printf("%s -%s-> %s\n", pulse.From, pulseStrMap[pulse.High], to)
				}

				if pulse.High {
					m.numHighPulsesSent++
				} else {
					m.numLowPulsesSent++
				}
				if _, ok := m.modules[to]; ok {
					toIndex := m.modules[to].Index
					output := m.modules[to].Process(m.numButtonPresses, pulse, m.state[toIndex])
					if output != nil {
						m.lastPulse[toIndex] = output.High
						next = append(next, output)
					} else if _, pulseOk := m.lastPulse[toIndex]; pulseOk {
						delete(m.lastPulse, toIndex)
					}
				}
			}
		}
		current = next

		if debug {
			fmt.Printf("----------------------------\n")
		}
	}
}

func (m *BroadCastSystem) stateHash() string {
	checksum := sha256.New()
	b := make([]byte, 8)
	for _, state := range m.state {
		binary.LittleEndian.PutUint64(b, *state)
		checksum.Write(b)
	}
	return fmt.Sprintf("%x", checksum.Sum(nil))
}

func (m *BroadCastSystem) FindCycle() (int, int) {
	seen := map[string]int{m.stateHash(): 0}
	cycle := 0
	for {
		cycle++
		m.ProcessButtonPress()
		after := m.stateHash()
		if _, ok := seen[after]; ok {
			return seen[after], cycle - seen[after]
		} else {
			seen[after] = cycle
		}
	}
}

func (m *BroadCastSystem) Reset() {
	m.numLowPulsesSent = 0
	m.numHighPulsesSent = 0
	m.numButtonPresses = 0
	m.queue = []Pulse{}
	m.state = make([]*uint64, len(m.modules))
	m.lastPulse = map[uint8]bool{}
	index := uint8(0)
	for _, module := range m.modules {
		module.Index = index
		module.Sources = []string{}
		module.numHighPulsesReceived = 0
		module.numLowPulsesReceived = 0
		module.SourceMask = 0
		m.state[module.Index] = new(uint64)
		*m.state[module.Index] = 0
		index++
	}
	for _, module := range m.modules {
		for _, dest := range module.Destinations {
			if _, ok := m.modules[dest]; ok {
				m.modules[dest].Sources = append(m.modules[dest].Sources, module.Station)
				m.modules[dest].SourceMask |= 1 << module.Index
			}
		}
	}
}

func (m *BroadCastSystem) GetAllSources(station string) []*Module {
	res := []*Module{m.modules[station]}

	stack := []string{station}
	processed := map[string]bool{}
	var current string
	for len(stack) > 0 {
		current, stack = stack[0], stack[1:]
		for _, source := range m.modules[current].Sources {
			if !processed[source] {
				stack = append(stack, source)
				res = append(res, m.modules[source])
				processed[source] = true
			}
		}
	}

	return res
}

func (m *BroadCastSystem) Stats() (int, int) {
	return m.numHighPulsesSent, m.numLowPulsesSent
}

func (m *BroadCastSystem) Print() {
	for _, module := range m.modules {
		fmt.Printf("%+v\n", module)
	}
}

func newBroadcastSystem() *BroadCastSystem {
	return &BroadCastSystem{
		modules: map[string]*Module{},
	}
}

func (m *Module) Process(num int, pulse *Pulse, state *uint64) *Pulse {
	output := false
	if pulse.High {
		m.numHighPulsesReceived++
	} else {
		m.numLowPulsesReceived++
	}

	switch m.Type {
	case Broadcast:
		output = pulse.High
	case Conjunction:
		mask := uint64(1 << pulse.FromIndex)
		if pulse.High {
			*state |= mask
		} else {
			*state &= ^mask
		}
		if *state != m.SourceMask {
			output = true
		}
	case FlipFlop:
		if !pulse.High {
			*state = (*state + 1) % 2
			output = *state == 1
		} else {
			return nil
		}
	case RX:
		if pulse.High {
			fmt.Printf("[%d] RX received: %s\n", num, pulseStrMap[pulse.High])
		}
		return nil
	}

	return &Pulse{output, m.Index, m.Station, m.Destinations}
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	broadCastSystem := newBroadcastSystem()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.Split(strings.ReplaceAll(scanner.Text(), ",", ""), " ")
		station, destinations := line[0], line[2:]
		moduleType := Broadcast
		if station[0] == '&' {
			moduleType = Conjunction
			station = station[1:]
		} else if station[0] == '%' {
			moduleType = FlipFlop
			station = station[1:]
		}
		broadCastSystem.AddModule(&Module{
			Station:      station,
			Type:         moduleType,
			Destinations: destinations,
		})
	}
	broadCastSystem.Reset()
	broadCastSystem.Print()

	for i := 0; i < 1000; i++ {
		broadCastSystem.ProcessButtonPress()
	}

	counterHigh, counterLow := broadCastSystem.Stats()
	fmt.Printf("Delivered %d high pulses and %d low pulses => %d\n", counterHigh, counterLow, counterHigh*counterLow)

	rxModule := &Module{
		Station:      "rx",
		Type:         RX,
		Destinations: []string{},
	}

	broadCastSystem.AddModule(rxModule)
	broadCastSystem.Reset()

	numPushes := 1
	joint := "bq"
	for _, sub := range broadCastSystem.modules[joint].Sources {
		sources := broadCastSystem.GetAllSources(sub)
		subSystem := newBroadcastSystem()
		subSystem.AddModule(broadCastSystem.modules[joint])
		subSystem.modules[joint].Type = RX
		for _, source := range sources {
			subSystem.AddModule(source)
		}
		subSystem.Reset()

		subIndex := subSystem.modules[sub].Index
		jointIndex := subSystem.modules[joint].Index

		offset, cycle := subSystem.FindCycle()
		fmt.Printf("[%s %d] found cycle length of %d (with offset %d) => %d %d %+v\n", sub, subIndex, cycle, offset, *subSystem.state[jointIndex], subSystem.modules[joint].SourceMask, subSystem.modules[joint].Sources)
		numPushes *= cycle
	}
	fmt.Printf("It should take %d button presses to turn on\n", numPushes)
}
