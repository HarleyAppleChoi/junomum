CREATE TABLE account(
    address TEXT UNIQUE PRIMARY KEY NOT NULL
);

CREATE TABLE locked_account(
    account_address TEXT UNIQUE PRIMARY KEY NOT NULL REFERENCES account(address),
    locked_address TEXT NOT NULL,
    balance NUMERIC NOT NULL,
    unlock_limit NUMERIC NOT NULL
);

CREATE TABLE delegator_account(
    account_address TEXT UNIQUE PRIMARY KEY NOT NULL REFERENCES account(address),
	delegator_id    INTEGER NOT NULL,
	delegator_node_id   TEXT NOT NULL,
	delegator_node_info TEXT NOT NULL
)