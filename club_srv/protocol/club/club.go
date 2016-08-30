package club

const (
	prefix            = "Club."
	ServiceGetClub    = prefix + "Get"
	ServiceCreateClub = prefix + "Create"
)

type Club struct {
	Avatar string  `json:"avatar"`
	Nick   string  `json:"nick"`
	Score  float64 `json:"score"`
	Rank   int     `json:""rank`
}

type GetRequest struct {
	ClubId int64 `json:"club_id"`
}

type GetResponse struct {
	TotalNum int64  `json:"total_num"`
	Data     []Club `json:"data"`
}
