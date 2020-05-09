DROP TABLE IF EXISTS `sms_black_white`;
CREATE TABLE `sms_black_white` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增id',
  `phone` varchar(50) DEFAULT NULL COMMENT '手机号',
  `bw_type` int(1) DEFAULT '0' COMMENT '1 白名单 2 黑名单',
  `bw_status` int(1) DEFAULT '0' COMMENT '1 有效 2 无效',
  `create_datetime` int(20) DEFAULT '0' COMMENT '创建时间',
  `update_datetime` int(20) DEFAULT '0' COMMENT '修改时间',
  `creator` int(10) DEFAULT '0' COMMENT '创建人id',
  `updator` int(10) DEFAULT '0' COMMENT '更新人id',
  PRIMARY KEY (`id`),
  KEY `idx_phone` (`phone`,`bw_type`,`bw_status`) USING BTREE
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb4;

DROP TABLE IF EXISTS `sms_carrier`;
CREATE TABLE `sms_carrier` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增id',
  `carrier_name` varchar(100) DEFAULT NULL COMMENT '运营商名字',
  `create_datetime` int(20) DEFAULT '0' COMMENT '创建时间',
  `update_datetime` int(20) DEFAULT '0' COMMENT '更新时间',
  `carrier_tpl_header` varchar(100) DEFAULT NULL COMMENT '运营商模板头',
  `carrier_tpl_suffix` varchar(200) DEFAULT NULL COMMENT '运营商模板后缀',
  `carrier_tpl_params` varchar(1000) DEFAULT NULL COMMENT '模板参数集合以,隔开',
  `carrier_gateway` varchar(1000) DEFAULT NULL COMMENT '运营商的api网关',
  `carrier_appid` varchar(100) DEFAULT NULL COMMENT '运营商认证账号',
  `carrier_sid` varchar(100) DEFAULT NULL,
  `carrier_token` varchar(500) DEFAULT NULL COMMENT '运营商认证密码',
  `carrier_status` tinyint(1) DEFAULT '0' COMMENT '1 有效 2 无效',
  `creator` int(10) DEFAULT '0' COMMENT '创建人',
  `updator` int(10) DEFAULT '0' COMMENT '更新人',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

DROP TABLE IF EXISTS `sms_tpl`;
CREATE TABLE `sms_tpl` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增id',
  `carrier_id` int(10) DEFAULT '0' COMMENT '运营商id （0 全部运营商支持）',
  `carrier_tpl_id` varchar(100) DEFAULT '0' COMMENT '运营商模板id',
  `carrier_tpl_content` varchar(1000) DEFAULT NULL COMMENT '运营商模板内容',
  `carrier_tpl_type` tinyint(1) DEFAULT '0' COMMENT '1 验证码 2 营销',
  `carrier_tpl_header` varchar(100) DEFAULT NULL COMMENT '运营商模板头',
  `carrier_tpl_suffix` varchar(200) DEFAULT NULL COMMENT '运营商模板后缀',
  `carrier_tpl` int(10) DEFAULT '0' COMMENT '内部系统模板id',
  `tpl_status` tinyint(1) DEFAULT '0' COMMENT '1 有效 2 无效',
  `create_datetime` int(20) DEFAULT '0' COMMENT '创建时间',
  `update_datetime` int(20) DEFAULT '0' COMMENT '更新时间',
  `creator` int(10) DEFAULT '0' COMMENT '创建人id',
  `updator` int(10) DEFAULT '0' COMMENT '更新人id',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb4;

DROP TABLE IF EXISTS `cj_sms_log_202005`;
CREATE TABLE `cj_sms_log_202005` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `phone` varchar(32) NOT NULL COMMENT '短信手机号',
  `content` varchar(500) NOT NULL DEFAULT '' COMMENT '短信内容',
  `prov_smsid` varchar(64) DEFAULT NULL COMMENT '短信提供方的消息id',
  `status` tinyint(4) NOT NULL DEFAULT '1' COMMENT '发送短信结果状态，1成功，0失败',
  `provider` varchar(32) DEFAULT NULL COMMENT '短信提供商，zhuoyun卓云，ucpaas云之讯',
  `create_datetime` bigint(20) DEFAULT NULL COMMENT '创建时间',
  `prov_time` varchar(64) DEFAULT NULL COMMENT '短息提供商记录的短信时间',
  `return_str` text COMMENT '发短信接口返回的字符串',
  `tpl` varchar(64) DEFAULT NULL COMMENT '短信模板关键字',
  `ip` varchar(32) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `smsid` (`prov_smsid`),
  KEY `idx_phone` (`phone`) USING BTREE
) ENGINE=InnoDB  DEFAULT CHARSET=utf8 COMMENT='短信记录表';