#!/bin/bash

# Auto Migration Script - Runs all SQL migrations from the migrations folder
# Usage: ./migrate.sh [up|down|reset]

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
MIGRATIONS_DIR="./migrations"
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-3306}"
DB_USER="${DB_USER:-root}"
DB_PASSWORD="${DB_PASSWORD:-root}"
DB_NAME="${DB_NAME:-assmi_super_shop}"
CONTAINER_NAME="assmi-super-shop-erp-backend-db-1"

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[✓]${NC} $1"
}

print_error() {
    echo -e "${RED}[✗]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[!]${NC} $1"
}

# Function to check if migrations directory exists
check_migrations_dir() {
    if [ ! -d "$MIGRATIONS_DIR" ]; then
        print_error "Migrations directory not found: $MIGRATIONS_DIR"
        exit 1
    fi
    print_success "Migrations directory found"
}

# Function to check if Docker container is running
check_docker_container() {
    if ! docker ps | grep -q "$CONTAINER_NAME"; then
        print_warning "Docker container '$CONTAINER_NAME' is not running. Starting containers..."
        docker-compose up -d
        sleep 5
    fi
    print_success "Docker container is running"
}

# Function to run UP migrations (*.up.sql files)
run_migrations_up() {
    print_status "Running UP migrations..."
    
    # Get all up migration files sorted by number
    local migration_files=$(find "$MIGRATIONS_DIR" -name "*.up.sql" | sort)
    
    if [ -z "$migration_files" ]; then
        print_warning "No migration files found"
        return
    fi
    
    for migration_file in $migration_files; do
        local filename=$(basename "$migration_file")
        print_status "Executing: $filename"
        
        # Run the migration using docker exec
        if docker exec -i "$CONTAINER_NAME" mysql -u "$DB_USER" -p"$DB_PASSWORD" "$DB_NAME" < "$migration_file" 2>&1 | grep -i error; then
            print_error "Failed to execute $filename"
            return 1
        fi
        print_success "Completed: $filename"
    done
    
    print_success "All UP migrations completed successfully!"
}

# Function to run DOWN migrations (*.down.sql files)
run_migrations_down() {
    print_status "Running DOWN migrations (rollback)..."
    
    # Get all down migration files sorted by number in reverse
    local migration_files=$(find "$MIGRATIONS_DIR" -name "*.down.sql" | sort -r)
    
    if [ -z "$migration_files" ]; then
        print_warning "No down migration files found"
        return
    fi
    
    for migration_file in $migration_files; do
        local filename=$(basename "$migration_file")
        print_status "Executing: $filename"
        
        # Run the migration using docker exec
        if docker exec -i "$CONTAINER_NAME" mysql -u "$DB_USER" -p"$DB_PASSWORD" "$DB_NAME" < "$migration_file" 2>&1 | grep -i "error\|failed"; then
            print_error "Failed to execute $filename"
            return 1
        fi
        print_success "Completed: $filename"
    done
    
    print_success "All DOWN migrations completed successfully!"
}

# Function to reset database (down then up)
reset_migrations() {
    print_warning "Resetting database... This will drop all tables!"
    read -p "Are you sure? (yes/no): " confirm
    
    if [ "$confirm" != "yes" ]; then
        print_status "Reset cancelled"
        return
    fi
    
    run_migrations_down
    run_migrations_up
    print_success "Database reset completed successfully!"
}

# Function to show migration status
show_status() {
    print_status "Migration Files Status:"
    echo ""
    
    print_status "UP Migrations:"
    find "$MIGRATIONS_DIR" -name "*.up.sql" | sort | while read file; do
        echo "  ✓ $(basename "$file")"
    done
    
    echo ""
    print_status "DOWN Migrations:"
    find "$MIGRATIONS_DIR" -name "*.down.sql" | sort | while read file; do
        echo "  ✓ $(basename "$file")"
    done
    
    echo ""
    print_status "Database Tables:"
    docker exec "$CONTAINER_NAME" mysql -u "$DB_USER" -p"$DB_PASSWORD" "$DB_NAME" -e "SELECT COUNT(*) as 'Table Count' FROM information_schema.tables WHERE table_schema='$DB_NAME';" 2>/dev/null || true
}

# Main logic
main() {
    local action="${1:-up}"
    
    print_status "Starting migration process..."
    check_migrations_dir
    check_docker_container
    
    case "$action" in
        up)
            run_migrations_up
            ;;
        down)
            run_migrations_down
            ;;
        reset)
            reset_migrations
            ;;
        status)
            show_status
            ;;
        *)
            print_error "Unknown action: $action"
            echo ""
            echo "Usage: $0 [command]"
            echo ""
            echo "Commands:"
            echo "  up      - Run all UP migrations (default)"
            echo "  down    - Run all DOWN migrations (rollback)"
            echo "  reset   - Run DOWN then UP migrations (full reset)"
            echo "  status  - Show migration status"
            exit 1
            ;;
    esac
    
    print_success "Migration process completed!"
}

# Run main function
main "$@"
