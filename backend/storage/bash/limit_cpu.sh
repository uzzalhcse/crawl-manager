#!/bin/bash

# Step 1: Check if a binary name is provided as an argument
if [ "$#" -lt 1 ]; then
    echo "Usage: $0 <binary_name> [cpu_percentage]"
    exit 1
fi

BINARY_NAME="$1"

# Step 2: Check if a CPU usage percentage is provided; default to 80% if not
CPU_PERCENTAGE=${2:-80}

# Validate the CPU percentage
if ! echo "$CPU_PERCENTAGE" | grep -qE '^[0-9]+$' || [ "$CPU_PERCENTAGE" -lt 0 ] || [ "$CPU_PERCENTAGE" -gt 100 ]; then
    echo "[ERROR] Invalid CPU percentage: $CPU_PERCENTAGE. Must be an integer between 0 and 100."
    exit 1
fi

# Step 3: Get total number of CPU cores
TOTAL_CPUS=$(nproc)
echo "[INFO] Total CPU cores detected: $TOTAL_CPUS"

# Step 4: Calculate the target CPUs based on the provided percentage
TARGET_CPUS=$(echo "scale=0; $TOTAL_CPUS * $CPU_PERCENTAGE / 100" | bc)
echo "[INFO] Target CPUs for $CPU_PERCENTAGE% usage: $TARGET_CPUS"

# Step 5: Calculate the quota for the specified CPU usage per core
# Assume a 100ms period (100000us)
QUOTA=$(echo "$TARGET_CPUS * 100000" | bc)
echo "[INFO] CPU quota calculated: $QUOTA microseconds for 100 milliseconds period"

# Step 6: Create the cgroup if it doesn't exist
CGROUP_PATH="/sys/fs/cgroup/mygroup"
if [ ! -d "$CGROUP_PATH" ]; then
    echo "[INFO] Creating cgroup directory: $CGROUP_PATH"
    sudo mkdir "$CGROUP_PATH"
else
    echo "[INFO] Cgroup directory already exists: $CGROUP_PATH"
fi

# Step 7: Set the CPU limit for the cgroup (specified CPU usage)
echo "$QUOTA 100000" | sudo tee "$CGROUP_PATH/cpu.max" > /dev/null
echo "[INFO] CPU limit set to: $QUOTA microseconds per 100 milliseconds"

# Step 8: Run the Go binary or any other program under the cgroup
echo "$$" | sudo tee "$CGROUP_PATH/cgroup.procs" > /dev/null
echo "[INFO] Running binary: $BINARY_NAME under cgroup"
./"$BINARY_NAME"

# Step 9: Log completion
if [ $? -eq 0 ]; then
    echo "[INFO] Successfully executed $BINARY_NAME"
else
    echo "[ERROR] Failed to execute $BINARY_NAME"
fi
