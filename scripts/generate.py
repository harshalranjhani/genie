import sys
from g4f.client import Client
import requests
from g4f.cookies import set_cookies

def download_image(image_url, prompt):
    # Create a filename based on the prompt (sanitize if necessary)
    filename = prompt.replace(" ", "_") + ".jpg"
    response = requests.get(image_url)
    if response.status_code == 200:
        with open(filename, 'wb') as f:
            f.write(response.content)
    return filename

def generate_and_download_image(prompt):
    # Your existing code to generate the image URL
    # For example, let's say image_url is the result
    client = Client()
    response = client.images.generate(
    model="gemini",
    prompt=prompt,
    )
    image_url = response.data[0].url
    
    # Download the image and get the filename
    filename = download_image(image_url, prompt)
    print(filename)

if __name__ == "__main__":
    if len(sys.argv) > 1:
        prompt = sys.argv[1]
        ssid = sys.argv[2]
        set_cookies(".google.com", {
        "__Secure-1PSID": "cookie value"
        })
        generate_and_download_image(prompt)
    else:
        print("Please provide a prompt.")
