-- +goose Up
-- +goose StatementBegin

-- ----------------------------
CREATE TABLE `users` (
                         `id` bigint unsigned NOT NULL AUTO_INCREMENT,
                         `created_at` datetime(3) DEFAULT NULL,
                         `updated_at` datetime(3) DEFAULT NULL,
                         `deleted_at` datetime(3) DEFAULT NULL,
                         `f_user_name` varchar(255) NOT NULL COMMENT '用户名',
                         `f_password` varchar(255) NOT NULL COMMENT '密码',
                         PRIMARY KEY (`id`),
                         INDEX `idx_users_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
-- +goose StatementEnd
