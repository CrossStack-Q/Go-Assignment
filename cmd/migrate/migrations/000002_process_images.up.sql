CREATE TABLE processed_images (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    csv_upload_id UUID REFERENCES csv_uploads(id) ON DELETE CASCADE,
    product_name TEXT NOT NULL,
    input_image_urls TEXT[] NOT NULL,
    output_image_urls TEXT[] NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);