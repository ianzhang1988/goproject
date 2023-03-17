package main

import (
	"encoding/json"
	"fmt"
)

var json_string = `
{
	"code": 1,
	"data": {
		"bangkok_ais": [
			"singapore_equinix"
		],
		"beijing11_dxt": [
			"beijing_zjy_gc_cnc",
			"beijing4_dxt"
		],
		"beijing2_bgctvnet": [
			"ixp_beijing5_cnc"
		],
		"beijing2_cnix": [
			"ixp_beijing5_cnc"
		],
		"beijing3_bgctvnet": [
			"ixp_beijing5_cnc"
		],
		"beijing3_cnc": [
			"beijing_zjy_gc_cnc",
			"jinan9_gc_cnc"
		],
		"beijing3_fbwn": [
			"beijing4_fbwn"
		],
		"beijing4_dxt": [
			"beijing_zjy_gc_cnc",
			"wuhan_lkg_gc_ct"
		],
		"beijing4_fbwn": [
			"beijing2_cnix",
			"ixp_beijing5_cnc"
		],
		"beijing5_cnc_proxy": [
			"wuhan_lkg_cnc_bus_proxy"
		],
		"beijing5_crtc": [
			"wuhan_lkg_gc_cmnet",
			"jinan9_cmnet_proxy"
		],
		"beijing5_ixp_bus_proxy": [
			"wuhan_lkg_gc_ct"
		],
		"beijing_zjy_cnc_proxy": [
			"beijing_zjy_gc_cnc"
		],
		"beijing_zjy_dxt_proxy": [
			"beijing_zjy_gc_cnc",
			"beijing4_dxt"
		],
		"beijing_zjy_gc_cnc": [
			"default"
		],
		"beijing_zjy_syncloud_proxy": [
			"wuhan_lkg_gc_cmnet"
		],
		"bj_ecs_ali_s3_gc_inner": [
			"beijing_zjy_gc_cnc"
		],
		"changchun2_fbwn": [
			"beijing2_cnix",
			"ixp_beijing5_cnc"
		],
		"changchun_crtc": [
			"jinan9_cmnet_proxy"
		],
		"changchun_fbwn": [
			"beijing4_fbwn"
		],
		"changchun_scc": [
			"luoyang5_cnc",
			"ixp_beijing5_cnc"
		],
		"changsha2_scc": [
			"foshan_kpl_qnet_proxy",
			"shanghai_qnet"
		],
		"changsha4_scc": [
			"wuhan_lkg_gc_cmnet",
			"jinan9_cmnet_proxy"
		],
		"changsha5_scc": [
			"foshan_kpl_qnet_proxy",
			"shanghai_qnet"
		],
		"changsha_cscatv": [
			"wuhan6_cnc"
		],
		"chengdu2_scc": [
			"chengdu3_scc",
			"ixp_beijing5_cnc"
		],
		"chengdu3_scc": [
			"beijing2_cnix",
			"wuhan6_cnc",
			"shanghai_qnet"
		],
		"chengdu_scc": [
			"chengdu3_scc",
			"ixp_beijing5_cnc"
		],
		"chongqing2_gwbn": [
			"beijing11_dxt"
		],
		"chongqing2_scc": [
			"ixp_beijing5_cnc",
			"shanghai_qnet",
			"beijing11_dxt"
		],
		"chongqing4_scc": [
			"ixp_beijing5_cnc",
			"shanghai_qnet",
			"beijing11_dxt"
		],
		"dalian_fbwn": [
			"beijing2_cnix",
			"ixp_beijing5_cnc"
		],
		"dongguan2_scc": [
			"dongguan3_scc"
		],
		"dongguan3_scc": [
			"wuhan6_cnc",
			"shanghai_qnet"
		],
		"dongguan_scc": [
			"dongguan3_scc"
		],
		"ecs_ali_s3_gc_inner": [
			"beijing_zjy_gc_cnc"
		],
		"ecs_ali_s3_gc_inner1": [
			"beijing_zjy_gc_cnc"
		],
		"foshan_kpl_ali_s3_gc_ct": [
			"default"
		],
		"foshan_kpl_cnc_bus_proxy": [
			"foshan_kpl_ali_s3_gc_ct"
		],
		"foshan_kpl_cnc_proxy": [
			"foshan_kpl_ali_s3_gc_ct"
		],
		"foshan_kpl_ct_bus_proxy": [
			"foshan_kpl_ali_s3_gc_ct"
		],
		"foshan_kpl_qnet_proxy": [
			"foshan_kpl_ali_s3_gc_ct",
			"wuhan_lkg_gc_cmnet"
		],
		"foshan_scc": [
			"dongguan_scc",
			"dongguan3_scc"
		],
		"fuzhou2_scc": [
			"ixp_beijing5_cnc",
			"shanghai_qnet"
		],
		"gs_beijing_ct": [
			"wuhan_lkg_gc_ct"
		],
		"guangzhou8_cmnet_proxy": [
			"foshan_kpl_ali_s3_gc_ct"
		],
		"guangzhou9_cmnet": [
			"guangzhou8_cmnet_proxy",
			"wuhan_cmnet"
		],
		"guangzhou9_cmnet_bus_proxy": [
			"guangzhou9_cmnet"
		],
		"guangzhou_gzcatv": [
			"foshan_kpl_qnet_proxy",
			"shanghai_qnet"
		],
		"guiyang2_scc": [
			"chongqing4_scc",
			"chongqing2_scc",
			"shanghai_qnet"
		],
		"guiyang3_scc": [
			"chongqing4_scc",
			"guiyang2_scc"
		],
		"guizhou_cqccn1_scc": [
			"guiyang2_scc",
			"guiyang3_scc"
		],
		"guizhou_cqccn2_scc": [
			"chongqing2_scc",
			"guiyang2_scc"
		],
		"guizhou_cqccn_scc": [
			"chongqing2_scc",
			"guiyang2_scc"
		],
		"haerbin2_cnc": [
			"beijing_zjy_gc_cnc",
			"jinan9_gc_cnc"
		],
		"haerbin2_scc": [
			"beijing4_dxt"
		],
		"haerbin3_cnc": [
			"beijing_zjy_gc_cnc",
			"jinan9_gc_cnc"
		],
		"hangzhou12_wasu": [
			"vcdn_Rcache_hangzhou10_wasu",
			"hangzhou9_wasu"
		],
		"hangzhou13_wasu": [
			"vcdn_Rcache_hangzhou10_wasu",
			"hangzhou9_wasu"
		],
		"hangzhou9_wasu": [
			"vcdn_Rcache_hangzhou10_wasu",
			"beijing2_cnix"
		],
		"hcache_beijing_bgctvnet": [
			"ixp_beijing5_cnc"
		],
		"hebei_cqccn_scc": [
			"chongqing2_scc",
			"guiyang2_scc"
		],
		"hefei2_scc": [
			"ixp_beijing5_cnc",
			"shanghai_qnet"
		],
		"hefei3_scc": [
			"ixp_beijing5_cnc",
			"shanghai_qnet"
		],
		"henan_cqccn_scc": [
			"chongqing4_scc",
			"chongqing2_scc"
		],
		"hongkong2_equinix": [
			"hongkong_gc_equinix"
		],
		"hongkong2_equinix_bus_proxy": [
			"hongkong_gc_equinix"
		],
		"hongkong_bn": [
			"hongkong_megaplus"
		],
		"hongkong_cds": [
			"singapore_equinix"
		],
		"hongkong_cmi": [
			"hongkong_megaplus"
		],
		"hongkong_gc_equinix": [
			"default"
		],
		"hongkong_hgc": [
			"hongkong_megaplus"
		],
		"hongkong_hkt": [
			"hongkong_megaplus"
		],
		"hongkong_megaplus": [
			"wuhan_lkg_qnet_proxy"
		],
		"hongkong_megaplus_bus_proxy": [
			"hongkong_megaplus"
		],
		"hongkong_zenlayer": [
			"hongkong_megaplus"
		],
		"huhehaote2_scc": [
			"huhehaote3_scc"
		],
		"huhehaote3_scc": [
			"ixp_beijing5_cnc"
		],
		"iptv_Rcache_chengdu2_cmnet": [
			"wuhan_cmnet"
		],
		"iptv_Rcache_chengdu_cmnet": [
			"wuhan_cmnet"
		],
		"iptv_Rcache_henan_cnc": [
			"iptv_Rcache_tianjin_cnc"
		],
		"iptv_Rcache_tianjin_cnc": [
			"beijing_zjy_gc_cnc",
			"jinan9_gc_cnc"
		],
		"ixp_beijing5_cnc": [
			"wuhan_lkg_gc_cmnet",
			"beijing_zjy_gc_cnc"
		],
		"ixp_beijing_office": [
			"wuhan_lkg_gc_cmnet"
		],
		"jinan10_p1_cnc": [
			"beijing_zjy_gc_cnc"
		],
		"jinan10_p2_cnc": [
			"beijing_zjy_gc_cnc",
			"jinan9_gc_cnc"
		],
		"jinan9_4k_cmnet": [
			"jinan9_gc_cnc"
		],
		"jinan9_cmnet_bus_proxy": [
			"wuhan_lkg_gc_cmnet",
			"jinan9_gc_cnc"
		],
		"jinan9_cmnet_proxy": [
			"wuhan_lkg_gc_cmnet",
			"jinan9_gc_cnc"
		],
		"jinan9_cnc": [
			"jinan9_gc_cnc"
		],
		"jinan9_cnc_bus_proxy": [
			"beijing_zjy_gc_cnc",
			"jinan9_gc_cnc"
		],
		"jinan9_cnc_proxy": [
			"beijing_zjy_gc_cnc",
			"jinan9_gc_cnc"
		],
		"jinan9_ct_bus_proxy": [
			"wuhan_lkg_gc_cmnet",
			"jinan9_gc_cnc"
		],
		"jinan9_ct_proxy": [
			"jinan9_gc_cnc",
			"wuhan_lkg_gc_ct"
		],
		"jinan9_gc_cnc": [
			"beijing_zjy_gc_cnc"
		],
		"jinan_cernet": [
			"jinan9_cmnet_proxy"
		],
		"jinan_scc": [
			"beijing2_cnix"
		],
		"lanzhou2_scc": [
			"beijing2_cnix",
			"wuhan6_cnc"
		],
		"luoyang5_cnc": [
			"beijing_zjy_gc_cnc",
			"jinan9_gc_cnc"
		],
		"luoyang6_cnc": [
			"beijing_zjy_gc_cnc",
			"jinan9_gc_cnc"
		],
		"migu_cmnet": [
			"wuhan_lkg_gc_cmnet"
		],
		"nanchang_scc": [
			"foshan_kpl_qnet_proxy",
			"shanghai_qnet"
		],
		"nanning_scc": [
			"wuhan_lkg_gc_cmnet",
			"jinan9_cmnet_proxy"
		],
		"shan3xi_scc": [
			"beijing4_dxt"
		],
		"shanghai2_colnet": [
			"beijing2_cnix",
			"shanghai_qnet"
		],
		"shanghai6_gwbn": [
			"shanghai_qnet",
			"beijing11_dxt"
		],
		"shanghai_qnet": [
			"wuhan_lkg_gc_cmnet",
			"ixp_beijing5_cnc"
		],
		"shenyang2_gwbn": [
			"beijing11_dxt"
		],
		"shenyang6_cnc": [
			"beijing_zjy_gc_cnc",
			"jinan9_gc_cnc"
		],
		"shenyang_scc": [
			"beijing4_dxt"
		],
		"shenzhen2_twnet": [
			"wuhan6_cnc",
			"foshan_kpl_qnet_proxy"
		],
		"shenzhen3_gwbn": [
			"beijing11_dxt"
		],
		"shijiazhuang4_scc": [
			"beijing4_dxt",
			"ixp_beijing5_cnc"
		],
		"singapore_equinix": [
			"hongkong_megaplus"
		],
		"singapore_huawei_proxy": [
			"singapore_equinix"
		],
		"taian2_scc": [
			"beijing2_cnix",
			"ixp_beijing5_cnc"
		],
		"tianjin2_scc": [
			"ixp_beijing5_cnc",
			"shanghai_qnet"
		],
		"tonghua_cnc": [
			"beijing_zjy_gc_cnc",
			"jinan9_gc_cnc"
		],
		"vcdn_Gcache_wuxi2_scc": [
			"ixp_beijing5_cnc",
			"shanghai_qnet"
		],
		"vcdn_Rcache_haikou_scc": [
			"beijing2_cnix",
			"wuhan_cmnet"
		],
		"vcdn_Rcache_hangzhou10_wasu": [
			"beijing2_cnix",
			"wuhan6_cnc"
		],
		"vcdn_Rcache_hangzhou11_wasu": [
			"beijing2_cnix",
			"wuhan_lkg_gc_cmnet"
		],
		"vcdn_Rcache_jinan7_4k_cmnet": [
			"jinan9_gc_cnc"
		],
		"vcdn_Rcache_wuhan5_scc": [
			"wuhan2_scc"
		],
		"vcdn_Rcache_zhengzhou4_scc": [
			"ixp_beijing5_cnc",
			"shanghai_qnet"
		],
		"wuhan2_scc": [
			"ixp_beijing5_cnc",
			"shanghai_qnet"
		],
		"wuhan3_scc": [
			"wuhan2_scc"
		],
		"wuhan4_scc": [
			"wuhan2_scc"
		],
		"wuhan5_ct": [
			"wuhan_lkg_gc_cmnet",
			"beijing_zjy_gc_cnc"
		],
		"wuhan6_cnc": [
			"wuhan_lkg_gc_cmnet",
			"beijing_zjy_gc_cnc"
		],
		"wuhan_cmnet": [
			"wuhan_lkg_gc_cmnet",
			"beijing_zjy_gc_cnc"
		],
		"wuhan_lkg2_ct_proxy": [
			"wuhan_lkg_gc_ct"
		],
		"wuhan_lkg_cmnet_bus_proxy": [
			"wuhan_lkg_gc_cmnet",
			"wuhan_lkg_gc_ct"
		],
		"wuhan_lkg_cnc_bus_proxy": [
			"wuhan_lkg_gc_cmnet",
			"beijing_zjy_gc_cnc",
			"wuhan_lkg_gc_ct"
		],
		"wuhan_lkg_ct_bus_proxy": [
			"wuhan_lkg_gc_cmnet",
			"beijing_zjy_gc_cnc",
			"wuhan_lkg_gc_ct"
		],
		"wuhan_lkg_ct_proxy": [
			"beijing_zjy_gc_cnc",
			"wuhan_lkg_gc_ct"
		],
		"wuhan_lkg_gc_cmnet": [
			"foshan_kpl_ali_s3_gc_ct",
			"beijing_zjy_gc_cnc"
		],
		"wuhan_lkg_gc_ct": [
			"foshan_kpl_ali_s3_gc_ct",
			"beijing_zjy_gc_cnc"
		],
		"wuhan_lkg_qnet_proxy": [
			"wuhan_lkg_gc_ct"
		],
		"wuhan_scc": [
			"wuhan2_scc"
		],
		"wulumuqi2_scc": [
			"ixp_beijing5_cnc",
			"wuhan_cmnet"
		],
		"wulumuqi_scc": [
			"wulumuqi2_scc"
		],
		"wuxi_scc": [
			"ixp_beijing5_cnc",
			"shanghai_qnet"
		],
		"wuxi_scc_bus_proxy": [
			"vcdn_Gcache_wuxi2_scc"
		],
		"wuxi_stb_scc": [
			"shanghai_qnet"
		],
		"xian_scc": [
			"beijing4_dxt"
		],
		"xianyang2_scc": [
			"xian_scc",
			"beijing4_dxt"
		],
		"xianyang_scc": [
			"xian_scc",
			"beijing4_dxt"
		],
		"yinchuan_scc": [
			"ixp_beijing5_cnc",
			"shanghai_qnet"
		],
		"zhengzhou2_scc": [
			"ixp_beijing5_cnc"
		],
		"zhengzhou3_scc": [
			"beijing2_cnix"
		]
	}
}`

type ConnectData struct {
	Code int32               `json:"code"`
	Data map[string][]string `json:"data"`
}

func main() {
	connect_data := ConnectData{}

	err := json.Unmarshal([]byte(json_string), &connect_data)
	if err != nil {
		fmt.Println(err)
	}
	for k, v := range connect_data.Data {
		fmt.Println("k", k)
		fmt.Println("v", v)
		break
	}
}
