# Let's use Amazon S3
import boto3
import json

ce_client = boto3.client('ce')

response = ce_client.get_cost_and_usage_with_resources(
    TimePeriod={
        'Start': '2022-08-2',
        'End': '2022-09-01'
    },
    Granularity='MONTHLY',
    Filter={
        'Dimensions': {
            'Key': 'SERVICE',
            'Values': [
                'EC2-Instances',
            ],
            # 'MatchOptions': [
            #     'EQUALS'|'ABSENT'|'STARTS_WITH'|'ENDS_WITH'|'CONTAINS'|'CASE_SENSITIVE'|'CASE_INSENSITIVE',
            # ]
        }
    },
    Metrics=[
        'BlendedCost',
    ],
    GroupBy=[
        {
            'Type': 'DIMENSION',
            'Key': 'RESOURCE_ID'
        },
    ]
)

print(json.dumps(response, indent=2))