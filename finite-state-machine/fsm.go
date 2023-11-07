package main

import (
	"errors"
	"fmt"
)

type Transition struct {
	from  string
	to    string
	event string
}

type StateMachine struct {
	state       string
	transitions []Transition
	handleEvent func(fromState, toState string, args []interface{}) error
}

func NewStateMachine(init string, transitions []Transition, handleEvent func(fromState, toState string, args []interface{}) error) *StateMachine {
	return &StateMachine{
		state:       init,
		transitions: transitions,
		handleEvent: handleEvent,
	}
}

func (m *StateMachine) changeState(state string) {
	m.state = state
}

func (m *StateMachine) findTransMatching(fromState, event string) *Transition {
	for _, v := range m.transitions {
		if v.from == fromState && v.event == event {
			return &v
		}
	}

	return nil
}

func (m *StateMachine) Trigger(event string, args ...interface{}) error {
	trans := m.findTransMatching(m.state, event)
	if trans == nil {
		return errors.New("状态转换失败：未找到 trans")
	}

	if trans.event == "" {
		return errors.New("未找到具体的 event")
	}

	err := m.handleEvent(m.state, trans.to, args)
	if err != nil {
		return errors.New("状态转换失败：未找到 handleEvent")
	}

	m.changeState(trans.to)
	return nil
}

func main() {
	transitions := make([]Transition, 0)
	transitions = append(transitions, Transition{
		from:  "create",
		to:    "running",
		event: "start",
	})
	transitions = append(transitions, Transition{
		from:  "create",
		to:    "end",
		event: "work",
	})

	fsm := NewStateMachine("create", transitions, func(fromState, toState string, args []interface{}) error {
		switch toState {
		case "end":
			fmt.Println("工作结束")
		}
		return nil
	})

	fsm.Trigger("start")
	fsm.Trigger("work")

	fmt.Println(fsm.state)
}
