DROP TABLE IF EXISTS "verify_emails" CASCADE;

ALTER TABLE "users" DROP COLUMN "is_email_verified";

-- CASCADE keywords make sure that if there are recode in other table that
-- reference this table ,they will all be deleted