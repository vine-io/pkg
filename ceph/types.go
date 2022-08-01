package ceph

type Pool struct {
	Name                 string `json:"pool"`
	Id                   int32  `json:"pool_id"`
	Size                 int32  `json:"size"`
	MinSize              int32  `json:"min_size"`
	PGNum                int32  `json:"pg_num"`
	PGPNum               int32  `json:"pgp_num"`
	CrushRule            string `json:"crush_rule"`
	HashPsPool           bool   `json:"hashpspool"`
	NoDelete             bool   `json:"nodelete"`
	NoPgChange           bool   `json:"nopgchange"`
	NoSizeChange         bool   `json:"nosizechange"`
	WriteFadviseDontNeed bool   `json:"write_fadvise_dontneed"`
	NoScrub              bool   `json:"noscrub"`
	NodeepScrube         bool   `json:"nodeep-scrube"`
	UseGmtHitset         bool   `json:"use_gmt_hitset"`
	FastRead             int64  `json:"fast_read"`
	PgAutoScaleMode      string `json:"pg_auto_scale_mode"`
}
