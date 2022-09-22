// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// document https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/costexplorer
package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer/types"
)

//
//aws ce get-cost-and-usage --time-period Start="2022-08-20",End="2022-09-01" --granularity DAILY --filter --metrics "BlendedCost" --group-by Type=DIMENSION,Key=RESOURCE_ID

//https://docs.aws.amazon.com/zh_cn/aws-cost-management/latest/APIReference/API_GetCostAndUsageWithResources.html
func main() {

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("configuration error, " + err.Error())
	}

	//https://aws.github.io/aws-sdk-go-v2/docs/making-requests/
	ceClient := costexplorer.NewFromConfig(cfg)

	dataInterval := types.DateInterval{
		End:   aws.String("2022-09-01"),
		Start: aws.String("2022-08-27"),
	}

	filter := types.Expression{
		Dimensions: &types.DimensionValues{
			Key: "RESOURCE_ID",
			//MatchOptions: types.MatchOptionEquals,
			Values: []string{"i-05221513bd7541ba7"},
		},
	}

	groupBy := types.GroupDefinition{
		Type: types.GroupDefinitionTypeDimension,
		Key:  aws.String("RESOURCE_ID"),
	}
	//groupBy := []&types.G

	//"Dimensions": { "Key": "REGION", "Values": [ "us-east-1", “us-west-1” ] }
	input := &costexplorer.GetCostAndUsageWithResourcesInput{
		Granularity: types.GranularityDaily,
		TimePeriod:  &dataInterval,
		Filter:      &filter,
		Metrics:     []string{"BlendedCost"},
		GroupBy:     []types.GroupDefinition{groupBy},
	}

	fmt.Println("GetCostAndUsageWithResources input is {}", input)

	fmt.Println("Data Interval result is {} ", 1)

	result, err := ceClient.GetCostAndUsageWithResources(context.TODO(), input)

	b, err := json.MarshalIndent(result.ResultsByTime, "", "  ")
	fmt.Println(string(b))
}

// go with ce
