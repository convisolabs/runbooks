import requests
import json
import getopt, sys
from string import Template

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

statuses = {
	'undefined': 12,
	'identified': 13,
	'in_progress': 14,
	'fix_accepted': 15,
	'fix_refused': 23,
	'waiting_validation': 22
}

help_string = f"""This is a tool to export vulnerabilities from Conviso Plaform into a Kanbanize Board.

Usage: ./kanbanize_vuln_sync.py -k <conviso_api_key> -p <company_id> -n <kanbanize_api_key> -b <board_id> -l <lane_id> -u <url> [-f] [-h]

	Mandatory Options:
	-k,--conviso_api_key=		Conviso Platform API key
	-p,--company_id=			Company Id (from Conviso Platform)
	-n,--kanbanize_api_key=		Kanbanize API key
	-b,--board_id=			Kanbanize Board Id
	-l,--lane_id= 			Kanbanize Lane Id
	-u,--url=				Kanbanize URL (eg: https://conviso.kanbanize.com)
	-f,--update				Flag for updating the columns in Kanbanize

	For (this) help:
	-h,--help
"""

def usage():
	print(help_string)
	exit(1)

env = 'app'

url = f'https://{env}.convisoappsec.com/graphql'

t_query = Template("""{
vulnerabilities(
	page: $page_num_template
	limit: 1000
	params: {
	projectScopeIdEq: $company_id_template
	}
	order: created_at
	orderType: DESC
) {
	collection {
	title
	id
	projectId
	companyId
	vulnerabilityStatus
	project {
		label
	}
	}
	metadata {
	currentPage
	totalPages
	}
}
}""")

def get_cards(kanbanize_api_key, board_id, kanbanize_url):
	my_url = f"{kanbanize_url}/api/v2/cards/"

	payload = {
		"boardid": board_id,
	}

	headers = {
		"apikey": kanbanize_api_key,
		"Content-Type": "application/json"
	}

	response = requests.get(my_url, data=json.dumps(payload), headers=headers)

	if response.status_code == 200:
		print("Cards retrieved successfully!")
		#print(response.content)
	else:
		print(f"{bcolors.WARNING}Failed to retrieve cards.{bcolors.ENDC}")
		print(response.content)
		exit(1)
		
	json_return = response.json()
	return json_return

def update_card(kanbanize_api_key, card_id, column_id, kanbanize_url):
	my_url = f"{kanbanize_url}/api/v2/cards/{card_id}"

	payload = {
		"column_id": column_id,
	}

	headers = {
		"apikey": kanbanize_api_key,
		"Content-Type": "application/json"
	}

	response = requests.patch(my_url, data=json.dumps(payload), headers=headers)

	if response.status_code == 200:
		print(f"- {bcolors.OKGREEN}Card updated successfully!{bcolors.ENDC}")
	else:
		print(f"- {bcolors.WARNING}Failed to update card.{bcolors.ENDC}")
		print(response.content)


def create_card(kanbanize_api_key, board_id, column_id, lane_id, title, description, kanbanize_url):
	my_url = f"{kanbanize_url}/api/v2/cards/"

	payload = {
		"boardid": board_id,
		"column_id": column_id,
		"lane_id": lane_id,
		"title": title,
		"description": description
	}

	headers = {
		"apikey": kanbanize_api_key,
		"Content-Type": "application/json"
	}

	response = requests.post(my_url, data=json.dumps(payload), headers=headers)

	if response.status_code == 200:
		print(f"- {bcolors.OKGREEN}Card created successfully!{bcolors.ENDC}")
	else:
		print(f"- {bcolors.WARNING}Failed to create card.{bcolors.ENDC}")
		print(response.content)
		
def check_value_in_json(json_data, search_value):
	for value in json_data['data'].values():
		if isinstance(value, list):
			for v_key in value:
				if isinstance(v_key, dict) and v_key.get('title') == search_value:
					return True
	return False

def search_id_column(cards, title):
	for value in cards['data'].values():
		if isinstance(value, list):
			for v_key in value:
				if title in v_key.values():
					return v_key['column_id'], v_key['card_id']
		
def fetch_vulns(company_id, page_num, url, conviso_api_key, page_max):
	query = t_query.substitute({'company_id_template': company_id, 'page_num_template': page_num})
	r = requests.post(url, json={'query': query}, headers={"x-api-key": conviso_api_key})

	if (r.status_code != 200):
		print(f'[ERROR] PROCESSING: PAGE {page_num} FROM {page_max} TOTAL PAGES - STATUS: {r.status_code} - SKIPPING')
		exit(1)


	json_data = json.loads(r.text)
	if json_data['data']['vulnerabilities'] is None:
		print(f'[ERROR] PROCESSING: PAGE {page_num} FROM {page_max} TOTAL PAGES - VULNERABILITY == None - SKIPPING')
		exit(1)

	page_max = json_data['data']['vulnerabilities']['metadata']['totalPages']
	vulns = json_data['data']['vulnerabilities']['collection']

	return vulns, page_max

def main():
	argument_list = sys.argv[1:]
	options = "k:p:n:b:l:u:fh"
	long_options = ["conviso_api_key=", "company_id=", "kanbanize_api_key=", "board_id=", "lane_id=", "url=", "update", "help"]

	conviso_api_key, company_id, kanbanize_api_key, board_id, lane_id, kanbanize_url = (None,) * 6
	update_flag = False

	try:
		arguments, values = getopt.getopt(argument_list, options, long_options)

		for arg, value in arguments:
			if arg in ("-k", "--conviso_api_key"):
				conviso_api_key = value
			elif arg in ("-p", "--company_id"):
				company_id = value
			elif arg in ("-n", "--kanbanize_api_key"):
				kanbanize_api_key = value
			elif arg in ("-b", "--board_id"):
				board_id = value
			elif arg in ("-l", "--lane_id"):
				lane_id = value
			elif arg in ("-u", "--url"):
				kanbanize_url = value
			elif arg in ("-f", "--update"):
				update_flag = True
			elif arg in ("-h", "--help"):
				usage()

	except getopt.error:
		usage()

	if any(var is None for var in [conviso_api_key, company_id, kanbanize_api_key, board_id, lane_id, kanbanize_url]):
		usage()

	page_num = 1
	page_max = 9999

	cards = get_cards(kanbanize_api_key, board_id, kanbanize_url)

	while (page_max >= page_num):

		vulns, page_max = fetch_vulns(company_id, page_num, url, conviso_api_key, page_max)

		for vuln in vulns:
			title = 'undefined' if vuln['title'] == None else vuln['title']
			companyId = company_id if vuln['companyId'] == None else vuln['companyId']
			projectId = 0 if vuln['projectId'] == None else vuln['projectId']
			projectName = 'undefined' if vuln['project']['label'] == None else vuln['project']['label']
			id = 'undefined' if vuln['id'] == None else vuln['id']
			status = 'undefined' if vuln['vulnerabilityStatus'] == None else vuln['vulnerabilityStatus']
			description = f"{title} (https://{env}.convisoappsec.com/scopes/{companyId}/projects/{projectId}/occurrences/{id}?locale=en)"
			
			print(f"[+] {projectName} - {title} ",end='')

			if not check_value_in_json(cards, title):
				create_card(kanbanize_api_key, board_id, statuses[status], lane_id, title, description, kanbanize_url)
			else:
				print(f"- {bcolors.WARNING}Card already exists{bcolors.ENDC} ", end='')
				column, card_id = search_id_column(cards, title)
				if column != statuses[status]:
					print(f"-  {bcolors.WARNING}Columns differ{bcolors.ENDC} ", end='')
					if (update_flag):
						update_card(kanbanize_api_key, card_id, statuses[status], kanbanize_url)
					else:
						print('')
				else:
					print(f"- {bcolors.OKGREEN}Columns OK{bcolors.ENDC}")
				
		page_num += 1

if __name__ == "__main__":
	main()
