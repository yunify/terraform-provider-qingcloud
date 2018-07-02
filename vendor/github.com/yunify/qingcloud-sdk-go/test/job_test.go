package main

import (
	"errors"
	"fmt"
	"github.com/DATA-DOG/godog"
	qc "github.com/yunify/qingcloud-sdk-go/service"
)

func JobFeatureContext(s *godog.Suite) {
	s.Step(`^initialize job service$`, initializeJobService)
	s.Step(`^the job service is initialized$`, theJobServiceIsInitialized)

	s.Step(`^describe jobs$`, describeJobs)
	s.Step(`^describe jobs should get (\d+) job at least$`, describeJobsShouldGetJobAtLeast)
	s.Step(`^describe jobs should have the jobs I just created$`, describeJobsShouldHaveTheJobsIJustCreated)
}

var jobService *qc.JobService
var describeJobOutput *qc.DescribeJobsOutput

// --------------------------------------------------------------------------

func initializeJobService() error {
	jobService, err = qcService.Job(tc.Zone)
	return err
}

func theJobServiceIsInitialized() error {
	if jobService == nil {
		return errors.New("Job sub service is not initialized")
	}
	return nil
}

// --------------------------------------------------------------------------

func describeJobs() error {
	describeJobOutput, err = jobService.DescribeJobs(nil)
	return err
}

func describeJobsShouldGetJobAtLeast(count int) error {
	if len(describeJobOutput.JobSet) >= count {
		return nil
	}
	return fmt.Errorf("DescribeJobs doesn't get \"%d\" job(s)", count)
}

func describeJobsShouldHaveTheJobsIJustCreated() error {
	okCount := 0
	for _, job := range describeJobOutput.JobSet {
		if qc.StringValue(job.JobID) == qc.StringValue(runInstanceOutput.JobID) {
			okCount++
		}
		if qc.StringValue(job.JobID) == qc.StringValue(terminateInstanceOutput.JobID) {
			okCount++
		}
	}

	if okCount == 2 {
		return nil
	}
	return errors.New("DescribeJobs doesn't get the jobs I just created")
}
