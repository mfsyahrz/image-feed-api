CREATE TABLE IF NOT EXISTS "comments"
(
  "id" BIGSERIAL PRIMARY KEY,
  "post_id" BIGSERIAL,
  "content" varchar(255) NOT NULL,
  "creator" varchar(100) NOT NULL,
  "created_date" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS "posts" 
(
  "id" BIGSERIAL PRIMARY KEY,
  "caption" varchar(255),
  "src_image" text DEFAULT NULL,
  "src_image_format" varchar(10) DEFAULT NULL,
  "display_image" text DEFAULT NULL,
  "display_image_format" varchar(10) DEFAULT NULL,
  "creator" varchar(100) NOT NULL,
  "comment_count" bigint NOT NULL DEFAULT 0,
  "created_date" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_posts_pagination ON posts (comment_count, created_date, id);
