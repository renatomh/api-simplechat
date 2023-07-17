CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "name" varchar NOT NULL,
  "username" varchar UNIQUE NOT NULL,
  "email" varchar UNIQUE,
  "avatar_url" varchar,
  "last_login_at" timestamp
);

CREATE TABLE "contacts" (
  "id" bigserial PRIMARY KEY,
  "from_user_id" bigint NOT NULL,
  "to_user_id" bigint NOT NULL,
  "status" varchar NOT NULL,
  "requested_at" timestamp NOT NULL DEFAULT (now()),
  "accepted_at" timestamp
);

CREATE TABLE "chats" (
  "id" bigserial PRIMARY KEY,
  "from_user_id" bigint NOT NULL,
  "to_user_id" bigint NOT NULL,
  "last_message_received_at" timestamp
);

CREATE TABLE "messages" (
  "id" bigserial PRIMARY KEY,
  "chat_id" bigint NOT NULL,
  "from_user_id" bigint NOT NULL,
  "to_user_id" bigint NOT NULL,
  "body" varchar NOT NULL,
  "sent_at" timestamp NOT NULL DEFAULT (now())
);

CREATE INDEX ON "users" ("username");

CREATE INDEX ON "users" ("email");

CREATE INDEX ON "contacts" ("from_user_id");

CREATE INDEX ON "contacts" ("to_user_id");

CREATE INDEX ON "contacts" ("status");

CREATE INDEX ON "chats" ("from_user_id");

CREATE INDEX ON "chats" ("to_user_id");

CREATE INDEX ON "chats" ("from_user_id", "to_user_id");

CREATE INDEX ON "messages" ("from_user_id");

CREATE INDEX ON "messages" ("to_user_id");

CREATE INDEX ON "messages" ("body");

CREATE INDEX ON "messages" ("from_user_id", "to_user_id");

COMMENT ON COLUMN "contacts"."from_user_id" IS 'The from/to order makes no difference here';

COMMENT ON COLUMN "contacts"."to_user_id" IS 'The from/to order makes no difference here';

COMMENT ON COLUMN "contacts"."status" IS 'Pending, Accepted or Rejected';

COMMENT ON COLUMN "chats"."from_user_id" IS 'The from/to order makes no difference here';

COMMENT ON COLUMN "chats"."to_user_id" IS 'The from/to order makes no difference here';

ALTER TABLE "contacts" ADD FOREIGN KEY ("from_user_id") REFERENCES "users" ("id");

ALTER TABLE "contacts" ADD FOREIGN KEY ("to_user_id") REFERENCES "users" ("id");

ALTER TABLE "chats" ADD FOREIGN KEY ("from_user_id") REFERENCES "users" ("id");

ALTER TABLE "chats" ADD FOREIGN KEY ("to_user_id") REFERENCES "users" ("id");

ALTER TABLE "messages" ADD FOREIGN KEY ("chat_id") REFERENCES "chats" ("id");

ALTER TABLE "messages" ADD FOREIGN KEY ("from_user_id") REFERENCES "users" ("id");

ALTER TABLE "messages" ADD FOREIGN KEY ("to_user_id") REFERENCES "users" ("id");
