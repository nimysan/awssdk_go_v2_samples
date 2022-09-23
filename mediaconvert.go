package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/mediaconvert"
	"github.com/aws/aws-sdk-go-v2/service/mediaconvert/types"
)

func main() {

	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		// 需要指定一个Endpoint ----
		return aws.Endpoint{
			PartitionID:   "aws",
			URL:           "https://lxlxpswfb.mediaconvert.us-east-1.amazonaws.com",
			SigningRegion: "us-east-1",
		}, nil

	})

	//配置更多日志
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithClientLogMode(aws.LogRequest), config.WithEndpointResolverWithOptions(customResolver))
	if err != nil {
		panic("configuration error, " + err.Error())
	}

	//fmt.Println(cfg)
	mediaConvertClient := mediaconvert.NewFromConfig(cfg)

	//需要转换的在S3中的视频文档
	inputFileInS3 := "s3://media.cuteworld.top/sample/210329_06B_Bali_1080p_013-r.mp4"

	fmt.Println("---- " + inputFileInS3)

	//最富在的部分， 创建转换任务
	job := mediaconvert.CreateJobInput{
		Role:  aws.String("arn:aws:iam::390468416359:role/service-role/MediaConvert_Default_Role"),
		Queue: aws.String("fordjs"),
		Settings: &types.JobSettings{
			TimecodeConfig: &types.TimecodeConfig{
				Source: "ZEROBASED",
			},
			Inputs: []types.Input{
				types.Input{
					FileInput: aws.String(inputFileInS3),
				},
			},
			OutputGroups: []types.OutputGroup{
				types.OutputGroup{
					CustomName: aws.String("mp4-720p-output"),
					Name:       aws.String("File Group"),
					OutputGroupSettings: &types.OutputGroupSettings{
						FileGroupSettings: &types.FileGroupSettings{
							Destination: aws.String("s3://media.cuteworld.top/output/"),
						},
						Type: "FILE_GROUP_SETTINGS",
					},
					Outputs: []types.Output{
						types.Output{
							ContainerSettings: &types.ContainerSettings{
								Container: "MP4",
							},
							VideoDescription: &types.VideoDescription{
								CodecSettings: &types.VideoCodecSettings{
									Codec: "H_264",

									H264Settings: &types.H264Settings{
										//Bitrate:         5000000,
										MaxBitrate:      5000000,
										RateControlMode: "QVBR",
										//ScanTypeConversionMode: "TRANSITION_DETECTION",
									},
								},
								Height: 720,
								Width:  1280,
							},
							//AudioDescriptions: []types.AudioDescription{
							//	types.AudioDescription{
							//		CodecSettings: &types.AudioCodecSettings{
							//			Codec: "AAC",
							//			AacSettings: &types.AacSettings{
							//				Bitrate:    96000,
							//				CodingMode: "CODING_MODE_2_0",
							//				SampleRate: 48000,
							//			},
							//		},
							//	},
							//},
							//NameModifier: aws.String("SAMPLE"),
						},
					},
				},
			},
		},
		//AccelerationSettings: &types.AccelerationSettings{
		//	Mode: "DISABLED",
		//},

		StatusUpdateInterval: types.StatusUpdateIntervalSeconds15,
		Priority:             4,
	}

	//打印配置
	//b, err := json.MarshalIndent(job, "", "  ")
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println(string(b))

	response, err := mediaConvertClient.CreateJob(context.TODO(), &job)

	if err != nil {
		fmt.Println("errorMsg is: ", err)
	}

	jobArn := *response.Job.Arn

	fmt.Println("the response is ", jobArn)

	jobStatusOutput, err := mediaConvertClient.GetJob(context.TODO(),
		&mediaconvert.GetJobInput{
			Id: response.Job.Id,
		})
	if err != nil {
		panic(err)
	}
	status := jobStatusOutput.Job.Status
	fmt.Println("job status is ", status)

	//cancel it - 是否需要开启取消测试
	if 1 == 2 {
		jobCancelOutput, err := mediaConvertClient.CancelJob(context.TODO(),
			&mediaconvert.CancelJobInput{Id: response.Job.Id},
		)
		if err != nil {
			panic(err)
		}

		j, err := json.Marshal(&jobCancelOutput.ResultMetadata)

		fmt.Println("cancel status is ", string(j))
	}
}
