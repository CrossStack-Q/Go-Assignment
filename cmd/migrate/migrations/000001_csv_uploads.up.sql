CREATE TABLE csv_uploads (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    unique_code TEXT NOT NULL,
    input_file_path TEXT NOT NULL,
    output_file_path TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);