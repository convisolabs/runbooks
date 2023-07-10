import requests
import re

base_url = 'https://app.conviso.com.br/scopes/REPLACECOMPANYID/assets/REPLACEASSETID?locale=en'
headers = {
    'authority': 'app.conviso.com.br',
    'accept': 'text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7',
    'accept-language': 'pt-BR,pt;q=0.9,en-US;q=0.8,en;q=0.7,es;q=0.6',
    'X-Armature-Api-Key': 'REPLACEAPIKEY',
    'if-none-match': 'W/"eb6d93c4f9c8a3febb0bff41f3525ab2"',
    'referer': 'https://app.conviso.com.br/scopes/REPLACECOMPANYID/integrations/fortify?utf8=%E2%9C%93&status=error&name=&button=',
    'sec-ch-ua': '"Not.A/Brand";v="8", "Chromium";v="114", "Google Chrome";v="114"',
    'sec-ch-ua-mobile': '?0',
    'sec-ch-ua-platform': '"Windows"',
    'sec-fetch-dest': 'document',
    'sec-fetch-mode': 'navigate',
    'sec-fetch-site': 'same-origin',
    'sec-fetch-user': '?1',
    'upgrade-insecure-requests': '1',
    'user-agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36',
}

params = {
    'button': '',
    'locale': 'en',
    'name': '',
    'page': 1,
    'status': 'success',
    'utf8': '✓'
}

response = requests.get(base_url, headers=headers, params=params)
if response.status_code == 200:
    data = response.text

    date_pattern = r'\b([0-9]{2}/){2}[0-9]{4} [0-9]{2}:[0-9]{2}:[0-9]{2} [AP]M\b'

    date_match = re.search(date_pattern, data)

    if date_match:
        date = date_match.group()
        print(f'{date}')
