//package handler_test
package handler

import (
	baymaxRPC "baymax/rpc"
	"baymax/user_srv/model"
	"flag"
	"fmt"
	"net/rpc"
	"github.com/jinzhu/configor"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/wawandco/fako"
	"net/rpc/jsonrpc"
	"sync"
	"testing"
	"time"
)

var once sync.Once

var Config = struct {
	Server struct {
		Addr string `default:":8080"`
	}
	Database struct{ Address string }
}{}

func init() {
	var (
		addr   string
		config string
	)

	flag.StringVar(&config, "config", "config.toml", "config file")
	flag.StringVar(&Config.Database.Address, "db_address", "", "database address")
	flag.StringVar(&addr, "addr", "", "addr, exmaple: 0.0.0.0:8080")
	flag.Parse()

	configor.Load(&Config)

	if addr != "" {
		Config.Server.Addr = addr
	}
}

func startRPCServer() {
	server := baymaxRPC.NewServer()
	server.RegisterName("User", &userHandler{})
	go server.Serve("tcp", Config.Server.Addr)
}

func connectDB() {
	model.Init(Config.Database.Address, true)
}

func initClient() *rpc.Client {
	var client *rpc.Client
	var err error
	for {
		client, err = jsonrpc.Dial("tcp", Config.Server.Addr)
		if err == nil {
			break
		}
	}
	return client
}

func init() {
	configor.Load(&Config)
}

type User struct {
	Id       string `fako:"digits"`
	Name string     `fako:""`
	Dob      time.Time
}

type UserDetailProtocol struct {
	Id       string `condition:"id = ?" fako:"digits"`
	Name  string `condition:"name = ?" fako:"user_name"`
	Dob time.Time `condition:"dob >= ?"`
}

type UserRequestTestSuite struct {
	suite.Suite
}

func (suite *UserRequestTestSuite) SetupTest() {
	// Do something setup the test environment
	connectDB()
}

func (suite *UserRequestTestSuite) TestBuildConditions() {
	request := UserDetailProtocol{
		Name: "show me",
	}
	conditions := *protocolToConditions(request)
	assert.Equal(suite.T(), 1, len(conditions), fmt.Sprintf("构建 Where 查询条件错误: %v", conditions))
	assert.Contains(suite.T(), conditions, "name = ?")
}

func (suite *UserRequestTestSuite) TestApplyFilter() {
	protocol := UserDetailProtocol{}
	fako.Fill(&protocol)
	user := User{}
	db := applyWhereFilter(model.DB, protocol)
	db.Find(&user)
}

func TestUserDetailRequest_Conditions(t *testing.T) {
	suite.Run(t, new(UserRequestTestSuite))
}
