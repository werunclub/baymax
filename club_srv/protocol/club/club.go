package club


const (
	prefix            = "Club."
	ServiceGetClub    = prefix + "Get"
	ServiceCreateClub = prefix + "Create"
)

type Club struct {
	ID uint `json:"id"`
	UserID string `json:"user_id"`
}

// 获取单条信息
type GetRequest struct {
	ClubId uint `json:"club_id"`
}

type GetResponse struct {
	Data     []Club `json:"data"`
}

// 创建新的俱乐部
type CreateRequest struct {
	Name string `json:"name"`
}

type CreateResponse struct {
	ClubId uint `json:"club_id"`
}
