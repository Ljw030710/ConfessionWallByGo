CREATE TABLE user (
    `id` BIGINT UNSIGNED AUTO_INCREMENT COMMENT '自增ID',
    `username` VARCHAR(20) NOT NULL COMMENT '用户名',
    `password` VARCHAR(255) NOT NULL COMMENT '密码',
    `nickname` VARCHAR(50) NOT NULL COMMENT '昵称',
    `created_at` TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
    `updated_at` TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
    `deleted_at` BIGINT NOT NULL DEFAULT 0 COMMENT '删除时间（软删除）',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='用户表';

CREATE TABLE `file_records` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增ID',
  `file_hash` VARCHAR(64) NOT NULL COMMENT '文件内容MD5指纹',
  `file_url` VARCHAR(255) NOT NULL COMMENT '七牛云完整外链',
  `created_at` TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `updated_at` TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_file_hash` (`file_hash`) -- 唯一索引是去重的关键
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='文件去重记录表';

CREATE TABLE `confessions` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增ID',
  `sender_id` BIGINT UNSIGNED NOT NULL COMMENT '发送者ID',
  `receiver_name` VARCHAR(64) NOT NULL COMMENT '接收者昵称/姓名',
  `content` TEXT NOT NULL COMMENT '表白正文内容',
  `image_url` VARCHAR(255) DEFAULT '' COMMENT '图片外链(来自file_records表)',
  `is_anonymous` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否匿名: 0-公开, 1-匿名',
  `status` TINYINT NOT NULL DEFAULT 1 COMMENT '状态: 1-正常显示, 2-审核中, 3-已删除',
  `created_at` TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `updated_at` TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_sender_id` (`sender_id`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='表白帖子主表';

CREATE TABLE `user_blocks` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增ID',
  `blocker_id` BIGINT UNSIGNED NOT NULL COMMENT '拉黑发起人ID',
  `blocked_id` BIGINT UNSIGNED NOT NULL COMMENT '被拉黑用户ID',
  `created_at` TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `updated_at` TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_block_pair` (`blocker_id`, `blocked_id`),
  KEY `idx_blocker_id` (`blocker_id`),
  KEY `idx_blocked_id` (`blocked_id`),
  CONSTRAINT `chk_not_self_block` CHECK (`blocker_id` <> `blocked_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户拉黑关系表';

CREATE TABLE `confession_comments` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增ID',
  `confession_id` BIGINT UNSIGNED NOT NULL COMMENT '表白ID',
  `username` VARCHAR(20) NOT NULL COMMENT '评论用户名',
  `content` VARCHAR(500) NOT NULL COMMENT '评论内容',
  `created_at` TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `updated_at` TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  `deleted_at` BIGINT NOT NULL DEFAULT 0 COMMENT '删除时间 (软删除)',
  PRIMARY KEY (`id`),
  KEY `idx_confession_id` (`confession_id`),
  KEY `idx_username` (`username`),
  CONSTRAINT `fk_comment_confession_id` FOREIGN KEY (`confession_id`) REFERENCES `confessions` (`id`),
  CONSTRAINT `fk_comment_username` FOREIGN KEY (`username`) REFERENCES `user` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='表白评论表';
