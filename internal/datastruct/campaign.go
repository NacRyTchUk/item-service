package datastruct

const CampaignTableName = "campaigns"

type Campaign struct {
	Id   int64  `db:"id"`
	Name string `db:"name"`
}
