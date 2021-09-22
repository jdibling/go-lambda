package golambda

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

type APIGatewayProxyRequest events.APIGatewayProxyRequest

// If the request is Base64 encoded, GetBody() will decode before returning.
func (r APIGatewayProxyRequest) GetBody() (string, error) {
	if !r.IsBase64Encoded {
		return r.Body, nil
	}
	dec, err := base64.StdEncoding.DecodeString(r.Body)
	if err != nil {
		return "", fmt.Errorf("base64 decoding string: %s; %w", r.Body, err)
	}
	return string(dec), nil
}

func (r APIGatewayProxyRequest) MustGetBody() string {
	dec, err := r.GetBody()
	if err != nil {
		panic(err)
	}
	return dec
}

// func (r APIGatewayProxyRequest) PathParameters() map[string]string {
// 	return r.PathParameters
// }

// func (r APIGatewayProxyRequest) QueryParameters() map[string]string {
// 	return r.QueryStringParameters
// }

// func (r APIGatewayProxyRequest) MultiValueQueryParameters() map[string][]string {
// 	return r.MultiValueQueryStringParameters
// }

// func (r APIGatewayProxyRequest) Headers() map[string]string {
// 	return r.Headers
// }

func (r APIGatewayProxyRequest) GetToken() string {
	authValue, ok := r.Headers["Authorization"]
	if ok != true {
		return ""
	}
	authParts := strings.Split(authValue, " ")
	if len(authParts) != 2 {
		return ""
	}
	return authParts[1]
}

// func (r APIGatewayProxyRequest) Body() string {
// 	return r.Body
// }

func (r APIGatewayProxyRequest) MapBodyStrings() (map[string]string, error) {
	m := make(map[string]string)
	err := json.Unmarshal([]byte(r.Body), &m)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal body; %w", err)
	}
	return m, nil
}

func (r APIGatewayProxyRequest) MapBodyObjects() (map[string]interface{}, error) {
	bodyMap := map[string]interface{}{}
	body, err := r.GetBody()
	if err != nil {
		return nil, fmt.Errorf("getting body; %w", err)
	}
	err = json.Unmarshal([]byte(body), &bodyMap)
	if err != nil {
		return nil, fmt.Errorf("unmarshal failed; %w", err)
	}
	return bodyMap, nil

}

func (r APIGatewayProxyRequest) ParseBody(v interface{}) error {
	body, err := r.GetBody()
	if err != nil {
		return fmt.Errorf("getting body; %w", err)
	}
	err = json.Unmarshal([]byte(body), v)
	if err != nil {
		return fmt.Errorf("body parse failed; %w", err)
	}
	return nil
}
