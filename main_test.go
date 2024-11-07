package main

import (
	"context"
	"flag"
	"os"
	"testing"

	componenttest "github.com/ONSdigital/dp-component-test"
	"github.com/ONSdigital/dp-identity-api/features/steps"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
)

var componentFlag = flag.Bool("component", false, "perform component tests")

type ComponentTest struct {
	MongoFeature *componenttest.MongoFeature
}

func (f *ComponentTest) InitializeScenario(godogCtx *godog.ScenarioContext) {
	component, err := steps.NewIdentityComponent()
	if err != nil {
		log.Fatal(context.Background(), "fatal error initialising a test scenario", err)
	}

	godogCtx.Before(func(ctx context.Context, _ *godog.Scenario) (context.Context, error) {
		component.Reset()
		return ctx, nil
	})

	godogCtx.After(func(ctx context.Context, _ *godog.Scenario, _ error) (context.Context, error) {
		if err := component.Close(); err != nil {
			log.Warn(context.Background(), "error closing identity component", log.FormatErrors([]error{err}))
		}
		return ctx, nil
	})

	component.RegisterSteps(godogCtx)
}

func (f *ComponentTest) InitializeTestSuite(_ *godog.TestSuiteContext) {}

func TestComponent(t *testing.T) {
	if *componentFlag {
		status := 0

		var opts = godog.Options{
			Output: colors.Colored(os.Stdout),
			Format: "pretty",
			Paths:  flag.Args(),
		}

		f := &ComponentTest{}

		status = godog.TestSuite{
			Name:                 "feature_tests",
			ScenarioInitializer:  f.InitializeScenario,
			TestSuiteInitializer: f.InitializeTestSuite,
			Options:              &opts,
		}.Run()

		if status > 0 {
			t.Fail()
		}
	} else {
		t.Skip("component flag required to run component tests")
	}
}
