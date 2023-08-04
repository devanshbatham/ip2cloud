#!/usr/bin/env python3

import json
import os

def read_ip_ranges_from_file(file_path):
    """
    Reads IP ranges from a text file and returns a list of IP ranges.
    
    Parameters:
        file_path (str): The path to the text file containing IP ranges.
        
    Returns:
        list: A list of IP ranges (strings).
    """
    with open(file_path, "r") as file:
        ip_ranges = [line.strip() for line in file if line.strip()]
    return ip_ranges

def create_cloud_data_json(file_paths, output_file):
    """
    Creates a JSON file containing cloud data with IP ranges from input files.
    
    Parameters:
        file_paths (list): A list of file paths (strings) where each file contains IP ranges for a cloud service.
        output_file (str): The path of the output JSON file to be created.
    """
    cloud_data = {}
    for file_path in file_paths:
        cloud_name = os.path.splitext(os.path.basename(file_path))[0]
        if os.path.exists(file_path):
            ip_ranges = read_ip_ranges_from_file(file_path)
            cloud_data[cloud_name] = ip_ranges
        else:
            print(f"File not found: {file_path}")

    # Get the directory path from the output_file
    output_dir = os.path.dirname(output_file)

    # Create the directory if it doesn't exist
    os.makedirs(output_dir, exist_ok=True)

    with open(output_file, "w") as output:
        json.dump(cloud_data, output, indent=4)

    print(f"cloud_data.json saved to: {output_file}")

if __name__ == "__main__":
    # Assuming the data folder contains text files with IP ranges for each cloud service.
    data_folder = "data"
    file_paths = [os.path.join(data_folder, file_name) for file_name in os.listdir(data_folder) if file_name.endswith(".txt")]

    # Define the output directory and file
    output_directory = os.path.expanduser("~/.config/.cloud-ranges")
    os.makedirs(output_directory, exist_ok=True)
    output_file = os.path.join(output_directory, "cloud_data.json")

    create_cloud_data_json(file_paths, output_file)
