import os
from dotenv import load_dotenv
import requests
from datetime import datetime, timedelta
from collections import defaultdict

# Carrega as variáveis de ambiente do arquivo .env
load_dotenv()

# Obtém a chave API da variável de ambiente
api_key = os.getenv("API_KEY")

# Verifica se a chave API foi carregada corretamente
if not api_key:
    raise ValueError("API_KEY não encontrada. Verifique se a variável de ambiente está definida corretamente.")

# Define o cabeçalho com a chave API obtida da variável de ambiente
headers = {
    "x-api-key": api_key,
    "Content-Type": "application/json"
}

# Define a URL da API
url = "https://app.convisoappsec.com/graphql"

# Define a query para obter os ativos com o campo updatedAt
query = """
{
  assets(companyId: "439", page: 1, limit: 100) {
    collection {
      id
      name
      assetsTagList
      updatedAt
    }
    metadata {
      currentPage
      totalPages
    }
  }
}
"""

# Define as tags específicas que você deseja considerar
specific_tags = ['AURA', 'Lab-IA', '4P Microservicos', '4P FAST DATA', '4P APIS', '4P BIG DATA', 'MLOps']

# Define as integrações específicas que você deseja considerar
integration_tags = {
    'FORTIFY': 'FORTIFY',
    'DEPENDENCY_TRACK': 'DEPENDENCY_TRACK'
}

# Pergunta ao usuário qual informação ele deseja validar
print("Selecione a operação desejada:")
print("1: Puxar o volume de ativos por tag")
print("2: Puxar o volume total de ativos por tags")
print("3: Puxar ativos que não foram sincronizados nos últimos 30 dias")
operation = input("Digite o número da operação desejada: ")

if operation == '1':
    print("Selecione uma das tags:")
    for i, tag in enumerate(specific_tags, 1):
        print(f"{i}: {tag}")
    tag_index = int(input("Digite o número da tag desejada: ")) - 1
    selected_tag = specific_tags[tag_index]
elif operation == '2':
    selected_tag = None
elif operation == '3':
    selected_tag = None
else:
    raise ValueError("Operação inválida selecionada.")

# Função para fazer a requisição e obter todos os ativos
def fetch_all_assets(url, headers, query):
    all_assets = []
    page = 1
    while True:
        response = requests.post(url, headers=headers, json={"query": query.replace("page: 1", f"page: {page}")})
        if response.status_code == 200:
            data = response.json()
            assets = data['data']['assets']['collection']
            all_assets.extend(assets)
            metadata = data['data']['assets']['metadata']
            if metadata['currentPage'] >= metadata['totalPages']:
                break
            page += 1
        else:
            print(f"Erro na requisição: {response.status_code}")
            break
    return all_assets

# Função para sincronizar um ativo
def sync_asset(asset_id, integration):
    mutation = """
    mutation SyncAsset($input: SyncAssetInput!) {
      syncAsset(input: $input) {
        assetId
        clientMutationId
        failureReason
        integration
      }
    }
    """
    variables = {
        "input": {
            "assetId": asset_id,
            "integration": integration
        }
    }
    response = requests.post(url, headers=headers, json={"query": mutation, "variables": variables})
    return response.json()

# Obtém todos os ativos
print("Buscando todos os ativos...")
assets = fetch_all_assets(url, headers, query)

# Filtra os ativos conforme a operação selecionada
filtered_assets = []
tag_groups = defaultdict(list)

if operation == '1':
    for asset in assets:
        if selected_tag in asset['assetsTagList']:
            filtered_assets.append(asset)
            tag_groups[selected_tag].append(asset)
elif operation == '2':
    for asset in assets:
        for tag in asset['assetsTagList']:
            if tag in specific_tags:
                filtered_assets.append(asset)
                tag_groups[tag].append(asset)
elif operation == '3':
    days_threshold = 30
    now = datetime.now().astimezone()
    for asset in assets:
        updated_at = datetime.fromisoformat(asset['updatedAt'])
        if now - updated_at > timedelta(days=days_threshold):
            filtered_assets.append(asset)
            for tag in asset['assetsTagList']:
                if tag in specific_tags:
                    tag_groups[tag].append(asset)

# Cria a mensagem para exibir no terminal
message_lines = []
if operation == '1':
    message_lines.append(f"TAG: {selected_tag}, Ativos:")
    for asset in tag_groups[selected_tag]:
        message_lines.append(f"  - Asset ID: {asset['id']}, Nome: {asset['name']}, Última Sincronização: {asset['updatedAt']}")
    message_lines.append(f"Total de ativos para a TAG {selected_tag}: {len(tag_groups[selected_tag])}")
elif operation == '2':
    for tag, assets in tag_groups.items():
        message_lines.append(f"TAG: {tag}, Quantidade de Ativos: {len(assets)}")
    total_assets_with_tags = len(filtered_assets)
    message_lines.append(f"\nTotal de todos os ativos que contêm as tags específicas: {total_assets_with_tags}")
elif operation == '3':
    for tag, assets in tag_groups.items():
        message_lines.append(f"TAG: {tag}, Ativos que não foram sincronizados nos últimos {days_threshold} dias:")
        for asset in assets:
            message_lines.append(f"  - Asset ID: {asset['id']}, Nome: {asset['name']}, Última Sincronização: {asset['updatedAt']}")
        message_lines.append(f"Total de ativos para a TAG {tag}: {len(assets)}")
    total_assets_with_tags = len(filtered_assets)
    message_lines.append(f"\nTotal de todos os ativos que não foram sincronizados nos últimos {days_threshold} dias: {total_assets_with_tags}")

message = "\n".join(message_lines)

# Verifica a mensagem criada
print("Mensagem criada:")
print(message)

# Variáveis para contar os ativos sincronizados por integração
fortify_count = 0
dependency_track_count = 0

# Pergunta ao usuário se deseja sincronizar os ativos
if operation in ['1', '2', '3']:
    user_input = input("Deseja sincronizar os ativos que estão vinculados nas tags selecionadas? (s/n): ").strip().lower()
    if user_input == 's':
        for asset in filtered_assets:
            for tag in asset['assetsTagList']:
                if 'FORTIFY' in integration_tags:
                    sync_response = sync_asset(asset['id'], integration_tags['FORTIFY'])
                    print(f"Verificando se o ativo com ID {asset['id']} tem integração com FORTIFY...")
                    
                    # Verifica se a resposta contém o campo esperado
                    sync_asset_data = sync_response.get('data', {}).get('syncAsset')
                    if sync_asset_data:
                        if sync_asset_data.get('failureReason'):
                            print(f"Ativo com ID {asset['id']} não tem integração com FORTIFY: {sync_asset_data['failureReason']}")
                        else:
                            errors = sync_response.get('errors', [])
                            if errors:
                                error_messages = [error.get('message', 'Erro desconhecido') for error in errors]
                                print(f"Erro ao sincronizar ativo {asset['id']} com FORTIFY: {', '.join(error_messages)}")
                            else:
                                fortify_count += 1
                    else:
                        print(f"Ativo com ID {asset['id']} sincronizado com FORTIFY.")

                if 'DEPENDENCY_TRACK' in integration_tags:
                    sync_response = sync_asset(asset['id'], integration_tags['DEPENDENCY_TRACK'])
                    print(f"Verificando se o ativo com ID {asset['id']} tem integração DEPENDENCY_TRACK...")
                    
                    # Verifica se a resposta contém o campo esperado
                    sync_asset_data = sync_response.get('data', {}).get('syncAsset')
                    if sync_asset_data:
                        if sync_asset_data.get('failureReason'):
                            print(f"Ativo com ID {asset['id']} não tem integração com DEPENDENCY_TRACK: {sync_asset_data['failureReason']}")
                        else:
                            errors = sync_response.get('errors', [])
                            if errors:
                                error_messages = [error.get('message', 'Erro desconhecido') for error in errors]
                                print(f"Erro ao sincronizar ativo {asset['id']} com DEPENDENCY_TRACK: {', '.join(error_messages)}")
                            else:
                                dependency_track_count += 1
                    else:
                        print(f"Ativo com ID {asset['id']} sincronizado com DEPENDENCY_TRACK.")
        
        # Exibe a contagem de ativos sincronizados por integração
        print(f"Sincronização concluída. Total de ativos sincronizados com FORTIFY: {fortify_count}")
        print(f"Total de ativos sincronizados com DEPENDENCY_TRACK: {dependency_track_count}")
