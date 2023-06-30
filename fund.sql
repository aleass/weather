
SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for df_fund_day_earnings
-- ----------------------------
DROP TABLE IF EXISTS `df_fund_day_earnings`;
CREATE TABLE `df_fund_day_earnings`  (
  `id` int NOT NULL AUTO_INCREMENT,
  `date` int NULL DEFAULT NULL COMMENT '日期',
  `code` varchar(255) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL DEFAULT NULL,
  `name` varchar(255) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL DEFAULT NULL,
  `unit_NV` decimal(8, 4) NULL DEFAULT NULL COMMENT '单位净值',
  `total_NV` decimal(8, 4) NULL DEFAULT NULL COMMENT '累计净值',
  `day_incre_val` decimal(6, 4) NULL DEFAULT NULL COMMENT '日增长值',
  `day_incre_rate` decimal(3, 2) NULL DEFAULT NULL COMMENT '日增长率',
  `buy_status` varchar(32) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL DEFAULT NULL COMMENT '申购状态',
  `sell_status` varchar(32) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL DEFAULT NULL COMMENT '赎回状态',
  `service_charge` decimal(3, 2) NULL DEFAULT NULL COMMENT '手续费',
  `type` varchar(255) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL DEFAULT NULL COMMENT '债券类型',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb3 COLLATE = utf8mb3_general_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for df_fund_earnings
-- ----------------------------
DROP TABLE IF EXISTS `df_fund_earnings`;
CREATE TABLE `df_fund_earnings`  (
  `code` varchar(8) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '基金代码',
  `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '基金简称',
  `date` datetime(0) NULL DEFAULT NULL COMMENT '日期',
  `nav_per_unit` decimal(10, 4) NULL DEFAULT NULL COMMENT '单位净值',
  `cumulative_nav` decimal(10, 4) NULL DEFAULT NULL COMMENT '累计净值',
  `daily_growth_rate` decimal(6, 4) NULL DEFAULT NULL COMMENT '日增长率',
  `past_1_week` decimal(8, 4) NULL DEFAULT NULL COMMENT '近1周增长率',
  `past_1_month` decimal(8, 4) NULL DEFAULT NULL COMMENT '近1个月增长率',
  `past_3_months` decimal(8, 4) NULL DEFAULT NULL COMMENT '近3个月增长率',
  `past_6_months` decimal(8, 4) NULL DEFAULT NULL COMMENT '近6个月增长率',
  `past_1_year` decimal(8, 4) NULL DEFAULT NULL COMMENT '近1年增长率',
  `past_2_years` decimal(8, 4) NULL DEFAULT NULL COMMENT '近2年增长率',
  `past_3_years` decimal(8, 4) NULL DEFAULT NULL COMMENT '近3年增长率',
  `this_year` decimal(8, 4) NULL DEFAULT NULL COMMENT '今年来增长率',
  `since_inception` decimal(8, 4) NULL DEFAULT NULL COMMENT '成立来增长率',
  `id` int NOT NULL AUTO_INCREMENT,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for df_fund_list
-- ----------------------------
DROP TABLE IF EXISTS `df_fund_list`;
CREATE TABLE `df_fund_list`  (
  `id` int NOT NULL AUTO_INCREMENT,
  `code` varchar(255) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL DEFAULT NULL COMMENT '代码',
  `name` varchar(255) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL DEFAULT NULL COMMENT '名字',
  `date` int NULL DEFAULT NULL COMMENT '日期',
  `pinyin` varchar(255) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb3 COLLATE = utf8mb3_general_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for df_fund_star
-- ----------------------------
DROP TABLE IF EXISTS `df_fund_star`;
CREATE TABLE `df_fund_star`  (
  `id` int NOT NULL AUTO_INCREMENT,
  `code` varchar(255) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL DEFAULT NULL,
  `name` varchar(255) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL DEFAULT NULL,
  `update_time` datetime(0) NULL DEFAULT NULL,
  `ZhaoShang_Securities_star` varchar(16) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '招商证券-星',
  `ZhaoShang_Securities_trend` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '招商证券-趋势 up down',
  `Shanghai_Securities_star` varchar(16) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '上海证券-星',
  `Shanghai_Securities_trend` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '上海证券-趋势 up down',
  `Jianan_Jinxin_star` varchar(16) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '济安金信-星',
  `Jianan_Jinxin_trend` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '济安金信-趋势 up down',
  `ZhaoShang_Securities_date` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '招商证券-更新时间',
  `Shanghai_Securities_date` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '上海证券-更新时间',
  `Jianan_Jinxin_Securities_date` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '济安金信-更新时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

SET FOREIGN_KEY_CHECKS = 1;
