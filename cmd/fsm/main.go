package main

import (
	"errors"
	"fmt"
	"strconv"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/multi"
)

type Event string
type Operator string

func NewRule(trigger Operator, eventVal Event) map[Operator]Event {
	return map[Operator]Event{trigger: eventVal}
}

var stateIdCntr = 0
var linkIdCntr = 1

type State struct {
	id    int64
	value any
}

func NewState(val any) State {
	state := State{id: int64(stateIdCntr), value: val}
	stateIdCntr++
	return state
}

type Link struct {
	id       int64
	from, to graph.Node
	rules    map[Operator]Event
}

func NewLink(from, to State, rules map[Operator]Event) Link {
	link := Link{id: int64(linkIdCntr), from: from, to: to, rules: rules}
	linkIdCntr++
	return link
}

func (n State) ID() int64 {
	return n.id
}

func (l Link) From() graph.Node {
	return l.from
}

func (l Link) To() graph.Node {
	return l.to
}

func (l Link) ID() int64 {
	return l.id
}

func (l Link) ReversedLine() graph.Line {
	return Link{from: l.to, to: l.from}
}

func (n State) String() string {
	switch n.value.(type) {
	case int:
		return strconv.Itoa(n.value.(int))
	case float32:
		return fmt.Sprintf("%f", n.value.(float32))
	case float64:
		return fmt.Sprintf("%f", n.value.(float64))
	case bool:
		return strconv.FormatBool(n.value.(bool))
	case string:
		return n.value.(string)
	default:
		return ""
	}
}

type StateMachine struct {
	CurrentState State
	g            *multi.DirectedGraph
}

func NewStateMachine() *StateMachine {
	s := &StateMachine{}
	s.g = multi.NewDirectedGraph()
	return s
}

func (s *StateMachine) Init(val any) State {
	s.CurrentState = NewState(val)
	s.g.AddNode(s.CurrentState)
	return s.CurrentState
}

func (s *StateMachine) MakeState(stateVal any) State {
	state := NewState(stateVal)
	s.g.AddNode(state)
	return state
}

func (s *StateMachine) LinkStates(s1, s2 State, rule map[Operator]Event) {
	s.g.SetLine(NewLink(s1, s2, rule))
}

func (s *StateMachine) FireEvent(e Event) error {
	curNode := s.CurrentState
	it := s.g.From(curNode.id)
	for it.Next() {
		n := s.g.Node(it.Node().ID()).(State)
		line := graph.LinesOf(s.g.Lines(curNode.id, n.id))[0].(Link)

		for key, val := range line.rules {
			k := string(key)
			switch k {
			case "eq":
				if val == e {
					s.CurrentState = n
					return nil
				}
			default:
				fmt.Printf("Sorry, the comparison operator '%s' is not supported\n", k)
				return errors.New("UNSUPPORTED_COMPARISON_OPERATOR")
			}
		}
	}
	return nil
}

func (s *StateMachine) Compute(events []string, printState bool) State {
	for _, e := range events {
		s.FireEvent(Event(e))
		if printState {
			fmt.Printf("%s\n", s.CurrentState.String())
		}
	}
	return s.CurrentState
}

func main() {
	sm := NewStateMachine()

	lockedState := sm.Init("locked")
	unlockedState := sm.MakeState("unlocked")

	coinRule := NewRule(Operator("eq"), Event("coin"))
	pushRule := NewRule(Operator("eq"), Event("push"))

	sm.LinkStates(lockedState, unlockedState, coinRule)
	sm.LinkStates(unlockedState, lockedState, pushRule)

	sm.LinkStates(lockedState, lockedState, pushRule)
	sm.LinkStates(unlockedState, unlockedState, coinRule)
	fmt.Printf("Initial state --- %s\n", sm.CurrentState.String())

	events := []string{"coin", "push"}
	sm.Compute(events, true)

	fmt.Printf("Final state --- %s\n", sm.CurrentState.String())
}
