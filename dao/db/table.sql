-- 数据库表的建表语句 --
-- admin_models --
CREATE TABLE `admin_models` (
                                `id` int unsigned NOT NULL AUTO_INCREMENT,
                                `name` varchar(255) NOT NULL,
                                `phone` varchar(255) NOT NULL,
                                `password` varchar(255) NOT NULL,
                                `level` int NOT NULL DEFAULT '1' COMMENT '1=普通, 2=超级, 3=只读',
                                `email` varchar(255) NOT NULL,
                                `direction` int NOT NULL COMMENT '1 go 2 java 3 前端',
                                `created_at` datetime NULL,
                                `updated_at` datetime NULL,
                                `deleted_at` datetime NULL,
                                PRIMARY KEY (`id`),
                                KEY `idx_admin_models_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
-- messages --
CREATE TABLE `messages` (
                            `id` bigint NOT NULL AUTO_INCREMENT,
                            `send_id` int unsigned NOT NULL,
                            `receive_id` int unsigned NOT NULL,
                            `title` varchar(100) NOT NULL,
                            `content` text NOT NULL,
                            `created_at` datetime NULL,
                            `is_read` tinyint(1) NOT NULL DEFAULT '0',
                            `type` int NOT NULL,
                            PRIMARY KEY (`id`),
                            KEY `idx_messages_send_id` (`send_id`),
                            KEY `idx_messages_receive_id` (`receive_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
-- interview_results --
CREATE TABLE `interview_results` (
                                     `id` int unsigned NOT NULL AUTO_INCREMENT,
                                     `user_id` int unsigned NOT NULL,
                                     `round` int NOT NULL,
                                     `status` int NOT NULL DEFAULT '0' COMMENT '1=通过, 0=不通过',
                                     `comment` text,
                                     `admin_id` int unsigned DEFAULT NULL,
                                     `created_at` datetime NULL,
                                     `updated_at` datetime NULL,
                                     PRIMARY KEY (`id`),
                                     KEY `idx_interview_results_user_id` (`user_id`),
                                     KEY `idx_interview_results_round` (`round`),
                                     KEY `idx_interview_results_admin_id` (`admin_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
-- interview_slots --
CREATE TABLE `interview_slots` (
                                   `id` int unsigned NOT NULL AUTO_INCREMENT,
                                   `round` int NOT NULL,
                                   `start_time` datetime NOT NULL,
                                   `end_time` datetime NOT NULL,
                                   `num` int NOT NULL DEFAULT '0',
                                   `max_num` int NOT NULL DEFAULT '50',
                                   PRIMARY KEY (`id`),
                                   KEY `idx_interview_slots_round` (`round`),
                                   KEY `idx_interview_slots_start_time` (`start_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
-- interview_assignments --
CREATE TABLE `interview_assignments` (
                                         `id` int unsigned NOT NULL AUTO_INCREMENT,
                                         `user_id` int unsigned NOT NULL,
                                         `slot_id` int unsigned NOT NULL,
                                         `round` int NOT NULL,
                                         `direction` int NOT NULL DEFAULT '0' COMMENT '0不确定，1为Go，2为Java，3为前端，4为后端',
                                         `deleted_at` datetime NULL,
                                         PRIMARY KEY (`id`),
                                         UNIQUE KEY `slot_user` (`user_id`, `slot_id`),
                                         KEY `idx_interview_assignments_round` (`round`),
                                         KEY `idx_interview_assignments_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
-- user_models --
CREATE TABLE `user_models` (
                               `id` int unsigned NOT NULL AUTO_INCREMENT,
                               `stu_id` varchar(32) NOT NULL,
                               `name` varchar(64) NOT NULL,
                               `password` varchar(255) NOT NULL COMMENT 'bcrypt hash',
                               `phone` varchar(32) NOT NULL,
                               `email` varchar(128) NOT NULL,
                               `first_pass` int NOT NULL DEFAULT '0',
                               `second_pass` int NOT NULL DEFAULT '0',
                               `direction` int NOT NULL,
                               `create_time` datetime NULL,
                               `update_time` datetime NULL,
                               `deleted_at` datetime NULL,
                               `gender` int NOT NULL DEFAULT '1',
                               PRIMARY KEY (`id`),
                               UNIQUE KEY `uix_user_models_stu_id` (`stu_id`),
                               UNIQUE KEY `uix_user_models_phone` (`phone`),
                               KEY `idx_user_models_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;