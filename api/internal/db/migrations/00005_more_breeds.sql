-- +goose Up

-- The pet form tells people to pick "Cross / Mixed breed" when unsure, so
-- make sure that option exists, with exactly that name, for cats and dogs.
UPDATE breeds
SET name = 'Cross / Mixed breed'
WHERE name = 'Mixed breed'
  AND species_id = (SELECT id FROM species WHERE name = 'Dog');

-- Breeds commonly seen in UK vet practices, including popular crosses that
-- owners (and vets) refer to by name. generic = true means "listed first".
-- ON CONFLICT ... DO NOTHING skips any that already exist, using the
-- UNIQUE (species_id, name) constraint from migration 00001.
INSERT INTO breeds (name, species_id, generic)
SELECT b.name, s.id, b.generic
FROM (VALUES
    -- Cats
    ('Cat', 'Cross / Mixed breed',     true),
    ('Cat', 'Domestic mediumhair',     true),
    ('Cat', 'Abyssinian',              false),
    ('Cat', 'Birman',                  false),
    ('Cat', 'British Longhair',        false),
    ('Cat', 'Burmese',                 false),
    ('Cat', 'Burmilla',                false),
    ('Cat', 'Cornish Rex',             false),
    ('Cat', 'Devon Rex',               false),
    ('Cat', 'Egyptian Mau',            false),
    ('Cat', 'Exotic Shorthair',        false),
    ('Cat', 'Norwegian Forest Cat',    false),
    ('Cat', 'Ocicat',                  false),
    ('Cat', 'Oriental Shorthair',      false),
    ('Cat', 'Persian',                 false),
    ('Cat', 'Russian Blue',            false),
    ('Cat', 'Selkirk Rex',             false),
    ('Cat', 'Siberian',                false),
    ('Cat', 'Somali',                  false),
    ('Cat', 'Sphynx',                  false),
    ('Cat', 'Tonkinese',               false),
    ('Cat', 'Turkish Angora',          false),
    ('Cat', 'Turkish Van',             false),
    -- Dogs
    ('Dog', 'Akita',                          false),
    ('Dog', 'Australian Shepherd',            false),
    ('Dog', 'Basset Hound',                   false),
    ('Dog', 'Beagle',                         false),
    ('Dog', 'Bernese Mountain Dog',           false),
    ('Dog', 'Bichon Frise',                   false),
    ('Dog', 'Border Terrier',                 false),
    ('Dog', 'Boxer',                          false),
    ('Dog', 'Bull Terrier',                   false),
    ('Dog', 'Bulldog',                        false),
    ('Dog', 'Cane Corso',                     false),
    ('Dog', 'Cavalier King Charles Spaniel',  false),
    ('Dog', 'Cavapoo',                        false),
    ('Dog', 'Chihuahua',                      false),
    ('Dog', 'Cockapoo',                       false),
    ('Dog', 'Dachshund',                      false),
    ('Dog', 'Dalmatian',                      false),
    ('Dog', 'Dobermann',                      false),
    ('Dog', 'English Springer Spaniel',       false),
    ('Dog', 'Golden Retriever',               false),
    ('Dog', 'Great Dane',                     false),
    ('Dog', 'Greyhound',                      false),
    ('Dog', 'Hungarian Vizsla',               false),
    ('Dog', 'Jack Russell Terrier',           false),
    ('Dog', 'Labradoodle',                    false),
    ('Dog', 'Lhasa Apso',                     false),
    ('Dog', 'Lurcher',                        false),
    ('Dog', 'Maltese',                        false),
    ('Dog', 'Miniature Schnauzer',            false),
    ('Dog', 'Newfoundland',                   false),
    ('Dog', 'Pomeranian',                     false),
    ('Dog', 'Poodle',                         false),
    ('Dog', 'Pug',                            false),
    ('Dog', 'Rottweiler',                     false),
    ('Dog', 'Rough Collie',                   false),
    ('Dog', 'Shar Pei',                       false),
    ('Dog', 'Shetland Sheepdog',              false),
    ('Dog', 'Shih Tzu',                       false),
    ('Dog', 'Siberian Husky',                 false),
    ('Dog', 'Sprocker Spaniel',               false),
    ('Dog', 'Weimaraner',                     false),
    ('Dog', 'West Highland White Terrier',    false),
    ('Dog', 'Whippet',                        false),
    ('Dog', 'Yorkshire Terrier',              false)
) AS b (species_name, name, generic)
JOIN species s ON s.name = b.species_name
ON CONFLICT (species_id, name) DO NOTHING;

-- +goose Down

-- Remove the breeds this migration added, keeping the original 00002 set and
-- any breed a pet is using (deleting those would violate the foreign key).
DELETE FROM breeds b
USING species s
WHERE s.id = b.species_id
  AND s.name IN ('Cat', 'Dog')
  AND b.name NOT IN (
      'Domestic shorthair', 'Domestic longhair', 'Bengal', 'British Shorthair',
      'Maine Coon', 'Ragdoll', 'Siamese',
      'Cross / Mixed breed', 'Border Collie', 'Cocker Spaniel', 'French Bulldog',
      'German Shepherd', 'Labrador Retriever', 'Staffordshire Bull Terrier'
  )
  AND NOT EXISTS (SELECT 1 FROM pets p WHERE p.breed_id = b.id);

DELETE FROM breeds b
USING species s
WHERE s.id = b.species_id
  AND s.name = 'Cat'
  AND b.name = 'Cross / Mixed breed'
  AND NOT EXISTS (SELECT 1 FROM pets p WHERE p.breed_id = b.id);

UPDATE breeds
SET name = 'Mixed breed'
WHERE name = 'Cross / Mixed breed'
  AND species_id = (SELECT id FROM species WHERE name = 'Dog');
