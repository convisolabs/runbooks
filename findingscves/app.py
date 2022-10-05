from fileinput import filename
import requests
import os
from datetime import *
from enum import Enum

CONVISO_PLATFORM_URI =  'https://appsec.convisoappsec.com/graphql'
CONVISO_PLATFORM_TOKEN = os.environ.get('CONVISO_PLATFORM_TOKEN','token')
CONVISO_PLATFORM_STATUS_CODE = 200
CONVISO_PLATFORM_HEADERS = {'x-api-key': CONVISO_PLATFORM_TOKEN}

clientId = input('Digite o código do cliente no Conviso Platform:') 
QUERYFINDINGS = """
query Findings($scopeId:Int){
  findings(page: 1, limit: 1000, params: {
          scopeId : $scopeId }
  ) {
    collection {
      affectedIp
      affectedLineNumbers
      affectedSourceFile
      applicationId
      applicationName
      businessSeverity
      category
      commitReference
      component
      createdAt
      cves
      descriptionEn
      evaluationStatus
      flowDeployId
      flowProjectId
      hashIssue
      id
      integrationName
      originalSeverity
      rawOutputFromScanner
      recommendationsEn
      references
      scanId
      scanType
      shortCommitRef
      source
      sourceFile
      tags
      technicalSeverity
      titleEn
      version
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

def runQuery(uri, query, statusCode, headers, variables):
    request = requests.post(uri, json={'query':query, 'variables':variables}, headers=headers)
    if request.status_code == statusCode:
        return request.json()
    else:
        raise Exception(f'Unexpected status code returned: {request.status_code}')


variables = {'scopeId':int(clientId)}
findings = runQuery(CONVISO_PLATFORM_URI, QUERYFINDINGS, CONVISO_PLATFORM_STATUS_CODE, CONVISO_PLATFORM_HEADERS, variables)
results=findings['data']['findings']['collection']
for i in results:
    if (i['cves']) != []:
        print("Foi identificado o",(i['cves']),"na aplicação",(i['applicationName']),"utilizando a ",(i['integrationName']))
