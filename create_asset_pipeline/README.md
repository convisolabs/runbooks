
## Sobre o script

Este script é utilizado para automatizar a criação de ativos (assets) na plataforma Conviso AppSec através de consultas GraphQL. Vou explicar cada parte do script em detalhes:

*Instalação do jq:*

apt-get update && apt-get install -y jq: Estes comandos são utilizados para atualizar os pacotes do sistema e instalar a ferramenta jq. O jq é uma ferramenta de linha de comando usada para processamento de JSON, que será utilizada posteriormente no script.

*create_asset:*

Este é um job definido em um pipeline de CI/CD. Ele é responsável por criar ou verificar a existência de um ativo na plataforma Conviso.

*Obtenção do nome do projeto:*

PROJECT_NAME=$(basename "$CI_PROJECT_DIR"): Este comando extrai o nome do projeto atual a partir do diretório do projeto no ambiente de integração contínua. O nome do projeto é importante para identificar o ativo na plataforma Conviso.

*Consulta para verificar a existência do ativo:*

RESPONSE=$(curl -s -X POST \ ... ): Este comando envia uma consulta GraphQL para a API da Conviso para verificar se já existe um ativo com o mesmo nome na plataforma. A consulta busca um ativo com o nome igual ao $PROJECT_NAME.
Se um ativo com o mesmo nome existir, o seu ID é armazenado em ASSET_ID.

*Criação do ativo:*

Se não for encontrado um ativo com o mesmo nome, o script envia uma nova consulta GraphQL para criar um novo ativo na plataforma Conviso.

O comando curl envia a consulta para a API da Conviso, especificando informações como o companyId (identificador da empresa), o nome do projeto, o impacto nos negócios, a classificação de dados e uma descrição.

Para classificar a severidade do ativo no campo businessImpact há 3 opções 

**VALUES**
```
LOW
MEDIUM
HIGH
```
Após a criação bem-sucedida, o ID do novo ativo é armazenado em ASSET_ID.

**Feedback de status:**

O script fornece feedback detalhado sobre se o ativo foi encontrado ou criado com sucesso, ou se houve falha na criação, exibindo mensagens adequadas.

Este script automatiza o processo de criação de ativos na plataforma Conviso AppSec, tornando-o adequado para integração em pipelines de CI/CD. Ele usa consultas GraphQL para interagir com a API da plataforma e realizar as operações desejadas.


## Opcional

Caso queira usar a o Conviso AST para criar o ativo e realizar o scan do código, basta remover o # do script.