package servers

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/betas-in/utils"
)

type Match struct {
	Attribute string `hcl:"attribute"`
	Operator  string `hcl:"operator"`
	Value     string `hcl:"value"`
}

func NewMatch(attribute, operator, value string) *Match {
	return &Match{
		Attribute: attribute,
		Operator:  operator,
		Value:     value,
	}
}

func (m *Match) CompareStringArray(list []string) (bool, error) {
	contains := utils.Array().Contains(list, m.Value, true)
	switch m.Operator {
	case "contains":
		if contains > -1 {
			return true, nil
		}
	case "not-contains":
		if contains == -1 {
			return true, nil
		}
	default:
		return false, fmt.Errorf("only ['contains', 'not-contains'] operators are supported")
	}
	return false, nil
}

func (m *Match) CompareString(name string) (bool, error) {
	switch m.Operator {
	case "=":
		if name == m.Value {
			return true, nil
		}
	case "!=":
		if name != m.Value {
			return true, nil
		}
	case "like":
		if strings.Contains(name, m.Value) {
			return true, nil
		}
	default:
		// #TODO Support nil match type
		return false, fmt.Errorf("only ['=', '!=', 'like'] operators are supported")
	}
	return false, nil
}

func (m *Match) CompareInt(name int) (bool, error) {
	valueInt, err := strconv.Atoi(m.Value)
	if err != nil {
		return false, nil
	}
	switch m.Operator {
	case "=":
		if name == valueInt {
			return true, nil
		}
	case "!=":
		if name != valueInt {
			return true, nil
		}
	case ">":
		if name > valueInt {
			return true, nil
		}
	case ">=":
		if name >= valueInt {
			return true, nil
		}
	case "<":
		if name < valueInt {
			return true, nil
		}
	case "<=":
		if name <= valueInt {
			return true, nil
		}
	default:
		return false, fmt.Errorf("only ['=', '!=', '>', '>=', '<', '<='] operators are supported")
	}
	return false, nil
}
