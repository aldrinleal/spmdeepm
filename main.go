package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/awsutil"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/tatsushid/go-fastping"
	"net"
	"time"
	"fmt"
	"os"
	"github.com/aws/aws-sdk-go/aws/session"
	"log"
)

func main() {
	host := os.Args[1]

	pinger := fastping.NewPinger()

	ra, err := net.ResolveIPAddr("ip4:icmp", host)

	if err != nil {
		panic(err)
	}

	pinger.AddIPAddr(ra)

	pinger.MaxRTT = time.Second * 15

	returnError := fmt.Errorf("Timed out")

	roundTripTime := time.Duration(0)

	pinger.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
		fmt.Printf("rtt: %016d ms\n", rtt.Nanoseconds() / 1000)
		roundTripTime = rtt
		returnError = nil
	}

	err = pinger.Run()

	if err == nil {
		publish_metrics(host, roundTripTime)
	} else {
		panic(err)
	}
}

func publish_metrics(host string, duration time.Duration) {
	var _ret []*cloudwatch.Dimension
	var metricData []*cloudwatch.MetricDatum

	region := "us-east-1"

	session := session.New(&aws.Config{Region: &region})

	cloudwatchService := cloudwatch.New(session)

	dim := cloudwatch.Dimension{
		Name: aws.String("gateway"),
		Value: aws.String(host),
	}

	_ret = append(_ret, &dim)

	metric := cloudwatch.MetricDatum{
		MetricName: aws.String("rtt"),
		Unit: aws.String("Microseconds"),
		Value: aws.Float64(float64(duration.Nanoseconds() / 1000)),
		Dimensions: _ret,
	}

	metricData = append(metricData, &metric)

	metricInput := &cloudwatch.PutMetricDataInput{
		MetricData: metricData,
		Namespace: aws.String("spmdeepm"),
	}

	resp, err := cloudwatchService.PutMetricData(metricInput)

	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			fmt.Errorf("[%s] %s", awsErr.Code, awsErr.Message)
			panic(awsErr)
		} else if err != nil {
			panic(err)
		}
	}

	log.Println(awsutil.StringValue(resp))
}
