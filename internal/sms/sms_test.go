package sms

import (
	log "git.dhgames.cn/svr_comm/gcore/glog"
	"glogin/util"
	"testing"
	"time"
)

func TestRand(t *testing.T) {
	for i := 0; i < 1000; i++ {
		code := CreateVerifyCode()
		log.Infow("TestRand, CreateVerifyCode", "code", code)
	}
}

func TestRandAccount(t *testing.T) {
	MinAccount := 100000000
	MaxAccount := 999999999
	for i := 0; i < 1000; i++ {
		accountId := util.Rand32Num(int32(MinAccount), int32(MaxAccount))
		log.Infow("TestRand, accountId", "accountId", accountId)
	}
}

func TestTime(t *testing.T) {
	for i := 0; i < 1000; i++ {
		timeStr := time.Now().Format("2006-01-02 15:04:05")
		localTime := time.Now().Local()
		log.Infow("TestRand, accountId", "timeStr", timeStr, "localTime", localTime)
	}
}
