/*
 Navicat Premium Data Transfer

 Source Server         : oldbai
 Source Server Type    : MySQL
 Source Server Version : 80027
 Source Host           : oldbai.top:3309
 Source Schema         : biz

 Target Server Type    : MySQL
 Target Server Version : 80027
 File Encoding         : 65001

 Date: 27/02/2023 15:43:10
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for lbblog_article
-- ----------------------------
DROP TABLE IF EXISTS `lbblog_article`;
CREATE TABLE `lbblog_article`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `created_at` int NULL DEFAULT NULL,
  `updated_at` int NULL DEFAULT NULL,
  `deleted_at` int NULL DEFAULT NULL,
  `title` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL,
  `desc` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL,
  `category_id` bigint UNSIGNED NULL DEFAULT NULL,
  `img` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL,
  `content` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 2 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_bin ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of lbblog_article
-- ----------------------------
INSERT INTO `lbblog_article` VALUES (1, 1677483225, 1677483225, 0, '引言', '该项目是使用Docker进行部署，接下来是对该项目的讲解。', 1, 'https://baifile-1309918034.cos.ap-guangzhou.myqcloud.com/public/link-info/assets/images/docker.png?q-sign-algorithm=sha1&q-ak=AKID4ERy3JypHZgP7moGwwPehVo96hFLcq0RxBtWKkvicg9cY3_2UV4UKbuWxsuCtMjt&q-sign-time=1677483342;1677486942&q-key-time=1677483342;1677486942&q-header-list=host&q-url-param-list=&q-signature=fcbc36d636a9fe80fbeb7926f2e7402464bb64e5&x-cos-security-token=aT2iHFbY4kP4UZzhhWbnirOkaLyuWuga682269288618cd36a442110c20b759f8R_MphXtRiLE2pfT1tfk_lJoJ3tAb0GHnIQEQ7kH4JyTaEF1Q9wQn36W3K-v6mhI0cRVKNnGHVb9IOj2klLtF3pMV6gDGl4db803pg0Yn7jEKKJbL6qc063CcAWlCKp0I6HoP14KMuh8XeDVJzljJqmuW3_PtKfQKGXrgAI5VQtok4IDfvcUlqL-HwnJysH9b', '<h2>技术框架：</h2>\n<ul>\n<li>后端：golang</li>\n<li>前端：vue, nuxt</li>\n<li>项目链接：\n<p><a href=\"https://github.com/oldbai555/bgg\">bgg</a></p>\n</li>\n</ul>');

-- ----------------------------
-- Table structure for lbblog_category
-- ----------------------------
DROP TABLE IF EXISTS `lbblog_category`;
CREATE TABLE `lbblog_category`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `created_at` int NULL DEFAULT NULL,
  `updated_at` int NULL DEFAULT NULL,
  `deleted_at` int NULL DEFAULT NULL,
  `name` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 2 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_bin ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of lbblog_category
-- ----------------------------
INSERT INTO `lbblog_category` VALUES (1, 1677483021, 1677483021, 0, 'Docker');

-- ----------------------------
-- Table structure for lbblog_comment
-- ----------------------------
DROP TABLE IF EXISTS `lbblog_comment`;
CREATE TABLE `lbblog_comment`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `created_at` int NULL DEFAULT NULL,
  `updated_at` int NULL DEFAULT NULL,
  `deleted_at` int NULL DEFAULT NULL,
  `article_id` bigint UNSIGNED NULL DEFAULT NULL,
  `user_id` bigint UNSIGNED NULL DEFAULT NULL,
  `user_email` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL,
  `content` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL,
  `status` bigint UNSIGNED NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_bin ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of lbblog_comment
-- ----------------------------

-- ----------------------------
-- Table structure for lbuser_user
-- ----------------------------
DROP TABLE IF EXISTS `lbuser_user`;
CREATE TABLE `lbuser_user`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `created_at` int NULL DEFAULT NULL,
  `updated_at` int NULL DEFAULT NULL,
  `deleted_at` int NULL DEFAULT NULL,
  `username` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL,
  `password` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL,
  `avatar` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL,
  `nickname` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL,
  `email` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL,
  `github` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL,
  `desc` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL,
  `role` int UNSIGNED NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_bin ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of lbuser_user
-- ----------------------------
INSERT INTO `lbuser_user` VALUES (1, 1677482452, 1677483274, 0, 'superadmin', '123456', 'https://baifile-1309918034.cos.ap-guangzhou.myqcloud.com/public/link-info/assets/images/20230131-162920.webp?q-sign-algorithm=sha1&q-ak=AKIDMdBgXmJYhQRDFl7jQpHSkkhEW8SWJ45pxWAo97i99jJ2h8LTQSDsGFjEQWo2rcZA&q-sign-time=1677482696;1677486296&q-key-time=1677482696;1677486296&q-header-list=host&q-url-param-list=&q-signature=f2e9596899d7ba53538d53aa8769bf1c200b8048&x-cos-security-token=cmQy9hKYWHE2Aky9IqB24lrnPcSjSxxa1e185b23d8a0b37adf32b2a10b8bf60ftMDG9SnQ3PrH8A865yHjc-fydov9_kpZF7BPkyZeDjRnJvoJ-Wpa65XPyeNPDsDxNxXmofHTRctr4MZ85Bpp737B50JU2SVcKLA1GoYqcuwZWenmQK-ZXaJcodc7Ws7kE767ySQJcIltv4TTc_QG9OZ8XK6d6hkfBi-Jo_NIOQDhL1UhFhNE4PAHLy-EYJpJ', '大白', '934105499@qq.com', 'https://github.com/oldbai555', '一个进击的程序员', 1);

SET FOREIGN_KEY_CHECKS = 1;
