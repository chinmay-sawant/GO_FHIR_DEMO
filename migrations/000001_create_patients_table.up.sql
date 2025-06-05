CREATE TABLE IF NOT EXISTS patients (
    id SERIAL PRIMARY KEY,
    fhir_data JSONB NOT NULL,
    active BOOLEAN,
    family VARCHAR(255),
    given VARCHAR(255),
    gender VARCHAR(20),
    birth_date DATE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_patients_active ON patients(active);
CREATE INDEX IF NOT EXISTS idx_patients_family ON patients(family);
CREATE INDEX IF NOT EXISTS idx_patients_given ON patients(given);
CREATE INDEX IF NOT EXISTS idx_patients_gender ON patients(gender);
CREATE INDEX IF NOT EXISTS idx_patients_birth_date ON patients(birth_date);
CREATE INDEX IF NOT EXISTS idx_patients_deleted_at ON patients(deleted_at);

-- Create a GIN index for JSONB data for efficient querying
CREATE INDEX IF NOT EXISTS idx_patients_fhir_data ON patients USING GIN(fhir_data);

-- Create updated_at trigger
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_patients_updated_at 
    BEFORE UPDATE ON patients 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();
