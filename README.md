<h1 align="center">
    ip2cloud
  <br>
</h1>

<h4 align="center">Check IP addresses against known cloud provider IP address ranges using interval tree</h4>

<p align="center">
  <a href="#installation">🔧 Installation</a>  
  <a href="#usage">⚙️ Usage</a>  
  <a href="#add-new-ip-ranges"> ⚡ Add New IP Ranges</a> 
  <br>
</p>

# Prerequisites

- Python 3.x
- intervaltree library (`pip install intervaltree`)
- tqdm library (`pip install tqdm`)

# Installation


```sh
git clone https://github.com/devanshbatham/ip2cloud
cd ip2cloud
chmod +x setup.sh
./setup.sh
```


Now, you can use `ip2cloud` from the command line to check IP addresses against known cloud provider IP address ranges.


# Usage



```sh
cat input_ips.txt | ip2cloud
```

The script will process the provided IP addresses and display the corresponding cloud provider for each IP address. Optionally, you can use the `-j` or `--json` flag to print the output in JSON format.



### Example

Suppose you have a file named `input_ips.txt` with the following IP addresses:

```
59.82.33.201
63.32.40.140   
63.33.205.240
64.4.8.90 
64.4.8.67
```

Run the script as follows:

```sh
(~)>> cat input_ips.txt | ip2cloud

[aliyun] : 59.82.33.201
[aws] : 63.32.40.140                         
[aws] : 63.33.205.240        
[azure] : 64.4.8.67 
[azure] : 64.4.8.90   
```

- with `--json` or `-j` flag: 


```sh
(~)>> cat input_ips.txt | ip2cloud -j

{
  "aliyun": [
    "59.82.33.201"
  ],
  "aws": [
    "63.32.40.140",
    "63.33.205.240"
  ],
  "azure": [
    "64.4.8.90",
    "64.4.8.67"
  ]
}
```

If no cloud provider is found for an IP address, it won't appear in the JSON or standard output.


# Add New IP Ranges

If you want to add new IP ranges for cloud providers, follow these steps:

1. Create a new text file in the `data` folder (e.g., `somecloud.txt`).
2. Add the IP address ranges in CIDR notation to the new text file. Each range should be on a separate line.
3. Run the setup script again:

```sh
./setup.sh
```

The new IP ranges will be updated in the `cloud_data.json` file, and `ip2cloud`` will use the updated data for IP lookups.


### Note

- The script only supports IPv4 addresses.
- Make sure to keep the `cloud_data.json` file updated with the latest IP address ranges for the cloud providers you want to check.
- `cloud_data.json` should be stored in `/.config/.cloud-ranges` folder after running `setup.sh`


# Supported Cloud Providers

The data folder contains ranges for the mentioned cloud providers, feel free to add any number of providers' ranges. 

- [x] Alibaba Cloud (Aliyun)
- [x] Amazon Web Services (AWS)
- [x] Microsoft Azure
- [x] Cloudflare
- [x] DigitalOcean
- [x] Fastly
- [x] Google Cloud
- [x] IBM Cloud
- [x] Linode
- [x] Oracle Cloud
- [x] Tencent Cloud
- [x] UCloud

