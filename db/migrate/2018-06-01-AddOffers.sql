/* 2018/06/01 */

CREATE TABLE tf_offer (
offer_id SERIAL NOT NULL,
offer_created_at timestamp,
offer_updated_at timestamp,
offer_name text,
offer_text text,
offer_thumbnail text,
offer_status integer NOT NULL DEFAULT 100,
offer_id_format integer NOT NULL DEFAULT 0,
offer_id_user integer,
offer_count_upvote integer,
offer_count_downvote integer,
offer_count_flag integer,
offer_count_comment integer
);

ALTER TABLE tf_offer OWNER TO techfront_user;

ALTER TABLE tf_upvote ADD COLUMN upvote_id_offer integer;
ALTER TABLE tf_downvote ADD COLUMN downvote_id_offer integer;
ALTER TABLE tf_flag ADD COLUMN flag_id_offer integer;

ALTER TABLE tf_user ADD COLUMN user_count_offer integer;

ALTER TABLE tf_user_favorite ADD COLUMN user_favorite_id_offer integer;