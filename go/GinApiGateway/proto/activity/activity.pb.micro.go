// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: activity.proto

package activity

import (
	fmt "fmt"
	proto "google.golang.org/protobuf/proto"
	_ "google.golang.org/protobuf/types/known/structpb"
	math "math"
)

import (
	context "context"
	client "go-micro.dev/v5/client"
	server "go-micro.dev/v5/server"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ client.Option
var _ server.Option

// Client API for Activity service

type ActivityService interface {
	GetActivity(ctx context.Context, in *GetActivityRequest, opts ...client.CallOption) (*GetActivityResponse, error)
	GetActivityById(ctx context.Context, in *GetActivityByIdRequest, opts ...client.CallOption) (*GetActivityByIdResponse, error)
	CreateActivity(ctx context.Context, in *CreateActivityRequest, opts ...client.CallOption) (*CreateActivityResponse, error)
	DeleteActivity(ctx context.Context, in *DeleteActivityRequest, opts ...client.CallOption) (*DeleteActivityResponse, error)
	UpdateActivity(ctx context.Context, in *UpdateActivityRequest, opts ...client.CallOption) (*UpdateActivityResponse, error)
	GetActivityUpdateHistoryById(ctx context.Context, in *GetActivityUpdateHistoryByIdRequest, opts ...client.CallOption) (*GetActivityUpdateHistoryByIdResponse, error)
}

type activityService struct {
	c    client.Client
	name string
}

func NewActivityService(name string, c client.Client) ActivityService {
	return &activityService{
		c:    c,
		name: name,
	}
}

func (c *activityService) GetActivity(ctx context.Context, in *GetActivityRequest, opts ...client.CallOption) (*GetActivityResponse, error) {
	req := c.c.NewRequest(c.name, "Activity.GetActivity", in)
	out := new(GetActivityResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *activityService) GetActivityById(ctx context.Context, in *GetActivityByIdRequest, opts ...client.CallOption) (*GetActivityByIdResponse, error) {
	req := c.c.NewRequest(c.name, "Activity.GetActivityById", in)
	out := new(GetActivityByIdResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *activityService) CreateActivity(ctx context.Context, in *CreateActivityRequest, opts ...client.CallOption) (*CreateActivityResponse, error) {
	req := c.c.NewRequest(c.name, "Activity.CreateActivity", in)
	out := new(CreateActivityResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *activityService) DeleteActivity(ctx context.Context, in *DeleteActivityRequest, opts ...client.CallOption) (*DeleteActivityResponse, error) {
	req := c.c.NewRequest(c.name, "Activity.DeleteActivity", in)
	out := new(DeleteActivityResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *activityService) UpdateActivity(ctx context.Context, in *UpdateActivityRequest, opts ...client.CallOption) (*UpdateActivityResponse, error) {
	req := c.c.NewRequest(c.name, "Activity.UpdateActivity", in)
	out := new(UpdateActivityResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *activityService) GetActivityUpdateHistoryById(ctx context.Context, in *GetActivityUpdateHistoryByIdRequest, opts ...client.CallOption) (*GetActivityUpdateHistoryByIdResponse, error) {
	req := c.c.NewRequest(c.name, "Activity.GetActivityUpdateHistoryById", in)
	out := new(GetActivityUpdateHistoryByIdResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Activity service

type ActivityHandler interface {
	GetActivity(context.Context, *GetActivityRequest, *GetActivityResponse) error
	GetActivityById(context.Context, *GetActivityByIdRequest, *GetActivityByIdResponse) error
	CreateActivity(context.Context, *CreateActivityRequest, *CreateActivityResponse) error
	DeleteActivity(context.Context, *DeleteActivityRequest, *DeleteActivityResponse) error
	UpdateActivity(context.Context, *UpdateActivityRequest, *UpdateActivityResponse) error
	GetActivityUpdateHistoryById(context.Context, *GetActivityUpdateHistoryByIdRequest, *GetActivityUpdateHistoryByIdResponse) error
}

func RegisterActivityHandler(s server.Server, hdlr ActivityHandler, opts ...server.HandlerOption) error {
	type activity interface {
		GetActivity(ctx context.Context, in *GetActivityRequest, out *GetActivityResponse) error
		GetActivityById(ctx context.Context, in *GetActivityByIdRequest, out *GetActivityByIdResponse) error
		CreateActivity(ctx context.Context, in *CreateActivityRequest, out *CreateActivityResponse) error
		DeleteActivity(ctx context.Context, in *DeleteActivityRequest, out *DeleteActivityResponse) error
		UpdateActivity(ctx context.Context, in *UpdateActivityRequest, out *UpdateActivityResponse) error
		GetActivityUpdateHistoryById(ctx context.Context, in *GetActivityUpdateHistoryByIdRequest, out *GetActivityUpdateHistoryByIdResponse) error
	}
	type Activity struct {
		activity
	}
	h := &activityHandler{hdlr}
	return s.Handle(s.NewHandler(&Activity{h}, opts...))
}

type activityHandler struct {
	ActivityHandler
}

func (h *activityHandler) GetActivity(ctx context.Context, in *GetActivityRequest, out *GetActivityResponse) error {
	return h.ActivityHandler.GetActivity(ctx, in, out)
}

func (h *activityHandler) GetActivityById(ctx context.Context, in *GetActivityByIdRequest, out *GetActivityByIdResponse) error {
	return h.ActivityHandler.GetActivityById(ctx, in, out)
}

func (h *activityHandler) CreateActivity(ctx context.Context, in *CreateActivityRequest, out *CreateActivityResponse) error {
	return h.ActivityHandler.CreateActivity(ctx, in, out)
}

func (h *activityHandler) DeleteActivity(ctx context.Context, in *DeleteActivityRequest, out *DeleteActivityResponse) error {
	return h.ActivityHandler.DeleteActivity(ctx, in, out)
}

func (h *activityHandler) UpdateActivity(ctx context.Context, in *UpdateActivityRequest, out *UpdateActivityResponse) error {
	return h.ActivityHandler.UpdateActivity(ctx, in, out)
}

func (h *activityHandler) GetActivityUpdateHistoryById(ctx context.Context, in *GetActivityUpdateHistoryByIdRequest, out *GetActivityUpdateHistoryByIdResponse) error {
	return h.ActivityHandler.GetActivityUpdateHistoryById(ctx, in, out)
}
