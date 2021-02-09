#! /bin/sh

# Wait for DB services
sh ./config/run/wait-for-services.sh

# Prepare DB (Migrate - If not? Create db & Migrate)
sh ./config/run/prepare-db.sh

# Pre-comple app assets
sh ./config/run/asset-pre-compile.sh

# Start Application
bundle exec puma -C config/puma.rb