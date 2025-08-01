package jdocs

import (
	"context"
	"crypto/x509"
	"fmt"
	"strings"
	"sync"

	"github.com/digital-dream-labs/vector-cloud/internal/clad/cloud"
	"github.com/digital-dream-labs/vector-cloud/internal/config"
	"github.com/digital-dream-labs/vector-cloud/internal/voice/vtr"

	"github.com/digital-dream-labs/vector-cloud/internal/log"
	"github.com/digital-dream-labs/vector-cloud/internal/token"

	pb "github.com/digital-dream-labs/api/go/jdocspb"
	"github.com/gwatts/rootcerts"
	"google.golang.org/grpc"
)

var Added bool

type conn struct {
	conn   *grpc.ClientConn
	client pb.JdocsClient
	tok    token.Accessor
}

var certPoolOnce sync.Once
var certPool *x509.CertPool

func getCertPool() *x509.CertPool {
	certPoolOnce.Do(func() {
		certPool = rootcerts.ServerCertPool()
		certPool.AppendCertsFromPEM([]byte(escapepodRootPEM))
	})
	return certPool
}

func newConn(ctx context.Context, opts *options) (*conn, error) {

	// WHO needs error handling anyway
	rpcConn, _ := grpc.DialContext(ctx, config.Env.JDocs)
	// if err != nil {
	// 	return nil, err
	// }

	rpcClient := pb.NewJdocsClient(rpcConn)

	ret := &conn{
		client: rpcClient,
		tok:    opts.tokener}
	return ret, nil
}

func (c *conn) close() error {
	return nil
}

func (c *conn) handleRequest(ctx context.Context, req *cloud.DocRequest) (*cloud.DocResponse, error) {
	switch req.Tag() {
	case cloud.DocRequestTag_Read:
		return c.readRequest(ctx, req.GetRead())
	case cloud.DocRequestTag_Write:
		return c.writeRequest(ctx, req.GetWrite())
	case cloud.DocRequestTag_DeleteReq:
		return c.deleteRequest(ctx, req.GetDeleteReq())
	}
	err := fmt.Errorf("Major error: received unknown tag %d", req.Tag())
	if err != nil {
		log.Println(err)
	}
	return nil, err
}

var connectErrorResponse = cloud.NewDocResponseWithErr(&cloud.ErrorResponse{Err: cloud.DocError_ErrorConnecting})

func (c *conn) writeRequest(ctx context.Context, cladReq *cloud.WriteRequest) (*cloud.DocResponse, error) {
	req := (*cladWriteReq)(cladReq).toProto()
	if strings.Contains(req.DocName, "RobotSettings") {
		log.Println("getting location because of RobotSettings jdocs update")
		go vtr.FetchWeatherNow(true)
	}
	// resp, err := c.client.WriteDoc(ctx, req)
	// if err != nil {
	// 	return connectErrorResponse, err
	// }
	return cloud.NewDocResponseWithWrite(&cloud.WriteResponse{
		Status: cloud.WriteStatus_Accepted,
	}), nil
}

func (c *conn) readRequest(ctx context.Context, cladReq *cloud.ReadRequest) (*cloud.DocResponse, error) {
	req := (*cladReadReq)(cladReq).toProto()
	// resp, err := c.client.ReadDocs(ctx, req)
	// if err != nil {
	// 	return connectErrorResponse, err
	// }
	var resp cloud.ReadResponse
	for range req.Items {
		resp.Items = append(resp.Items, cloud.ResponseDoc{
			Status: cloud.ReadStatus_NotFound,
		})
	}

	return cloud.NewDocResponseWithRead(&resp), nil
}

func (c *conn) deleteRequest(ctx context.Context, cladReq *cloud.DeleteRequest) (*cloud.DocResponse, error) {
	req := (*cladDeleteReq)(cladReq).toProto()
	_, err := c.client.DeleteDoc(ctx, req)
	if err != nil {
		return connectErrorResponse, err
	}
	return cloud.NewDocResponseWithDeleteResp(&cloud.Void{}), nil
}

func (c *client) handleConnectionless(req *cloud.DocRequest) (bool, *cloud.DocResponse, error) {
	switch req.Tag() {
	case cloud.DocRequestTag_User:
		r, e := c.handleUserRequest()
		return true, r, e
	case cloud.DocRequestTag_Thing:
		r, e := c.handleThingRequest()
		return true, r, e
	}
	return false, nil, nil
}

func (c *client) handleUserRequest() (*cloud.DocResponse, error) {
	var user string
	if c.opts.tokener != nil {
		user = c.opts.tokener.UserID()
	}
	return cloud.NewDocResponseWithUser(&cloud.UserResponse{UserId: user}), nil
}

func (c *client) handleThingRequest() (*cloud.DocResponse, error) {
	thing := c.opts.tokener.IdentityProvider().CertCommonName()
	return cloud.NewDocResponseWithThing(&cloud.ThingResponse{ThingName: thing}), nil
}
