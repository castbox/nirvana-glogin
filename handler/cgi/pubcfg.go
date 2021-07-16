package cgi

import (
	"github.com/gin-gonic/gin"
	"glogin/internal/configure"
)

//load pub_cfg: Url "http://127.0.0.1:2000/v1/kv/app_dynamic_cfg/lwk_dev/pub_cfg/?recurse=true"
//load pub_cfg: [#{<<"CreateIndex">> => 16921314,<<"Flags">> => 0,<<"Key">> => <<"app_dynamic_cfg/lwk_dev/pub_cfg/com.dh.bpc.gp">>,<<"LockIndex">> => 0,<<"ModifyIndex">> => 16921314,<<"Value">> => <<"ewogICAgImNvbm5fY2ZnIjp7CiAgICAgICAgInN0YWJsZSI6ewoJCQkJCQkiY2x1c3Rlcl90eXBlIiA6IDEsCiAgICAgICAgICAgICJ1Z2F0ZV9hZGRyIjoiMTAuMC4yNDAuMjM0OjE4ODg5IiwKICAgICAgICAgICAgInZzbl9hZGRyIjoiMTAuMC4yNDAuMjM0OjE5ODg5IiwKICAgICAgICAgICAgInVwYXlfYWRkciI6IjEwLjAuMjQwLjE5OjgwODgiLAogICAgICAgICAgICAidWNoYXRfYWRkciI6ICJhb2QtZGV2LXVjaGF0LmRoZ2FtZXMuY246MTg4ODciLAogICAgICAgICAgICAiY29tbXVuaXR5X3dlYl9hZGRyIjoiaHR0cDovLzEwLjAuMC4xOTo3NzcwL2FvZC92MS9pbmRleC5odG1sIy9ob21lIiwKICAgICAgICAgICAgImNvbW11bml0eV9zcnZfYWRkciI6Imh0dHA6Ly8xODIuMTUwLjIyLjYxOjI3Nzc4IiwKICAgICAgICAgIAkiYWljc193c19hZGRyIjogIndzczovL2FpY3MtY2xpLmRldi1kaC5jb20vZGhfd3Mvd3NfYXBwIiwKICAgICAgICAgICAgImFpY3NfaHR0cF9hZGRyIjogImh0dHBzOi8vYWljcy1jbGkuZGV2LWRoLmNvbSIKICAgICAgICB9LAogICAgICAgICJjaGVjayI6ewogICAgICAgICAgICAiY2x1c3Rlcl90eXBlIjogMiwKICAgICAgICAgICAgInVnYXRlX2FkZHIiOiIiLAogICAgICAgICAgICAidnNuX2FkZHIiOiAiIiwKICAgICAgICAgICAgInVwYXlfYWRkciI6IiIsCiAgICAgICAgICAgICJ1Y2hhdF9hZGRyIjogIiIKICAgICAgICB9CiAgICB9LAogICAgImNoZWNrX3ZzbiI6IiIKfQ==">>},#{<<"CreateIndex">> => 16821432,<<"Flags">> => 0,<<"Key">> => <<"app_dynamic_cfg/lwk_dev/pub_cfg/com.droidhang.bpc.ios">>,<<"LockIndex">> => 0,<<"ModifyIndex">> => 16821481,<<"Value">> => <<"ewogICAgImNvbm5fY2ZnIjp7CiAgICAgICAgInN0YWJsZSI6ewoJCQkJCQkiY2x1c3Rlcl90eXBlIiA6IDEsCiAgICAgICAgICAgICJ1Z2F0ZV9hZGRyIjoiMTAuMC4yNDAuMjM0OjE4ODg5IiwKICAgICAgICAgICAgInZzbl9hZGRyIjoiMTAuMC4yNDAuMjM0OjE5ODg5IiwKICAgICAgICAgICAgInVwYXlfYWRkciI6Imh0dHA6Ly8xMC4wLjI0MC4xOTo4MDg4IiwKICAgICAgICAgICAgInVjaGF0X2FkZHIiOiAiYW9kLWRldi11Y2hhdC5kaGdhbWVzLmNuOjE4ODg3IiwKICAgICAgICAgICAgImNvbW11bml0eV93ZWJfYWRkciI6Imh0dHA6Ly8xMC4wLjAuMTk6Nzc3MC9hb2QvdjEvaW5kZXguaHRtbCMvaG9tZSIsCiAgICAgICAgICAgICJjb21tdW5pdHlfc3J2X2FkZHIiOiJodHRwOi8vMTgyLjE1MC4yMi42MToyNzc3OCIsCiAgICAgICAgICAJImFpY3Nfd3NfYWRkciI6ICJ3c3M6Ly9haWNzLWNsaS5kZXYtZGguY29tL2RoX3dzL3dzX2FwcCIsCiAgICAgICAgICAgICJhaWNzX2h0dHBfYWRkciI6ICJodHRwczovL2FpY3MtY2xpLmRldi1kaC5jb20iCiAgICAgICAgfSwKICAgICAgICAiY2hlY2siOnsKICAgICAgICAgICAgImNsdXN0ZXJfdHlwZSI6IDIsCiAgICAgICAgICAgICJ1Z2F0ZV9hZGRyIjoiIiwKICAgICAgICAgICAgInZzbl9hZGRyIjogIiIsCiAgICAgICAgICAgICJ1cGF5X2FkZHIiOiIiLAogICAgICAgICAgICAidWNoYXRfYWRkciI6ICIiCiAgICAgICAgfQogICAgfSwKICAgICJjaGVja192c24iOiIiCn0=">>}]s
//pub_cfg AppName <<"com.droidhang.bpc.ios">>, info #{<<"check_vsn">> => <<>>,<<"conn_cfg">> => #{<<"check">> => #{<<"cluster_type">> => 2,<<"uchat_addr">> => <<>>,<<"ugate_addr">> => <<>>,<<"upay_addr">> => <<>>,<<"vsn_addr">> => <<>>},<<"stable">> => #{<<"aics_http_addr">> => <<"https://aics-cli.dev-dh.com">>,<<"aics_ws_addr">> => <<"wss://aics-cli.dev-dh.com/dh_ws/ws_app">>,<<"cluster_type">> => 1,<<"community_srv_addr">> => <<"http://182.150.22.61:27778">>,<<"community_web_addr">> => <<"http://10.0.0.19:7770/aod/v1/index.html#/home">>,<<"uchat_addr">> => <<"aod-dev-uchat.dhgames.cn:18887">>,<<"ugate_addr">> => <<"10.0.240.234:18889">>,<<"upay_addr">> => <<"http://10.0.240.19:8088">>,<<"vsn_addr">> => <<"10.0.240.234:19889">>}}}
//pub_cfg AppName <<"com.dh.bpc.gp">>, info #{<<"check_vsn">> => <<>>,<<"conn_cfg">> => #{<<"check">> => #{<<"cluster_type">> => 2,<<"uchat_addr">> => <<>>,<<"ugate_addr">> => <<>>,<<"upay_addr">> => <<>>,<<"vsn_addr">> => <<>>},<<"stable">> => #{<<"aics_http_addr">> => <<"https://aics-cli.dev-dh.com">>,<<"aics_ws_addr">> => <<"wss://aics-cli.dev-dh.com/dh_ws/ws_app">>,<<"cluster_type">> => 1,<<"community_srv_addr">> => <<"http://182.150.22.61:27778">>,<<"community_web_addr">> => <<"http://10.0.0.19:7770/aod/v1/index.html#/home">>,<<"uchat_addr">> => <<"aod-dev-uchat.dhgames.cn:18887">>,<<"ugate_addr">> => <<"10.0.240.234:18889">>,<<"upay_addr">> => <<"10.0.240.19:8088">>,<<"vsn_addr">> => <<"10.0.240.234:19889">>}}}
func CfgHandler(ctx *gin.Context) {
	loginInfo := &configure.CfgRequest{}
	err := ctx.Bind(loginInfo)
	if err != nil {
		ParseRequestError(ctx, err)
		return
	}
	cfgRsp := configure.GetCfg(loginInfo.AppName, loginInfo.Vsn, ctx.ClientIP())
	if cfgRsp != nil {
		ctx.JSON(200, cfgRsp)
		return
	}
	ctx.JSON(200, configure.ClusterCfg{})
}
