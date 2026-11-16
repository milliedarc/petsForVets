-- +goose Up

-- 1. A pet's breed must belong to the pet's species (no "Labrador" cats).
--    A foreign key on two columns checks the pair (breed_id, species_id)
--    exists together in breeds. A NULL breed_id skips the check, so pets
--    without a breed are still fine.
--    The target of a foreign key needs a UNIQUE constraint on exactly those
--    columns, hence the first statement.
ALTER TABLE breeds
    ADD CONSTRAINT breeds_id_species_id_key UNIQUE (id, species_id);

ALTER TABLE pets
    DROP CONSTRAINT pets_breed_id_fkey,
    ADD CONSTRAINT pets_breed_matches_species_fkey
        FOREIGN KEY (breed_id, species_id) REFERENCES breeds (id, species_id);

-- 2. Gender values match the radio buttons in the frontend.
UPDATE pets SET gender = lower(gender);
ALTER TABLE pets
    ADD CONSTRAINT pets_gender_check CHECK (gender IN ('', 'male', 'female'));

-- 3. A view is a saved SELECT you can query like a table. Every query that
--    returns pets selects from this, so the JOINs are written only once.
CREATE VIEW pet_details AS
SELECT p.id,
       p.owner_id,
       p.name,
       p.species_id,
       p.breed_id,
       p.date_of_birth,
       p.gender,
       p.neutered,
       p.colour,
       p.imported,
       p.import_country_code,
       p.microchip_number,
       p.avatar_path,
       s.name        AS species_name,
       s.pet_type_id AS pet_type_id,
       b.name        AS breed_name
FROM pets p
JOIN species s     ON s.id = p.species_id
LEFT JOIN breeds b ON b.id = p.breed_id;

-- +goose Down
DROP VIEW pet_details;

ALTER TABLE pets DROP CONSTRAINT pets_gender_check;

ALTER TABLE pets
    DROP CONSTRAINT pets_breed_matches_species_fkey,
    ADD CONSTRAINT pets_breed_id_fkey FOREIGN KEY (breed_id) REFERENCES breeds (id);

ALTER TABLE breeds DROP CONSTRAINT breeds_id_species_id_key;
