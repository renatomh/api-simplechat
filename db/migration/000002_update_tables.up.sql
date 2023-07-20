-- Apply changes to "users" table
ALTER TABLE "users" RENAME COLUMN "name" TO "full_name";
ALTER TABLE "users" ADD COLUMN "hash_pass" varchar NOT NULL;
ALTER TABLE "users" ADD COLUMN "password_changed_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z';
ALTER TABLE "users" ALTER COLUMN "created_at" SET DATA TYPE timestamptz;
ALTER TABLE "users" ALTER COLUMN "last_login_at" SET DATA TYPE timestamptz;

-- Apply changes to "contacts" table
ALTER TABLE "contacts" ALTER COLUMN "requested_at" SET DATA TYPE timestamptz;
ALTER TABLE "contacts" ALTER COLUMN "accepted_at" SET DATA TYPE timestamptz;

-- Apply changes to "chats" table
ALTER TABLE "chats" ALTER COLUMN "last_message_received_at" SET DATA TYPE timestamptz;

-- Apply changes to "messages" table
ALTER TABLE "messages" ALTER COLUMN "sent_at" SET DATA TYPE timestamptz;
