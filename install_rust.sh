#!/bin/bash

# A script to automate the installation of Rust via rustup.
# It checks for dependencies and runs the installer non-interactively.

# Exit immediately if a command exits with a non-zero status.
set -e

# --- Helper Functions ---

# Function to check if a command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# --- Main Script Logic ---

echo "Starting Rust installation script..."

# 1. Check for dependencies
echo "Checking for required dependencies (curl, gcc)..."
if ! command_exists curl; then
    echo "Error: 'curl' is not installed. Please install it to continue." >&2
    exit 1
fi

if ! command_exists gcc; then
    echo "Error: A C compiler like 'gcc' is not found." >&2
    echo "Please install a C compiler. For example:" >&2
    echo "  - On Debian/Ubuntu: sudo apt update && sudo apt install build-essential" >&2
    echo "  - On Fedora/CentOS: sudo dnf groupinstall \"Development Tools\"" >&2
    echo "  - On Arch Linux:    sudo pacman -S base-devel" >&2
    exit 1
fi
echo "Dependencies found."

# 2. Download and run the rustup installer non-interactively
echo "Downloading and running the rustup installer..."
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y

echo ""
echo "Rust was installed successfully!"
echo "To get started, you need to configure your current shell. Run:"
echo "  source \"\$HOME/.cargo/env\""
echo "Or, simply open a new terminal session."