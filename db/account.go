package db

import (
	log "git.dhgames.cn/svr_comm/gcore/glog"
	"git.dhgames.cn/svr_comm/gcore/gmongo"
	"glogin/config"
	"glogin/db/db_core"
	"go.mongodb.org/mongo-driver/bson"
	"math/rand"
)

const (
	AccountTableName = "glogin_account"
	MinAccount       = 100000000
	MaxAccount       = 999999999
	MaxTryTime       = 10
)

//
//-spec create_third(binary(), binary(), binary(), binary(), binary(), map()) -> {ok, integer()} | {error, any()}.
//create_third(Plat, OpenId, BundleId, IP, DhID, #{<<"unite">> := IsUnite} = Game) when IsUnite =:= true ->
//Account = #{Plat => OpenId, <<"create">> => #{<<"time">> => util:now(), <<"ip">> => IP, <<"sm_id">> => DhID, <<"bundle_id">> => BundleId}},
//lager:info("req create_third account unite ex, ~p ", [Account]),
//create(Account, BundleId, IP, DhID, Game);
//create_third(Plat, ThirdUid, BundleId, IP, DhID, Game) ->
//Account = #{Plat => ThirdUid, <<"bundle_id">> => BundleId, <<"create">> => #{<<"time">> => util:now(), <<"ip">> => IP}},
//lager:info("req create_third account , ~p ", [Account]),
//create(Account, BundleId, IP, DhID, Game).

func CheckNotExist(filter interface{}) bool {
	count, err := gmongo.CountDocuments(config.All{}.MongoUrl, config.All{}.MongoDb, AccountTableName, filter)
	if err != nil {
		return false
	}
	if count == 0 {
		return true
	}
	return false
}

func Load(filter interface{}) (result db_core.AccountData, err error) {
	doc, errFind := gmongo.FindOne(config.All{}.MongoUrl, config.All{}.MongoDb, AccountTableName, filter)
	if errFind != nil {
		err = errFind
		return
	}
	err = doc.Decode(&result)
	return
}

func CreateDhId(document bson.M) (accountId int32, err error) {
	i := 0
	for ; i < MaxTryTime; i++ {
		accountId = rand.Int31n(MaxAccount-MinAccount) + MinAccount
		document["_id"] = accountId
		_, errInsert := gmongo.InsertOne(config.All{}.MongoUrl, config.All{}.MongoDb, AccountTableName, document)
		if errInsert == nil {
			log.Infow("new account ok", "account", "times", accountId, i)
			return
		}
	}
	return
}
