-- 创建用户表
CREATE TABLE IF NOT EXISTS `usr` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `username` varchar(20) NOT NULL COMMENT '用户名',
  `password` varchar(32) NOT NULL COMMENT '密码',
  `email` varchar(100) NOT NULL COMMENT '邮箱',
  `mobile` varchar(20) DEFAULT NULL COMMENT '手机号',
  `address` varchar(255) DEFAULT NULL COMMENT '地址',
  `token` text DEFAULT NULL COMMENT '用户令牌',
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_username` (`username`),
  UNIQUE KEY `idx_email` (`email`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- 插入测试数据
INSERT INTO `usr` (`username`, `password`, `email`, `mobile`, `address`) VALUES
('admin', '21232f297a57a5a743894a0e4a801fc3', 'admin@example.com', '+8613800138000', '北京市海淀区'),
('test', '098f6bcd4621d373cade4e832627b4f6', 'test@example.com', '+8613900139000', '上海市浦东新区');