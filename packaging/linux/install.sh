#!/bin/bash
set -e

INSTALL_DIR="/usr/local/bin"

# Copy binary
if [ -f "ope" ]; then
    sudo cp ope "$INSTALL_DIR/ope"
    sudo chmod +x "$INSTALL_DIR/ope"
    echo "Installed ope to $INSTALL_DIR/ope"
else
    echo "Error: ope binary not found in current directory"
    exit 1
fi

# Register URL scheme
ope install

echo "Done! Try: xdg-open 'ope:///tmp'"
