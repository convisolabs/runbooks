import requests
import json
import os
import csv


from enum import Enum
from datetime import datetime,date


class HTTP_METHOD(Enum):
    GET = 0
    POST = 1

CONVISO_PLATFORM_URI =  'https://app.convisoappsec.com/graphql'
CONVISO_PLATFORM_TOKEN = os.environ.get('CONVISO_PLATFORM_TOKEN','token')
CONVISO_PLATFORM_STATUS_CODE = 200
CONVISO_PLATFORM_STATUS_CODE_ERROR = 500
CONVISO_PLATFORM_HEADERS = {'x-api-key':CONVISO_PLATFORM_TOKEN}

CLICKUP_URI = 'https://api.clickup.com/api/v2/'
CLICKUP_TOKEN = os.environ.get('CLICKUP_TOKEN','token')
CLICKUP_STATUS_CODE = 200
CLICKUP_HEADERS = {'Content-Type':'application/json','Authorization':CLICKUP_TOKEN}
CLICKUP_LIST_PROJECTS = 163176498

ClickupIdCustomFieldCustomer = 0
ClickupIdCustomFieldDemandType = 0
ClickupIdCustomFieldConvisoLink = 0


dictClickUpCustomer = dict()
dictClickUpDemandType = dict()
dictConvisoPlatformProjectTypes = dict()
dictConvisoPlatformRequirements = dict()


CONVISO_PLATFORM_REQUIREMENTS_QUERY = '''
query Playbooks($companyId:ID!, $page:Int){
  playbooks(id: $companyId, page: $page, limit: 100, params: {}) {
    collection {
      checklistTypeId
      companyId
      createdAt
      deletedAt
      description
      id
      label
      updatedAt
    }
    metadata {
      currentPage
      limitValue
      totalCount
      totalPages
    }
  }
}

'''

CONVISO_PLATFORM_PROJECT_TYPES_QUERY = '''
query ProjectTypes($page:Int){
  projectTypes(page: $page, limit: 100, params: {}) {
    collection {
      code
      description
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
'''

CONVISO_PLATFORM_PROJECT_MUTATION = '''
mutation CreateProject($input:CreateProjectInput!){
  createProject(
    input: $input
  ) {
    clientMutationId
    errors
    project {
      apiCode
      apiResponseReview
      closeComments
      companyId
      connectivity
      continuousDelivery
      contractedHours
      createdAt
      deploySendFrequency
      dueDate
      endDate
      environmentInvaded
      estimatedDays
      estimatedHours
      executiveSummary
      freeRetest
      hasOpenRetest
      hoursOrDays
      id
      integrationDeploy
      inviteToken
      isOpen
      isPublic
      label
      language
      mainRecommendations
      microserviceFolder
      negativeScope
      notificationList
      objective
      pid
      plannedStartedAt
      playbookFinishedAt
      playbookStartedAt
      receiveDeploys
      repositoryUrl
      sacCode
      sacProjectId
      scope
      secretId
      sshPublicKey
      startDate
      status
      students
      subScopeId
      totalAnalysisLines
      totalChangedLines
      totalNewLines
      totalPublishedVulnerabilities
      totalRemovedLines
      type
      updatedAt
      userableId
      userableType
      waiting
    }
  }
}
'''

BANNER = '''
  ____  _       _    __                       ____ _ _      _    _   _       
 |  _ \| | __ _| |_ / _| ___  _ __ _ __ ___  / ___| (_) ___| | _| | | |_ __  
 | |_) | |/ _` | __| |_ / _ \| '__| '_ ` _ \| |   | | |/ __| |/ / | | | '_ \ 
 |  __/| | (_| | |_|  _| (_) | |  | | | | | | |___| | | (__|   <| |_| | |_) |
 |_|   |_|\__,_|\__|_|  \___/|_|  |_| |_| |_|\____|_|_|\___|_|\_\\___/| .__/ 
                                                                      |_|    

'''

def runCallHttpMethod(method, uri, headers, statusCode,jsonData):
    teste=json.dumps(jsonData)
    if method == HTTP_METHOD.GET:
        request = requests.get(uri, headers=headers)
    else:
        request = requests.post(uri, data=teste,headers=headers)
    if request.status_code == statusCode:
        return request.json()
    else:
        raise Exception(f'Unexpected status code returned: {request.status_code}')


def runQuery(uri, query, statusCode, headers, variables):
    request = requests.post(uri, json={'query':query, 'variables':variables}, headers=headers)
    if request.status_code == statusCode:
        return request.json()
    else:
        raise Exception(f'Unexpected status code returned: {request.status_code}')


def runMutationError(uri, query, headers, variables):
    request = requests.post(uri, json={'query':query, 'variables':variables}, headers=headers)
    teste = ""



def SelectTypeId():
    for projectType in dictConvisoPlatformProjectTypes:
        print(f'{projectType};{dictConvisoPlatformProjectTypes[projectType]}')
    return input('Conviso Platform Project TypeId(only number):')

def SelectCustomFieldCustomer():
    cont = 1
    for customer in dictClickUpCustomer:
        print(f'{cont};{dictClickUpCustomer[customer]}')
        cont = cont + 1
    value = input('Digite o código do cliente:')
    return list(dictClickUpCustomer.keys())[int(value) - 1]

def SelectCustomFieldDemandType():
    cont = 1
    for customer in dictClickUpDemandType:
        print(f'{cont};{dictClickUpDemandType[customer]}')
        cont = cont + 1
    value = input('Digite o código da demanda:')
    return list(dictClickUpDemandType.keys())[int(value) - 1]

def SelectProject():
    return input('Conviso Platform Project ID:')

def FillConvisoPlatformFieldsAndClickUpCsv():
    arquivoCSV = input('Digite o caminho completo do arquivo CSV: ')
    projectId = SelectProject()
    typeId = SelectTypeId() 
    custom_field_customer = SelectCustomFieldCustomer()
    custom_field_demand_type = SelectCustomFieldDemandType()

    header = True
    with open(arquivoCSV, 'r') as arqCSV:
        reader = csv.reader(arqCSV)
        for linha in reader:
            if header == True:
                header = False
            elif linha[0] != '':
                label = linha[0]
                goal = linha[1]
                scope = linha[2]
                estimatedHours = linha[3]
                startDate = linha[4]
                finishDate = linha[5]
                objStartDate = datetime.strptime(startDate, '%Y-%m-%d')
                objFinishDate = datetime.strptime(finishDate, '%Y-%m-%d')
                graphQLObj = {
                'companyId': int(projectId),
                'label': label,
                'goal': goal,
                'playbooksIds': [1],
                'scope': scope,
                'typeId': int(typeId),
                'startDate': f'{objStartDate:%Y%m%d}',
                'estimatedHours':estimatedHours,
                'students':20
                }
                clickupObject = {
                    'name':f'[EPIC] - {label}',
                    'description':scope,
                    'status':'backlog',
                    'due_date':objFinishDate.timestamp()*1000,
                    'due_date_time':True,
                    'time_estimate':int(estimatedHours) * 60 * 60 * 1000,
                    'start_date':objStartDate.timestamp()*1000,
                    'start_date_time':True,
                    'notify_all':True,
                    'check_required_custom_fields':True,
                    'custom_fields':[
                        {
                            'id': ClickupIdCustomFieldCustomer,
                            'value':custom_field_customer
                        },
                        {
                            'id': ClickupIdCustomFieldDemandType,
                            'value':custom_field_demand_type
                        },
                        {
                            'id':ClickupIdCustomFieldConvisoLink,
                            'value':'https://app.convisoappsec.com/'
                        }
                    ]
                }         
                variables = {'input': graphQLObj}
                runMutationError(CONVISO_PLATFORM_URI, CONVISO_PLATFORM_PROJECT_MUTATION, CONVISO_PLATFORM_HEADERS, variables)
                uri = f'{CLICKUP_URI}list/{CLICKUP_LIST_PROJECTS}/task'
                runCallHttpMethod(HTTP_METHOD.POST, uri, CLICKUP_HEADERS, CLICKUP_STATUS_CODE, clickupObject)            

def FillConvisoPlatformFieldsAndClickUp():
    projectId = SelectProject()
    label = input('Conviso Platform Project Label:')
    goal = input('Conviso Platform Projetct Goal:')
    scope = input('Conviso Platform Project Scope:')
    typeId = SelectTypeId() 
    estimatedHours = input('Conviso Platform Estimated Hours(only number):')
    startDate = input('Start Date: (YYYY-MM-DD):')
    finishDate = input('Finish Date: (YYYY-MM-DD):')
    custom_field_customer = SelectCustomFieldCustomer()
    custom_field_demand_type = SelectCustomFieldDemandType()
    objStartDate = datetime.strptime(startDate, '%Y-%m-%d')
    objFinishDate = datetime.strptime(finishDate, '%Y-%m-%d')
    objStartDateConvisoPlatform = date.today()
    if (objStartDateConvisoPlatform < datetime.date(objStartDate)):
        objStartDateConvisoPlatform = objStartDate    
    graphQLObj = {
       'companyId': int(projectId),
       'label': label,
       'goal': goal,
       'playbooksIds': [1],
       'scope': scope,
       'typeId': int(typeId),
       'startDate': f'{objStartDateConvisoPlatform:%Y%m%d}',
       'estimatedHours':estimatedHours,
    }
    clickupObject = {
        'name':f'[EPIC] - {label}',
        'description':scope,
        'status':'backlog',
        'due_date':objFinishDate.timestamp()*1000,
        'due_date_time':True,
        'time_estimate':int(estimatedHours) * 60 * 60 * 1000,
        'start_date':objStartDate.timestamp()*1000,
        'start_date_time':True,
        'notify_all':True,
        'check_required_custom_fields':True,
        'custom_fields':[
            {
                'id': ClickupIdCustomFieldCustomer,
                'value':custom_field_customer
            },
            {
                'id': ClickupIdCustomFieldDemandType,
                'value':custom_field_demand_type
            },
            {
                'id':ClickupIdCustomFieldConvisoLink,
                'value':'https://app.convisoappsec.com/'
            }
        ]
    }

    variables = {'input': graphQLObj}
    runMutationError(CONVISO_PLATFORM_URI, CONVISO_PLATFORM_PROJECT_MUTATION, CONVISO_PLATFORM_HEADERS, variables)
    uri = f'{CLICKUP_URI}list/{CLICKUP_LIST_PROJECTS}/task'
    runCallHttpMethod(HTTP_METHOD.POST, uri, CLICKUP_HEADERS, CLICKUP_STATUS_CODE, clickupObject)

def CustomFieldsClickUp():
    uri = f'{CLICKUP_URI}list/{CLICKUP_LIST_PROJECTS}/field'
    request = requests.get(uri, headers=CLICKUP_HEADERS)
    if request.status_code == CLICKUP_STATUS_CODE:
        fields = request.json()
        for field in fields['fields']:
            if field['name'].lower().strip() == 'customer':
                global ClickupIdCustomFieldCustomer 
                ClickupIdCustomFieldCustomer = field['id']
                for option in field['type_config']['options']:
                    dictClickUpCustomer[option['id']] = option['name'] 
            elif field['name'].lower().strip() == 'demand type':
                global ClickupIdCustomFieldDemandType 
                ClickupIdCustomFieldDemandType = field['id']
                for option in field['type_config']['options']:
                    dictClickUpDemandType[option['id']] = option['name'] 
            elif field['name'].lower().strip() == 'conviso link':
                global ClickupIdCustomFieldConvisoLink 
                ClickupIdCustomFieldConvisoLink = field['id']
    else:
        raise Exception(f'Unexpected status code returned: {request.status_code}')


def RequirementsConvisoPlatform(projectId):
    cont = 1
    while True:
        variables = {'companyId':projectId,'page': cont}
        requirements = runQuery(CONVISO_PLATFORM_URI, CONVISO_PLATFORM_REQUIREMENTS_QUERY, CONVISO_PLATFORM_STATUS_CODE, CONVISO_PLATFORM_HEADERS, variables)
        for requirement in requirements['data']['playbooks']['collection']:
            dictConvisoPlatformRequirements[requirement['id']] = requirement['label']
        cont = cont + 1
        if cont > requirements['data']['playbooks']['metadata']['totalPages']:
            break


def ProjectTypesConvisoPlatform():
    cont = 1
    while True:
        variables = {'page': cont}
        projectTypes = runQuery(CONVISO_PLATFORM_URI, CONVISO_PLATFORM_PROJECT_TYPES_QUERY, CONVISO_PLATFORM_STATUS_CODE, CONVISO_PLATFORM_HEADERS, variables)
        for projectType in projectTypes['data']['projectTypes']['collection']:
            dictConvisoPlatformProjectTypes[projectType['id']] = projectType['label']
        cont = cont + 1
        if cont > projectTypes['data']['projectTypes']['metadata']['totalPages']:
            break

def menu():
     while True:
        print('-----Menu-----')
        print('0 - Sair')
        print('1 - Criar Tarefa ConvisoPlatform/ClickUp')
        print('2 - Criar Tarefa ConvisoPlatform/ClickUp By CSV')
        key = input('Escolha a opção desejada:')
        if key == '0':
            print('programa finalizado!')
            break
        elif key == '1':
            FillConvisoPlatformFieldsAndClickUp()
            print('Atividade cadastrada com sucesso!')
        elif key == '2':
            FillConvisoPlatformFieldsAndClickUpCsv()
            print('Atividade cadastrada com sucesso!')            
        else:
            print('Nenhuma opção válida teclada!!!')
            menu()


if __name__ == '__main__':
    CustomFieldsClickUp()
    ProjectTypesConvisoPlatform()
    print(BANNER)
    menu()