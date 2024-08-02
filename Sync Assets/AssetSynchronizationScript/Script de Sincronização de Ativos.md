# Script de Sincronização de Ativos

## Visão Geral

Este script é projetado para interagir com a API da Conviso AppSec para buscar e sincronizar ativos com base em vários critérios. Ele suporta buscar ativos por tags, obter o número total de ativos com tags específicas e identificar ativos que não foram sincronizados nos últimos 30 dias.

## Recursos

- Buscar e listar ativos por tags específicas.
- Obter o número total de ativos associados a tags específicas.
- Identificar ativos que não foram sincronizados nos últimos 30 dias.
- Sincronizar ativos com integrações específicas (FORTIFY e DEPENDENCY_TRACK).

## Requisitos

- Python 3.x
- Biblioteca `requests`
- Biblioteca `python-dotenv`

## Instalação

1. Clone o repositório ou baixe o script.

2. Instale as bibliotecas necessárias usando pip:

    ```bash
    pip install requests python-dotenv
    ```

3. Crie um arquivo `.env` no diretório raiz do projeto e adicione sua chave de API:

    ```env
    API_KEY=sua_chave_de_api_aqui
    ```



## Company ID

Na query é necessário adicionar CompanyID, exemplo.

`assets(companyId: "439", page: 1, limit: 100)`

## Define as tags específicas que você deseja considerar

No código adicione as tags que a API deve consultar na CP da Conviso, exemplo.

`specific_tags = ['AURA', 'Lab-IA', '4P Microservicos', '4P FAST DATA', '4P APIS', '4P BIG DATA', 'MLOps']`

## Uso
1. Execute o script:

    ```bash
    python sync.py
    ```

2. Siga os prompts para selecionar uma operação:
    - **1**: Buscar o volume de ativos por tag.
    - **2**: Buscar o volume total de ativos por tags.
    - **3**: Buscar ativos que não foram sincronizados nos últimos 30 dias.

3. Se solicitado, insira a tag ou o número da operação conforme necessário.

4. O script exibirá os resultados no terminal e pedirá para sincronizar os ativos. Se você escolher sincronizar, ele processará os ativos e fornecerá o status da sincronização.

## Explicação do Código

- **Carregamento de Variáveis de Ambiente**:

  O script usa `python-dotenv` para carregar chaves de API e outras configurações a partir de um arquivo `.env`.

- **Busca de Ativos**:

  A função `fetch_all_assets` recupera todos os ativos da API paginando através dos resultados.

- **Sincronização de Ativos**:

  A função `sync_asset` envia solicitações de sincronização para a API para as integrações especificadas.

- **Operações**:

  - **Operação 1**: Lista ativos filtrados pela tag selecionada.
  - **Operação 2**: Conta ativos associados a cada tag especificada.
  - **Operação 3**: Lista ativos que não foram sincronizados nos últimos 30 dias.

## Tratamento de Erros

Se houver um problema com a sincronização, o script imprimirá mensagens de erro no terminal. Certifique-se de que sua chave de API está correta e que você tem acesso adequado na CP.


