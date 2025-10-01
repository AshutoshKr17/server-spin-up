-- Initialize database with extensions and basic data

-- Create extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_trgm";

-- Create a default admin user (optional)
INSERT INTO users (email, name, api_key, is_active, created_at, updated_at)
VALUES (
    'admin@gpu-cloud-manager.com',
    'System Administrator',
    'admin-api-key-' || uuid_generate_v4(),
    true,
    NOW(),
    NOW()
) ON CONFLICT (email) DO NOTHING;

-- Create indexes for better performance
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_instances_user_provider_status 
    ON instances(user_id, provider, status);

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_instances_provider_status_created 
    ON instances(provider, status, created_at);

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_user_providers_user_provider_enabled 
    ON user_providers(user_id, provider, is_enabled);

-- Create full-text search index on instance names
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_instances_name_gin 
    ON instances USING gin(name gin_trgm_ops);

-- Add comments to tables
COMMENT ON TABLE users IS 'System users who can manage GPU instances';
COMMENT ON TABLE user_providers IS 'User configurations for different GPU providers';
COMMENT ON TABLE instances IS 'GPU instances managed by the system';

-- Grant permissions (adjust as needed for your setup)
-- GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO your_app_user;
-- GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO your_app_user;
