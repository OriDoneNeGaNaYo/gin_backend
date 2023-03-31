from bs4 import BeautifulSoup
import requests
import re 

url = "https://www.ansan.go.kr/stat/common/cntnts/selectContents.do?cntnts_id=C0001015"
req = requests.get(url).content

a = []
soup = BeautifulSoup(req, "html.parser")
for table in soup.find_all("table", {"class": "table tl_c"}):
    for tbody in table.find_all("tbody"):
        for t in tbody.find_all("td"):
            str_ = t.text.replace(" ", "")
            pattern = re.compile(r'[가-힣]+').findall(str_)
            if pattern:
                for i in pattern:
                    a.append(i)

import pandas as pd 
data = pd.DataFrame(a, columns=[a[0]])
d = data.drop(data[data["경기도"] == "경기도"].index)

d.to_csv("gyoung.csv", index=False, index_label=False)