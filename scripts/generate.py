import sys
from g4f.client import Client
import requests
from g4f.cookies import set_cookies

def download_image(image_url, prompt):
    filename = prompt.replace(" ", "_") + ".jpg"
    response = requests.get(image_url)
    if response.status_code == 200:
        with open(filename, 'wb') as f:
            f.write(response.content)
    return filename

def generate_and_download_image(prompt):
    client = Client()
    response = client.images.generate(
    model="gemini",
    prompt=prompt,
    )
    image_url = response.data[0].url
    
    filename = download_image(image_url, prompt)
    print(filename)

if __name__ == "__main__":
    try:
        if len(sys.argv) > 1:
            prompt = sys.argv[1]
            ssid = sys.argv[2]
            set_cookies(".google.com", {
            "__Secure-1PSID": ssid
            })
            generate_and_download_image(prompt)
        else:
            print("Please provide a prompt.")
    except Exception as e:
        print("An error occurred:", e)
