CREATE TABLE integration
(
    id      varchar primary key default gen_random_uuid(),
    person  varchar,
    login   varchar,
    request varchar
);

CREATE TABLE test
(
    id      varchar primary key default gen_random_uuid(),
    person  varchar,
    login   varchar,
    request varchar
);

CREATE TABLE traffic
(
    id      varchar primary key default gen_random_uuid(),
    person  varchar,
    login   varchar,
    request varchar
);

CREATE TABLE other
(
    id      varchar primary key default gen_random_uuid(),
    person  varchar,
    login   varchar,
    request varchar
);

CREATE TABLE signature
(
    id              varchar primary key default gen_random_uuid(),
    publicKey       varchar,
    body            varchar,
    login           varchar,
    timeOfSignature varchar
)