#!/usr/bin/env python3

import json
import ipaddress
import sys
from tqdm import tqdm
from intervaltree import Interval, IntervalTree
import argparse
import os

def build_interval_tree(cloud_data):
    """
    Build an interval tree to efficiently store IP address ranges for each cloud provider.

    Args:
        cloud_data (dict): A dictionary where keys are cloud provider names, and values are lists of IP address ranges.

    Returns:
        IntervalTree: An interval tree containing the IP address ranges for each cloud provider.

    Complexity:
        - Time Complexity: O(N * M), where N is the number of cloud providers, and M is the average number of IP address ranges per provider.
        - Space Complexity: O(K), where K is the total number of IP address ranges across all cloud providers.
    """
    interval_tree = IntervalTree()
    for cloud, ip_ranges in cloud_data.items():
        for ip_range in ip_ranges:
            try:
                # Convert the IP range string into an IPv4Network object
                ip_range_obj = ipaddress.IPv4Network(ip_range)

                # Check if the IP range represents a single IP address (/32 subnet mask)
                if ip_range_obj.prefixlen == 32:
                    start_ip = end_ip = ip_range_obj.network_address
                else:
                    # Sort IP addresses to make sure start is less than or equal to end
                    start_ip, end_ip = sorted((ip_range_obj.network_address, ip_range_obj.broadcast_address))

                # Add the IP address range interval to the interval tree with the associated cloud provider name
                interval_tree.add(Interval(start_ip._ip, end_ip._ip, cloud))
            except ValueError as e:
                # If there is an error in parsing the IP range, skip and continue
                pass
    return interval_tree

def find_cloud(ip, interval_tree):
    """
    Find the cloud provider associated with a given IP address using the interval tree.

    Args:
        ip (str): The IP address to check.
        interval_tree (IntervalTree): The interval tree containing IP address ranges for each cloud provider.

    Returns:
        str or None: The cloud provider name associated with the IP address, or None if no match is found.

    Complexity:
        - Time Complexity: O(log K), where K is the total number of IP address ranges in the interval tree.
        - Space Complexity: O(1).
    """
    ip_obj = ipaddress.IPv4Address(ip)
    intervals = interval_tree[ip_obj._ip]
    for interval in intervals:
        # If the interval has associated data (cloud provider name), return the name
        if interval.data:
            return interval.data
    return None

if __name__ == "__main__":
    # Set up argparse to handle command-line arguments
    parser = argparse.ArgumentParser(description="Check IP addresses against cloud providers.")
    parser.add_argument("-j", "--json", action="store_true", help="Print output in JSON format.")
    args = parser.parse_args()

    # Assuming your JSON data is stored in a file called "cloud_data.json" in ~/.config/.cloud-ranges/
    config_dir = os.path.expanduser("~/.config/.cloud-ranges")
    json_file = os.path.join(config_dir, "cloud_data.json")

    # Read the JSON data from the specified file
    with open(json_file, "r") as f:
        cloud_data = json.load(f)

    # Build the interval tree using the provided cloud_data
    interval_tree = build_interval_tree(cloud_data)

    # Read IPs from stdin (one IP per line)
    ips_to_check = [line.strip() for line in sys.stdin]

    # If the JSON output option is provided, use a dictionary to collect the results
    if args.json:
        result_dict = {}
        for ip_to_check in ips_to_check:
            cloud = find_cloud(ip_to_check, interval_tree)
            if cloud:
                result_dict.setdefault(cloud, []).append(ip_to_check)

        # Print the result in JSON format
        print(json.dumps(result_dict, indent=4))
    else:
        # Use tqdm to display a progress bar while checking each IP address
        with tqdm(total=len(ips_to_check), desc="Checking IPs", unit="IP", leave=False) as pbar:
            for ip_to_check in ips_to_check:
                cloud = find_cloud(ip_to_check, interval_tree)
                if cloud:
                    tqdm.write(f"[{cloud}] : {ip_to_check}")
                pbar.update(1)
