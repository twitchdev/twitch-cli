create table events (
  id text not null primary key, event text not null, 
  json text not null, from_user text not null, 
  to_user text not null, transport text not null, 
  timestamp text not null
);
create table categories(
  id text not null primary key, category_name text not null
);
create table users(
  id text not null primary key, 
  user_login text not null, 
  display_name text not null, 
  email text not null, 
  user_type text, 
  broadcaster_type text, 
  user_description text, 
  created_at text not null, 
  category_id text, 
  modified_at text, 
  stream_language text not null default 'en', 
  title text not null default '', 
  delay int not null default 0, 
  foreign key (category_id) references categories(id)
);
create table follows (
  broadcaster_id text not null, 
  user_id text not null, 
  created_at text not null, 
  primary key (broadcaster_id, user_id), 
  foreign key (broadcaster_id) references users(id), 
  foreign key (user_id) references users(id)
);
create table blocks (
  broadcaster_id text not null, 
  user_id text not null, 
  created_at text not null, 
  primary key (broadcaster_id, user_id), 
  foreign key (broadcaster_id) references users(id), 
  foreign key (user_id) references users(id)
);
create table bans (
  broadcaster_id text not null, 
  user_id text not null, 
  created_at text not null, 
  expires_at text, 
  primary key (broadcaster_id, user_id), 
  foreign key (broadcaster_id) references users(id), 
  foreign key (user_id) references users(id)
);
create table ban_events (
  id text not null primary key, 
  event_timestamp text not null, 
  event_type text not null, 
  event_version text not null default '1.0', 
  broadcaster_id text not null, 
  user_id text not null, 
  expires_at text, 
  foreign key (broadcaster_id) references users(id), 
  foreign key (user_id) references users(id)
);
create table moderators (
  broadcaster_id text not null, 
  user_id text not null, 
  created_at text not null, 
  primary key (broadcaster_id, user_id), 
  foreign key (broadcaster_id) references users(id), 
  foreign key (user_id) references users(id)
);
create table moderator_actions (
  id text not null primary key, 
  event_timestamp text not null, 
  event_type text not null, 
  event_version text not null default '1.0', 
  broadcaster_id text not null, 
  user_id text not null, 
  foreign key (broadcaster_id) references users(id), 
  foreign key (user_id) references users(id)
);
create table editors (
  broadcaster_id text not null, 
  user_id text not null, 
  created_at text not null, 
  primary key (broadcaster_id, user_id), 
  foreign key (broadcaster_id) references users(id), 
  foreign key (user_id) references users(id)
);
create table channel_points_rewards(
  id text not null primary key, 
  broadcaster_id text not null, 
  reward_image text, 
  background_color text, 
  is_enabled boolean not null default false, 
  cost int not null default 0, 
  title text not null, 
  reward_prompt text, 
  is_user_input_required boolean default false, 
  stream_max_enabled boolean default false, 
  stream_max_count int default 0, 
  stream_user_max_enabled boolean default false, 
  stream_user_max_count int default 0, 
  global_cooldown_enabled boolean default false, 
  global_cooldown_seconds int default 0, 
  is_paused boolean default false, 
  is_in_stock boolean default true, 
  should_redemptions_skip_queue boolean default false, 
  redemptions_redeemed_current_stream int, 
  cooldown_expires_at text, 
  foreign key (broadcaster_id) references users(id)
);
create table channel_points_redemptions(
  id text not null primary key, 
  reward_id text not null, 
  broadcaster_id text not null, 
  user_id text not null, 
  user_input text, 
  redemption_status text not null, 
  redeemed_at text, 
  foreign key (reward_id) references channel_points_rewards(id), 
  foreign key (broadcaster_id) references users(id), 
  foreign key (user_id) references users(id)
);
create table streams(
  id text not null primary key, 
  broadcaster_id id text not null, 
  stream_type text not null default 'live', 
  viewer_count int not null, 
  started_at text not null, 
  is_mature boolean not null default false, 
  foreign key (broadcaster_id) references users(id)
);
create table tags(
  id text not null primary key, is_auto boolean not null default false, 
  tag_name text not null
);
create table stream_tags(
  user_id text not null, 
  tag_id text not null, 
  primary key(user_id, tag_id), 
  foreign key(user_id) references users(id), 
  foreign key(tag_id) references tags(id)
);
create table teams(
  id text not null primary key, background_image_url text, 
  banner text, created_at text not null, 
  updated_at text, info text, thumbnail_url text, 
  team_name text, team_display_name text
);
create table team_members(
  team_id text not null, 
  user_id text not null, 
  primary key (team_id, user_id) foreign key (team_id) references teams(id), 
  foreign key (user_id) references users(id)
);
create table videos(
  id text not null primary key, 
  stream_id text, 
  broadcaster_id text not null, 
  title text not null, 
  video_description text not null, 
  created_at text not null, 
  published_at text, 
  viewable text not null, 
  view_count int not null default 0, 
  duration text not null, 
  video_language text not null default 'en', 
  category_id text, 
  type text default 'archive', 
  foreign key (stream_id) references streams(id), 
  foreign key (broadcaster_id) references users(id), 
  foreign key (category_id) references categories(id) 
  );
create table stream_markers(
  id text not null primary key, 
  video_id text not null, 
  position_seconds int not null, 
  created_at text not null, 
  description text not null, 
  broadcaster_id text not null, 
  foreign key (broadcaster_id) references users(id), 
  foreign key (video_id) references videos(id)
);
create table video_muted_segments (
  video_id text not null, 
  video_offset int not null, 
  duration int not null, 
  primary key (video_id, video_offset), 
  foreign key (video_id) references videos(id)
);
create table subscriptions (
  broadcaster_id text not null, 
  user_id text not null, 
  is_gift boolean not null default false, 
  gifter_id text, 
  tier text not null default '1000', 
  created_at text not null, 
  primary key (broadcaster_id, user_id), 
  foreign key (broadcaster_id) references users(id), 
  foreign key (user_id) references users(id), 
  foreign key (gifter_id) references users(id)
);
create table drops_entitlements(
  id text not null primary key, 
  benefit_id text not null, 
  timestamp text not null, 
  user_id text not null, 
  game_id text not null, 
  foreign key (user_id) references users(id), 
  foreign key (game_id) references categories(id)
);
create table clients (
  id text not null primary key, secret text not null, 
  is_extension boolean default false, 
  name text not null
);
create table authorizations (
  id integer not null primary key AUTOINCREMENT, 
  client_id text not null, 
  user_id text, 
  token text not null unique, 
  expires_at text not null, 
  scopes text, 
  foreign key (client_id) references clients(id)
);
create table polls (
  id text not null primary key, 
  broadcaster_id text not null, 
  title text not null, 
  bits_voting_enabled boolean default false, 
  bits_per_vote int default 10, 
  channel_points_voting_enabled boolean default false, 
  channel_points_per_vote int default 10, 
  status text not null, 
  duration int not null, 
  started_at text not null, 
  ended_at text, 
  foreign key (broadcaster_id) references users(id)
);
create table poll_choices (
  id text not null primary key, 
  title text not null, 
  votes int not null default 0, 
  channel_points_votes int not null default 0, 
  bits_votes int not null default 0, 
  poll_id text not null, 
  foreign key (poll_id) references polls(id)
);
create table predictions (
  id text not null primary key, 
  broadcaster_id text not null, 
  title text not null, 
  winning_outcome_id text, 
  prediction_window int, 
  status text not null, 
  created_at text not null, 
  ended_at text, 
  locked_at text, 
  foreign key (broadcaster_id) references users(id)
);
create table prediction_outcomes (
  id text not null primary key, 
  title text not null, 
  users int not null default 0, 
  channel_points int not null default 0, 
  color text not null, 
  prediction_id text not null, 
  foreign key (prediction_id) references predictions(id)
);
create table prediction_predictions (
  prediction_id text not null, 
  user_id text not null, 
  amount int not null, 
  outcome_id text not null, 
  primary key(prediction_id, user_id), 
  foreign key(user_id) references users(id), 
  foreign key(prediction_id) references predictions(id), 
  foreign key(outcome_id) references prediction_outcomes(id)
);
create table clips (
  id text not null primary key, 
  broadcaster_id text not null, 
  creator_id text not null, 
  video_id text not null, 
  game_id text not null, 
  title text not null, 
  view_count int default 0, 
  created_at text not null, 
  duration real not null, 
  foreign key (broadcaster_id) references users(id), 
  foreign key (creator_id) references users(id)
);
