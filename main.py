import os, requests, sys
import subprocess
import importlib.util
from colorama import init, Fore, Style

init(autoreset=True)

owner = 'Tux-MacG1v'
repo = 'A-ULTIMATE-FINDER'

def import_module(module_name):
    try:
        importlib.import_module(module_name)
        return True
    except ImportError:
        return False
def install_package(package_name):
    try:
        subprocess.run([sys.executable, "-m", "pip", "install", package_name], check=True)
        return True
    except subprocess.CalledProcessError:
        return False
required_libraries = ['requests', 'colorama']
for lib in required_libraries:
    if not import_module(lib):
        print(f"Installing {lib}...")
        if install_package(lib):
            print(f"{lib} installed successfully.")
        else:
            print(f"Failed to install {lib}. Exiting.")
            exit(1)

def download_github_release(owner, repo):
    tag_name_url = "https://raw.githubusercontent.com/Tux-MacG1v/A-ULTIMATE-FINDER/main/version.db"
    tag_name = requests.get(tag_name_url).text.strip()
    url = f"https://api.github.com/repos/{owner}/{repo}/releases/tags/{tag_name}"
    response = requests.get(url)
    release_info = response.json()
    current_dir = os.getcwd()
    for asset in release_info['assets']:
        asset_url = asset['browser_download_url']
        asset_name = asset['name']
        download_path = os.path.join(current_dir, asset_name)
        print(Fore.GREEN + f"DOWNLOADING ULTIMATE-FINDER..." + Style.RESET_ALL)
        response = requests.get(asset_url)
        with open(download_path, 'wb') as f:
            f.write(response.content)

    print(Fore.GREEN + "LATEST TOOL DOWNLOAD COMPLETE..." + Style.RESET_ALL)
    print()
    executable_path = os.path.join(current_dir, asset_name)
    print(Fore.YELLOW + "TOOL IS RUNNING..." + Style.RESET_ALL)
    subprocess.run(executable_path, shell=True)

if __name__ == "__main__":
    print(Fore.CYAN + "SEARCHING FOR LATEST..." + Style.RESET_ALL)
    download_github_release(owner, repo)
