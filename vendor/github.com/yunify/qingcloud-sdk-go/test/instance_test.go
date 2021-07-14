package main

import (
	"errors"
	"fmt"
	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	qcErrors "github.com/yunify/qingcloud-sdk-go/request/errors"
	qc "github.com/yunify/qingcloud-sdk-go/service"
	"strconv"
	"time"
)

var instanceService *qc.InstanceService
var runInstanceInput *qc.RunInstancesInput
var runInstanceOutput *qc.RunInstancesOutput

func InstanceFeatureContext(s *godog.Suite) {
	s.Step(`^initialize instance service$`, initializeInstanceService)
	s.Step(`^the instance service is initialized$`, theInstanceServiceIsInitialized)

	s.Step(`^instance configuration:$`, instanceConfiguration)
	s.Step(`^run instances$`, runInstances)

	s.Step(`^run instances should get a job ID$`, runInstancesShouldGetAJobID)
	s.Step(`^run instances will be finished$`, runInstancesWillBeFinished)

	s.Step(`^terminate instances$`, terminateInstances)
	s.Step(`^terminate instances should get a job ID$`, terminateInstancesShouldGetAJobID)
	s.Step(`^terminate instances will be finished$`, terminateInstancesWillBeFinished)
}

func initializeInstanceService() error {
	instanceService, err = qcService.Instance(tc.Zone)
	return err
}

func theInstanceServiceIsInitialized() error {
	if instanceService == nil {
		return errors.New("Instance sub service is not initialized")
	}
	return nil
}

// --------------------------------------------------------------------------

func instanceConfiguration(configuration *gherkin.DataTable) error {
	count, err := strconv.Atoi(configuration.Rows[1].Cells[2].Value)
	if err != nil {
		return err
	}

	runInstanceInput = &qc.RunInstancesInput{
		ImageID:      qc.String(configuration.Rows[1].Cells[0].Value),
		InstanceType: qc.String(configuration.Rows[1].Cells[1].Value),
		Count:        qc.Int(count),
		LoginMode:    qc.String(configuration.Rows[1].Cells[3].Value),
		LoginPasswd:  qc.String(configuration.Rows[1].Cells[4].Value),
	}
	return nil
}

func runInstances() error {
	runInstanceOutput, err = instanceService.RunInstances(runInstanceInput)
	return err
}

// --------------------------------------------------------------------------

func runInstancesShouldGetAJobID() error {
	if runInstanceOutput.JobID != nil {
		return nil
	}
	return errors.New("RunInstances don't get a job ID")
}

func runInstancesWillBeFinished() error {
	retries := 0
	for retries < tc.MaxRetries {
		describeJobOutput, err := jobService.DescribeJobs(
			&qc.DescribeJobsInput{
				Jobs: []*string{runInstanceOutput.JobID},
			},
		)
		if err != nil {
			return err
		}
		if qc.StringValue(describeJobOutput.JobSet[0].Status) == "successful" {
			return nil
		}
		retries++
		time.Sleep(time.Second * time.Duration(tc.RetryWaitTime))
	}
	return nil
}

// --------------------------------------------------------------------------

var terminateInstanceOutput *qc.TerminateInstancesOutput

func terminateInstances() error {
	retries := 0
	for retries < tc.MaxRetries {
		terminateInstanceOutput, err = instanceService.TerminateInstances(
			&qc.TerminateInstancesInput{
				Instances: runInstanceOutput.Instances,
			},
		)
		if err != nil {
			switch e := err.(type) {
			case *qcErrors.QingCloudError:
				fmt.Println(e)
				if e.RetCode != 1400 {
					return e
				}
			default:
				return err
			}
		} else {
			return nil
		}
		retries++
		time.Sleep(time.Second * time.Duration(tc.RetryWaitTime))
	}
	return nil
}

func terminateInstancesShouldGetAJobID() error {
	if terminateInstanceOutput.JobID != nil {
		return nil
	}
	return errors.New("TerminateInstances doesn't get a job ID")
}

func terminateInstancesWillBeFinished() error {
	retries := 0
	for retries < tc.MaxRetries {
		describeJobOutput, err := jobService.DescribeJobs(
			&qc.DescribeJobsInput{
				Jobs: []*string{terminateInstanceOutput.JobID},
			},
		)
		if err != nil {
			return err
		}
		if qc.StringValue(describeJobOutput.JobSet[0].Status) == "successful" {
			return nil
		}
		retries++
		time.Sleep(time.Second * time.Duration(tc.RetryWaitTime))
	}
	return nil
}
