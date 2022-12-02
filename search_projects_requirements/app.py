import requests
import os
import sys

CONVISO_PLATFORM_URI =  'https://app.convisoappsec.com/graphql'
CONVISO_PLATFORM_TOKEN = os.environ.get('CONVISO_PLATFORM_TOKEN','token')
CONVISO_PLATFORM_STATUS_CODE = 200
CONVISO_PLATFORM_HEADERS = {'x-api-key':CONVISO_PLATFORM_TOKEN}


QUERYCOMPANIES = """
 {
  companies(page: 1, limit: 1000, params: {}, order: label, orderType: DESC) {
    collection {
      id
      label
    }
    metadata {
      currentPage
      limitValue
      totalCount
      totalPages
    }
  }
}
"""

QUERYPROJECTS = """
	query Projects($scopeIdEq:Int, $page:Int)
	{
	  projects(page: $page, limit: 1000, params: {scopeIdEq:$scopeIdEq}, sortBy: "string", descending: true) {
	    collection {
        id
	      companyId,
	      label
	      playbooks{
          id,
          label
	      }
	    }
	    metadata {
	      currentPage
	      total
	      totalAnalysis
	      totalCount
	      totalDone
	      totalEstimate
	      totalFixing
	      totalPages
	      totalPaused
	      totalPlanned
	      totalVulnsCount
	      totalVulnsCritical
	      totalVulnsHigh
	      totalVulnsLow
	      totalVulnsMedium
	      totalVulnsNotification
	    }
	  }
	}

"""



def runQuery(uri, query, statusCode, headers, variables):
  request = requests.post(uri, json={'query':query, 'variables':variables}, headers=headers)
  if request.status_code == statusCode:
    return request.json()
  else:
    raise Exception(f'Unexpected status code returned: {request.status_code}')
    
def runSearchCompanyProjectRequirement():
  projs = []
  fileName = input('Nome do Arquivo:')
  companies = runQuery(CONVISO_PLATFORM_URI, QUERYCOMPANIES, CONVISO_PLATFORM_STATUS_CODE, CONVISO_PLATFORM_HEADERS,None)
  contCompany = 1
  for company in companies["data"]["companies"]["collection"]:
    print("Empresas - {0}/{1}".format(contCompany,companies["data"]["companies"]["metadata"]["totalCount"]))
    cont = 1
    while True:
      variables = {"page":cont,"scopeIdEq":int(company["id"])}
      projects = runQuery(CONVISO_PLATFORM_URI, QUERYPROJECTS, CONVISO_PLATFORM_STATUS_CODE, CONVISO_PLATFORM_HEADERS, variables)
      contProject = 1
      for project in projects["data"]["projects"]["collection"]:
        print("Projetos - {0}/{1}".format(contProject,projects["data"]["projects"]["metadata"]["total"]))
        if (("gap" in project["label"].lower()) or (len(project["playbooks"]) > 0 and int(project["playbooks"][0]["id"]) == 164)):
          texto = "{0};{1};{2}\n".format(company["id"], company["label"],project["label"])
          projs.append(texto)
        contProject = contProject + 1
      cont = cont + 1
      if cont > projects["data"]["projects"]["metadata"]["totalPages"]:
        break
    contCompany = contCompany + 1

  if len(projs) > 0:
    saveFile(fileName, projs)
    print("Arquivo gerado com sucesso!!!")
  else:
    "Nenhuma empresa/projeto encontrada!"          
    
def saveFile(fileName, lines):
  with open(fileName, 'w', encoding='UTF8') as f:
    f.writelines(lines)
        

if __name__ == '__main__':
  if sys.version_info.major >= 3:     
    while True:
      print('-----Menu-----')
      print('0 - Sair')
      print('1 - Gerar csv Clientes/Projetos com o Requirement Desejado:')
      key = input('Escolha a opção desejada:')
      if key == '0':
          print('programa finalizado!')
          break
      elif key == '1':
          runSearchCompanyProjectRequirement()
