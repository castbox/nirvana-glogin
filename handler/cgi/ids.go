package cgi

import (
	"fmt"
	log "git.dhgames.cn/svr_comm/gcore/glog"
	"github.com/gin-gonic/gin"
	"glogin/config"
	"glogin/constant"
	"glogin/internal/account"
	"glogin/internal/ids"
	"glogin/util"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

type OldAccountReq struct {
	NewAccountId int32  `json:"new_account_id" binding:"required"`
	OldBundleId  string `json:"old_bundle_id" binding:"required"`
}

type OldAccountRsp struct {
	OldAccount AccountRsp `json:"old_account" binding:"required"`
}
type AccountRsp struct {
	Platform     string `json:"platform" binding:"required"`
	OldAccountID string `json:"old_account_id" binding:"required"`
	OldBundleId  string `json:"old_bundle_id" binding:"required"`
}

type OldFacebookReq struct {
	NewAccountId int32  `json:"new_account_id" binding:"required"`
	OldBundleId  string `json:"old_bundle_id" binding:"required"`
}

type OldFacebookRsp struct {
	OldFacebookId string `json:"old_facebook_id" binding:"required"`
}

func IdsHandler(ctx *gin.Context) {
	action := ctx.Param("action")
	switch action {
	case constant.IdsGetOldAccount:
		getOLdAccount(ctx)
		break
	case constant.IdsGetFaceBookID:
		getFacebookId(ctx)
		break
	default:
		err := fmt.Errorf("request not found: %v", action)
		ParseRequestError(ctx, 500, err)
	}
}

// 参数 新的账号ID
// 返回老表中老的账号ID
func getOLdAccount(ctx *gin.Context) {
	oldAccountReq := &OldAccountReq{}
	err := ctx.Bind(oldAccountReq)
	log.Infow("got new query IdsHandler getOLdAccount request", "request", oldAccountReq)
	if err != nil {
		ParseRequestError(ctx, 500, err)
		return
	}
	accountRsp := &AccountRsp{}
	response := &OldAccountRsp{}
	before := time.Now().UnixNano()
	defer func() {
		if err := recover(); err != nil {
			log.Errorw("got panic", "err", err)
		}
		log.Infow("query IdsHandler getOLdAccount rsp", "response", response, "time_cost", (time.Now().UnixNano()-before)/1000000)
	}()

	// 账号不存在
	if account.CheckNotExist(bson.M{"_id": oldAccountReq.NewAccountId}) {
		err := fmt.Errorf("get_account_info account is not exist: %v", oldAccountReq.NewAccountId)
		ParseRequestError(ctx, 500, err)
		return
	}

	// 账号不正确
	dhAccount, err := account.Load(bson.M{"_id": oldAccountReq.NewAccountId})
	if err != nil || dhAccount.ID == 0 {
		err := fmt.Errorf("get_account_info account is not exist: %v", oldAccountReq.NewAccountId)
		ParseRequestError(ctx, 500, err)
		return
	}

	// facebook账号
	if dhAccount.Facebook != "" {
		log.Infow("this account is facebook", "FacebookID", dhAccount.Facebook)
		accessToken, err1 := ids.GetAccessToken(dhAccount.BundleID)
		if err1 != nil {
			errR1 := fmt.Errorf("getOLdAccount ids.GetAccessToken: %v", oldAccountReq.NewAccountId)
			ParseRequestError(ctx, 500, errR1)
			return
		}
		facebookId, err2 := ids.GetIds(dhAccount.Facebook, oldAccountReq.OldBundleId, accessToken)
		if err2 != nil {
			errR2 := fmt.Errorf("getOLdAccount ids.GetIds: %v", oldAccountReq.NewAccountId)
			ParseRequestError(ctx, 500, errR2)
			return
		}
		log.Infow("get_facebook_ids get_ids ok ", "FaceBookId", facebookId)

		tableName := config.Field("old_account_table_name").String()
		dbName := config.Field("mongo_old_gpdb").String()
		dhAccount, err := account.LoadSpecify(bson.M{"facebook": facebookId}, dbName, tableName)
		if err != nil || dhAccount.ID == 0 {
			//err := fmt.Errorf("get_account_info account is not exist: %v", oldAccountReq.NewAccountId)
			log.Infow("this old account is facebook login no account ", "FaceBookId", facebookId)
			accountRsp.Platform = "facebook"
			accountRsp.OldAccountID = "noaccount"
			response.OldAccount = *accountRsp
			ctx.JSON(200, response)
			return
		}
		accountRsp.Platform = "facebook"
		accountRsp.OldAccountID = util.Int2String(int(dhAccount.ID))
		accountRsp.OldBundleId = dhAccount.BundleID
		response.OldAccount = *accountRsp

	} else if dhAccount.Google != "" {
		log.Infow("this account is google login ", "Google", dhAccount.Google)
		tableName := config.Field("old_account_table_name").String()
		dbName := config.Field("mongo_old_gpdb").String()
		dhAccount, err := account.LoadSpecify(bson.M{"google": dhAccount.Google}, dbName, tableName)
		if err != nil || dhAccount.ID == 0 {
			err := fmt.Errorf("get_account_info account is not exist: %v", oldAccountReq.NewAccountId)
			ParseRequestError(ctx, 500, err)
			return
		}
		accountRsp.Platform = "google"
		accountRsp.OldAccountID = util.Int2String(int(dhAccount.ID))
		accountRsp.OldBundleId = dhAccount.BundleID
		response.OldAccount = *accountRsp
	}

	ctx.JSON(200, response)
}

// 参数 新的账号ID 旧的BundleId
// 返回值 旧的FacebookID
func getFacebookId(ctx *gin.Context) {
	oldFacebookReq := &OldFacebookReq{}
	err := ctx.Bind(oldFacebookReq)
	log.Infow("got new query IdsHandler getFacebookId request", "request", oldFacebookReq)
	if err != nil {
		ParseRequestError(ctx, 500, err)
		return
	}

	response := &OldFacebookRsp{}
	before := time.Now().UnixNano()
	defer func() {
		if err := recover(); err != nil {
			log.Errorw("got panic", "err", err)
		}
		log.Infow("query IdsHandler getOLdAccount rsp", "response", response, "time_cost", (time.Now().UnixNano()-before)/1000000)
	}()

	// 账号不存在
	if account.CheckNotExist(bson.M{"_id": oldFacebookReq.NewAccountId}) {
		err := fmt.Errorf("getFacebookId account is not exist: %v", oldFacebookReq.NewAccountId)
		ParseRequestError(ctx, 500, err)
		return
	}

	// 账号不正确
	dhAccount, err := account.Load(bson.M{"_id": oldFacebookReq.NewAccountId})
	if err != nil || dhAccount.ID == 0 {
		err := fmt.Errorf("getFacebookId account  error : %v", oldFacebookReq.NewAccountId)
		ParseRequestError(ctx, 500, err)
		return
	}

	// facebook账号
	if dhAccount.Facebook != "" {
		log.Infow("this account is facebook", "FacebookID", dhAccount.Facebook)
		accessToken, err1 := ids.GetAccessToken(dhAccount.BundleID)
		if err1 != nil {
			errR1 := fmt.Errorf("getOLdAccount ids.GetAccessToken: %v", oldFacebookReq.NewAccountId)
			ParseRequestError(ctx, 500, errR1)
			return
		}
		facebookId, err2 := ids.GetIds(dhAccount.Facebook, oldFacebookReq.OldBundleId, accessToken)
		if err2 != nil {
			errR2 := fmt.Errorf("getOLdAccount ids.GetIds: %v", oldFacebookReq.NewAccountId)
			ParseRequestError(ctx, 500, errR2)
			return
		}
		response.OldFacebookId = facebookId

	} else if dhAccount.Google != "" {
		log.Infow("this account is google login ", "Google", dhAccount.Google)
	}

	ctx.JSON(200, response)
}
