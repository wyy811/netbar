-- 网吧计费管理系统数据库初始化脚本

-- 用户表（员工/管理员）
CREATE TABLE IF NOT EXISTS `users` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `username` VARCHAR(50) NOT NULL COMMENT '登录名',
  `password_hash` VARCHAR(255) NOT NULL COMMENT '密码哈希',
  `real_name` VARCHAR(50) COMMENT '真实姓名',
  `role` TINYINT DEFAULT 1 COMMENT '角色：1=普通员工，2=管理员，3=超级管理员',
  `status` TINYINT DEFAULT 1 COMMENT '状态：1=启用，0=禁用',
  `last_login_at` DATETIME COMMENT '最后登录时间',
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_users_username` (`username`),
  KEY `idx_users_role` (`role`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='员工/管理员表';

-- 会员表
CREATE TABLE IF NOT EXISTS `members` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `card_number` VARCHAR(20) NOT NULL COMMENT '会员卡号',
  `name` VARCHAR(50) NOT NULL COMMENT '会员姓名',
  `phone` VARCHAR(15) COMMENT '手机号',
  `id_card` VARCHAR(18) COMMENT '身份证号',
  `balance` DECIMAL(10,2) DEFAULT 0.00 COMMENT '账户余额',
  `total_spent` DECIMAL(10,2) DEFAULT 0.00 COMMENT '累计消费',
  `discount_level` TINYINT DEFAULT 0 COMMENT '折扣等级：0=无，1=9折，2=8折',
  `status` TINYINT DEFAULT 1 COMMENT '状态：1=正常，0=冻结',
  `registered_at` DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '注册时间',
  `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_members_card_number` (`card_number`),
  KEY `idx_members_phone` (`phone`),
  KEY `idx_members_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='会员表';

-- 电脑机位表
CREATE TABLE IF NOT EXISTS `computers` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `machine_number` VARCHAR(10) NOT NULL COMMENT '机位编号',
  `area` VARCHAR(20) COMMENT '区域',
  `ip_address` VARCHAR(15) COMMENT '内网IP',
  `status` TINYINT DEFAULT 0 COMMENT '状态：0=空闲，1=使用中，2=维护中，3=预约中',
  `hourly_rate` DECIMAL(8,2) NOT NULL COMMENT '每小时单价',
  `current_session_id` BIGINT COMMENT '当前上机记录ID',
  `last_online_at` DATETIME COMMENT '最后上线时间',
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_computers_machine_number` (`machine_number`),
  KEY `idx_computers_status` (`status`),
  KEY `idx_computers_area` (`area`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='电脑机位表';

-- 上机记录表（核心业务表）
CREATE TABLE IF NOT EXISTS `sessions` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `computer_id` BIGINT UNSIGNED NOT NULL COMMENT '机位ID',
  `member_id` BIGINT COMMENT '会员ID',
  `user_id` BIGINT UNSIGNED NOT NULL COMMENT '操作员工ID',
  `start_time` DATETIME NOT NULL COMMENT '上机开始时间',
  `end_time` DATETIME COMMENT '下机时间',
  `duration_minutes` INT COMMENT '总时长（分钟）',
  `total_amount` DECIMAL(10,2) COMMENT '总费用',
  `discount_amount` DECIMAL(10,2) DEFAULT 0.00 COMMENT '优惠金额',
  `paid_amount` DECIMAL(10,2) COMMENT '实付金额',
  `payment_method` TINYINT COMMENT '支付方式：1=现金，2=会员余额，3=支付宝，4=微信',
  `status` TINYINT DEFAULT 0 COMMENT '状态：0=进行中，1=已结束，2=异常中断',
  `note` VARCHAR(255) COMMENT '备注',
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_sessions_computer_id` (`computer_id`),
  KEY `idx_sessions_member_id` (`member_id`),
  KEY `idx_sessions_start_time` (`start_time`),
  KEY `idx_sessions_status` (`status`),
  CONSTRAINT `fk_sessions_computer` FOREIGN KEY (`computer_id`) REFERENCES `computers` (`id`),
  CONSTRAINT `fk_sessions_member` FOREIGN KEY (`member_id`) REFERENCES `members` (`id`),
  CONSTRAINT `fk_sessions_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='上机记录表';

-- 计费规则表
CREATE TABLE IF NOT EXISTS `rate_rules` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `rule_name` VARCHAR(50) NOT NULL COMMENT '规则名称',
  `day_of_week` TINYINT COMMENT '星期几：1-7，NULL表示每天',
  `start_time` TIME NOT NULL COMMENT '时段开始',
  `end_time` TIME NOT NULL COMMENT '时段结束',
  `hourly_rate` DECIMAL(8,2) NOT NULL COMMENT '该时段单价',
  `is_overnight` TINYINT DEFAULT 0 COMMENT '是否通宵场',
  `priority` TINYINT DEFAULT 1 COMMENT '优先级（数值越大越优先）',
  `enabled` TINYINT DEFAULT 1 COMMENT '是否启用',
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_rate_rules_day_time` (`day_of_week`, `start_time`, `end_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='计费规则表';

-- 充值记录表
CREATE TABLE IF NOT EXISTS `recharges` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `member_id` BIGINT UNSIGNED NOT NULL COMMENT '会员ID',
  `user_id` BIGINT UNSIGNED NOT NULL COMMENT '操作员工ID',
  `amount` DECIMAL(10,2) NOT NULL COMMENT '充值金额',
  `bonus_amount` DECIMAL(10,2) DEFAULT 0.00 COMMENT '赠送金额',
  `payment_method` TINYINT NOT NULL COMMENT '支付方式：1=现金，2=支付宝，3=微信',
  `before_balance` DECIMAL(10,2) NOT NULL COMMENT '充值前余额',
  `after_balance` DECIMAL(10,2) NOT NULL COMMENT '充值后余额',
  `remark` VARCHAR(255) COMMENT '备注',
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_recharges_member_id` (`member_id`),
  KEY `idx_recharges_created_at` (`created_at`),
  CONSTRAINT `fk_recharges_member` FOREIGN KEY (`member_id`) REFERENCES `members` (`id`),
  CONSTRAINT `fk_recharges_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='充值记录表';

-- 操作日志表
CREATE TABLE IF NOT EXISTS `audit_logs` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id` BIGINT UNSIGNED NOT NULL COMMENT '操作人ID',
  `action` VARCHAR(50) NOT NULL COMMENT '操作类型',
  `target_type` VARCHAR(30) COMMENT '目标类型',
  `target_id` BIGINT COMMENT '目标ID',
  `detail` JSON COMMENT '详细数据',
  `ip_address` VARCHAR(45) COMMENT '操作IP',
  `user_agent` VARCHAR(255) COMMENT '客户端信息',
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_audit_logs_user_id` (`user_id`),
  KEY `idx_audit_logs_created_at` (`created_at`),
  KEY `idx_audit_logs_action` (`action`),
  CONSTRAINT `fk_audit_logs_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='操作日志表';