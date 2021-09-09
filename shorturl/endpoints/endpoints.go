package endpoints

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/septemhill/dh/shorturl/service"
	"github.com/septemhill/dh/shorturl/transport"
)

type Endpoints struct {
	UploadURL endpoint.Endpoint
	AccessURL endpoint.Endpoint
}

func MakeEndpoints(srv service.ShortURLService) Endpoints {
	return Endpoints{
		UploadURL: makeUploadURLEndpoint(srv),
		AccessURL: makeAccessURLEndpoint(srv),
	}
}

func makeUploadURLEndpoint(srv service.ShortURLService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(transport.UploadURLRequest)
		id, url, err := srv.UploadURL(req.URL, req.ExpiredAt)
		if err != nil {
			return nil, err
		}
		return transport.UploadURLResponse{ID: id, ShortURL: url}, nil
	}
}

func makeAccessURLEndpoint(srv service.ShortURLService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(transport.AccessURLRequest)
		url, err := srv.AccessURL(req.ID)
		if err != nil {
			return nil, err
		}
		return transport.AccessURLResponse{URL: url}, nil
	}
}
