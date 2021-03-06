package awsgateway

import (
	"encoding/json"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/iotaledger/iota.go/trinary"
)

const (
	// Public
	// https://docs.aws.amazon.com/lambda/latest/dg/limits.html
	// 6MB payload, 300 sec execution time, 1000 concurrent exectutions.
	// Limit to 1000 POSTs and 20 chunks per request.

	// MaxConcurrency is the number of lambdas running concurrently
	MaxConcurrency = 5000

	// MaxChunksLen is the number of chunks sent to each lambda
	MaxChunksLen = 300 // .3 MB

	// private
	hooknodeFnNameDev  = "arn:aws:lambda:us-east-2:174232317769:function:lambda-node-dev-hooknode"
	hooknodeFnNameProd = "arn:aws:lambda:us-east-2:174232317769:function:lambda-node-production-hooknode"
	hooknodeRegion     = "us-east-2"
)

// HooknodeChunk is the chunk object sent to lambda
type HooknodeChunk struct {
	Address string         `json:"address"`
	Value   int            `json:"value"`
	Message trinary.Trytes `json:"message"`
	Tag     trinary.Trytes `json:"tag"`
}

// HooknodeReq is the payload sent to lambda
type HooknodeReq struct {
	Provider string           `json:"provider"`
	Chunks   []*HooknodeChunk `json:"chunks"`
}

var (
	sess = session.Must(session.NewSession(&aws.Config{Region: aws.String(hooknodeRegion)}))
)

// InvokeHooknode will invoke lambda to do PoW for the chunks in HooknodeReq
func InvokeHooknode(req *HooknodeReq) error {
	hooknodeFnName := hooknodeFnNameProd
	if os.Getenv("LAMBDA_ENV") == "dev" {
		hooknodeFnName = hooknodeFnNameDev
	}

	// Serialize params
	payload, err := json.Marshal(*req)
	if err != nil {
		return err
	}

	// Invoke lambda.
	client := lambda.New(sess)
	_, err = client.Invoke(&lambda.InvokeInput{
		FunctionName: aws.String(hooknodeFnName),
		Payload:      payload,
	})

	return err

	// fmt.Println("=========RESPONSE START=======")
	// fmt.Println("LAMBDA RETURNED")
	// bodyString := string(res.Payload)
	// fmt.Println(bodyString)
	// fmt.Println("========= RESPONSE END =======")

}
