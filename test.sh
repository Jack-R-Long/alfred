#!/bin/bash

# Alfred API Test Runner Script

set -e

echo "🧪 Alfred API Test Suite"
echo "========================"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    local color=$1
    local message=$2
    echo -e "${color}${message}${NC}"
}

# Clean previous test artifacts
print_status $YELLOW "Cleaning previous test artifacts..."
rm -f coverage.out coverage.html
go clean -testcache

# Run tests with coverage
print_status $YELLOW "Running tests with coverage..."
if go test -v -race -coverprofile=coverage.out ./...; then
    print_status $GREEN "✅ All tests passed!"
    
    # Generate coverage report
    if [ -f coverage.out ]; then
        print_status $YELLOW "Generating coverage report..."
        go tool cover -html=coverage.out -o coverage.html
        
        # Show coverage summary
        echo ""
        print_status $YELLOW "Coverage Summary:"
        go tool cover -func=coverage.out | tail -1
        
        print_status $GREEN "📊 Coverage report generated: coverage.html"
    fi
else
    print_status $RED "❌ Tests failed!"
    exit 1
fi

echo ""
print_status $GREEN "🎉 Test suite completed successfully!"