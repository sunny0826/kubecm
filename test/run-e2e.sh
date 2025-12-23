#!/usr/bin/env bash

# E2E Test Runner for kubecm
# This script builds kubecm and runs end-to-end tests

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"
BIN_DIR="${PROJECT_ROOT}/bin"
KUBECM_BIN="${BIN_DIR}/kubecm"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}================================================${NC}"
echo -e "${GREEN}  KubeCM E2E Test Runner${NC}"
echo -e "${GREEN}================================================${NC}"
echo ""

# Parse arguments
VERBOSE=false
BUILD=true

while [[ $# -gt 0 ]]; do
  case $1 in
    -v|--verbose)
      VERBOSE=true
      shift
      ;;
    --no-build)
      BUILD=false
      shift
      ;;
    -h|--help)
      echo "Usage: $0 [OPTIONS]"
      echo ""
      echo "Options:"
      echo "  -v, --verbose    Enable verbose output"
      echo "  --no-build       Skip building kubecm"
      echo "  -h, --help       Show this help message"
      exit 0
      ;;
    *)
      echo -e "${RED}Unknown option: $1${NC}"
      exit 1
      ;;
  esac
done

# Build kubecm
if [ "$BUILD" = true ]; then
  echo -e "${YELLOW}Building kubecm...${NC}"
  mkdir -p "${BIN_DIR}"
  
  cd "${PROJECT_ROOT}"
  go build -o "${KUBECM_BIN}" .
  
  if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ Build successful${NC}"
    echo ""
  else
    echo -e "${RED}✗ Build failed${NC}"
    exit 1
  fi
else
  # Check if binary exists
  if [ ! -f "${KUBECM_BIN}" ]; then
    echo -e "${RED}✗ kubecm binary not found at ${KUBECM_BIN}${NC}"
    echo -e "${YELLOW}  Please build kubecm first or run without --no-build${NC}"
    exit 1
  fi
  echo -e "${YELLOW}Skipping build (using existing binary)${NC}"
  echo ""
fi

# Display kubecm version
echo -e "${YELLOW}kubecm version:${NC}"
"${KUBECM_BIN}" version
echo ""

# Run e2e tests
echo -e "${YELLOW}Running e2e tests...${NC}"
cd "${PROJECT_ROOT}/test/e2e"

export KUBECM_BIN="${KUBECM_BIN}"

if [ "$VERBOSE" = true ]; then
  go test -v -timeout 10m ./...
else
  go test -timeout 10m ./...
fi

TEST_RESULT=$?

echo ""
if [ $TEST_RESULT -eq 0 ]; then
  echo -e "${GREEN}================================================${NC}"
  echo -e "${GREEN}  ✓ All E2E tests passed!${NC}"
  echo -e "${GREEN}================================================${NC}"
  exit 0
else
  echo -e "${RED}================================================${NC}"
  echo -e "${RED}  ✗ E2E tests failed${NC}"
  echo -e "${RED}================================================${NC}"
  exit 1
fi
