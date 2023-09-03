package main

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/client/callopt"
	etcd "github.com/kitex-contrib/registry-etcd"
	relation "github.com/lyujsally/MinTikTok-lyujsally/kitex_gen/relation"
	"github.com/lyujsally/MinTikTok-lyujsally/kitex_gen/relation/relationservice"
	"github.com/lyujsally/MinTikTok-lyujsally/middlewares"
	"github.com/lyujsally/MinTikTok-lyujsally/service"
	"github.com/lyujsally/MinTikTok-lyujsally/settings"
)

const CtxUserIDKey = "userID"

// RelationServiceImpl implements the last service interface defined in the IDL.
type RelationServiceImpl struct{}

// IsFollow implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) IsFollow(ctx context.Context, req *relation.DouyinRelationIsfollowRequest) (resp *relation.DouyinRelationIsfollowResponse, err error) {
	// 从请求中获取用户ID和目标ID
	userId := req.UserId
	targetId := req.TargetId

	// 调用IsFollow业务逻辑
	followServiceImp := service.NewFollowInstance()
	isfollow, err := followServiceImp.IsFollow(userId, targetId)
	if err != nil {
		return &relation.DouyinRelationIsfollowResponse{
			StatusCode: int32(middlewares.CodeIsfollowFailed),
			StatusMsg:  middlewares.CodeIsfollowFailed.Msg(),
		}, err
	}

	// 返回响应
	return &relation.DouyinRelationIsfollowResponse{
		StatusCode: 0,
		StatusMsg:  "success",
		IsFollow:   isfollow,
	}, nil
}

// GetFolloweeCount implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) GetFolloweeCount(ctx context.Context, req *relation.DouyinRelationFolloweeCountRequest) (resp *relation.DouyinRelationFolloweeCountResponse, err error) {
	// 从请求中获取用户ID
	userId := req.UserId
	// 调用GetFolloweeCount业务逻辑
	followServiceImp := service.NewFollowInstance()
	followeeCount, err := followServiceImp.GetFolloweeCount(userId)
	if err != nil {
		return &relation.DouyinRelationFolloweeCountResponse{
			StatusCode: int32(middlewares.CodeGetFolloweeCountFailed),
			StatusMsg:  middlewares.CodeGetFolloweeCountFailed.Msg(),
		}, err
	}
	// 返回响应
	return &relation.DouyinRelationFolloweeCountResponse{
		StatusCode:    0,
		StatusMsg:     "success",
		FolloweeCount: followeeCount,
	}, nil
}

// GetFollowerCount implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) GetFollowerCount(ctx context.Context, req *relation.DouyinRelationFollowerCountRequest) (resp *relation.DouyinRelationFollowerCountResponse, err error) {
	// 从请求中获取用户ID
	userId := req.UserId
	// 调用GetFollowerCount业务逻辑
	followServiceImp := service.NewFollowInstance()
	followerCount, err := followServiceImp.GetFollowerCount(userId)
	if err != nil {
		return &relation.DouyinRelationFollowerCountResponse{
			StatusCode: int32(middlewares.CodeGetFollowerCountFailed),
			StatusMsg:  middlewares.CodeGetFollowerCountFailed.Msg(),
		}, err
	}
	// 返回响应
	return &relation.DouyinRelationFollowerCountResponse{
		StatusCode:    0,
		StatusMsg:     "success",
		FollowerCount: followerCount,
	}, nil
}

// Action implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) Action(ctx context.Context, req *relation.DouyinRelationActionRequest) (resp *relation.DouyinRelationActionResponse, err error) {
	// 获取当前请求的参数
	userId := req.UserId
	targetId := req.ToUserId
	actType := req.ActionType

	if actType < 1 || actType > 2 {
		return &relation.DouyinRelationActionResponse{
			StatusCode: int32(middlewares.CodeInvalidParam),
			StatusMsg:  middlewares.CodeInvalidParam.Msg(),
		}, errors.New(middlewares.CodeInvalidParam.Msg())
	}

	// 关注-取关业务逻辑
	followServiceImp := service.NewFollowInstance()
	if actType == 1 {
		go followServiceImp.Follow(userId, targetId)
	} else {
		go followServiceImp.UnFollow(userId, targetId)
	}

	return &relation.DouyinRelationActionResponse{
		StatusCode: 0,
		StatusMsg:  "success",
	}, nil
}

// GetFolloweeList implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) GetFolloweeList(ctx context.Context, req *relation.DouyinRelationFollowListRequest) (resp *relation.DouyinRelationFollowListResponse, err error) {
	log.Printf("进入GetFolloweeList方法")
	// 获取当前用户id
	userId := req.UserId
	// 调用GetFolloweeList业务逻辑
	followServiceImp := service.NewFollowInstance()
	users, err := followServiceImp.GetFolloweeList(userId)
	if err != nil {
		return &relation.DouyinRelationFollowListResponse{
			StatusCode: int32(middlewares.CodeGetFolloweeListFailed),
			StatusMsg:  middlewares.CodeGetFolloweeListFailed.Msg(),
		}, err
	}

	return &relation.DouyinRelationFollowListResponse{
		StatusCode: 0,
		StatusMsg:  "success",
		UserList:   users,
	}, nil
}

// GetFollowerList implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) GetFollowerList(ctx context.Context, req *relation.DouyinRelationFollowerListRequest) (resp *relation.DouyinRelationFollowerListResponse, err error) {
	// 获取当前用户id
	userId := req.UserId
	// 调用GetFollowerList业务逻辑
	followServiceImp := service.NewFollowInstance()
	users, err := followServiceImp.GetFollowerList(userId)
	if err != nil {
		return &relation.DouyinRelationFollowerListResponse{
			StatusCode: int32(middlewares.CodeGetFollowerListFailed),
			StatusMsg:  middlewares.CodeGetFollowerListFailed.Msg(),
		}, err
	}

	return &relation.DouyinRelationFollowerListResponse{
		StatusCode: 0,
		StatusMsg:  "success",
		UserList:   users,
	}, nil
}

// GetFriendList implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) GetFriendList(ctx context.Context, req *relation.DouyinRelationFriendListRequest) (resp *relation.DouyinRelationFriendListResponse, err error) {
	// 获取当前用户id
	userId := req.UserId
	// 调用GetFriendList业务逻辑
	followServiceImp := service.NewFollowInstance()
	friends, err := followServiceImp.GetFriendList(userId)
	if err != nil {
		return &relation.DouyinRelationFriendListResponse{
			StatusCode: int32(middlewares.CodeGetFriendListFailed),
			StatusMsg:  middlewares.CodeGetFriendListFailed.Msg(),
		}, err
	}

	r, err := etcd.NewEtcdResolver(settings.Conf.Endpoints)
	if err != nil {
		log.Printf("error:%v", err)
		return &relation.DouyinRelationFriendListResponse{
			StatusCode: int32(middlewares.CodeGetFriendListFailed),
			StatusMsg:  middlewares.CodeGetFriendListFailed.Msg(),
		}, err
	}

	cli, err := relationservice.NewClient(settings.Conf.ServiceName, client.WithResolver(r))
	if err != nil {
		log.Printf("Create Kitex client failed")
		return &relation.DouyinRelationFriendListResponse{
			StatusCode: int32(middlewares.CodeGetFriendListFailed),
			StatusMsg:  middlewares.CodeGetFriendListFailed.Msg(),
		}, err
	}

	reqToM := &relation.DouyinRelationGetFriendListRequest{
		UserId:     userId,
		FriendList: friends,
	}

	// 发送请求
	respFromM, err := cli.GetFriendsList(ctx, reqToM, callopt.WithRPCTimeout(3*time.Second))
	if err != nil {
		log.Printf("Redirect Message Server failed")
		return &relation.DouyinRelationFriendListResponse{
			StatusCode: int32(middlewares.CodeGetFriendListFailed),
			StatusMsg:  middlewares.CodeGetFriendListFailed.Msg(),
		}, err
	}

	resp.StatusCode = respFromM.StatusCode
	resp.StatusMsg = respFromM.StatusMsg
	resp.UserList = respFromM.FriendList

	return resp, nil
}

// GetFriendsList implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) GetFriendsList(ctx context.Context, req *relation.DouyinRelationGetFriendListRequest) (resp *relation.DouyinRelationGetFriendListResponse, err error) {
	// TODO: Your code here...
	return
}
