package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type PrayerJournalStackProps struct {
	awscdk.StackProps
}

func NewPrayerJournalStack(scope constructs.Construct, id string, props *PrayerJournalStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	// TODO: Add S3, CloudFront, and Lambda here

	return stack
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	NewPrayerJournalStack(app, "PrayerJournalStack", &PrayerJournalStackProps{
		StackProps: awscdk.StackProps{
			Env: &awscdk.Environment{
				Region: jsii.String("us-east-1"),
			},
		},
	})

	app.Synth(nil)
}
