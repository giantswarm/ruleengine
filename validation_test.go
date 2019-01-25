package ruleengine

import (
	"context"
	"strconv"
	"testing"

	"github.com/giantswarm/microerror"
)

func Test_ExecuteValidation(t *testing.T) {
	testCases := []struct {
		Name         string
		Rules        []ValidationRule
		ErrorMatcher func(err error) bool
	}{
		{
			Name:         "case 0 ensures that not executing any rule does not return an error",
			Rules:        []ValidationRule{},
			ErrorMatcher: nil,
		},
		{
			Name: "case 1 ensures that executing one rule not returning an error does not return an error",
			Rules: []ValidationRule{
				{
					Validation: func() error {
						return nil
					},
				},
			},
			ErrorMatcher: nil,
		},
		{
			Name: "case 2 ensures that executing two rules not returning an error does not return an error",
			Rules: []ValidationRule{
				{
					Validation: func() error {
						return nil
					},
				},
				{
					Validation: func() error {
						return nil
					},
				},
			},
			ErrorMatcher: nil,
		},
		{
			Name: "case 3 ensures that executing one rule returning an error does return an error",
			Rules: []ValidationRule{
				{
					Validation: func() error {
						return microerror.Mask(validationError)
					},
				},
			},
			ErrorMatcher: IsValidation,
		},
		{
			Name: "case 4 ensures that executing two rules returning an error does return an error",
			Rules: []ValidationRule{
				{
					Validation: func() error {
						return microerror.Mask(validationError)
					},
				},
				{
					Validation: func() error {
						return microerror.Mask(validationError)
					},
				},
			},
			ErrorMatcher: IsValidation,
		},
		{
			Name: "case 5 ensures that executing multiple rules of which one returns an error does return an error",
			Rules: []ValidationRule{
				{
					Validation: func() error {
						return nil
					},
				},
				{
					Validation: func() error {
						return microerror.Mask(validationError)
					},
				},
				{
					Validation: func() error {
						return nil
					},
				},
			},
			ErrorMatcher: IsValidation,
		},
		{
			Name: "case 6 ensures that one false condition does not return an error",
			Rules: []ValidationRule{
				{
					Conditions: []func() bool{
						func() bool { return false },
					},
					Validation: func() error {
						return microerror.Mask(validationError)
					},
				},
			},
			ErrorMatcher: nil,
		},
		{
			Name: "case 7 ensures that two false conditions do not return an error",
			Rules: []ValidationRule{
				{
					Conditions: []func() bool{
						func() bool { return false },
						func() bool { return false },
					},
					Validation: func() error {
						return microerror.Mask(validationError)
					},
				},
			},
			ErrorMatcher: nil,
		},
		{
			Name: "case 8 ensures that one true condition does not return an error",
			Rules: []ValidationRule{
				{
					Conditions: []func() bool{
						func() bool { return false },
						func() bool { return true },
					},
					Validation: func() error {
						return microerror.Mask(validationError)
					},
				},
			},
			ErrorMatcher: nil,
		},
		{
			Name: "case 9 ensures that one true condition does return an error",
			Rules: []ValidationRule{
				{
					Conditions: []func() bool{
						func() bool { return true },
					},
					Validation: func() error {
						return microerror.Mask(validationError)
					},
				},
			},
			ErrorMatcher: IsValidation,
		},
		{
			Name: "case 10 ensures that two true conditions do return an error",
			Rules: []ValidationRule{
				{
					Conditions: []func() bool{
						func() bool { return true },
						func() bool { return true },
					},
					Validation: func() error {
						return microerror.Mask(validationError)
					},
				},
			},
			ErrorMatcher: IsValidation,
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			err := ExecuteValidation(context.Background(), tc.Rules)

			switch {
			case err == nil && tc.ErrorMatcher == nil:
				// correct; carry on
			case err != nil && tc.ErrorMatcher == nil:
				t.Fatalf("error == %#v, want nil", err)
			case err == nil && tc.ErrorMatcher != nil:
				t.Fatalf("error == nil, want non-nil")
			case !tc.ErrorMatcher(err):
				t.Fatalf("error == %#v, want matching", err)
			}
		})
	}
}
