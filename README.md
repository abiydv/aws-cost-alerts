# AWS Cost Alerts

This project sets up following resources in an AWS account to keep on top of unexpected costs.

- `AWS::CE::AnomalyMonitor`
- `AWS::CE::AnomalySubscription`
- `AWS::Budgets::Budget`

## Alerts

To modify the alert configs, please refer to the [app.go](./app.go) file. 

### Anomaly Monitor Alert

Once deployed, Anomaly monitor watches the per service spend and sends an email alert when both following conditions are met - 

 * Cost impact is 25% or higher than usual
 * Cost impact is USD 100 or more

### AWS Budgets Alert

Once deployed, Budget watches the total spend of an account, and sends an email alert when the spend is greater than 90% of the configured budget (USD 1000).

## CDK 

The `cdk.json` file tells the CDK Toolkit how to execute your app.

## Useful commands

 * `cdk deploy`      deploy this stack to your default AWS account/region
 * `cdk diff`        compare deployed stack with current state
 * `cdk synth`       emits the synthesized CloudFormation template
 * `go test`         run unit tests
