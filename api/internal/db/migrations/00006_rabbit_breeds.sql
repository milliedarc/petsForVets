-- +goose Up

-- Rabbit breeds recognised in the UK, plus a generic option listed first.
-- The pet form shows the breed question for any species that has breeds,
-- so adding rows here is all another species would need.
INSERT INTO breeds (name, species_id, generic)
SELECT b.name, s.id, b.generic
FROM (VALUES
    ('Cross / Mixed breed', true),
    ('Angora',              false),
    ('Belgian Hare',        false),
    ('British Giant',       false),
    ('Californian',         false),
    ('Cashmere Lop',        false),
    ('Chinchilla',          false),
    ('Continental Giant',   false),
    ('Dutch',               false),
    ('Dwarf Lop',           false),
    ('English Lop',         false),
    ('English Spot',        false),
    ('Flemish Giant',       false),
    ('French Lop',          false),
    ('Harlequin',           false),
    ('Havana',              false),
    ('Himalayan',           false),
    ('Holland Lop',         false),
    ('Lionhead',            false),
    ('Mini Lop',            false),
    ('Mini Rex',            false),
    ('Netherland Dwarf',    false),
    ('New Zealand White',   false),
    ('Polish',              false),
    ('Rex',                 false),
    ('Satin',               false),
    ('Thrianta',            false)
) AS b (name, generic)
JOIN species s ON s.name = 'Rabbit'
ON CONFLICT (species_id, name) DO NOTHING;

-- +goose Down

-- Remove rabbit breeds, except any a pet is using.
DELETE FROM breeds b
USING species s
WHERE s.id = b.species_id
  AND s.name = 'Rabbit'
  AND NOT EXISTS (SELECT 1 FROM pets p WHERE p.breed_id = b.id);
