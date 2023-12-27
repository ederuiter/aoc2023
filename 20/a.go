package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

var debug = false

type ModuleType int

const (
	Broadcast ModuleType = iota
	Conjunction
	FlipFlop
	RX
)

var pulseStrMap = map[bool]string{true: "-high-", false: "-low-"}
var typeStrMap = map[ModuleType]string{Broadcast: "BR", Conjunction: "CJ", FlipFlop: "FF", RX: "RX"}

type Module struct {
	Station               string
	Type                  ModuleType
	Destinations          []string
	Sources               []string
	State                 map[string]bool
	numHighPulsesReceived int
	numLowPulsesReceived  int
}

type Pulse struct {
	High bool
	From string
	To   []string
}

type BroadCastSystem struct {
	modules           map[string]*Module
	queue             []Pulse
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
	m.numButtonPresses++
	current := []*Pulse{{false, "button", []string{"broadcaster"}}}
	for len(current) > 0 {
		next := []*Pulse{}
		numPulses := 0
		for _, pulse := range current {
			if pulse.High {
				m.numHighPulsesSent++
			} else {
				m.numLowPulsesSent++
			}

			for _, to := range pulse.To {
				output := m.modules[to].Process(pulse)
				if output != nil {
					next = append(next, output)
					numPulses += len(output.To)
				}
			}
		}
		current = next

		if debug {
			fmt.Printf("----------------------------\n")
		}
	}
}

func (m *BroadCastSystem) FindCycle(outputModule string, output bool) int {
	//TODO: how do we detect a cycle
}

func (m *BroadCastSystem) Reset() {
	m.numLowPulsesSent = 0
	m.numHighPulsesSent = 0
	m.numButtonPresses = 0
	m.queue = []Pulse{}
	for _, module := range m.modules {
		module.State = map[string]bool{}
		module.Sources = []string{}
		module.numHighPulsesReceived = 0
		module.numLowPulsesReceived = 0
	}
	for _, module := range m.modules {
		for _, dest := range module.Destinations {
			if _, ok := m.modules[dest]; ok {
				m.modules[dest].Sources = append(m.modules[dest].Sources, module.Station)
			}
		}
	}
}

func (m *BroadCastSystem) GetAllSources(station string) []*Module {
	res := []*Module{}

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

func (m *Module) Process(pulse *Pulse) *Pulse {
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
		m.State[pulse.From] = pulse.High
		if m.State["&&"] {
			m.State["&&"] = pulse.High
			output = !m.State["&&"]
		} else if !pulse.High {
			output = true
			m.State["&&"] = false
		} else {
			for _, source := range m.Sources {
				val, ok := m.State[source]
				if !ok || !val {
					output = true
					break
				}
			}
			m.State["&&"] = !output
		}

	case FlipFlop:
		triggered := !pulse.High
		m.State["triggered"] = triggered
		if triggered {
			m.State["ff"] = !m.State["ff"]
			output = m.State["ff"]
		} else {
			return nil
		}
	case RX:
		return nil
	}

	return &Pulse{output, m.Station, m.Destinations}
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

	//for i := 0; i < 1000; i++ {
	//	broadCastSystem.PressButton()
	//	broadCastSystem.Process()
	//}

	//counterHigh, counterLow := broadCastSystem.Stats()
	//fmt.Printf("Delivered %d high pulses and %d low pulses => %d\n", counterHigh, counterLow, counterHigh*counterLow)

	rxModule := &Module{
		Station:      "rx",
		Type:         RX,
		Destinations: []string{},
	}

	broadCastSystem.AddModule(rxModule)
	broadCastSystem.Reset()

	numPushes := 1
	for _, sub := range broadCastSystem.modules["bq"].Sources {
		sources := broadCastSystem.GetAllSources(sub)
		subSystem := newBroadcastSystem()
		for _, source := range sources {
			subSystem.AddModule(source)
		}
		subSystem.Reset()
		numPushes *= subSystem.FindCycle(sub, true)
	}
	fmt.Printf("It should take %d button presses to turn on\n", numPushes)

	//for rxModule.numLowPulsesReceived != 1 {
	//	rxModule.numLowPulsesReceived = 0
	//	rxModule.numHighPulsesReceived = 0
	//	broadCastSystem.ProcessButtonPress()
	//
	//	if broadCastSystem.numButtonPresses%100_000 == 0 {
	//		fmt.Printf("[%d] not yet ..\n", broadCastSystem.numButtonPresses)
	//	}
	//}
	//
	//fmt.Printf("It took %d button presses to turn on\n", broadCastSystem.numButtonPresses)
}
