CREATE TABLE account(
    address TEXT UNIQUE PRIMARY KEY NOT NULL,
    balance BIGINT NOT NULL,
    code TEXT NOT NULL,
    keys_list JSONB,
    contract_map JSONB
);

CREATE TABLE locked_account_delegator
(  address TEXT  NOT NULL NOT NULL REFERENCES account(address),
  locked_address TEXT  NOT NULL UNIQUE,
  node_id TEXT  NOT NULL UNIQUE ,
  delegator_id BIGINT  NOT NULL UNIQUE
);

CREATE TABLE locked_account_staker
(
    address TEXT  NOT NULL NOT NULL REFERENCES account(address),
    node_id TEXT  NOT NULL 
);

CREATE TABLE locked_account_balance(
    locked_address TEXT NOT NULL REFERENCES locked_account_delegator(locked_address),
    balance BIGINT NOT NULL,
    unlock_limit BIGINT NOT NULL,
    height BIGINT NOT NULL
);

CREATE INDEX locked_account_balance_index ON locked_account_balance (height);


CREATE TABLE delegator_account(
    account_address TEXT UNIQUE PRIMARY KEY NOT NULL REFERENCES account(address),
	delegator_id    BIGINT NOT NULL,
	delegator_node_id   TEXT NOT NULL
);
