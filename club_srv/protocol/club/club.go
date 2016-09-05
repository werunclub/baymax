package club

import (
	"time"
)

const (
	prefix = "Club."

	// 获取指定的俱乐部信息
	SrvGetOneClub = prefix + "GetOne"
	// 创建新的俱乐部
	SrvCreateClub = prefix + "Create"
	// 获取多个俱乐部
	SrvGetBatchClub = prefix + "GetBatch"
	// 修改俱乐部信息
	SrvUpdateClub = prefix + "Update"
	// 查询俱乐部
	SrvSearchClub = prefix + "Search"
)


type Club struct {
	ID         int      `json:"id"`
	UserID     string    `json:"user_id"`
	Name       string    `json:"name"`
	Icon       string    `json:"icon"`
	Des        string    `json:"des"`
	ShortUrl   string    `json:"short_url"`
	SortNum    int       `json:"sort_num"`
	State      bool      `json:"state"`
	Authorized bool      `json:"authorized"`
	DataBody   string    `json:"data_body"`
	Source     int       `json:"source"`
	CityCode   string    `json:"city_code"`
	IndustryID int       `json:"industry_id"`
	CommonByte int       `json:"common_byte"`
	CreateTime time.Time `json:"create_time"`
}

// 获取指定 id 俱乐部
type GetOneArgs struct {
	ClubID int
}
type GetOneReply struct {
	Data Club
}

// 获取多条 id 对应的 club 信息
type GetBatchArgs struct {
	ClubIDS []int
}
type GetBatchReply struct {
	Total int
	Data []Club
}

// 根据条件查询 club
// TODO 如何检测类型的范围比如这里的 limit 以及 offset
type SearchArgs struct {
	Name string
	Limit int
	Offset int
}
type SearchReply struct {
	Total int
	Data []Club
}

// 创建新的 club
type CreateArgs struct {
	Club Club
}
type CreateReply struct {
	Club Club
}

// 修改 club
type UpdateArgs struct {
	ClubID int
	NewClub Club
}
type UpdateReply struct {
	Club Club
}
