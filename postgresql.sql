CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "name" varchar NOT NULL,
  "username" varchar UNIQUE NOT NULL,
  "email" varchar UNIQUE,
  "avatar_url" varchar,
  "last_login_at" timestamp
);

CREATE TABLE "contact" (
  "id" bigserial PRIMARY KEY,
  "from_user_id" bigint NOT NULL,
  "to_user_id" bigint NOT NULL,
  "status" varchar NOT NULL,
  "requested_at" timestamp NOT NULL DEFAULT (now()),
  "accepted_at" timestamp
);

CREATE TABLE "chat" (
  "id" bigserial PRIMARY KEY,
  "from_user_id" bigint NOT NULL,
  "to_user_id" bigint NOT NULL,
  "last_message_received_at" timestamp
);

CREATE TABLE "message" (
  "id" bigserial PRIMARY KEY,
  "chat_id" bigint NOT NULL,
  "from_user_id" bigint NOT NULL,
  "to_user_id" bigint NOT NULL,
  "body" varchar NOT NULL,
  "sent_at" timestamp NOT NULL DEFAULT (now())
);

CREATE INDEX ON "users" ("username");

CREATE INDEX ON "users" ("email");

CREATE INDEX ON "contact" ("from_user_id");

CREATE INDEX ON "contact" ("to_user_id");

CREATE INDEX ON "contact" ("status");

CREATE INDEX ON "chat" ("from_user_id");

CREATE INDEX ON "chat" ("to_user_id");

CREATE INDEX ON "chat" ("from_user_id", "to_user_id");

CREATE INDEX ON "message" ("from_user_id");

CREATE INDEX ON "message" ("to_user_id");

CREATE INDEX ON "message" ("body");

CREATE INDEX ON "message" ("from_user_id", "to_user_id");

COMMENT ON COLUMN "contact"."from_user_id" IS 'The from/to order makes no difference here';

COMMENT ON COLUMN "contact"."to_user_id" IS 'The from/to order makes no difference here';

COMMENT ON COLUMN "contact"."status" IS 'Pending, Accepted or Rejected';

COMMENT ON COLUMN "chat"."from_user_id" IS 'The from/to order makes no difference here';

COMMENT ON COLUMN "chat"."to_user_id" IS 'The from/to order makes no difference here';

ALTER TABLE "contact" ADD FOREIGN KEY ("from_user_id") REFERENCES "users" ("id");

ALTER TABLE "contact" ADD FOREIGN KEY ("to_user_id") REFERENCES "users" ("id");

ALTER TABLE "chat" ADD FOREIGN KEY ("from_user_id") REFERENCES "users" ("id");

ALTER TABLE "chat" ADD FOREIGN KEY ("to_user_id") REFERENCES "users" ("id");

ALTER TABLE "message" ADD FOREIGN KEY ("chat_id") REFERENCES "chat" ("id");

ALTER TABLE "message" ADD FOREIGN KEY ("from_user_id") REFERENCES "users" ("id");

ALTER TABLE "message" ADD FOREIGN KEY ("to_user_id") REFERENCES "users" ("id");
