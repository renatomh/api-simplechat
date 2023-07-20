-- Revert changes to "users" table
ALTER TABLE "users" ALTER COLUMN "last_login_at" SET DATA TYPE timestamp;
ALTER TABLE "users" ALTER COLUMN "created_at" SET DATA TYPE timestamp;
ALTER TABLE "users" DROP COLUMN "password_changed_at";
ALTER TABLE "users" DROP COLUMN "hash_pass";
ALTER TABLE "users" RENAME COLUMN "full_name" TO "name";

-- Revert changes to "contacts" table
ALTER TABLE "contacts" ALTER COLUMN "accepted_at" SET DATA TYPE timestamp;
ALTER TABLE "contacts" ALTER COLUMN "requested_at" SET DATA TYPE timestamp;

-- Revert changes to "chats" table
ALTER TABLE "chats" ALTER COLUMN "last_message_received_at" SET DATA TYPE timestamp;

-- Revert changes to "messages" table
ALTER TABLE "messages" ALTER COLUMN "sent_at" SET DATA TYPE timestamp;
