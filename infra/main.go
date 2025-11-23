package main

import (
	"os"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigatewayv2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigatewayv2integrations"
	"github.com/aws/aws-cdk-go/awscdk/v2/awscertificatemanager"
	"github.com/aws/aws-cdk-go/awscdk/v2/awscloudfront"
	"github.com/aws/aws-cdk-go/awscdk/v2/awscloudfrontorigins"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsroute53"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsroute53targets"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3assets"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3deployment"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

const (
	DomainName   = "faithforge.academy"
	AppSubdomain = "prayer." + DomainName
	ApiSubdomain = "prayerapi." + DomainName
	TableName    = "prayer-journal-data"
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

	// -----------------------------------------------------------------------
	// 0. Context & Lookups
	// -----------------------------------------------------------------------

	// Verify Auth0 Env Vars are present for the Lambda
	auth0Issuer := os.Getenv("AUTH0_ISSUER")     // e.g. https://your-tenant.us.auth0.com/
	auth0Audience := os.Getenv("AUTH0_AUDIENCE") // e.g. https://prayerapi.faithforge.academy
	if auth0Issuer == "" || auth0Audience == "" {
		panic("❌ Missing Environment Variables: AUTH0_ISSUER or AUTH0_AUDIENCE must be set for CDK deploy")
	}

	// Lookup Hosted Zone
	zone := awsroute53.HostedZone_FromLookup(stack, jsii.String("Zone"), &awsroute53.HostedZoneProviderProps{
		DomainName: jsii.String(DomainName),
	})

	// Create HTTPS Certificate (applies to *.faithforge.academy)
	cert := awscertificatemanager.NewCertificate(stack, jsii.String("SiteCert"), &awscertificatemanager.CertificateProps{
		DomainName: jsii.String(DomainName),
		SubjectAlternativeNames: jsii.Strings(
			AppSubdomain,
			ApiSubdomain,
		),
		Validation: awscertificatemanager.CertificateValidation_FromDns(zone),
	})

	// -----------------------------------------------------------------------
	// 1. Database (DynamoDB)
	// -----------------------------------------------------------------------
	table := awsdynamodb.NewTable(stack, jsii.String("PrayerTable"), &awsdynamodb.TableProps{
		TableName:     jsii.String(TableName),
		PartitionKey:  &awsdynamodb.Attribute{Name: jsii.String("pk"), Type: awsdynamodb.AttributeType_STRING},
		SortKey:       &awsdynamodb.Attribute{Name: jsii.String("sk"), Type: awsdynamodb.AttributeType_STRING},
		BillingMode:   awsdynamodb.BillingMode_PAY_PER_REQUEST,
		RemovalPolicy: awscdk.RemovalPolicy_DESTROY, // Change to RETAIN for production
	})

	// -----------------------------------------------------------------------
	// 2. Backend (Lambda + API Gateway)
	// -----------------------------------------------------------------------

	// Go Lambda Function
	backendFunc := awslambda.NewFunction(stack, jsii.String("APIHandler"), &awslambda.FunctionProps{
		Runtime: awslambda.Runtime_PROVIDED_AL2023(),
		Handler: jsii.String("bootstrap"),
		Code: awslambda.Code_FromAsset(jsii.String("../"), &awss3assets.AssetOptions{
			// Only upload the compiled binary, ignore source files
			Exclude: jsii.Strings("web", "infra", "internal", "cmd", "go.*", "*.md", "Makefile"),
		}),
		Architecture: awslambda.Architecture_ARM_64(),
		Environment: &map[string]*string{
			"AUTH0_ISSUER":   jsii.String(auth0Issuer),
			"AUTH0_AUDIENCE": jsii.String(auth0Audience),
			"TABLE_NAME":     table.TableName(),
		},
	})

	// Grant Lambda permissions to DynamoDB
	table.GrantReadWriteData(backendFunc)

	// Custom Domain for API Gateway
	apiDomain := awsapigatewayv2.NewDomainName(stack, jsii.String("ApiDomain"), &awsapigatewayv2.DomainNameProps{
		DomainName:  jsii.String(ApiSubdomain),
		Certificate: cert,
	})

	// HTTP API
	httpApi := awsapigatewayv2.NewHttpApi(stack, jsii.String("HttpApi"), &awsapigatewayv2.HttpApiProps{
		DefaultDomainMapping: &awsapigatewayv2.DomainMappingOptions{
			DomainName: apiDomain,
		},
		CorsPreflight: &awsapigatewayv2.CorsPreflightOptions{
			AllowOrigins: jsii.Strings(
				"http://localhost:5173",
				"https://"+AppSubdomain,
			),
			AllowMethods: &[]awsapigatewayv2.CorsHttpMethod{
				awsapigatewayv2.CorsHttpMethod_GET,
				awsapigatewayv2.CorsHttpMethod_POST,
				awsapigatewayv2.CorsHttpMethod_PUT,
				awsapigatewayv2.CorsHttpMethod_PATCH,
				awsapigatewayv2.CorsHttpMethod_DELETE,
				awsapigatewayv2.CorsHttpMethod_OPTIONS,
			},
			AllowHeaders: jsii.Strings("Authorization", "Content-Type"),
			MaxAge:       awscdk.Duration_Days(jsii.Number(1)),
		},
	})

	// Add Route (Proxy to Lambda)
	httpApi.AddRoutes(&awsapigatewayv2.AddRoutesOptions{
		Path:        jsii.String("/{proxy+}"),
		Integration: awsapigatewayv2integrations.NewHttpLambdaIntegration(jsii.String("LambdaIntegration"), backendFunc, &awsapigatewayv2integrations.HttpLambdaIntegrationProps{}),
	})

	// DNS Record for API
	awsroute53.NewARecord(stack, jsii.String("ApiAliasRecord"), &awsroute53.ARecordProps{
		Zone:       zone,
		RecordName: jsii.String(ApiSubdomain),
		Target:     awsroute53.RecordTarget_FromAlias(awsroute53targets.NewApiGatewayv2DomainProperties(apiDomain.RegionalDomainName(), apiDomain.RegionalHostedZoneId())),
	})

	// -----------------------------------------------------------------------
	// 3. Frontend (S3 + CloudFront)
	// -----------------------------------------------------------------------

	// S3 Bucket
	bucket := awss3.NewBucket(stack, jsii.String("WebBucket"), &awss3.BucketProps{
		BlockPublicAccess: awss3.BlockPublicAccess_BLOCK_ALL(),
		RemovalPolicy:     awscdk.RemovalPolicy_DESTROY,
		AutoDeleteObjects: jsii.Bool(true),
	})

	// Create Origin Access Identity (OAI) so CloudFront can read the private bucket
	oai := awscloudfront.NewOriginAccessIdentity(stack, jsii.String("OAI"), nil)
	bucket.GrantRead(oai, nil)

	// CloudFront Distribution
	dist := awscloudfront.NewDistribution(stack, jsii.String("WebDist"), &awscloudfront.DistributionProps{
		DefaultBehavior: &awscloudfront.BehaviorOptions{
			Origin: awscloudfrontorigins.NewS3Origin(bucket, &awscloudfrontorigins.S3OriginProps{
				OriginAccessIdentity: oai,
			}),
			ViewerProtocolPolicy: awscloudfront.ViewerProtocolPolicy_REDIRECT_TO_HTTPS,
		},
		DomainNames:       jsii.Strings(AppSubdomain),
		Certificate:       cert,
		DefaultRootObject: jsii.String("index.html"),
		ErrorResponses: &[]*awscloudfront.ErrorResponse{
			{
				HttpStatus:         jsii.Number(403),
				ResponseHttpStatus: jsii.Number(200),
				ResponsePagePath:   jsii.String("/index.html"),
				Ttl:                awscdk.Duration_Minutes(jsii.Number(0)),
			},
			{
				HttpStatus:         jsii.Number(404),
				ResponseHttpStatus: jsii.Number(200),
				ResponsePagePath:   jsii.String("/index.html"),
				Ttl:                awscdk.Duration_Minutes(jsii.Number(0)),
			},
		},
	})

	// Deploy Vue assets to S3
	awss3deployment.NewBucketDeployment(stack, jsii.String("DeployWeb"), &awss3deployment.BucketDeploymentProps{
		Sources: &[]awss3deployment.ISource{
			awss3deployment.Source_Asset(jsii.String("../web/dist"), nil),
		},
		DestinationBucket: bucket,
		Distribution:      dist, // Invalidate cache on deploy
		DistributionPaths: jsii.Strings("/*"),
	})

	// DNS Record for Frontend
	awsroute53.NewARecord(stack, jsii.String("AppAliasRecord"), &awsroute53.ARecordProps{
		Zone:       zone,
		RecordName: jsii.String(AppSubdomain),
		Target:     awsroute53.RecordTarget_FromAlias(awsroute53targets.NewCloudFrontTarget(dist)),
	})

	// -----------------------------------------------------------------------
	// Outputs
	// -----------------------------------------------------------------------
	awscdk.NewCfnOutput(stack, jsii.String("ApiUrl"), &awscdk.CfnOutputProps{
		Value: jsii.String("https://" + ApiSubdomain),
	})
	awscdk.NewCfnOutput(stack, jsii.String("AppUrl"), &awscdk.CfnOutputProps{
		Value: jsii.String("https://" + AppSubdomain),
	})

	return stack
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	// Fix: HostedZone_FromLookup requires an explicit Account ID.
	// CDK_DEFAULT_ACCOUNT is automatically populated by the CDK CLI when running 'cdk deploy'.
	account := os.Getenv("CDK_DEFAULT_ACCOUNT")
	if account == "" {
		// Fallback if running manually/locally without CDK CLI
		account = os.Getenv("AWS_ACCOUNT_ID")
	}

	if account == "" {
		panic("❌ Account ID is missing. HostedZone lookup requires explicit account configuration. Ensure you are running via 'cdk deploy' or set AWS_ACCOUNT_ID.")
	}

	NewPrayerJournalStack(app, "PrayerJournalStack", &PrayerJournalStackProps{
		StackProps: awscdk.StackProps{
			Env: &awscdk.Environment{
				Region:  jsii.String("us-east-1"),
				Account: jsii.String(account),
			},
		},
	})

	app.Synth(nil)
}
