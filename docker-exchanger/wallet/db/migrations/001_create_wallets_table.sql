-- PostgreSQL Script

CREATE SCHEMA IF NOT EXISTS mydb;

SET search_path TO mydb;

-- -----------------------------------------------------
-- Table: currencies
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS currencies (
  id SERIAL PRIMARY KEY,
  currency VARCHAR(45) NOT NULL
);

-- -----------------------------------------------------
-- Table: users
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS users (
  id SERIAL PRIMARY KEY,
  username VARCHAR(45) NOT NULL,
  email VARCHAR(256) NOT NULL,
  password VARCHAR(256) NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- -----------------------------------------------------
-- Table: wallets
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS wallets (
  id SERIAL PRIMARY KEY,
  user_id INTEGER NOT NULL,
  name VARCHAR(45) NOT NULL,
  CONSTRAINT user_fk
    FOREIGN KEY (user_id)
    REFERENCES users (id)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION
);

CREATE INDEX user_idx ON wallets (user_id);

-- -----------------------------------------------------
-- Table: balances
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS balances (
  id SERIAL PRIMARY KEY,
  balance INTEGER NOT NULL DEFAULT 0,
  wallet_id INTEGER NOT NULL,
  currency_id INTEGER NOT NULL,
  CONSTRAINT currency_fk
    FOREIGN KEY (currency_id)
    REFERENCES currencies (id)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  CONSTRAINT wallet_fk
    FOREIGN KEY (wallet_id)
    REFERENCES wallets (id)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION
);

CREATE INDEX currency_idx ON balances (currency_id);
CREATE INDEX wallet_idx ON balances (wallet_id);

INSERT INTO mydb.currencies (currency) VALUES ('RUB');
INSERT INTO mydb.currencies (currency) VALUES ('USD');
INSERT INTO mydb.currencies (currency) VALUES ('EUR');
