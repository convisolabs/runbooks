from fileinput import filename
import requests
import os
from datetime import *
from enum import Enum




class METHOD(Enum):
    PERCRITICY = 0
    PERCRITICYANDDATE = 1


CONVISO_PLATFORM_URI =  'https://app.convisoappsec.com/graphql'
CONVISO_PLATFORM_TOKEN = os.environ.get('CONVISO_PLATFORM_TOKEN','token')
CONVISO_PLATFORM_STATUS_CODE = 200
CONVISO_PLATFORM_HEADERS = {'x-api-key':CONVISO_PLATFORM_TOKEN}

QUERYVULNERABILITIES = """
    query Vulnerabilities($criticityEq:String, $page:Int, $projectScopeIdEq:Int){
        vulnerabilities(
            page: $page
            limit: 100
            params: {
                projectScopeIdEq:$projectScopeIdEq,
                criticityEq:$criticityEq
            }
            order: project_id
            orderType: DESC
        ) {
            collection {
                projectId,
                vulnerabilityStatus,
                description,
                title,
                criticity,
                createdAt,
          	    updatedAt,
                project{
                id,
                label,
                assets{
                    id,
                    name
                }
                }
            }
            metadata {
            currentPage
            limitValue
            totalCount
            totalCritical
            totalHigh
            totalLow
            totalMedium
            totalNotifications
            totalPages
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

def saveFile(fileName, lines):
    with open(fileName, 'w', encoding='UTF8') as f:
        f.writelines(lines)

def loadingVulns(kindOf):
    cont = 1
    vulns = []   
    clientId = input('Digite o código do cliente no Conviso Platform:') 
    fileName = input('Nome do Arquivo:') 
    criticity = input('Digite a criticidade desejada (low,medium,high,critical):') 
    objDate = date.today()
    if kindOf == METHOD.PERCRITICYANDDATE:
        dateString = input('Digite uma data de corte (YYYY-MM-DD):')
        objDate = datetime.strptime(dateString, '%Y-%m-%d')

    while True:
        variables = {'criticityEq':criticity,'projectScopeIdEq':int(clientId),'page': cont}
        vulnerabilities = runQuery(CONVISO_PLATFORM_URI, QUERYVULNERABILITIES, CONVISO_PLATFORM_STATUS_CODE, CONVISO_PLATFORM_HEADERS, variables)
        for vulnerability in vulnerabilities['data']['vulnerabilities']['collection']:
            if len(vulnerability['project']['assets']) > 0 and ((vulnerability['vulnerabilityStatus'] == 'identified' or vulnerability['vulnerabilityStatus'] == 'fix_refused')):
                if kindOf == METHOD.PERCRITICY or (kindOf == METHOD.PERCRITICYANDDATE and datetime.strptime(vulnerability['updatedAt'].split('T')[0], '%Y-%m-%d') >= objDate):
                    texto = "{0};{1};{2}\n".format(vulnerability['project']['assets'][0]['name'],vulnerability['project']['label'], vulnerability['title'])
                    vulns.append(texto)
        cont = cont + 1
        if cont > vulnerabilities['data']['vulnerabilities']['metadata']['totalPages']:
            break
    if len(vulns) > 0:
        saveFile(fileName, vulns)
        print("Arquivo gerado com sucesso!!!")
    else:
        "Nenhum vulnerabilidade encontrada com o filtro selecionado!"

if __name__ == '__main__':
     while True:
        print('-----Menu-----')
        print('0 - Sair')
        print('1 - Gerar csv com vulnerabilidades por data e criticidade:')
        print('2 - Gerar csv com vulnerabilidades por criticidade:')
        key = input('Escolha a opção desejada:')
        if key == '0':
            print('programa finalizado!')
            break
        elif key == '1':
            loadingVulns(METHOD.PERCRITICYANDDATE)
        elif key == '2':
            loadingVulns(METHOD.PERCRITICY)

