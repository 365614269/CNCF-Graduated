import os
import requests
import yaml

sesn = requests.session()


def getsuffix(filename):  # Get the file suffix.
    for i in range(len(filename) - 1, -1, -1):
        if filename[i] == ".":
            return filename[i + 1:]


def upload(text, kbId):  # Upload the text to the website.
    url = "https://fastgpt.daocloud.io/api/openapi/kb/pushData"

    payload = {
        "kbId": kbId,
        "mode": "index",
        "prompt": "",
        "data": [
            {
                "q": text,
            }
        ]
    }

    headers = {
        "apikey": secrets["apikey"],
        "Content-Type": "application/json"
    }

    response = sesn.post(url, headers=headers, json=payload)

    print(response.text)


if __name__ == "__main__":
    repos = os.listdir("modules")
    for repo in repos:
        f = open(f'kbId/{repo}/secret.yaml', 'r', encoding='utf-8')
        secrets = yaml.load(f.read(), Loader=yaml.FullLoader)
        kbId = secrets["kbId"]
        for root, dirs, files in os.walk(f"modules/{repo}"):
            for file in files:
                if getsuffix(file) == "md":  # We only care about .md files
                    cwd = root + '/' + file

                    lines = []
                    with open(cwd, 'rb') as f:
                        bintext = f.readlines()

                    for i in bintext:
                        try:  # Ignore the characters that cannot be decoded
                            lines.append(i.decode('gb2312'))
                        except:
                            pass

                    text = "".join(lines)
                    upload(text, kbId)