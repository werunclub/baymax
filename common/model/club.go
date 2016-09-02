package model

import (
	"github.com/jinzhu/gorm"
	"strings"
	"time"
)

func init() {
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		defaultTableName = strings.Replace(defaultTableName, "_", "", -1)
		return "club_" + defaultTableName
	}
}

// 俱乐部
type Club struct {
	ID          uint   `gorm:"primary_key"`
	UserID      string `gorm:"size:36"`
	Name        string `gorm:"size:32;unique"`
	Icon        string `gorm:"size:128;default:''"`
	Description string `gorm:"type:text;column:des"`
	ShortUrl    string `gorm:"type:"`
	PersonCount uint   `gorm:"default:0"`
	SortNum     int    `gorm:"default:0"`
	State       bool   `gorm:"default:true"`
	Authorized  bool   `gorm:"default:true"`
	// 保留数据集合
	DataBody string `gorm:"size:1024;default:''"`
	Source   int    `gorm:"default:0"`
	// 俱乐部城市
	CityCode string `gorm:"size:8;default:0"`
	// 所属行业
	IndustryID int
	// -1位表示 是否参与排序  -2位表示 是否已经验证身份
	CommonByte int `gorm:"default:0"`

	CreateTime time.Time

	Teams                 []Team
	Persons               []ClubPerson
	Roles                 []ClubTeamRole
	Boards                []ClubBoard
	Topics                []Topic
	Activities            []Activity
	Matchs                []Match
	Challenges            []Challenge
	Messahes              []ClubMessage
	ClubPersonWeekSorts   []ClubPersonWeekSort
	UserMonthSteps        []UserMonthSteps
	ClubAuthorizations    []ClubAuthorization
	ClubGpsRecords        []ClubGpsRecord
	ClubIndustryWeekSorts []ClubIndustryWeekSort
}

// 部门
type Team struct {
	ID     uint `gorm:"primary_key"`
	ClubID int  `gorm:"index"`
	Parent int  `gorm:"index"`

	Name     string `gorm:"size:32"`
	DataBody string `gorm:"size:512;default:''"`
	Count    int    `gorm:"default:0"`

	CreateTime time.Time

	Teams   []Team
	Persons []ClubPerson
	Roles   []ClubTeamRole
	Matchs  []Match
}

// 俱乐部成员
type ClubPerson struct {
	ID     uint   `gorm:"primary_key"`
	UserID string `gorm:"size:36;default:'';index:idx_user_club"`
	// FIXME 第一个 index 是否必须
	ClubID int `gorm:"index;index:idx_user_club"`
	TeamId int `gorm:"index"`

	Role   int    `gorm:"default:0"`
	Name   string `gorm:"size:16;default:''"`
	Mobile string `gorm:"size:24;default:''"`
	Email  string `gorm:"size:64;default:''"`
	State  bool   `gorm:"default:false"`

	UpdateTime time.Time
	CreateTime time.Time

	Priv     string `gorm:"size:10;default:111"`
	DataBody string `gorm:"size:512;default:''"`

	Roles                 []ClubTeamRole
	ClubPersonFixedDatas  []ClubPersonFixedData
	ClubSummaryStepsRanks []ClubSummaryStepsRank
}

// 俱乐部部门权限
type ClubTeamRole struct {
	ID     uint `gorm:"primary_key"`
	TeamID int  `gorm:"index"`
	ClubID int  `gorm:"index"`
	Person int  `gorm:"index;unique_index"`

	// FIXME role_choices = ((1, u'管理员'),)
	Role      int `gorm:"default:'0'"`
	LastMsgID int `gorm:"default:'0'"`

	UpdateTime time.Time
	CreateTime time.Time
}

// 俱乐部成员数据修复
type ClubPersonFixedData struct {
	ID       uint      `gorm:"primary_key"`
	PersonID int       `gorm:"index;index:'idx_person_curday'"`
	Curday   time.Time `gorm:"index:'idx_person_curday"`
	Steps    int       `gorm:"default:0"`
	Meters   int       `gorm:"default:0"`
	DataBody string    `gorm:"size:512;default:''"`
}

// 俱乐部公告
type ClubBoard struct {
	ID         uint   `gorm:"primary_key"`
	UserID     string `gorm:"size:36"`
	ClubID     int    `gorm:"index"`
	Notice     string `gorm:"type:text;default:''"`
	CreateTime time.Time
}

// 俱乐部话题
type Topic struct {
	ID     uint   `gorm:"primary_key"`
	ClubID int    `gorm:"index"`
	UserId string `gorm:"size:36"`

	// 主题
	Title string `gorm:"size:64"`
	// 内容
	Content string `gorm:"type:text"`
	// 置顶标志
	SetTop bool `gorm:"default:false"`
	// 置顶时间
	SetTopTime time.Time
	// 最后回复时间
	TatestreplyTime time.Time
	// 最后修改时间
	LatestmodifyTime time.Time
	// 发布时间
	CreateTime time.Time
	// 客户端 IP
	ClientIP string `gorm:"size:50"`
	// 审核状态
	AdminState string `gorm:"default:false"`
	// 有效状态
	State bool `gorm:"default:true"`
}

// 活动
type Activity struct {
	ID     uint   `gorm:"primary_key"`
	ClubID int    `gorm:"index"`
	UserID string `gorm:"size:36"`

	Title       string `gorm:"size:64"`
	Content     string `gorm:"type:text"`
	StartTime   time.Time
	EndTime     time.Time
	Poster      string `gorm:"size:120"`
	Location    string `gorm:"size:64"`
	MapPosition string `gorm:"size:64"`
	ClientIP    string `gorm:"size:50"`
	PersonCount string `gorm:"default:0"`
	CreateTime  time.Time
	SortNum     int  `gorm:"default:0"`
	State       bool `gorm:"default:true"`
	AdminState  bool `gorm:"default:false"`

	Persons []ActivityPerson
}

// 活动参与人员
type ActivityPerson struct {
	ID         uint   `gorm:"primary_key"`
	ActivityID int    `gorm:"index;index:idx_activity_user"`
	UserID     string `gorm:"size:36;index:idx_activity_user"`
	JoinName   string `gorm:"size:32"`
	JoinPhone  string `gorm:"size:64"`
	CreateTime time.Time
}

// 竞赛
type Match struct {
	// FIXME
	// type_choices = ((11, u'多人竞赛'), (12, u'团队竞赛'),)
	// data_type_choices = ((0, u'追踪器(走跑类)'), (1, u'码表(骑行类)'), (2, u'手机咕咚(路线类)'))
	// data_use_type_choices = ((0, u'步'), (1, u'米'),)

	ID           uint   `gorm:"primary_key"`
	ClubID       int    `gorm:"index"`
	UserID       string `gorm:"size:36;index"`
	TeamId       int    `gorm:"index"`
	CommentCount int    `gorm:"default:0"`
	Icon         string `gorm:"size:256;default=''"`
	Type         int
	Title        string `gorm:"size:36"`
	Content      string `gorm:"type:text"`
	StartTime    time.Time
	EndTime      time.Time
	TeamLimit    int `gorm:"default:2"`
	PersonLimit  int `gorm:"default:0"`
	CalcoinChip  int `gorm:"default:0"`
	CalcoinTotal int `gorm:"default:0"`
	DateType     int `gorm:"default:0"`
	DataUseType  int `gorm:"default:0"`
	// 大于 0 就为目标竞赛
	DataTotalLimit int `gorm:"default:0"`
	PersonCount    int `gorm:"default:0"`
	// 0 无效   1 有效   2 创建中  -1 已归档
	State      int `gorm:"default:1"`
	CreateTime time.Time
	DataBody   string `gorm:"size:512"`

	MatchTeams     []MatchTeam
	MatchPersons   []MatchPerson
	MatchTeamInfos []MatchTeamInfo
}

// 竞赛团队
type MatchTeam struct {
	ID          uint   `gorm:"primary_key"`
	MatchID     int    `gorm:"index;index:index:idx_match_name"`
	Name        string `gorm:"size:32;default:'';index:idx_match_name"`
	UpdateTime  time.Time
	DataTotal   int  `gorm:"default:0"`
	SuccessFlag bool `gorm:"default:false"`
	CreateTime  time.Time
	DataBody    string `gorm:"size:512"`

	MatchPersons []MatchPerson
}

// 竞赛成员
type MatchPerson struct {
	ID      uint   `gorm:"primary_key"`
	UserID  string `gorm:"size:32;index;index:idx_match_user"`
	MatchID int    `gorm:"index;index:idx_match_user"`
	TeamID  int    `gorm:"index"`

	TeamPosition int  `gorm:"default:0"`
	CaptainFlag  bool `gorm:"default:false"`

	UpdateTime  time.Time
	DataTotal   int  `gorm:"default:0"`
	DataModify  int  `gorm:"default:0"`
	SuccessFlag bool `gorm:"default:false"`
	// 分配到得卡币数
	CalcoinDistributed int `gorm:"default:0"`
	CreateTime         time.Time
	DataBody           string `gorm:"size:512"`
}

// 挑战
type Challenge struct {
	// FIXME
	// type_choices = ((0, u'团队方式挑战'), (12, u'个人挑战'),)
	// data_type_choices = ((0, u'追踪器(走跑类)'), (1, u'码表(骑行类)'), (2, u'手机咕咚
	// (路线类)'))

	ID          uint   `gorm:"primary_key"`
	ClubID      int    `gorm:"index"`
	MapID       int    `gorm:"index"`
	UserID      string `gorm:"size:36"`
	Type        int
	Title       string `gorm:"size:36"`
	Content     string `gorm:"type:text"`
	StartTime   time.Time
	EndTime     time.Time
	PersonCount int  `gorm:"default:0"`
	IsOver      bool `gorm:"default:false"`
	// 结束时间
	OverTime   time.Time
	IsPublic   bool `gorm:"default:false"`
	DataType   int
	Icon       string `gorm:"size:128;default:''"`
	CreateTime time.Time
	State      bool   `gorm:"default:true"`
	DataBody   string `gorm:"size:512"`

	Teams   []ChallengeTeam
	Persons []ChallengePerson
}

// 竞赛团队
type ChallengeTeam struct {
	ID          uint   `gorm:"primary_key"`
	Name        string `gorm:"size:32:index:idx_challenge_name"`
	ChallengeID int    `gorm:"index;index:idx_challenge_name"`
	UpdateTime  time.Time
	DataTotal   float32 `gorm:"default:0"`
	SuccessFlag bool    `gorm:"default:false"`
	// 结束时间
	OverTime   time.Time
	CreateTime time.Time
	DataBody   string `gorm:"size:512"`

	Persons []ChallengePerson
}

// 挑战成员表
type ChallengePerson struct {
	ID          uint   `gorm:"primary_key"`
	UserID      string `gorm:"size:36;index:idx_challenge_user"`
	ChallengeID int    `gorm:"index;index:idx_challenge_user"`
	TeamID      int    `gorm:"index"`

	TeamPosition int  `gorm:"default:0"`
	CaptainFlag  bool `gorm:"default:false"`
	UpdateTime   time.Time
	DataTotal    float32 `gorm:"default:0;index"`
	CreateTime   time.Time
	DataBody     string `gorm:"size:512"`
}

// 挑战者对应的线路
type Map struct {
	ID            uint   `gorm:"primary_key"`
	Name          string `gorm:"size:32"`
	TotalDistance float32
	Img           string `gorm:"size:128"`
	Thumbnail     string `gorm:"size:128"`
	TemplateName  string `gorm:"size:32"`
	Desc          string `gorm:"default=''"`
	CreateTime    time.Time
	DataBody      string `gorm:"size:512"`

	Points []Point
}

// 线路对应的点
type Point struct {
	ID    uint `gorm:"primary_key"`
	MapID int  `gorm:"index"`
	// 距离起点距离
	Distance float32
	// 样式,决定位置
	Style string `gorm:"size:128"`
	// 有名字的代表地图上的实点
	Name       string `gorm:"size:32"`
	CreateTime time.Time
	DataBody   string `gorm:"size:512"`
}

// 俱乐部推荐阅读
type Book struct {
	ID         uint   `gorm:"primary_key"`
	Summary    string `gorm:"type:text"`
	Img        string `gorm:"size:128"`
	Url        string `gorm:"size:128"`
	SortNum    int    `gorm:"default:0"`
	State      bool   `gorm:"default:true"`
	CreateTime time.Time
	DataBody   string `gorm:"size:512"`
}

// 俱乐部申请
type ClubApply struct {
	ID       uint   `gorm:"primary_key"`
	Name     string `gorm:"size:128"`
	Email    string `gorm:"size:128"`
	Nick     string `gorm:"size:128"`
	Mobile   string `gorm:"size:11"`
	Position string `gorm:"size:128"`
	Notes    string `gorm:"size:128"`
	ClubID   int    `gorm:"default:0"`
	// 0 申请  1 通过  2 拒绝  -1 删除
	State      int `gorm:"default:0"`
	UpdateTime time.Time
	CreateTime time.Time
	DataBody   string `gorm:"size:512"`
	Source     int    `gorm:"default:nil"`
}

// 运营俱乐部中奖名单
type ClubApplyLottery struct {
	ID     uint   `gorm:"primary_key"`
	Name   string `gorm:"size:128"`
	Mobile string `gorm:"size:11"`
	Email  string `gorm:"size:128"`
	Data   string `gorm:"size:512"`
}

// 俱乐部消息中心
type ClubMessage struct {
	// FIXME msg_type_choices = ((0, u'俱乐部昨天加入成员'), (1, u'活动蒋开始消息'), (2, u'活动将结束通知'))
	ID         uint   `gorm:"primary_key"`
	ClubID     int    `gorm:"index"`
	TeamID     int    `gorm:"default:0"`
	MsgType    int    `gorm:"default:0"`
	MsgValue   int    `gorm:"default:0"`
	Content    string `gorm:"size:256"`
	State      int    `gorm:"default:0"`
	CreateTime time.Time
}

// 俱乐部平均步数
type ClubWeekSort struct {
	ID          uint      `gorm:"primary_key"`
	ClubID      int       `gorm:"index;index:idx_club_curday"`
	Curday      time.Time `gorm:"index;index:idx_club_curday"`
	Steps       int       `gorm:"default:0"`
	Ranking     int       `gorm:"default:0;index"`
	Progress    int       `gorm:"default:0"`
	PersonCount int       `gorm:"default:0"`
}

// 俱乐部内用户周数据
type ClubPersonWeekSort struct {
	ID       uint   `gorm:"primary_key"`
	ClubID   int    `gorm:"index"`
	TeamID   int    `gorm:"default:0;index"`
	UserID   string `gorm:"size:36;index"`
	Name     string `gorm:"size:36"`
	Curday   time.Time
	Steps    int `gorm:"default:0"`
	Ranking  int `gorm:"default:0;index"`
	Progress int `gorm:"default:0"`
}

// 用户周数据
type UserWeekSteps struct {
	ID     uint      `gorm:"primary_key"`
	UserID string    `gorm:"size:36;index;index:idx_user_curday"`
	Curday time.Time `gorm:"index;index:idx_user_curday"`
	Steps  int       `gorm:"default:0"`
}

// 用户月数据
type UserMonthSteps struct {
	ID     uint      `gorm:"primary_key"`
	ClubID int       `gorm:"default:0;index"`
	UserID string    `gorm:"size:36;index"`
	Curday time.Time `gorm:"index"`
	Steps  int       `gorm:"default:0"`
}

// 俱乐部帮助中心
type ClubHelpMsg struct {
	// FIXME msg_type_choices = ((0, u'常见问题'), (1, u'申请俱乐部'), (2, u'管理俱乐部'), (3, u'活动'), (4, u'统计'),)
	ID      uint   `gorm:"primary_key"`
	MsgType int    `gorm:"default:0"`
	Title   string `gorm:"size:128;"`
	Summary string `gorm:"type:text;default:''"`
	Content string `gorm:"type:text"`
	SortNum string `gorm:"default:0"`
	// 0 普通  1 置顶
	State      int `gorm:"default:0"`
	UpdateTime time.Time
	CreateTime time.Time
}

// 创建竞赛后竞赛部门和俱乐部部门对应表
type MatchTeamInfo struct {
	ID          uint `gorm:"primary_key"`
	MatchID     int  `gorm:"index"`
	TeamID      int  `gorm:"default:0"`
	MatchTeamID int  `gorm:"default:0"`
	CreateTime  time.Time
}

// 俱乐部API权限对应表
type ClubAuthorization struct {
	ID     uint   `gorm:"primary_key"`
	ClubID int    `gorm:"default:0;index"`
	Token  string `gorm:"size:128"`
	// 0：普通；1：调试中
	Status     int `gorm:"default:0"`
	CreateTime time.Time
}

// 俱乐部用户GPS数据搜集表
type ClubGpsRecord struct {
	ID         uint      `gorm:"primary_key"`
	ClubID     string    `gorm:"size:36;index"`
	RouteId    string    `gorm:"size:36"`
	UserID     string    `gorm:"size:36;index"`
	UpdateTime time.Time `gorm:"index"`
}

// 俱乐部用户计步数据搜集表
type ClubStepRecord struct {
	ID         uint   `gorm:"primary_key"`
	ClubID     string `gorm:"size:36;index"`
	UserID     string `gorm:"size:36;index"`
	UploadTime time.Time
	TargetTime time.Time
}

// 俱乐部广告
type ClubAdvertisement struct {
	// FIXME vtype_choices = ((0, u'非俱乐部用户可见'), (1, u'俱乐部用户可见'), (2, u'全部用户可见'))
	ID         uint      `gorm:"primary_key"`
	ImageUrl   string    `gorm:"size:256;default=''"`
	NextUrl    string    `gorm:"size:256;default=''"`
	StartTime  time.Time `gorm:"index"`
	EndTime    time.Time `gorm:"index"`
	Status     bool      `gorm:"default:false;index"`
	SortNum    int       `gorm:"default:0"`
	UserId     string    `gorm:"size:36;defaut:''"`
	CreateTime time.Time
	UpdateTime time.Time
	Vtype      int `gorm:"default:0"`
}

// 用户俱乐部汇总表-总表 自用户加入俱乐部后开始
type ClubSummaryStepsRank struct {
	ID          uint      `gorm:"primary_key"`
	Person      int       `gorm:"index"`
	TodayCurday time.Time `gorm:"index"`
	TodaySteps  time.Time `gorm:"default:0"`
	MonthCurday time.Time `gorm:"index"`
	MonthSteps  int       `gorm:"default:0"`
	TotalSteps  int       `gorm:"default:0"`
	// 计算总数据辅助日期
	AssistCurday time.Time
	// 计算总数据辅助数据  不包含 assist_curday 当天的数据
	AssistSteps int `gorm:"default:0"`
	PraiseNum   int `gorm:"default:0"`
	CreateTime  time.Time
	UpdateTime  time.Time
}

// 俱乐部行业
type ClubIndustry struct {
	ID         uint   `gorm:"primary_key"`
	Name       string `gorm:"size:256"`
	Parent     int    `gorm:index`
	SortNum    int    `gorm:"default:0"`
	CreateTime time.Time

	ClubIndustrys         []ClubIndustry
	ClubIndustryWeekSorts []ClubIndustryWeekSort
}

// 俱乐部行业周排名
type ClubIndustryWeekSort struct {
	ID     uint      `gorm:"primary_key"`
	ClubId int       `gorm:"index"`
	Curday time.Time `gorm:index`
	// 俱乐部的父行业
	IndustryId int `gorm:"index"`
	AvgSteps   int `gorm:"default:0"`
	CreateTime time.Time
}

// 俱乐部认证信息
type ClubAuthenticateInfo struct {
	ID     uint `gorm:"primary_key"`
	ClubId int  `gorm:"index"`
	// 公司全称
	Name string `gorm:"size:256"`
	// 法人
	LegalPerson string `gorm:"size:20;default:''"`
	Phone       string `gorm:"size:16;default:''"`
	// 隶属企业
	SubjectionName string `gorm:"size:256;default:''"`
	// 机构代码证
	OccUrl string `gorm:"size:256;default:''"`
	// 工商营业执照
	BleUrl string `gorm:"size:256;default:''"`
	// 是否通过
	State int `gorm:"default:0"`
	//  默认0  1通过  2失败  -1等待审核
	CheckState int `gorm:"default:0"`
	// 审核人
	CheckUserId string `gorm:"size:36;default:''"`
	// 失败原因
	FailReason string `gorm:"size:1024;default:''"`
	CheckTime  time.Time
	CreateTime time.Time
}

// 用户周数据V2版 , 加入了有户加入俱乐部日期的判断
type UserWeekStepsV2 struct {
	ID     uint   `gorm:"primary_key"`
	ClubId int    `gorm:"default:0;index;index:idx_club_user_curday"`
	UserId string `gorm:"size:36;index;index:idx_club_user_curday"`
	// 每周 1
	Curday     time.Time `gorm:"index;index:idx_club_user_curday"`
	Steps      int       `gorm:"default:0"`
	UpdateTime time.Time
	CreateTime time.Time
}
