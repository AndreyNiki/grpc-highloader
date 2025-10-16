package proto

import (
	"context"
	"errors"
	"fmt"

	"github.com/golang/protobuf/jsonpb"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/dynamic"
	"github.com/jhump/protoreflect/dynamic/grpcdynamic"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	"github.com/AndreyNiki/grpc-highloader/internal/entity"
	"github.com/AndreyNiki/grpc-highloader/internal/logger"
	"github.com/AndreyNiki/grpc-highloader/internal/metrics"
	"github.com/AndreyNiki/grpc-highloader/internal/templates"
)

// Requester send dynamic GRPC requests.
type Requester struct {
	methodDesc *desc.MethodDescriptor
	conn       *grpc.ClientConn
	stub       grpcdynamic.Stub
	metrics    *metrics.Metrics
	req        *entity.RequestParams
	parser     *ProtoParser
	tb         *templates.TemplateBuilder
}

// NewRequester create a new Requester.
func NewRequester(req *entity.RequestParams, metrics *metrics.Metrics) (*Requester, error) {
	p := NewProtoParser()
	r := &Requester{
		metrics: metrics,
		req:     req,
		parser:  p,
		tb:      templates.NewTemplateBuilder(),
	}

	methodDesc, err := r.parser.GetMethodDescriptor(req.Proto.FilePath, req.Method, req.Service)
	if err != nil {
		return nil, err
	}
	r.methodDesc = methodDesc

	conn, err := r.newConn(req.Host)
	if err != nil {
		return nil, err
	}
	stub := grpcdynamic.NewStub(conn)
	r.conn = conn
	r.stub = stub
	return r, nil
}

// SendUnaryRPCRequest send one unary rpc request.
func (r *Requester) SendUnaryRPCRequest(ctx context.Context) error {
	log := logger.LoggerFromContext(ctx)
	md := r.methodDesc.GetInputType()
	msg := dynamic.NewMessage(md)

	err := r.makeMessage(msg, r.req.Message)
	if err != nil {
		return err
	}

	r.metrics.IncrementRequestCount()
	resp, err := r.stub.InvokeRpc(ctx, r.methodDesc, msg)
	if err != nil {
		return err
	}
	log.Info("Response", resp.String())

	statusErr, ok := status.FromError(err)
	if !ok {
		return errors.New("failed to get status from response")
	}
	r.metrics.IncrementResponseStatus(statusErr.Code())

	return nil
}

// Close requester.
func (r *Requester) Close() {
	if r.conn != nil {
		r.conn.Close()
	}
}

// makeMessage make dynamic message for request.
func (r *Requester) makeMessage(message *dynamic.Message, data string) error {
	msg, err := r.tb.Process(data)
	if err != nil {
		return fmt.Errorf("processing template message failed: %w", err)
	}
	err = jsonpb.UnmarshalString(msg, message)
	if err != nil {
		return err
	}

	return nil
}

// newConn create a new connection for grpc.
func (r *Requester) newConn(host string) (*grpc.ClientConn, error) {
	ctx := context.Background()
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	return grpc.DialContext(ctx, host, opts...)
}
