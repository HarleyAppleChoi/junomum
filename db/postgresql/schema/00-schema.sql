CREATE TABLE block
(
    height           BIGINT UNIQUE PRIMARY KEY,
    id               TEXT NOT NULL UNIQUE,
    parent_id        TEXT NOT NULL,
    collection_guarantees JSONB NOT NULL,
    timestamp        TIMESTAMP WITHOUT TIME ZONE NOT NULL
);

CREATE INDEX block_index ON block (height);
CREATE INDEX block_id_index ON block (id);


CREATE TABLE block_seal
(
    height BIGINT NOT NULL REFERENCES block (height),
    execution_receipt_id TEXT UNIQUE,
    execution_receipt_signatures TEXT[][]
);

CREATE INDEX block_seal_index ON block_seal (height);
CREATE INDEX block_seal_execution_receipt_id_index ON block_seal (execution_receipt_id);


CREATE TABLE collection
(  height BIGINT  NOT NULL REFERENCES block (height),
  id TEXT  NOT NULL,
  processed BOOLEAN  NOT NULL ,
  transaction_id TEXT  NOT NULL UNIQUE
);

CREATE INDEX collection_index ON collection (height);
CREATE INDEX collection_transaction_id_index ON collection (transaction_id);


CREATE TABLE transaction
(
		height BIGINT NOT NULL REFERENCES block (height),
        transaction_id TEXT NOT NULL REFERENCES collection (transaction_id),

		script TEXT ,
		arguments TEXT[],
		reference_block_id TEXT,
		gas_limit BIGINT,
		proposal_key TEXT,
		payer TEXT,
		authorizers TEXT[],
		payload_signature JSONB,
		envelope_signatures JSONB
);
CREATE INDEX transaction_index ON transaction (height);


CREATE TABLE transaction_result
(  height BIGINT  NOT NULL REFERENCES block (height),
  transaction_id TEXT  NOT NULL REFERENCES collection (transaction_id),
  status TEXT  NOT NULL ,
  error TEXT 
);

CREATE INDEX transaction_result_index ON transaction_result (height);



CREATE TABLE event
(
    height BIGINT NOT NULL REFERENCES block (height),
    type TEXT,
    transaction_id TEXT REFERENCES collection (transaction_id),
    transaction_index TEXT,
    event_index BIGINT,
    value TEXT
);

CREATE INDEX event_index ON event (height);


CREATE TABLE pruning
(
    last_pruned_height BIGINT NOT NULL
);

