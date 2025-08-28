#!/bin/bash

# Migration script: run from parent directory containing both repos
# Usage: ./migrate.sh

OLD_REPO="agents-at-scale"
NEW_REPO="agents-at-scale-ark"

# Get current branch from old repo
cd $OLD_REPO
CURRENT_BRANCH=$(git branch --show-current)

if [[ "$CURRENT_BRANCH" == "main" ]]; then
    echo "Error: Switch to your feature branch first"
    exit 1
fi

# Create patch of your changes
git diff main..$CURRENT_BRANCH > ../migration.patch

# Switch to new repo and create branch
cd ../$NEW_REPO
NEW_BRANCH="feature/migrated-$(echo $CURRENT_BRANCH | sed 's/.*\///')"
git checkout -b $NEW_BRANCH

# Apply patch with 3-way merge for conflict resolution
git apply --3way ../migration.patch

# Stage changes
git add .

echo "Migration complete. Next steps:"
echo "1. git commit -m 'feat: migrate feature from old repo'"
echo "2. git push -u origin $NEW_BRANCH"
echo "3. Create PR"

# Cleanup
rm ../migration.patch