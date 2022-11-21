import sys
import pandas as pd
import requests
import os

from pandas import ExcelWriter #pip install openpyxl

CONVISO_PLATFORM_URI =  os.environ.get('CONVISO_PLATFORM_URI')
CONVISO_PLATFORM_TOKEN = os.environ.get('CONVISO_PLATFORM_TOKEN')
CONVISO_PLATFORM_STATUS_CODE = 200
CONVISO_PLATFORM_STATUS_CODE_ERROR = 500
CONVISO_PLATFORM_HEADERS = {'x-api-key':CONVISO_PLATFORM_TOKEN}

CONVISO_PLATFORM_REQUIREMENTS_QUERY = '''
query Projects($id:ID!){
    project(id: $id) {
      activities {
        id
        title
        reference
        status
        description
        justify
      }
      playbooks {
        checklistTypeId
        companyId
      }
    }
}
'''
CONVISO_PLATFORM_VULNERABILITIES_QUERY = '''
query Assets($id:ID!){
	assets(id: $id, page:1, limit:1000){
    metadata{
      currentPage
      limitValue
      totalCount
      totalPages
    }
    collection{
      id
      name
      projects{
        pid
        ...VulnByAsset
      }
    }
  }
}
fragment VulnByAsset on Project{
  vulnerabilities{
    title
    vulnerabilityTemplate{
      categoryList
      reference
    }
  }
}
'''
# Online Request to Conviso Platform
#
def runQuery(uri, query, statusCode, headers, variables):
    request = requests.post(uri, json={'query':query, 'variables':variables}, headers=headers)
    if request.status_code == statusCode:
        return request.json()
    else:
        raise Exception(f'Unexpected status code returned: {request.status_code}')

# Get Requirements from Conviso Platform by Project ID
#
def get_requeriments(project):
    var_requi = {'id':project}
    resultado = runQuery(CONVISO_PLATFORM_URI, CONVISO_PLATFORM_REQUIREMENTS_QUERY, CONVISO_PLATFORM_STATUS_CODE, CONVISO_PLATFORM_HEADERS, var_requi)
    requirements = []
    # Get ASVS code
    for item in resultado['data']['project']['activities']:
        title = item['title'].split('-')
        requirements.append(title[1].strip())
    return requirements

# Get CWE from Vulnerabilities templates
#
def get_vulnerabilities(company):
    var_vuln = {'id':company}
    response = runQuery(CONVISO_PLATFORM_URI, CONVISO_PLATFORM_VULNERABILITIES_QUERY, CONVISO_PLATFORM_STATUS_CODE, CONVISO_PLATFORM_HEADERS, var_vuln)
    # Get CWE code
    vulnerabilities = []
    for asset in response['data']['assets']['collection']:
        projects = asset['projects']
        if projects is not None and len(projects) > 0:
            for project in projects:
                vulns = project['vulnerabilities']
                if vulns is not None:
                    for vul in vulns:
                        if vul is not None:
                            template = vul['vulnerabilityTemplate']
                            if template is not None:
                                category = template['categoryList']
                                if category is not None and category != 'N/A':
                                    vulnerabilities.append(category[4:category.find(" ")])
    return vulnerabilities                                

# Get ASVS from file - only cols from asvs code and cwe
#
def get_owasp_asvs():
    return pd.read_csv('OWASP ASVS 4.0.2-en.csv', usecols=['req_id','cwe','level1'])

# write excel file
#
def writeFile(data, header):    
    dataframe = pd.DataFrame(data, columns=header)
    writer = ExcelWriter('Result.xlsx')
    dataframe.to_excel(writer,'Sheet1',index=False)
    writer.save()

def main():
    try:
        # Loading the security "requeriments"
        requeriments = get_requeriments(input('Enter the Project ID:'))

        # Loading "Vulnerabilities" from cia's assets
        vulnerabilities = get_vulnerabilities(input('Enter the Company ID:'))

        # Loading "ASVS 4.0.2" (currenty version in Conviso Platform)
        asvs = get_owasp_asvs()
                
        data = []
        for column, req in enumerate(requeriments):

            for field in asvs.values:
                # ASVS -> CWE
                if req == field[0]:
                    cwe = str(int(field[2]))
                    # CWE -> is in vulnerability list?
                    if cwe in vulnerabilities:
                        # Case vulnerability found => requeriment is not according
                        result = 'Not according'
                    elif str(field[1]) == '✓':
                        # Case vulnerability NOT found and level 1 => requeriment okay
                        result = 'Done'
                    else:
                        # Case vulnerability NOT found, but level 2 or 3 => manual analysis necessary
                        result = 'non-automated'

                    data.append([req, cwe, result])
                    print('Requeriment:', req, '- CWE:', cwe, ' -> ', result)

        # generate spreadsheet with result
        header = ['ASVS', 'CWE', 'Result']
        writeFile(data, header)

    except ValueError as ve:
        return str(ve)

if __name__ == '__main__':
  if sys.version_info.major == 3: # mandatory python 3
    sys.exit(main())
