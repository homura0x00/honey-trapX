CREATE DATABASE IF NOT EXISTS `honey_db`;

USE `honey_db`;

# 用户表
CREATE TABLE IF NOT EXISTS `users`
(
    `id` BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT 'ID',
    `user_account` VARCHAR(256) NOT NULL COMMENT '账户',
    `user_password` VARCHAR(512) NOT NULL COMMENT '用户密码',
    `user_name` VARCHAR(256) NOT NULL COMMENT '用户名',
    `user_role` VARCHAR(255) NOT NULL COMMENT '权限角色: admin/user',
    `last_login_at` VARCHAR(256) COMMENT '最后登陆时间',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` DATETIME DEFAULT NULL COMMENT '删除时间（软删除标记）',
    index idx_user_account(`user_account`)  # 快查字段，检索用户名
) COMMENT '用户表' COLLATE = utf8mb4_unicode_ci;

# 系统日志表
CREATE TABLE IF NOT EXISTS `logs`
(
    `id` BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT 'ID',
    `user_id` BIGINT COMMENT '用户id',
    `user_name` VARCHAR(32) COMMENT '用户名称',
    `user_account` VARCHAR(32) COMMENT '用户账户',
    `ip` VARCHAR(32) COMMENT '用户ip',
    `addr` VARCHAR(255) COMMENT '当前所在地址',
    `login_status` TINYINT DEFAULT 0 COMMENT '当前登陆状态：0-未登陆/1-已登陆',
    `service_title` VARCHAR(255) COMMENT '业务操作名',
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间（操作时间）'
) COMMENT '系统日志表' COLLATE = utf8mb4_unicode_ci;

# 蜜罐节点表
CREATE TABLE IF NOT EXISTS `honey_posts`
(
    `id` BIGINT PRIMARY KEY COMMENT 'ID',
    `pot_title` VARCHAR(64) NOT NULL COMMENT '节点名称',
    `pot_type` VARCHAR(32) NOT NULL COMMENT '蜜罐类型（SSH/FTP/Redis',
    `status` TINYINT DEFAULT 0 NOT NULL COMMENT '节点状态：0-close/1-open',
    `listen_port` INT NOT NULL COMMENT '监听IP',
    `created_at` DATETIME NOT NULL COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL COMMENT '最后心跳时间（监听存活）'
) COMMENT '蜜罐节点表' COLLATE = utf8mb4_unicode_ci;

# 攻击事件表
CREATE TABLE IF NOT EXISTS `attack_events`
(
    `id` BIGINT PRIMARY KEY COMMENT 'ID',
    `pot_id` BIGINT NOT NULL COMMENT '关联蜜罐节点（外键）',
    `attacker_ip` varchar(32) COMMENT '攻击者ip',
    `attacker_port` INT COMMENT '攻击者端口',
    `attack_type` varchar(32) COMMENT '攻击类型（暴力破解/注入）',
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间（攻击时间）'
) COMMENT '攻击事件表' COLLATE = utf8mb4_unicode_ci;

# 攻击者的用户画像表
CREATE TABLE IF NOT EXISTS `attack_profile`
(
    `id` BIGINT PRIMARY KEY COMMENT 'ID',
    `attacker_name` VARCHAR(32) NOT NULL COMMENT '攻击者名称',
    `create_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE  CURRENT_TIMESTAMP COMMENT '更新时间'
) COMMENT 'hacker画像表' COLLATE = utf8mb4_unicode_ci;

# 业务容器表，在容器中所挂载的业务
CREATE TABLE IF NOT EXISTS `service`
(
    `id` BIGINT PRIMARY KEY COMMENT 'ID',
    `title` VARCHAR(32) COMMENT '业务名',
    `host` VARCHAR(32) COMMENT '地址',
    `port` VARCHAR(32) COMMENT '监听端口',
    `status` TINYINT DEFAULT 0 COMMENT '当前状态：0-closed/1-open',
    INDEX idx_title(`title`)    # 快查字段，检索业务名
) COMMENT '业务容器表' COLLATE = utf8mb4_unicode_ci;