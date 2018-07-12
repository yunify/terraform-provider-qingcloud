// Code generated by private/model/cli/gen-api/main.go. DO NOT EDIT.

package cloudsearchdomain

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/client/metadata"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/signer/v4"
	"github.com/aws/aws-sdk-go/private/protocol/restjson"
)

// CloudSearchDomain provides the API operation methods for making requests to
// Amazon CloudSearch Domain. See this package's package overview docs
// for details on the service.
//
// CloudSearchDomain methods are safe to use concurrently. It is not safe to
// modify mutate any of the struct's properties though.
type CloudSearchDomain struct {
	*client.Client
}

// Used for custom client initialization logic
var initClient func(*client.Client)

// Used for custom request initialization logic
var initRequest func(*request.Request)

// Service information constants
const (
	ServiceName = "cloudsearchdomain"  // Name of service.
	EndpointsID = ServiceName          // ID to lookup a service endpoint with.
	ServiceID   = "CloudSearch Domain" // ServiceID is a unique identifer of a specific service.
)

// New creates a new instance of the CloudSearchDomain client with a session.
// If additional configuration is needed for the client instance use the optional
// aws.Config parameter to add your extra config.
//
// Example:
//     // Create a CloudSearchDomain client from just a session.
//     svc := cloudsearchdomain.New(mySession)
//
//     // Create a CloudSearchDomain client with additional configuration
//     svc := cloudsearchdomain.New(mySession, aws.NewConfig().WithRegion("us-west-2"))
func New(p client.ConfigProvider, cfgs ...*aws.Config) *CloudSearchDomain {
	var c client.Config
	if v, ok := p.(client.ConfigNoResolveEndpointProvider); ok {
		c = v.ClientConfigNoResolveEndpoint(cfgs...)
	} else {
		c = p.ClientConfig(EndpointsID, cfgs...)
	}
	if c.SigningNameDerived || len(c.SigningName) == 0 {
		c.SigningName = "cloudsearch"
	}
	return newClient(*c.Config, c.Handlers, c.Endpoint, c.SigningRegion, c.SigningName)
}

// newClient creates, initializes and returns a new service client instance.
func newClient(cfg aws.Config, handlers request.Handlers, endpoint, signingRegion, signingName string) *CloudSearchDomain {
	svc := &CloudSearchDomain{
		Client: client.New(
			cfg,
			metadata.ClientInfo{
				ServiceName:   ServiceName,
				ServiceID:     ServiceID,
				SigningName:   signingName,
				SigningRegion: signingRegion,
				Endpoint:      endpoint,
				APIVersion:    "2013-01-01",
				JSONVersion:   "1.1",
			},
			handlers,
		),
	}

	// Handlers
	svc.Handlers.Sign.PushBackNamed(v4.SignRequestHandler)
	svc.Handlers.Build.PushBackNamed(restjson.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(restjson.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(restjson.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(restjson.UnmarshalErrorHandler)

	// Run custom client initialization if present
	if initClient != nil {
		initClient(svc.Client)
	}

	return svc
}

// newRequest creates a new request for a CloudSearchDomain operation and runs any
// custom request initialization.
func (c *CloudSearchDomain) newRequest(op *request.Operation, params, data interface{}) *request.Request {
	req := c.NewRequest(op, params, data)

	// Run custom request initialization if present
	if initRequest != nil {
		initRequest(req)
	}

	return req
}
