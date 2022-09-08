import sys
import pandas as pd
import requests

from pandas import ExcelWriter #pip install openpyxl

#CONVISO_PLATFORM_URI =  'https://app.convisoappsec.com/graphql'
CONVISO_PLATFORM_URI =  'https://homologa.convisoappsec.com/graphql'
CONVISO_PLATFORM_TOKEN = input('Please, API Key>')
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
#query Projects($id:ID!, $page:Int, $limit:Int){
#assets(id: $id, page: $page, limit: $limit){
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
# Requisição online (Consulta na Conviso Platform)
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
    #print(resultado)
    #print("---------------------------------------------------")
    requirements = []
    # Etapa obter código ASVS
    for item in resultado['data']['project']['activities']:
        title = item['title'].split('-')
        #new_row = [item['id'], item['title'], title[1]]
        requirements.append(title[1].strip())
        #print(f'{new_row}')
    return requirements

# Get CWE from Vulnerabilities templates
#
def get_vulnerabilities(company):
    var_vuln = {'id':company}
    response = runQuery(CONVISO_PLATFORM_URI, CONVISO_PLATFORM_VULNERABILITIES_QUERY, CONVISO_PLATFORM_STATUS_CODE, CONVISO_PLATFORM_HEADERS, var_vuln)
    #print(vulnerabilities)
    #print("---------------------------------------------------")
    # Etapa obter código CWE
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
                                    #new_row = [asset['id'], asset['name'], category[4:category.find(" ")]]
                                    vulnerabilities.append(category[4:category.find(" ")])
                                    #print(f'{new_row}')
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
        # Carregando os "Requisitos" de segurança da informação
        requeriments = get_requeriments(input('Enter the Project ID:'))

        # Carregando as "Vulnerabilidades" dos ativos da companhia
        vulnerabilities = get_vulnerabilities(input('Enter the Company ID:'))

        # Carregar "ASVS 4.0.2" (versão compatível com Conviso Platform)
        asvs = get_owasp_asvs()
                
        data = []
        # Etapa de Análise
        for column, req in enumerate(requeriments):

            for field in asvs.values:
                # ASVS -> CWE
                if req == field[0]:
                    cwe = str(int(field[2]))
                    # CWE -> está na lista de vulnerabilidades?
                    if cwe in vulnerabilities:
                        # Se temos vulnerabilidade encontrada => requisito não conforme
                        result = 'Not according'
                    elif str(field[1]) == '✓':
                        # Se não temos vulnerabilidade encontrada e é do nível 1 => requisito em conformidade
                        result = 'Done'
                    else:
                        # Se não temos vulnerabilidade encontrada, mas não é do nível 1 => requer análise manual
                        result = 'non-automated'

                    data.append([req, cwe, result])
                    print('Requito:', req, '- CWE:', cwe, ' -> ', result)

        # Etapa de gerar planilha com resultado
        header = ['ASVS', 'CWE', 'Result']
        writeFile(data, header)

    except ValueError as ve:
        return str(ve)

if __name__ == '__main__':
    sys.exit(main())