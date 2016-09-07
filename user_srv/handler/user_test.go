package handler

import (
	userProtocol "baymax/user_srv/protocol/user"
	"testing"
	"net/rpc"
	"baymax/user_srv/model"
	"github.com/stretchr/testify/suite"
	"github.com/stretchr/testify/assert"
	"github.com/pborman/uuid"
	"time"
	"github.com/icrowley/fake"
)


//func TestUserHandler_Get(t *testing.T) {
//	once.Do(startRPCServer)
//
//	client := initClient()
//	defer client.Close()
//
//	var user userProtocol.UserDetailReply
//	err := client.Call(userProtocol.UserDetail, "accc", &user)
//	if err != nil {
//		t.Errorf("获取用户详情失败: %v", err)
//	}
//
//	if user == (userProtocol.UserDetailReply{}) {
//		t.Error("获取用户信息失败, 返回用户为空")
//	}
//}
//
//func TestUserHandler_List(t *testing.T) {
//	once.Do(startRPCServer)
//	client := initClient()
//	defer client.Close()
//
//	users := userProtocol.UserListReply{}
//	err := client.Call(userProtocol.UserList, []string{"acaacc"}, &users)
//	if err != nil {
//		t.Fatalf("批量获取用户信息失败: %v", err)
//	}
//
//	if users == (userProtocol.UserListReply{}) {
//		t.Fatal("返回用户信息列表为空")
//	}
//}
//
//func TestUserHandler_List_emptyReply(t *testing.T) {
//	once.Do(startRPCServer)
//	client := initClient()
//	defer client.Close()
//
//	var users userProtocol.UserListReply
//	err := client.Call(userProtocol.UserList, []string{}, &users)
//	if err != nil {
//		t.Fatalf("批量获取用户信息失败: %v", err)
//	}
//
//	if users != (userProtocol.UserListReply{}) {
//		t.Fatalf("返回用户信息应该为空: %v", users)
//	}

//}

func createUser() *model.User {
	t := time.Now()
	user := model.User{Id: uuid.New(), Dob: &t, Name: fake.UserName()}
	model.DB.Create(&user)
	return &user
}

type UserHandlerTestSuite struct {
	suite.Suite
	client *rpc.Client
}

func (suite *UserHandlerTestSuite) SetupTest() {
	once.Do(startRPCServer)
	connectDB()
	suite.client = initClient()
}

func (suite *UserHandlerTestSuite) TestCreateUser() {
	t := time.Now()
	args := userProtocol.UserInfo{
		Name: fake.UserName(),
		Dob: &t,
	}
	reply := userProtocol.UserDetailReply{}
	err := suite.client.Call(userProtocol.UserCreate, &args, &reply)
	assert.Nil(suite.T(), err, "创建用户失败")
	assert.Equal(suite.T(), args.Name, reply.Name, "返回昵称不一样")
}

func (suite *UserHandlerTestSuite) TestUserDetail() {
	user := createUser()
	reply := userProtocol.UserDetailReply{}
	args := userProtocol.UserDetailRequest{Id: user.Id}

	err := suite.client.Call(userProtocol.UserDetail, &args, &reply)
	assert.Nil(suite.T(), err, "获取用户详情出错")
	assert.Equal(suite.T(), user.Id, reply.Id, "获取到错误用户")
}

func (suite *UserHandlerTestSuite) TestUserDetailNotFound() {
	reply := userProtocol.UserDetailReply{}

	err := suite.client.Call(userProtocol.UserDetail, uuid.New(), &reply)
	assert.NotNil(suite.T(), err, "用户不存在应该返回错误")
}

func (suite *UserHandlerTestSuite) TestUserList() {
	user := createUser()
	reply := userProtocol.UserListReply{}
	err := suite.client.Call(userProtocol.UserList, []string{user.Id}, &reply)
	assert.Nil(suite.T(), err, "批量获取用户失败")
	assert.NotEmpty(suite.T(), *reply.Users, "未获取到用户信息")
}

func (suite *UserHandlerTestSuite) TestUserListWithEmpty()  {
	reply := userProtocol.UserListReply{}
	err := suite.client.Call(userProtocol.UserList, []string{}, &reply)
	assert.Nil(suite.T(), err, "批量获取用户信息失败")

	assert.Empty(suite.T(), *reply.Users, "返回用户信息应该为空")
}

func (suite *UserHandlerTestSuite) TestPatchUser() {
	user := createUser()
	args := userProtocol.PatchUserRequest{
		Id: user.Id,
		Payload: &userProtocol.UserInfo{
			Id: user.Id,
			DNDEnabled: true,
			DNDStart: 10,
			Name: "new name",
		},
	}
	reply := userProtocol.UserDetailReply{}
	err := suite.client.Call(userProtocol.UserPatch, args, &reply)
	assert.Nil(suite.T(), err, "修改用户信息出错")
	assert.Equal(suite.T(), user.Id, reply.Id)
	assert.Equal(suite.T(), "new name", reply.Name)
}

func TestUserHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(UserHandlerTestSuite))
}
