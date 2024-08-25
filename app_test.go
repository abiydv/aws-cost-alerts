package main

import (
	"testing"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/assertions"
	"github.com/aws/jsii-runtime-go"
)

func TestAppStack(t *testing.T) {
	// GIVEN
	app := awscdk.NewApp(nil)

	// WHEN
	stack := NewCostAlertStack(app, "TestCostAlertStack", &CostAlertStackProps{
		Email:    "alert-email@example.com",
		Currency: "USD",
		AnomalyThresholds: map[string][]AnomalyCondition{
			"And": {
				{Dimensions: Dimensions{
					Key:          "ANOMALY_TOTAL_IMPACT_PERCENTAGE",
					Values:       []int{25},
					MatchOptions: []string{GtEq},
				}},
				{Dimensions: Dimensions{
					Key:          "ANOMALY_TOTAL_IMPACT_ABSOLUTE",
					Values:       []int{100},
					MatchOptions: []string{GtEq},
				}},
			}},
		SpendAlertPercentThreshold: 1000,
		StackProps: awscdk.StackProps{
			Env: env(),
		},
	})

	// THEN
	template := assertions.Template_FromStack(stack, nil)
	template.ResourceCountIs(jsii.String("AWS::CE::AnomalyMonitor"), jsii.Number(1))
	template.ResourceCountIs(jsii.String("AWS::CE::AnomalySubscription"), jsii.Number(1))
	template.ResourceCountIs(jsii.String("AWS::Budgets::Budget"), jsii.Number(1))
}
