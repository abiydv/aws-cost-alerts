package main

import (
	"encoding/json"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsbudgets"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsce"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type CostAlertStackProps struct {
	Email                      string
	Currency                   string
	AnomalyThresholds          map[string][]AnomalyCondition
	SpendAmount                float64
	SpendAlertPercentThreshold int
	awscdk.StackProps
}

type AppStackProps struct {
	awscdk.StackProps
}

type Dimensions struct {
	Key          string   `json:"Key"`
	Values       []int    `json:"Values"`
	MatchOptions []string `json:"MatchOptions"`
}
type AnomalyCondition struct {
	Dimensions Dimensions `json:"Dimensions"`
}

const GtEq string = "GREATER_THAN_OR_EQUAL"
const Gt string = "GREATER_THAN"

func NewCostAlertStack(scope constructs.Construct, id string, props *CostAlertStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)
	// AWS Cost Anomaly Monitor
	svcAnomalyMonitor := awsce.NewCfnAnomalyMonitor(stack, jsii.String("ServiceAnomalyMonitor"), &awsce.CfnAnomalyMonitorProps{
		MonitorName:      jsii.String("ServiceAnomalyMonitor"),
		MonitorType:      jsii.String("DIMENSIONAL"),
		MonitorDimension: jsii.String("SERVICE"),
	})
	anomalyThresholdExpr, _ := json.Marshal(props.AnomalyThresholds)
	awsce.NewCfnAnomalySubscription(stack, jsii.String("ServiceAnomalyMonitorEmailAlert"), &awsce.CfnAnomalySubscriptionProps{
		Frequency:      jsii.String("DAILY"),
		MonitorArnList: jsii.Strings(*svcAnomalyMonitor.AttrMonitorArn()),
		Subscribers: []awsce.CfnAnomalySubscription_SubscriberProperty{
			{
				Address: &props.Email,
				Type:    jsii.String("EMAIL"),
			},
		},
		SubscriptionName:    jsii.String("ServiceAnomalyMonitorEmailAlert"),
		ThresholdExpression: jsii.String(string(anomalyThresholdExpr)),
	})
	// AWS Budget
	awsbudgets.NewCfnBudget(stack, jsii.String("CostBudget"), &awsbudgets.CfnBudgetProps{
		Budget: awsbudgets.CfnBudget_BudgetDataProperty{
			BudgetType: jsii.String("COST"),
			TimeUnit:   jsii.String("MONTHLY"),
			BudgetLimit: awsbudgets.CfnBudget_SpendProperty{
				Amount: &props.SpendAmount,
				Unit:   jsii.String("USD"),
			},
			BudgetName: jsii.String("CostBudget"),
		},
		NotificationsWithSubscribers: []awsbudgets.CfnBudget_NotificationWithSubscribersProperty{
			{
				Notification: awsbudgets.CfnBudget_NotificationProperty{
					ComparisonOperator: jsii.String("EMAIL"),
					NotificationType:   jsii.String("ACTUAL"),
					Threshold:          jsii.Number(props.SpendAlertPercentThreshold),
					ThresholdType:      jsii.String("PERCENTAGE"),
				},
				Subscribers: []awsbudgets.CfnBudget_SubscriberProperty{
					{
						Address:          &props.Email,
						SubscriptionType: jsii.String("EMAIL"),
					},
				},
			},
		},
	})
	return stack
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	NewCostAlertStack(app, "CostAlertStack", &CostAlertStackProps{
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

	app.Synth(nil)
}

// env determines the AWS environment (account+region) in which our stack is to
// be deployed. For more information see: https://docs.aws.amazon.com/cdk/latest/guide/environments.html
func env() *awscdk.Environment {
	return nil
}
