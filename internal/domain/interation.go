package domain

import (
	"context"
)

const (
	BizPost = "post"
)

type Interaction struct {
	Model
	// idx_biz select * from . where biz ==
	// idx_bizid_biz 联合索引 (bizID区分度高)
	BizID      int64  `gorm:"uniqueIndex:idx_interaction_bizID_biz"`
	Biz        string `gorm:"uniqueIndex:idx_interaction_bizID_biz;type:varchar(255)"` // MYSQL 写法
	ReadCnt    int
	LikeCnt    int
	CollectCnt int // 3个cnt 相比较 type+cnt 在读性能友好, 每次只需要读一行
}

func (Interaction) TableName() string {
	return `interaction`
}

//go:generate mockgen -source=./interaction.go -destination=./mock/interaction.go -package=domain_mock
type InteractionService interface {
	IncrReadCount(c context.Context, biz string, id int64) error
	Like(c context.Context, biz string, bizID, userID int64) error
	CancelLike(c context.Context, biz string, bizID, userID int64) error
	Collect(c context.Context, biz string, bizID, userID, collectionID int64) error
	CancelCollect(c context.Context, biz string, bizID, userID, collectionID int64) error
	Stat(c context.Context, biz string, bizID, userID int64) (Interaction, UserInteractionStat, error)

	// IncrReadCount(c context.Context, req InteractionIncrReadCntRequest) (Response, error)
	// Like(c context.Context, req InteractionLikeRequest) (Response, error)
	// CancelLike(c context.Context, req InteractionLikeRequest) (Response, error)
	// Collect(c context.Context, req InteractionCollectRequest) (Response, error)
	// CancelCollect(c context.Context, req InteractionCollectRequest) (Response, error)
	// Stat(c context.Context, req InteractionStatRequest) (Response, error)

	//	TODO 展示用户收藏、点赞列表 select bizID from 。。。 where biz and userID
	//GetByIDs(c context.Context, biz string, bizIDs []int64) (map[int64]Interaction, error)
}

type InteractionRepository interface {
	IncrReadCount(c context.Context, biz string, id int64) error
	BatchIncrReadCount(c context.Context, biz []string, id []int64) error
	Like(c context.Context, biz string, bizID, userID int64) error
	CancelLike(c context.Context, biz string, bizID, userID int64) error
	Collect(c context.Context, biz string, bizID, userID, collectionID int64) error
	CancelCollect(c context.Context, biz string, bizID, userID, collectionID int64) error
	GetByIDs(c context.Context, biz string, ids []int64) (map[int64]Interaction, error)
	Stat(c context.Context, biz string, bizID, userID int64) (Interaction, UserInteractionStat, error)
}

type UserLike struct {
	Model
	UserID int64  `gorm:"uniqueIndex:idx_userLike_userID_bizID_biz"`
	BizID  int64  `gorm:"uniqueIndex:idx_userLike_userID_bizID_biz"`
	Biz    string `gorm:"uniqueIndex:idx_userLike_userID_bizID_biz;type:varchar(255)"`
	Status bool   // true 点赞 false 取消点赞
	// `gorm:"uniqueIndex:idx_userID_bizID_biz"`
	//	具体索引顺序，需要根据业务需求规定，此外还需根据字段区分度
	//1. 查询用户喜欢的东西 select * from user_like where user_id = ? and biz = ?
	//2. 查询某个东西的点赞数 select * from user_like where bizID = ? and biz = ?
}

func (UserLike) TableName() string {
	return `user_like`
}

type UserCollect struct {
	Model
	UserID       int64  `gorm:"uniqueIndex:idx_userCollect_userID_bizID_biz"`
	BizID        int64  `gorm:"uniqueIndex:idx_userCollect_userID_bizID_biz"`
	Biz          string `gorm:"uniqueIndex:idx_userCollect_userID_bizID_biz;type:varchar(255);"`
	CollectionID int64  `gorm:"index"`
	Status       bool
}

func (UserCollect) TableName() string {
	return `user_collect`
}

type UserInteractionStat struct {
	Liked     bool
	Collected bool
}

type InteractionIncrReadCntRequest struct {
	Biz   string `json:"biz"`
	BizID int64  `json:"biz_id"`
}

type InteractionLikeRequest struct {
	Biz    string `json:"biz"`
	BizID  int64  `json:"biz_id"`
	UserID int64  `json:"user_id"`
}

type InteractionCollectRequest struct {
	Biz       string `json:"biz"`
	BizID     int64  `json:"biz_id"`
	UserID    int64  `json:"user_id"`
	CollectID int64  `json:"collect_id"`
}

type InteractionStatRequest struct {
	Biz    string `json:"biz"`
	BizID  int64  `json:"biz_id"`
	UserID int64  `json:"user_id"`
}
