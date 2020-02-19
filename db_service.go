package common


import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

// DBService : For creating dynamodb connection.
type DBService struct {
	TableName string
}

// MyDynamo : For creating dynamodb connection.
type MyDynamo struct {
	Db dynamodbiface.DynamoDBAPI
}

// Dyna - exported
var Dyna *MyDynamo

// TransactItem : For creating TransactItem.
type TransactItem struct {
	TransactType              string
	TableName                 string
	KeyDetails                map[string]interface{}
	ItemDetails               interface{}
	ExpressionAttributeValues map[string]interface{}
	UpdateExpression          string
}

// ConfigureDynamoDB : For creating dynamodb connection.
func ConfigureDynamoDB() {
	Dyna = new(MyDynamo)

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(ConfigurationObj.AWS.Region)},
	)
	if err != nil {
		fmt.Println("error in session creation: ", err)
	}

	// Create DynamoDB client
	svc := dynamodb.New(sess)

	Dyna.Db = dynamodbiface.DynamoDBAPI(svc)
}

// PutItem : Inserts an item into Dynamodb
func (db DBService) PutItem(transactItem TransactItem) {
	av, err := dynamodbattribute.MarshalMap(transactItem.ItemDetails)
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(db.TableName),
	}

	_, err = Dyna.Db.PutItem(input)

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("Successfully created the record")
}

//UpdateItem : Updates an item in dynamodb
func (db DBService) UpdateItem(keyDetails map[string]string,
	itemDetails map[string]string, UpdateExpression string) {

	av, merr := dynamodbattribute.MarshalMap(itemDetails)
	fmt.Println("error: ", merr)
	keyAv, merr := dynamodbattribute.MarshalMap(keyDetails)
	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: av,
		TableName:                 aws.String(db.TableName),
		Key:                       keyAv,
		ReturnValues:              aws.String("UPDATED_NEW"),
		UpdateExpression:          aws.String(UpdateExpression),
	}

	_, err := Dyna.Db.UpdateItem(input)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("Successfully updated the record")
}

// GetItem gets data from dynamodb
func (db DBService) GetItem(keyDetails map[string]string, item interface{}) error {
	keyAv, err := dynamodbattribute.MarshalMap(keyDetails)
	fmt.Println(keyAv)
	result, err := Dyna.Db.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(db.TableName),
		Key:       keyAv,
	})
	if err != nil {
		fmt.Println(err.Error())
	}
	err = dynamodbattribute.UnmarshalMap(result.Item, item)
	return err
}

// QuerytItem gets data from dynamodb (only HK returns multiple records)
func (db DBService) QuerytItem(keyDetails map[string]string, items interface{}, keyEx string) error {
	keyAv, err := dynamodbattribute.MarshalMap(keyDetails)
	result, err := Dyna.Db.Query(&dynamodb.QueryInput{
		TableName:                 aws.String(db.TableName),
		KeyConditionExpression:    aws.String(keyEx),
		ExpressionAttributeValues: keyAv,
	})
	if err != nil {
		fmt.Println(err.Error())
	}

	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, items)
	return err
}

// BatchGetItem gets data from dynamodb (with multiple PK returns multiple records)
func (db DBService) BatchGetItem(keyDetails []map[string]string, items interface{}) error {
	var keys []map[string]*dynamodb.AttributeValue
	for _, element := range keyDetails {
		keyAv, err := dynamodbattribute.MarshalMap(element)
		if err != nil {
			fmt.Println(err.Error())
		}
		keys = append(keys, keyAv)
	}
	input := &dynamodb.BatchGetItemInput{
		RequestItems: map[string]*dynamodb.KeysAndAttributes{
			db.TableName: {
				Keys: keys,
				// ProjectionExpression: aws.String("AlbumTitle"),
			},
		},
	}
	result, err := Dyna.Db.BatchGetItem(input)
	if err != nil {
		fmt.Println(err.Error())
	}

	err = dynamodbattribute.UnmarshalListOfMaps(result.Responses[db.TableName], items)
	return err
}

// BatchWriteItemDelete delete data from dynamodb (with multiple PK delete in single go)
func (db DBService) BatchWriteItemDelete(keyDetails []map[string]string, items interface{}) error {
	var writeItems []*dynamodb.WriteRequest

	for _, element := range keyDetails {
		keyAv, err := dynamodbattribute.MarshalMap(element)
		if err != nil {
			fmt.Println(err.Error())
		}
		writeItems = append(writeItems, &dynamodb.WriteRequest{DeleteRequest: &dynamodb.DeleteRequest{Key: keyAv}})
	}
	input := &dynamodb.BatchWriteItemInput{

		RequestItems: map[string][]*dynamodb.WriteRequest{
			db.TableName: writeItems,
		},
	}
	result, err := Dyna.Db.BatchWriteItem(input)
	if err != nil {
		fmt.Println(err.Error())
	}
	if len(result.UnprocessedItems[db.TableName]) > 0 {
		fmt.Println("some of are not deleted")
	}
	return err
}

// QueryOnGSI gets data from dynamodb by querying on GSI
func (db DBService) QueryOnGSI(keyDetails map[string]string, items interface{}, keyEx string, indexName string) error {
	keyAv, err := dynamodbattribute.MarshalMap(keyDetails)
	result, err := Dyna.Db.Query(&dynamodb.QueryInput{
		TableName:                 aws.String(db.TableName),
		IndexName:                 aws.String(indexName),
		KeyConditionExpression:    aws.String(keyEx),
		ExpressionAttributeValues: keyAv,
	})
	if err != nil {
		fmt.Println(err.Error())
	}

	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, items)
	return err
}

// TransactWriteItems write to dynamodb in a single transaction
func TransactWriteItems(transactItems []TransactItem) map[string]interface{} {
	var result = make(map[string]interface{})
	transactWriteItems := []*dynamodb.TransactWriteItem{}
	for _, item := range transactItems {
		if item.TransactType == "put" {
			av, err := dynamodbattribute.MarshalMap(item.ItemDetails)
			if err != nil {
				ZapLoggerObj.Error(err.Error())
			}
			actualItem := &dynamodb.TransactWriteItem{
				Put: &dynamodb.Put{
					TableName: aws.String(item.TableName),
					Item:      av,
				},
			}
			transactWriteItems = append(transactWriteItems, actualItem)

		} else if item.TransactType == "update" {
			av, err := dynamodbattribute.MarshalMap(item.ExpressionAttributeValues)
			keyAv, err := dynamodbattribute.MarshalMap(item.KeyDetails)
			if err != nil {
				ZapLoggerObj.Error(err.Error())
			}
			actualItem := &dynamodb.TransactWriteItem{
				Update: &dynamodb.Update{
					TableName:                 aws.String(item.TableName),
					ExpressionAttributeValues: av,
					Key:                       keyAv,
					UpdateExpression:          aws.String(item.UpdateExpression),
				},
			}
			transactWriteItems = append(transactWriteItems, actualItem)

		} else if item.TransactType == "delete" {
			keyAv, err := dynamodbattribute.MarshalMap(item.KeyDetails)
			if err != nil {
				ZapLoggerObj.Error(err.Error())
			}
			actualItem := &dynamodb.TransactWriteItem{
				Delete: &dynamodb.Delete{
					TableName: aws.String(item.TableName),
					Key:       keyAv,
				},
			}
			transactWriteItems = append(transactWriteItems, actualItem)

		}

	}
	_, err4 := Dyna.Db.TransactWriteItems(&dynamodb.TransactWriteItemsInput{
		TransactItems: transactWriteItems,
	})

	if err4 != nil {
		ZapLoggerObj.Error("Error in Dynamo DB call")
		result["err"] = err4
		fmt.Println(err4)
	}

	return result
}
