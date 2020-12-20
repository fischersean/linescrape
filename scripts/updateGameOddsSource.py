"""
Update every item in the game odds table to have a source of 'mybookie'
"""

import boto3
import pandas as pd


def main():

    # Read in item list
    key_list = pd.read_csv("item_list.csv")

    dynamodb = boto3.resource("dynamodb")

    table = dynamodb.Table("game-odds")
    print(table.creation_date_time)

    for i, row in key_list.iterrows():

        league = row["league"]
        time_stamp = row["time_stamp"]

        table.update_item(
            Key={"league": league, "time_stamp": time_stamp,},
            UpdateExpression="SET #S = :val1",
            ExpressionAttributeNames={"#S": "source",},
            ExpressionAttributeValues={":val1": "mybookie"},
        )

        print(f"Udated {league}, {time_stamp}")


if __name__ == "__main__":
    main()
