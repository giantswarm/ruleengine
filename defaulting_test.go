package ruleengine

import (
	"context"
	"reflect"
	"strconv"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type defaultingMutation struct {
	Value string
}

func Test_ExecuteDefaulting(t *testing.T) {
	testCases := []struct {
		Name             string
		ScopeFunc        func() ([]DefaultingRule, *defaultingMutation)
		ExpectedMutation *defaultingMutation
	}{
		{
			Name: "case 0 ensures that not executing any rule does not mutate the mutation structure",
			ScopeFunc: func() ([]DefaultingRule, *defaultingMutation) {
				return nil, nil
			},
			ExpectedMutation: nil,
		},
		{
			Name: "case 1 ensures that executing one rule does mutate the mutation structure",
			ScopeFunc: func() ([]DefaultingRule, *defaultingMutation) {
				mutation := &defaultingMutation{
					Value: "foo",
				}

				rules := []DefaultingRule{
					{
						Defaulting: func() {
							mutation.Value = "bar"
						},
					},
				}

				return rules, mutation
			},
			ExpectedMutation: &defaultingMutation{
				Value: "bar",
			},
		},
		{
			Name: "case 2 ensures that executing two rules does mutate the mutation structure",
			ScopeFunc: func() ([]DefaultingRule, *defaultingMutation) {
				mutation := &defaultingMutation{
					Value: "foo",
				}

				rules := []DefaultingRule{
					{
						Defaulting: func() {
							mutation.Value = "bar"
						},
					},
					{
						Defaulting: func() {
							mutation.Value = "baz"
						},
					},
				}

				return rules, mutation
			},
			ExpectedMutation: &defaultingMutation{
				Value: "baz",
			},
		},
		{
			Name: "case 3 ensures that one false condition does not mutate the mutation structure",
			ScopeFunc: func() ([]DefaultingRule, *defaultingMutation) {
				mutation := &defaultingMutation{
					Value: "foo",
				}

				rules := []DefaultingRule{
					{
						Conditions: []func() bool{
							func() bool { return false },
						},
						Defaulting: func() {
							mutation.Value = "bar"
						},
					},
				}

				return rules, mutation
			},
			ExpectedMutation: &defaultingMutation{
				Value: "foo",
			},
		},
		{
			Name: "case 4 ensures that two false conditions do not mutate the mutation structure",
			ScopeFunc: func() ([]DefaultingRule, *defaultingMutation) {
				mutation := &defaultingMutation{
					Value: "foo",
				}

				rules := []DefaultingRule{
					{
						Conditions: []func() bool{
							func() bool { return false },
							func() bool { return false },
						},
						Defaulting: func() {
							mutation.Value = "bar"
						},
					},
				}

				return rules, mutation
			},
			ExpectedMutation: &defaultingMutation{
				Value: "foo",
			},
		},
		{
			Name: "case 5 ensures that multiple conditions of which one is true do not mutate the mutation structure",
			ScopeFunc: func() ([]DefaultingRule, *defaultingMutation) {
				mutation := &defaultingMutation{
					Value: "foo",
				}

				rules := []DefaultingRule{
					{
						Conditions: []func() bool{
							func() bool { return false },
							func() bool { return true },
						},
						Defaulting: func() {
							mutation.Value = "bar"
						},
					},
				}

				return rules, mutation
			},
			ExpectedMutation: &defaultingMutation{
				Value: "foo",
			},
		},
		{
			Name: "case 6 ensures that multiple conditions of which all are true do mutate the mutation structure",
			ScopeFunc: func() ([]DefaultingRule, *defaultingMutation) {
				mutation := &defaultingMutation{
					Value: "foo",
				}

				rules := []DefaultingRule{
					{
						Conditions: []func() bool{
							func() bool { return true },
							func() bool { return true },
						},
						Defaulting: func() {
							mutation.Value = "bar"
						},
					},
				}

				return rules, mutation
			},
			ExpectedMutation: &defaultingMutation{
				Value: "bar",
			},
		},
		{
			Name: "case 7 ensures that one true condition does mutate the mutation structure",
			ScopeFunc: func() ([]DefaultingRule, *defaultingMutation) {
				mutation := &defaultingMutation{
					Value: "foo",
				}

				rules := []DefaultingRule{
					{
						Conditions: []func() bool{
							func() bool { return true },
						},
						Defaulting: func() {
							mutation.Value = "bar"
						},
					},
				}

				return rules, mutation
			},
			ExpectedMutation: &defaultingMutation{
				Value: "bar",
			},
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			rules, mutation := tc.ScopeFunc()
			ExecuteDefaulting(context.Background(), rules)
			if !reflect.DeepEqual(mutation, tc.ExpectedMutation) {
				t.Fatalf("want matching\n\n%s\n", cmp.Diff(mutation, tc.ExpectedMutation))
			}
		})
	}
}
