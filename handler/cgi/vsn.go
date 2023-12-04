/**
 * @Author: littlesea
 * @copyright (C) 2022, <DroideLite>
 * @doc
 * 前后端比对vsn
 * @end
 * @File:  vsn
 * @Created: 2022/07/04 02:07
 */

package cgi

import (
	log "github.com/castbox/nirvana-gcore/glog"
	"github.com/gin-gonic/gin"
	"glogin/config"
	"glogin/db"
	"go.mongodb.org/mongo-driver/bson"
	"strconv"
	"strings"
	"time"
)

// vsn请求
type SetVsnReq struct {
	Token    string `json:"token" bson:"token"`
	Platform string `json:"platform" bson:"platform"`

	AppVersion                      string `json:"app_version" bson:"app_version"`
	ClientVersion                   string `json:"client_version" bson:"client_version"`
	ClientHotversion                string `json:"client_hotversion" bson:"client_hotversion"`
	ClientForceupdate               string `json:"client_forceupdate" bson:"client_forceupdate"`
	ClientDownloadurl               string `json:"client_downloadurl" bson:"client_downloadurl"`
	ClientUnderreviewVersion        string `json:"client_underreview_version" bson:"client_underreview_version"`
	ServerVersion                   string `json:"server_version" bson:"server_version"`
	ServerForceequal                string `json:"server_forceequal" bson:"server_forceequal"`
	ClientSpecifyForceupdateVersion string `json:"client_specify_forceupdate_version" bson:"client_specify_forceupdate_version"`
}

type SetVsnRsp struct {
	Errno int32  `json:"errno" bson:"errno"`
	Info  string `json:"info" bson:"info"`
}

type GetVsnReq struct {
	Token         string `json:"token" bson:"token"`
	Platform      string `json:"platform" bson:"platform"`
	ClientVersion string `json:"client_version" bson:"client_version"`
}

type GetVsnRsp struct {
	Errno                           int32  `json:"errno" bson:"errno"`
	Info                            string `json:"info" bson:"info"`
	AppVersion                      string `json:"app_version" bson:"app_version"`
	ClientVersion                   string `json:"client_version" bson:"client_version"`
	ClientHotversion                string `json:"client_hotversion" bson:"client_hotversion"`
	ClientForceupdate               string `json:"client_forceupdate" bson:"client_forceupdate"`
	ClientDownloadurl               string `json:"client_downloadurl" bson:"client_downloadurl"`
	ClientUnderreviewVersion        string `json:"client_underreview_version" bson:"client_underreview_version"`
	ServerVersion                   string `json:"server_version" bson:"server_version"`
	ServerForceequal                string `json:"server_forceequal" bson:"server_forceequal"`
	ClientSpecifyForceupdateVersion string `json:"client_specify_forceupdate_version" bson:"client_specify_forceupdate_version"`
}

type DBVsn struct {
	AppVersion                      string                 `json:"app_version" bson:"app_version"`
	ClientVersionMap                map[string]interface{} `json:"client_version_map" bson:"client_version_map"`
	ClientForceupdate               string                 `json:"client_forceupdate" bson:"client_forceupdate"`
	ClientDownloadurl               string                 `json:"client_downloadurl" bson:"client_downloadurl"`
	ClientUnderreviewVersion        string                 `json:"client_underreview_version" bson:"client_underreview_version"`
	ServerVersion                   string                 `json:"server_version" bson:"server_version"`
	ServerForceequal                string                 `json:"server_forceequal" bson:"server_forceequal"`
	ClientSpecifyForceupdateVersion string                 `json:"client_specify_forceupdate_version" bson:"client_specify_forceupdate_version"`
}

func SetVsn(ctx *gin.Context) {
	vsnReq := &SetVsnReq{}
	vsnRsp := SetVsnRsp{}
	defer func() {
		if err := recover(); err != nil {
			log.Errorw("got panic", "err", err)
		}
	}()
	err := ctx.Bind(vsnReq)
	if err != nil {
		log.Errorw("SetVsn", "err", err)
		panic(err)
	}
	log.Infow("SetVsn request", "request", vsnReq)
	if err != nil {
		ParseRequestError(ctx, 500, err)
		return
	}
	// 判断platform是否为""
	if vsnReq.Platform == "" {
		vsnRsp.Errno = -1
		vsnRsp.Info = "Platform == nil"
		ctx.JSON(404, vsnRsp)
	} else {
		//todo 存储字段
		err = DbSetVsn(vsnReq)
		if err != nil {
			ParseRequestError(ctx, 500, err)
			return
		}

		vsnRsp.Errno = 0
		vsnRsp.Info = "success"
		ctx.JSON(200, vsnRsp)
	}
}

func GetVsn(ctx *gin.Context) {
	startTime := time.Now().Unix()
	log.Infow("GetVsn request", "start", startTime)
	VsnReq := &GetVsnReq{}
	VsnRsp := &GetVsnRsp{}
	defer func() {
		log.Infow("GetVsn request end", "time", time.Now().Unix()-startTime)
		if err := recover(); err != nil {
			log.Errorw("got panic", "err", err)
		}
	}()
	err := ctx.Bind(VsnReq)
	log.Infow("GetVsn request", "request", VsnReq)
	if err != nil {
		log.Infow("GetVsn request", "err", err)
		ParseRequestError(ctx, 500, err)
		return
	}
	//todo 存储字段
	dbVsn, err := DbGetVsn(VsnReq)
	if err != nil {
		log.Infow("DbGetVsn request", "err", err)
		ParseRequestError(ctx, 500, err)
		return
	}

	hot, ok := dbVsn.ClientVersionMap[VsnReq.ClientVersion]
	if ok {
		VsnRsp.ClientHotversion = hot.(string)
	}
	clientVersion := VsnReq.ClientVersion
	result := strings.Split(clientVersion, ".")
	num, _ := strconv.Atoi(result[1])
	if num == 11 {
		clientVersion = result[0] + "." + "12" + "." + result[2]
	}
	VsnRsp.AppVersion = dbVsn.AppVersion
	VsnRsp.ClientVersion = clientVersion
	VsnRsp.ClientForceupdate = dbVsn.ClientForceupdate
	VsnRsp.ClientDownloadurl = dbVsn.ClientDownloadurl
	VsnRsp.ClientUnderreviewVersion = dbVsn.ClientUnderreviewVersion
	VsnRsp.ServerVersion = dbVsn.ServerVersion
	VsnRsp.ServerForceequal = dbVsn.ServerForceequal
	VsnRsp.ClientSpecifyForceupdateVersion = dbVsn.ClientSpecifyForceupdateVersion
	VsnRsp.Errno = 0
	VsnRsp.Info = "success"
	ctx.JSON(200, VsnRsp)
}

const VsnTable = "glogin_vsn"

func VsnTableName() string {
	tableName := config.Field("vsn_table_name").String()
	if tableName == "" {
		tableName = VsnTable
	}
	return tableName
}

func DbSetVsn(request *SetVsnReq) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Errorw("got panic", "err", err)
		}
	}()

	// 首先查询
	doc := &DBVsn{}
	_ = db.LoadOne(bson.M{"_id": request.Platform}, doc, VsnTableName())

	upData := bson.M{}
	upData["client_version_map"] = bson.M{}
	if request.AppVersion != "" {
		upData["app_version"] = request.AppVersion
	}
	if request.ClientVersion != "" && request.ClientHotversion != "" {
		if doc.ClientVersionMap == nil {
			tmp := bson.M{}
			tmp[request.ClientVersion] = request.ClientHotversion
			upData["client_version_map"] = tmp
		} else {
			doc.ClientVersionMap[request.ClientVersion] = request.ClientHotversion
			upData["client_version_map"] = doc.ClientVersionMap
		}
	} else {
		upData["client_version_map"] = doc.ClientVersionMap
	}
	if request.ClientForceupdate != "" {
		upData["client_forceupdate"] = request.ClientForceupdate
	}
	if request.ClientDownloadurl != "" {
		upData["client_downloadurl"] = request.ClientDownloadurl
	}
	if request.ClientUnderreviewVersion != "" {
		upData["client_underreview_version"] = request.ClientUnderreviewVersion
	}
	if request.ServerVersion != "" {
		upData["server_version"] = request.ServerVersion
	}
	if request.ServerForceequal != "" {
		upData["server_forceequal"] = request.ServerForceequal
	}
	if request.ClientSpecifyForceupdateVersion != "" {
		upData["client_specify_forceupdate_version"] = request.ClientSpecifyForceupdateVersion
	}

	setData := bson.M{"$set": upData}
	upDataErr := db.UpdateOne_Upsert(bson.M{"_id": request.Platform}, setData, VsnTableName())
	//res, err := gmongo.UpdateOne(uri, db, config.BuildTableName(roleDigestDBName), filter, set, options.Update().SetUpsert(true))
	if upDataErr != nil {
		return upDataErr
	}
	return nil
}

func DbGetVsn(request *GetVsnReq) (rsp *DBVsn, err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Errorw("got panic", "err", err)
		}
	}()
	doc := &DBVsn{}
	err = db.LoadOne(bson.M{"_id": request.Platform}, doc, VsnTableName())
	if err != nil {
		return nil, err
	}
	return doc, nil
}
