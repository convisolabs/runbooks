import requests
import json
import getopt, sys
from string import Template
import hashlib

conviso_plat_url = "https://homologa.convisoappsec.com/findings"

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
					$references
				],
				"severity": "$severity",
				"component": "$component",
				"version": "$version",
				"cve": [
					"$cve"
				],
				"plugin": "dependabot",
				"solution": "$solution",
				"hash_issue": "$hash"
			}
		]
	}
}""")

help_string = f"""This is a tool to import Dependabot findings into Conviso Platform.

Usage: ./insert_dependabot_findings.py -k <conviso_api_key> -p <project_id> [-g <github_api_key> -o <github_owner> -r <github_repo>]/[-f <json_input_file>]
	
	Mandatory Options:
	-k,--conviso_api_key =		Conviso Platform API key
	-p,--project_id =		Project key ID (from Project page inside Conviso Platform)
	
	For fetching alerts directly from dependabot:
	-g,--github_api_key =		GitHub Authentication key
	-o,--github_owner =		GitHub repository owner (ex: https://github.com/{bcolors.WARNING}Company{bcolors.ENDC}/Repository)
	-r,--github_repo =		GitHub repository name (ex: https://github.com/Company/{bcolors.WARNING}Repository{bcolors.ENDC})

	For fetching the alerts from a .json file:
	-f,--file			Name of the .json File (ex: alerts.json)

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
	r = requests.post(conviso_plat_url, data=json_request, headers=conviso_headers)
	if r.status_code == 200:
		print(f"{bcolors.OKCYAN}{r.status_code}{bcolors.ENDC}")
	else:
		print(f"{bcolors.WARNING}{r.status_code}{bcolors.ENDC}")
		print(json_request)

def generate_json_request(issue, conviso_project_id):
	title = json_safe(issue['security_advisory']['summary'])
	path = json_safe(issue['dependency']['manifest_path'])
	description = json_safe(issue['security_advisory']['description'])
	
	references = ''
	for index, reference in enumerate(issue['security_advisory']['references']):
		references += '"'
		references += json_safe(reference['url'])
		if index != len(issue['security_advisory']['references']) - 1:
			references += '",'
		else:
			references += '"'

	severity = json_safe(issue['security_advisory']['severity'])
	component = json_safe(issue['security_vulnerability']['package']['name'])
	version = json_safe(issue['security_vulnerability']['vulnerable_version_range'])
	
	cve = 'null'
	for reference in issue['security_advisory']['identifiers']:
		if reference['type'] == 'CVE':
			cve = json_safe(reference['value'])

	solution = 'Update to (at least) version '
	solution += json_safe(issue['security_vulnerability']['first_patched_version']['identifier'])

	vuln_id = json_safe(issue['security_advisory']['ghsa_id'])
	hash = generate_sha256_hash(vuln_id)

	json_request = create_findings.substitute({'conviso_project_id': conviso_project_id, 'title': title, 
	'description': description, 'path': path, 'severity': severity, 'component': component, 'version': version,	
	'cve': cve, 'hash': hash, 'references': references, 'solution': solution})

	return json_request

def fetch_dependabot_alerts(github_api_key, github_owner, github_repo):
	github_headers = {
		'Authorization': f'Bearer {github_api_key}',
		'Accept': 'application/vnd.github+json'
	}
	github_url = f"https://api.github.com/repos/{github_owner}/{github_repo}/dependabot/alerts"
	print("[-] Fetching alerts from Dependabot. Status: ",end='')
	response = requests.get(github_url, headers=github_headers)

	if response.status_code == 200:
		alerts = response.json()
		print(f"{bcolors.OKCYAN}{response.status_code}{bcolors.ENDC}")
	else:
		print(f"{bcolors.WARNING}{response.status_code}{bcolors.ENDC}")
		print(f"{bcolors.FAIL}[X] Error: Received status code while getting dependabot alerts.{bcolors.ENDC}")
		exit(1)
	
	return alerts

def read_json_file(input_file):
	print("[-] Fetching alerts from file.")
	with open(input_file) as file:
		alerts = json.load(file)
	
	return alerts

def main():
	argument_list = sys.argv[1:]
	options = "k:p:g:o:r:f:h"
	long_options = ["conviso_api_key=", "project_id=", "github_api_key=", "github_owner=", "github_repo=", "input_file=", "help"]

	conviso_api_key,github_api_key,conviso_project_id,github_owner,github_repo,input_file = None,None,None,None,None,None

	try:
		arguments, values = getopt.getopt(argument_list, options, long_options)

		for arg, value in arguments:
			if arg in ("-k", "--conviso_api_key"):
				conviso_api_key = value
			elif arg in ("-p", "--project_id"):
				conviso_project_id = value
			elif arg in ("-g", "--github_api_key"):
				github_api_key = value
			elif arg in ("-o", "--github_owner"):
				github_owner = value
			elif arg in ("-r", "--github_repo"):
				github_repo = value
			elif arg in ("-f", "--input_file"):
				input_file = value
			elif arg in ("-h", "--help"):
				usage()

	except getopt.error:
		usage()

	if conviso_api_key == None or conviso_project_id == None or (input_file == None and (github_api_key == None or 
		github_owner == None or github_repo == None)):
		usage()

	if input_file == None:
		alerts = fetch_dependabot_alerts(github_api_key, github_owner, github_repo)
	else:
		alerts = read_json_file(input_file)

	for issue in alerts:
		title = json_safe(issue['security_advisory']['summary'])
		json_request = generate_json_request(issue, conviso_project_id)
		create_finding(title, json_request, conviso_api_key)

if __name__ == "__main__":
    main()
