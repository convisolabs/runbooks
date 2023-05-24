import requests
import getopt, sys
from string import Template
import hashlib
import csv
import re

conviso_plat_url = "https://app.convisoappsec.com/findings"

class bcolors:
	HEADER = '\033[95m'
	OKBLUE = '\033[94m'
	OKCYAN = '\033[96m'
	OKGREEN = '\033[92m'
	WARNING = '\033[93m'
	FAIL = '\033[91m'
	ENDC = '\033[0m'
	BOLD = '\033[1m'
	UNDERLINE = '\033[4m'

create_findings = Template("""{
	"flow_project_id": "$conviso_project_id",
	"report": {
		"type": "sca",
		"issues": [ 
			{
				"path": "$path",
				"title": "$title",
				"description": "$description",
				"references": [
					"$references"
				],
				"severity": "$severity",
				"component": "$component",
				"version": "$version",
				"cve": [
					"$cve"
				],
				"plugin": "conviso",
				"solution": "$solution",
				"hash_issue": "$hash"
			}
		]
	}
}""")

help_string = f"""This is a tool to import SCA findings from AWS script into Conviso Platform.

Usage: ./findings_from_aws.py -k <conviso_api_key> -p <project_id> -f <csv_input_file>
	
	Mandatory Options:
	-k,--conviso_api_key =		Conviso Platform API key
	-p,--project_id =		Project key ID (from Project page inside Conviso Platform)
	-f,--file			Name of the .csv File (ex: findings.csv)

	For (this) help:
	-h,--help
"""

def usage():
	print(help_string)
	exit(1)

def generate_sha256_hash(string):
	sha256_hash = hashlib.sha256()
	sha256_hash.update(string.encode('utf-8'))
	return sha256_hash.hexdigest()

def json_safe(input_string):
	safe_string = input_string.replace('\r','').replace('\n','').replace('"','\\"')
	return safe_string

def create_finding(title, json_request, conviso_api_key):
	print(f"[+] Create Finding: {title} - Status: ",end='')
	conviso_headers = {
		"Content-Type": 'application/json',
		"Accept": 'application/json',
		"x-api-key": conviso_api_key
	}
	#print(json_request)
	r = requests.post(conviso_plat_url, data=json_request, headers=conviso_headers)
	if r.status_code == 200:
		print(f"{bcolors.OKCYAN}{r.status_code}{bcolors.ENDC}")
	else:
		print(f"{bcolors.WARNING}{r.status_code}{bcolors.ENDC}")
		print(json_request)

def generate_json_request(vulnerability, issue, conviso_project_id):
	identifier = f"{conviso_project_id};{vulnerability['cve']};{issue['package_name']};{issue['cve_id']};{issue['aws_account']};{issue['repository_name']}"
	hash = generate_sha256_hash(identifier)

	reference = f"https://nvd.nist.gov/vuln/detail/{vulnerability['cve']}".lower()
	title = f"{issue['cve_id']} - {issue['package_name']}"
	path = f"{issue['aws_account']} - {issue['repository_name']}"

	json_request = create_findings.substitute({'conviso_project_id': conviso_project_id, 'title': title, 
	'description': vulnerability['description'], 'path': path, 'severity': issue['severity_level'], 'component': issue['package_name'], 'version': 'N/A',	
	'cve': vulnerability['cve'], 'hash': hash, 'references': reference, 'solution': 'Update package'})

	return json_request

def parse_vulnerabilities(text):
    pattern = r"(CVE-\d{4}-\d+): (.+?)(?=(CVE-\d{4}-\d+):|$)"
    matches = re.findall(pattern, text, re.DOTALL)
    
    vulnerabilities = []
    for match in matches:
        cve_match = match[0].strip()
        description = match[1].strip()
        
        vulnerability = {
            "cve": cve_match,
            "description": description
        }
        vulnerabilities.append(vulnerability)
    
    return vulnerabilities

def read_csv_file(input_file):
	findings = []
	print("[-] Fetching findings from file.")
	with open(input_file, 'r') as file:
		csv_reader = csv.DictReader(file, delimiter=',')
		for row in csv_reader:
			aws_account = row['aws_account']
			repository_name = row['repository_name']
			cve_id = row['cve_id']
			severity_level = row['severity_level']
			package_name = row['package_name']
			description = row['description']
			created_at = row['createdAt']
			timing = row['timing']
			
			record = {
				'aws_account': aws_account,
				'repository_name': repository_name,
				'cve_id': cve_id,
				'severity_level': severity_level,
				'package_name': package_name,
				'description': description,
				'created_at': created_at,
				'timing': timing
			}
			
			findings.append(record)
	return findings

def main():
	argument_list = sys.argv[1:]
	options = "k:p:f:h"
	long_options = ["conviso_api_key=", "project_id=", "input_file=", "help"]

	conviso_api_key, conviso_project_id, input_file = None, None, None

	try:
		arguments, values = getopt.getopt(argument_list, options, long_options)

		for arg, value in arguments:
			if arg in ("-k", "--conviso_api_key"):
				conviso_api_key = value
			elif arg in ("-p", "--project_id"):
				conviso_project_id = value
			elif arg in ("-f", "--input_file"):
				input_file = value
			elif arg in ("-h", "--help"):
				usage()

	except getopt.error:
		usage()

	if conviso_api_key == None or conviso_project_id == None or input_file == None:
		usage()

	findings = read_csv_file(input_file)
	#print(findings)

	for issue in findings:
		#print(issue)
		vulnerabilities = parse_vulnerabilities(issue['description'])
		for vulnerability in vulnerabilities:
			#print(vulnerability)
			title = vulnerability['cve']
			#print(title)
			json_request = generate_json_request(vulnerability, issue, conviso_project_id)
			create_finding(title, json_request, conviso_api_key)

if __name__ == "__main__":
	main()
