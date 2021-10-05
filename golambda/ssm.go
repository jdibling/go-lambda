package golambda

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
)

type SSMParameterReader interface {
	ReadString(name string) (string, error)
	ReadStringWithContext(ctx context.Context, name string) (string, error)
}

type SSMParameterWriter interface {
	WriteString(name string, value string) error
	WriteStringWithContext(ctx context.Context, name string, value string) error
}

type SSMParameterReadWriter interface {
	SSMParameterReader
	SSMParameterWriter
}

func NewSSMParameterReadWriter() (SSMParameterReadWriter, error) {
	return &ssmParamReader{
		ssmSvc:    ssm.New(session.Must(session.NewSession())),
		encrypted: true,
	}, nil
}

type SSMParameterReaderOption func(ssmParamReader)

type ssmParamReader struct {
	ssmSvc    ssmiface.SSMAPI
	encrypted bool
}

type ssmParam struct {
	name      string
	encrypted bool
}

func (r *ssmParamReader) ReadStringWithContext(ctx context.Context, name string) (string, error) {
	// make sure name isn't empty
	if len(name) <= 1 {
		return "", fmt.Errorf("getting parameter; invalid name")
	}
	// append leading / if not present (otherwise we will get permission denied error)
	if !strings.HasPrefix(name, "/") {
		name = fmt.Sprintf("/%s", name)
	}
	// read parameter
	p := ssmParam{
		name:      name,
		encrypted: r.encrypted,
	}
	get, err := r.ssmSvc.GetParameterWithContext(ctx, &ssm.GetParameterInput{
		Name:           aws.String(p.name),
		WithDecryption: aws.Bool(p.encrypted),
	})
	if err != nil {
		return "", fmt.Errorf("getting parameter; %w", err)
	}
	if get.Parameter.Value == nil {
		return "", fmt.Errorf("getting parameter; nil")
	}
	return *get.Parameter.Value, nil
}

func (r *ssmParamReader) ReadString(name string) (string, error) {
	return r.ReadStringWithContext(context.Background(), name)
}

func (r *ssmParamReader) WriteStringWithContext(ctx context.Context, name string, value string) error {
	return fmt.Errorf("not implemented")
}

func (r *ssmParamReader) WriteString(name string, value string) error {
	return r.WriteStringWithContext(context.Background(), name, value)
}
