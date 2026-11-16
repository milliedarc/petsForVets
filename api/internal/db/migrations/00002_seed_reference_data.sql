-- +goose Up
-- Reference data every environment needs (not test users or pets).

INSERT INTO pet_types (id, name) VALUES
    ('cat',          'Cat'),
    ('dog',          'Dog'),
    ('small_animal', 'Small Animal'),
    ('reptile',      'Reptile');

INSERT INTO species (name, pet_type_id) VALUES
    ('Cat',            'cat'),
    ('Dog',            'dog'),
    ('Rabbit',         'small_animal'),
    ('Guinea pig',     'small_animal'),
    ('Hamster',        'small_animal'),
    ('Ferret',         'small_animal'),
    ('Bearded dragon', 'reptile'),
    ('Leopard gecko',  'reptile'),
    ('Tortoise',       'reptile'),
    ('Corn snake',     'reptile');

-- INSERT ... SELECT lets us look up the species id by name instead of
-- hard-coding generated uuids.
INSERT INTO breeds (name, species_id, generic)
SELECT b.name, s.id, b.generic
FROM (VALUES
    ('Cat', 'Domestic shorthair', true),
    ('Cat', 'Domestic longhair',  true),
    ('Cat', 'Bengal',             false),
    ('Cat', 'British Shorthair',  false),
    ('Cat', 'Maine Coon',         false),
    ('Cat', 'Ragdoll',            false),
    ('Cat', 'Siamese',            false),
    ('Dog', 'Mixed breed',        true),
    ('Dog', 'Border Collie',      false),
    ('Dog', 'Cocker Spaniel',     false),
    ('Dog', 'French Bulldog',     false),
    ('Dog', 'German Shepherd',    false),
    ('Dog', 'Labrador Retriever', false),
    ('Dog', 'Staffordshire Bull Terrier', false)
) AS b (species_name, name, generic)
JOIN species s ON s.name = b.species_name;

INSERT INTO medicines (name, description) VALUES
    ('Bravecto',  'Flea and tick treatment'),
    ('Drontal',   'Wormer'),
    ('Metacam',   'Anti-inflammatory pain relief'),
    ('Vetmedin',  'Heart medication'),
    ('Apoquel',   'Anti-itch medication');

-- +goose Down
DELETE FROM medicines;
DELETE FROM breeds;
DELETE FROM species;
DELETE FROM pet_types;
