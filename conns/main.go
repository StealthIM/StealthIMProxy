package conns

import (
	"StealthIMProxy/config"
	"StealthIMProxy/connpool"
	"context"
	"time"

	pb_fileapi "StealthIMProxy/StealthIM.FileAPI"
	pb_groupuser "StealthIMProxy/StealthIM.GroupUser"
	pb_msap "StealthIMProxy/StealthIM.MSAP"
	pb_session "StealthIMProxy/StealthIM.Session"
	pb_user "StealthIMProxy/StealthIM.User"

	"google.golang.org/grpc"
)

var (
	// PoolFileAPI 文件API连接池
	PoolFileAPI *connpool.Pool
	// PoolGroup 群连接池
	PoolGroup *connpool.Pool
	// PoolMSAPR MSAP读连接池
	PoolMSAPR *connpool.Pool
	// PoolMSAPW MSAP写连接池
	PoolMSAPW *connpool.Pool
	// PoolSession 会话连接池
	PoolSession *connpool.Pool
	// PoolUser 用户连接池
	PoolUser *connpool.Pool
)

// Init 初始化连接池
func Init(_ config.Config) {
	PoolFileAPI = connpool.NewPool("FileAPI", config.LatestConfig.Fileapi, func(conn *grpc.ClientConn) (bool, error) {
		cli := pb_fileapi.NewStealthIMFileAPIClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		_, err := cli.Ping(ctx, &pb_fileapi.PingRequest{})
		return true, err
	})
	PoolGroup = connpool.NewPool("Group", config.LatestConfig.Group, func(conn *grpc.ClientConn) (bool, error) {
		cli := pb_groupuser.NewStealthIMGroupUserClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		_, err := cli.Ping(ctx, &pb_groupuser.PingRequest{})
		return true, err
	})
	PoolMSAPR = connpool.NewPool("MSAP(R)", config.LatestConfig.MsapRead, func(conn *grpc.ClientConn) (bool, error) {
		cli := pb_msap.NewStealthIMMSAPSyncClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		_, err := cli.Ping(ctx, &pb_msap.PingRequest{})
		return true, err
	})
	PoolMSAPW = connpool.NewPool("MSAP(W)", config.LatestConfig.MsapWrite, func(conn *grpc.ClientConn) (bool, error) {
		cli := pb_msap.NewStealthIMMSAPClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		_, err := cli.Ping(ctx, &pb_msap.PingRequest{})
		return true, err
	})
	PoolSession = connpool.NewPool("Session", config.LatestConfig.Session, func(conn *grpc.ClientConn) (bool, error) {
		cli := pb_session.NewStealthIMSessionClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		_, err := cli.Ping(ctx, &pb_session.PingRequest{})
		return true, err
	})
	PoolUser = connpool.NewPool("User", config.LatestConfig.User, func(conn *grpc.ClientConn) (bool, error) {
		cli := pb_user.NewStealthIMUserClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		_, err := cli.Ping(ctx, &pb_user.PingRequest{})
		return true, err
	})
}
