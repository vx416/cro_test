-- +goose Up
CREATE TABLE IF NOT EXISTS `users` (
  `id` int NOT NULL AUTO_INCREMENT,
  `email` varchar(128) UNIQUE NOT NULL,
  `password_hash` varchar(100) NOT NULL,
  `created_ati` int NOT NULL,
  `deleted_ati` int,
  PRIMARY KEY (`id`),
  KEY `idx_created_ati` (`created_ati`)
) ENGINE=InnoDB AUTO_INCREMENT=624 DEFAULT CHARSET=utf8;


CREATE TABLE IF NOT EXISTS `wallets` (
  `id` int NOT NULL AUTO_INCREMENT,
  `serial_number` varchar(128) UNIQUE NOT NULL,
  `user_id` int NOT NULL,
  `amount` decimal(55,10) NOT NULL DEFAULT '0',
  `created_ati` int NOT NULL,
  `updated_ati` int NOT NULL,
  `deleted_ati` int,
  PRIMARY KEY (`id`),
  FOREIGN KEY(user_id) REFERENCES users(id),
  KEY `idx_created_ati` (`created_ati`)
) ENGINE=InnoDB AUTO_INCREMENT=624 DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `transactions` (
  `id` int NOT NULL AUTO_INCREMENT,
  `kind` tinyint NOT NULL COMMENT '1: deposit, 2: withdrawal, 3:transfer',
  `from_wallet_id` int NOT NULL DEFAULT 0,
  `to_wallet_id` int NOT NULL DEFAULT 0,
  `tx_amount` decimal(55,10) NOT NULL DEFAULT '0',
  `created_ati` int NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_from_wallet_id_created_ati` (`from_wallet_id`, `created_ati`),
  KEY `idx_to_wallet_id_created_ati` (`to_wallet_id`, `created_ati`)
) ENGINE=InnoDB AUTO_INCREMENT=624 DEFAULT CHARSET=utf8;

-- +goose Down
DROP TABLE IF EXISTS `transactions`;
DROP TABLE IF EXISTS `wallets`;
DROP TABLE IF EXISTS `users`;