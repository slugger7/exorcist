-- TODO: remove the newest records that will prevent this constraint from being added
alter table media_tag add constraint unique_media_tag unique (media_id, tag_id);

alter table media_person add constraint unique_media_person unique (media_id, person_id);
