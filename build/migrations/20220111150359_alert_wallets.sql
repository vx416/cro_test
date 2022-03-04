-- +goose Up
ALTER TABLE `wallets` ADD COLUMN currency varchar(10) NOT NULL;
CREATE UNIQUE INDEX idx_user_id_currency ON wallets(user_id, currency);
-- +goose Down
ALTER TABLE `wallets` DROP COLUMN currency;
