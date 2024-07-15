from sys import argv

FILE_TO_OPEN = None
API_KEY = None
COMPANY_ID = None


def create_asset(asset_name: str):
    import requests as req

    global API_KEY
    global COMPANY_ID

    headers = {
        'x-api-key': API_KEY
    }

    json_data = {
        'operationName': 'CreateAsset',
        'variables': {
            'asset': {
                'name': asset_name,
                'teamIds': [],
                'companyId': int(COMPANY_ID),
            },
        },
        'query': 'mutation CreateAsset($asset: CreateAssetInput!) {\n  createAsset(input: $asset) {\n    asset {\n      id\n    }\n    errors\n  }\n}',
    }

    response = req.post('https://app.convisoappsec.com/graphql', headers=headers, json=json_data)
    print(response.text)


if 5 > argv.__len__() < 4:
    print(f"[-] Please execute {argv[0]} file.txt company_id api_key")
    exit(1)

FILE_TO_OPEN = argv[1]
COMPANY_ID = argv[2]
API_KEY = argv[3]

if not API_KEY:
    print('[-] Api-key not defined.')

if not COMPANY_ID:
    print('[-] Company ID not defined.')

if FILE_TO_OPEN:
    f = open(argv[1], 'r')
    lines = f.readlines()
    f.close()

    for l in lines:
        create_asset(l)
