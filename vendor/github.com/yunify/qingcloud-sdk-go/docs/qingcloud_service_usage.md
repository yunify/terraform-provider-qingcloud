# QingCloud Service Usage Guide

Import and initialize QingCloud service with a config, and you are ready to use the initialized service.

Each API function take a Input struct and return an Output struct. The Input struct consists of request params, request headers and request elements, and the Output holds the HTTP status code, response headers, response elements and error (if error occurred).

``` go
import (
	qc "github.com/yunify/qingcloud-sdk-go/service"
)
```

### Code Snippet

Initialize the QingCloud service with a configuration

``` go
qcService, _ := qc.Init(configuration)
```

Initialize the instance service in a zone

``` go
pek3aInstance, _ := qcService.Instance("pek3a")
```

DescribeInstances

``` go
iOutput, _ := pek3aInstance.DescribeInstances(
	&qc.DescribeInstancesInput{
		Instances: qc.StringSlice([]string{"i-xxxxxxxx"}),
	},
)

// Print the return code.
fmt.Println(qc.IntValue(iOutput.RetCode))

// Print the first instance ID.
fmt.Println(qc.StringValue(iOutput.InstanceSet[0].InstanceID))
```

RunInstances

``` go
iOutput, _ := pek3aInstance.RunInstances(
	&qc.RunInstancesInput{
		ImageID:      qc.String("centos7x64d"),
		CPU:          qc.Int(1),
		Memory:       qc.Int(1024),
		LoginMode:    qc.String("keypair"),
		LoginKeyPair: qc.String("kp-xxxxxxxx"),
	},
)

// Print the return code.
fmt.Println(qc.IntValue(iOutput.RetCode))

// Print the job ID.
fmt.Println(qc.StringValue(iOutput.JobID))
```

Initialize the volume service in a zone

``` go
pek3aVolume, _ := qcService.Volume("pek3a")
```

DescribeVolumes

``` go
volOutput, _ := pek3aVolume.DescribeVolumes(&qc.DescribeVolumesInput{
	Volumes: qc.StringValue([]string{"vol-xxxxxxxx"}),
})

// Print the return code.
fmt.Println(qc.IntValue(volOutput.RetCode))

// Print the first volume name.
fmt.Println(qc.StringValue(volOutput.VolumeSet[0].VolumeName))
```

CreateVolumes

``` go
volOutput, _ := pek3aVolume.CreateVolumes(
	&qc.CreateVolumesInput{
		Size:  qc.Int(10),
		Count: qc.Int(2),
	},
)

// Print the return code.
fmt.Println(qc.IntValue(volOutput.RetCode))

// Print the job ID.
fmt.Println(qc.StringValue(volOutput.JobID))
```

