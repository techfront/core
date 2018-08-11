/* SQL migration AddTweetedMailedAt */
alter table topics add column tw_posted_at timestamp;
alter table topics add column vk_posted_at timestamp;
alter table topics add column fb_posted_at timestamp;
alter table topics add column newsletter_at timestamp;