SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for df_date_note
-- ----------------------------
DROP TABLE IF EXISTS `df_date_note`;
CREATE TABLE `df_date_note`  (
  `date` int NOT NULL AUTO_INCREMENT,
  PRIMARY KEY (`date`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 20230628 CHARACTER SET = utf8mb3 COLLATE = utf8mb3_general_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for df_fund
-- ----------------------------
DROP TABLE IF EXISTS `df_fund`;
CREATE TABLE `df_fund`  (
  `id` int NOT NULL AUTO_INCREMENT,
  `code` varchar(255) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL DEFAULT NULL COMMENT '代码',
  `name` varchar(255) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL DEFAULT NULL COMMENT '名字',
  `date` int NULL DEFAULT NULL COMMENT '日期',
  `pinyin` varchar(255) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 5326 CHARACTER SET = utf8mb3 COLLATE = utf8mb3_general_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for df_fund_earings
-- ----------------------------
DROP TABLE IF EXISTS `df_fund_earings`;
CREATE TABLE `df_fund_earings`  (
  `id` int NOT NULL AUTO_INCREMENT,
  `date` int NULL DEFAULT NULL COMMENT '日期',
  `code` varchar(255) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL DEFAULT NULL,
  `name` varchar(255) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL DEFAULT NULL,
  `unit_NV` decimal(8, 4) NULL DEFAULT NULL COMMENT '单位净值',
  `total_NV` decimal(8, 4) NULL DEFAULT NULL COMMENT '累计净值',
  `day_incre_val` decimal(3, 2) NULL DEFAULT NULL COMMENT '日增长值',
  `day_incre_rate` decimal(3, 2) NULL DEFAULT NULL COMMENT '日增长率',
  `buy_status` varchar(32) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL DEFAULT NULL COMMENT '申购状态',
  `sell_status` varchar(32) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL DEFAULT NULL COMMENT '赎回状态',
  `service_charge` decimal(3, 2) NULL DEFAULT NULL COMMENT '手续费',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 5301 CHARACTER SET = utf8mb3 COLLATE = utf8mb3_general_ci ROW_FORMAT = DYNAMIC;

SET FOREIGN_KEY_CHECKS = 1;
