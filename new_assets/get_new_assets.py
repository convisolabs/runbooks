from bs4 import BeautifulSoup
import re
import requests

scope_id = 430 # Vivo - Telefonica scopeId
cookie = "" # get cookie from app.conviso.com.br site
url = 'https://app.conviso.com.br/scopes/{scope_id}/integrations/fortify/select_projects?page={{page_number}}'.format(scope_id=scope_id)
headers = {
  "accept": "*/*;q=0.5, text/javascript, application/javascript, application/ecmascript, application/x-ecmascript",
    "accept-language": "pt-BR,pt;q=0.9,en-US;q=0.8,en;q=0.7",
    "if-none-match": "W/\"de422fb04ff7bbf12eaa311b3ab29048\"",
    "sec-fetch-dest": "empty",
    "sec-fetch-mode": "cors",
    "sec-fetch-site": "same-origin",
    "sec-gpc": "1",
    "x-csrf-token": "Z7PHbAZLfhkxir+R7J38fRhuFWOocgNa0GzMsRqAoXZ73cfauhWelNAVHsa+1JWh0DsWYLlfypIIxy0RdFMCog==",
    "x-requested-with": "XMLHttpRequest",
    "cookie": cookie,
    "Referer": "https://app.conviso.com.br/scopes/{}/integrations/fortify/select_projects?locale=en".format(scope_id),
    "Referrer-Policy": "strict-origin-when-cross-origin"
}

new_assets = []

# remove jquery structures
def format_soup(data):
  # regex: ^\$\(\'\#projects\'\).html\(\".*\);$
  lfilter, rfilter = ("$('#projects').html(\"\\n", ");")
  formated_data = (data[data.index(lfilter)+len(lfilter):data.index(rfilter)]).replace('\\"', '"').replace("\\/", "/")
  return BeautifulSoup(formated_data, "html.parser")
  
def get_project_name(div_str):
  return re.search("(.*)Project: ([a-zA-Z0-9\-\_ ]*)", div_str).group(2).strip()

def get_new_assets(data):
  assets = data.findAll("div", {"class": "col-lg-6 col-md-6 col-xs-12"})
  # print(data.prettify())
  for asset in assets:
    has_new_tag = asset.find("p", {"class": "tag"})
    if has_new_tag is not None:
      project_div = asset.find("strong").getText()
      project_name = get_project_name(project_div)
      new_assets.append(project_name)
      print("[{}] New asset: {}".format(len(new_assets), project_name))

def get_page(page_number):
  f_url = url.format(page_number=page_number)
  response = requests.get(f_url, headers=headers)
  return response.text, response.status_code

def dump_to_file():
  with open('new_assets.txt', 'w') as f:
    try:
      for new_asset in new_assets:
        f.write(new_asset + "\n")
    except Exception as e:
      print("Error on file write: ", e)

if __name__ == "__main__":
  page_number = 1
  try:
    while True:
      print("Page number: {}".format(page_number))
      content, status_code = get_page(page_number)
      if status_code == 500:
        break
      soup = format_soup(content)
      get_new_assets(soup)
      page_number += 1 
  except Exception as e:
    print("Error on page {}".format(page_number), e)
  finally:
    dump_to_file()
