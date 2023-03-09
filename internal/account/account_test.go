package account

import (
	log "github.com/castbox/nirvana-gcore/glog"
	"github.com/castbox/nirvana-gcore/gmongo"
	"glogin/config"
	"glogin/db/db_core"
	"glogin/util"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
)

func TestAccount(t *testing.T) {
	var IOSStr = "001286.b6c90ca4961f4dc39791e2e63ea9d134.0148_cnofficial"
	var BundleString = "com.droidhang.aod.ios"
	t.Logf("src %v %v", IOSStr, BundleString)

	//IOSBson := util.StrToIntArray(IOSStr)
	IOSInferceArray := util.StrToInterfaceArray(IOSStr)
	filter := bson.M{"ios": IOSInferceArray, "bundle_id": BundleString}
	//"mongo_db": "ulogin-account-test",
	//	"code_table_name":"ulogin_verifycode",
	//	"account_table_name":"ulogin_account",
	count, err := gmongo.CountDocuments(config.MongoUrl(), "ulogin-account-test", "ulogin_account", filter)
	if err != nil {
		t.Logf("CountDocuments %v %v", err, count)
	}
	t.Logf("CountDocuments %v %v", err, count)
	if count >= 1 {
		accountData, err := LoadTest(config.MongoUrl(), "ulogin-account-test", "ulogin_account", filter)
		if err == nil {
			t.Logf("LoadTest %v", accountData)
		}
	}

	t.Logf("LoadTest finish !!!!!!!")

}

func LoadTest(mongoUrl string, mongoDb string, tableName string, filter interface{}) (result db_core.AccountData, err error) {
	doc, errFind := gmongo.FindOne(mongoUrl, mongoDb, tableName, filter)
	if errFind != nil {
		log.Warnw("AccountTable Load", "err", err)
		err = errFind
		return
	}
	err = doc.Decode(&result)
	return
}
