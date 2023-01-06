#!/bin/bash
input="wordlist"
NOW="$(date)"
echo ""
echo ""
echo "---------------Verificação Iniciada $NOW-----------------"
#Sincronização de ativos
while IFS= read -r line
do
SYNC=$(echo "$line")
echo "-------------------------------------------------------------------"
echo "-------------------------------------------------------------------"
sleep 2
POST=$(echo | curl 'https://app.conviso.com.br/scopes/439/assets/'$SYNC'/sync_findings?integration_type=fortify&locale=en' \
  -H 'Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9' \
  -H 'Accept-Language: pt-BR,pt;q=0.9,en-US;q=0.8,en;q=0.7,es;q=0.6' \
  -H 'Cache-Control: max-age=0' \
  -H 'Connection: keep-alive' \
  -H 'Content-Type: application/x-www-form-urlencoded' \
  -H 'Cookie: ADICIONE O COOKIE AQUI' \
  -H 'Origin: https://app.conviso.com.br' \
  -H 'Referer: https://app.conviso.com.br/scopes/439/assets/'$SYNC'?locale=en' \
  -H 'Sec-Fetch-Dest: document' \
  -H 'Sec-Fetch-Mode: navigate' \
  -H 'Sec-Fetch-Site: same-origin' \
  -H 'Sec-Fetch-User: ?1' \
  -H 'Upgrade-Insecure-Requests: 1' \
  -H 'User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36' \
  -H 'sec-ch-ua: "Not?A_Brand";v="8", "Chromium";v="108", "Google Chrome";v="108"' \
  -H 'sec-ch-ua-mobile: ?0' \
  -H 'sec-ch-ua-platform: "Windows"' \
  --data-raw 'authenticity_token='ADICIONE AQUI O TOKEN''\
  --compressed)

echo "RESULTADOS POST : $POST"


sleep 10
done <"$input"
