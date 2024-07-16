DROP DATABASE IF EXISTS datapoint;

CREATE DATABASE datapoint;

\c datapoint

SET CLIENT_ENCODING = 'UTF-8';

CREATE TABLE source (
                        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                        name TEXT NOT NULL,
                        host TEXT,
                        port INTEGER NOT NULL,
                        username TEXT NOT NULL,
                        password TEXT,
                        database_name TEXT NOT NULL,
                        driver TEXT NOT NULL
);

CREATE TABLE widget (
                        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                        name TEXT NOT NULL,
                        type TEXT NOT NULL,
                        parent_id UUID REFERENCES widget (id) ON DELETE CASCADE,
                        props JSONB,
                        query JSONB
);
