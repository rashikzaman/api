create or replace function trigger_set_updated_at_timestamp()
returns trigger as $$
begin
  new.updated_at = NOW();
  return new;
end;
$$ language plpgsql;
