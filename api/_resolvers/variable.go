package resolvers

import (
	"context"
	"errors"
	"fmt"
	"sort"

	"github.com/mjm/graphql-go"
	"github.com/mjm/graphql-go/relay"

	"github.com/mjm/speedrungql/speedrun"
)

type Variable struct {
	speedrun.Variable
	client *speedrun.Client
}

func (v *Variable) ID() graphql.ID {
	return relay.MarshalID("variable", v.Variable.ID)
}

func (v *Variable) RawID() string {
	return v.Variable.ID
}

func (v *Variable) Game(ctx context.Context) (*Game, error) {
	gameURI := speedrun.FindLink(v.Links, "game")
	if gameURI == "" {
		return nil, nil
	}

	game, err := v.client.GetGame(ctx, gameURI)
	if err != nil {
		return nil, err
	}

	if game == nil {
		return nil, nil
	}
	return &Game{*game, v.client}, nil
}

func (v *Variable) Category(ctx context.Context) (*Category, error) {
	if v.CategoryID == "" {
		return nil, nil
	}

	c, err := v.client.GetCategory(ctx, v.CategoryID)
	if err != nil {
		return nil, err
	}

	if c == nil {
		return nil, nil
	}

	return &Category{*c, v.client}, nil
}

func (v *Variable) Scope() VariableScopeType {
	return VariableScopeType(v.Variable.Scope.Type)
}

func (v *Variable) Values() []*VariableValue {
	var vals []*VariableValue
	for valID, val := range v.Variable.Values.Values {
		vals = append(vals, &VariableValue{val, valID, v})
	}

	sort.Slice(vals, func(i, j int) bool {
		return vals[i].id < vals[j].id
	})
	return vals
}

func (v *Variable) DefaultValue() *VariableValue {
	valID := v.Variable.Values.Default
	if valID == "" {
		return nil
	}

	val, ok := v.Variable.Values.Values[valID]
	if !ok {
		return nil
	}

	return &VariableValue{val, valID, v}
}

func (v *Variable) Value(args struct {
	ID graphql.ID
}) *VariableValue {
	valID := string(args.ID)
	val, ok := v.Variable.Values.Values[valID]
	if !ok {
		return nil
	}

	return &VariableValue{val, valID, v}
}

type VariableScopeType speedrun.VariableScopeType

func (VariableScopeType) ImplementsGraphQLType(name string) bool {
	return name == "VariableScopeType"
}

func (v VariableScopeType) String() string {
	switch speedrun.VariableScopeType(v) {
	case speedrun.ScopeGlobal:
		return "GLOBAL"
	case speedrun.ScopeFullGame:
		return "FULL_GAME"
	case speedrun.ScopeAllLevels:
		return "ALL_LEVELS"
	case speedrun.ScopeSingleLevel:
		return "SINGLE_LEVEL"
	default:
		return ""
	}
}

func (v *VariableScopeType) UnmarshalGraphQLType(input interface{}) error {
	s, ok := input.(string)
	if !ok {
		return errors.New("VariableScopeType value was not a string")
	}

	switch s {
	case "GLOBAL":
		*v = VariableScopeType(speedrun.ScopeGlobal)
	case "FULL_GAME":
		*v = VariableScopeType(speedrun.ScopeFullGame)
	case "ALL_LEVELS":
		*v = VariableScopeType(speedrun.ScopeAllLevels)
	case "SINGLE_LEVEL":
		*v = VariableScopeType(speedrun.ScopeSingleLevel)
	default:
		return fmt.Errorf("unknown VariableScopeType value %q", s)
	}

	return nil
}

type VariableValue struct {
	speedrun.VariableValue
	id       string
	variable *Variable
}

func (v *VariableValue) ID() graphql.ID {
	return graphql.ID(v.id)
}

func (v *VariableValue) Variable() *Variable {
	return v.variable
}

func (v *VariableValue) Rules() *string {
	if v.VariableValue.Rules == "" {
		return nil
	}
	return &v.VariableValue.Rules
}
