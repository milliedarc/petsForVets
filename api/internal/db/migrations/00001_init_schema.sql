-- +goose Up

-- Users log in with email + password. We never store the password itself,
-- only a bcrypt hash of it (see internal/httpapi/auth.go).
CREATE TABLE users (
    id            uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    email         text NOT NULL UNIQUE,          -- stored lower-cased by the API
    password_hash text NOT NULL,
    name          text NOT NULL DEFAULT '',
    surname       text NOT NULL DEFAULT '',
    app_mode      text NOT NULL DEFAULT 'petOwner'
                  CHECK (app_mode IN ('petOwner', 'clinicTeam')),
    veterinary    boolean NOT NULL DEFAULT false,
    created_at    timestamptz NOT NULL DEFAULT now()
);

-- One row per logged-in browser. The cookie holds a random token; we store
-- only its SHA-256 hash, so a leaked database can't be used to log in.
CREATE TABLE sessions (
    token_hash bytea PRIMARY KEY,
    user_id    uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    expires_at timestamptz NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX sessions_user_id_idx ON sessions (user_id);

-- Reference data: the four big buttons on the pet form (cat, dog, ...).
-- A readable text key instead of a uuid, so the frontend can hard-code it.
CREATE TABLE pet_types (
    id   text PRIMARY KEY,
    name text NOT NULL
);

CREATE TABLE species (
    id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name        text NOT NULL UNIQUE,
    pet_type_id text NOT NULL REFERENCES pet_types (id)
);

CREATE TABLE breeds (
    id         uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name       text NOT NULL,
    species_id uuid NOT NULL REFERENCES species (id),
    generic    boolean NOT NULL DEFAULT false,   -- e.g. "Mixed breed", listed first
    UNIQUE (species_id, name)
);

CREATE TABLE pets (
    id                  uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_id            uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    name                text NOT NULL CHECK (name <> ''),
    species_id          uuid NOT NULL REFERENCES species (id),
    breed_id            uuid REFERENCES breeds (id),        -- optional
    date_of_birth       date,                               -- optional
    gender              text NOT NULL DEFAULT '',
    neutered            boolean NOT NULL DEFAULT false,
    colour              text NOT NULL DEFAULT '',
    imported            boolean NOT NULL DEFAULT false,
    import_country_code text NOT NULL DEFAULT '',
    microchip_number    text NOT NULL DEFAULT '',
    avatar_path         text,                               -- file name on disk, if any
    created_at          timestamptz NOT NULL DEFAULT now(),
    updated_at          timestamptz NOT NULL DEFAULT now()
);
-- "Show me my pets" is the most common query, so index the owner column.
CREATE INDEX pets_owner_id_idx ON pets (owner_id);

CREATE TABLE medicines (
    id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name        text NOT NULL UNIQUE,
    description text NOT NULL DEFAULT ''
);

CREATE TABLE prescriptions (
    id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    pet_id      uuid NOT NULL REFERENCES pets (id) ON DELETE CASCADE,
    medicine_id uuid NOT NULL REFERENCES medicines (id),
    dosage      text NOT NULL,
    quantity    integer NOT NULL CHECK (quantity > 0),
    created_at  timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX prescriptions_pet_id_idx ON prescriptions (pet_id);

CREATE TABLE prescription_requests (
    id               uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    pet_id           uuid NOT NULL REFERENCES pets (id) ON DELETE CASCADE,
    prescription_id  uuid NOT NULL REFERENCES prescriptions (id) ON DELETE CASCADE,
    status           text NOT NULL DEFAULT 'Pending'
                     CHECK (status IN ('Pending', 'Approved', 'Ready', 'Rejected')),
    comments         text NOT NULL DEFAULT '',
    rejection_reason text NOT NULL DEFAULT '',
    approved_by      uuid REFERENCES users (id),
    rejected_by      uuid REFERENCES users (id),
    created_at       timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX prescription_requests_pet_id_idx ON prescription_requests (pet_id);

-- +goose Down
DROP TABLE prescription_requests;
DROP TABLE prescriptions;
DROP TABLE medicines;
DROP TABLE pets;
DROP TABLE breeds;
DROP TABLE species;
DROP TABLE pet_types;
DROP TABLE sessions;
DROP TABLE users;
