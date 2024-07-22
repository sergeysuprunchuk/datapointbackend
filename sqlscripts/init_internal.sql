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

CREATE TABLE dashboard (
                           id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                           name TEXT NOT NULL
);

CREATE TABLE dashboard_widget (
                                  dashboard_id UUID REFERENCES dashboard (id) ON DELETE CASCADE,
                                  widget_id UUID REFERENCES widget (id) ON DELETE CASCADE,
                                  x SMALLINT NOT NULL,
                                  y SMALLINT NOT NULL,
                                  w SMALLINT NOT NULL,
                                  h SMALLINT NOT NULL
);

WITH RECURSIVE r AS (
    SELECT d.id "d_id", d.name "d_name",
           w.id "w_id", w.name "w_name", w.type, w.props, w.query, w.parent_id, dw.x, dw.y, dw.w, dw.h
    FROM dashboard d
             LEFT JOIN dashboard_widget dw ON d.id = dw.dashboard_id
             LEFT JOIN widget w ON dw.widget_id = w.id
    UNION
    SELECT r.d_id, r.d_name, w.id, w.name, w.type, w.props, w.query, w.parent_id, 0::SMALLINT, 0::SMALLINT, 0::SMALLINT, 0::SMALLINT
    FROM widget w
             JOIN r ON r.w_id = w.parent_id
) SELECT
      r.d_id, r.d_name,
      COALESCE(r.w_id, 'f6ea48bc-b204-408a-9826-92d791381ec6'), r.w_name, r.type, r.parent_id, r.props, r.query,
      r.x, r.y, r.w, r.h
FROM r;
;

WITH RECURSIVE r AS (
    SELECT d.id d_id, d.name d_name,
           w.id w_id, w.name w_name, w.type, w.parent_id, w.props, w.query,
           dw.x, dw.y, dw.w, dw.h
    FROM dashboard d
             LEFT JOIN dashboard_widget dw ON dw.dashboard_id = d.id
             LEFT JOIN widget w ON w.id = dw.widget_id
    UNION
    SELECT r.d_id, r.d_name, w.id, w.name, w.type, w.parent_id, w.props, w.query, 0::SMALLINT, 0::SMALLINT, 0::SMALLINT, 0::SMALLINT
    FROM widget w
             JOIN r ON r.w_id = w.parent_id
) SELECT r.d_id, r.d_name, COALESCE(r.w_id, '00000000-0000-0000-0000-000000000000'),
         COALESCE(r.w_name, ''), COALESCE(r.type,''), r.parent_id, r.props, r.query,
         COALESCE(r.x, 0), COALESCE(r.y, 0), COALESCE(r.w, 0), COALESCE(r.h, 0)
  FROM r