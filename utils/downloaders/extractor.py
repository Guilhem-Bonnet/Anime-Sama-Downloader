from utils.var import print_status, Colors
import re
from bs4 import BeautifulSoup
import requests, time, re

def extract_sendvid_video_source(html_content):
    if not html_content:
        return None
    video_source_pattern = r'var\s+video_source\s*=\s*"([^"]+)"'
    match = re.search(video_source_pattern, html_content)
    if match:
        return match.group(1)
    print_status("Could not extract video source from SendVid", "warning")
    return None

def extract_sibnet_video_source(html_content):
    if not html_content:
        return None
    soup = BeautifulSoup(html_content, 'html.parser')
    scripts = soup.find_all('script', type='text/javascript')
    for script in scripts:
        if 'player.src' in script.text:
            match = re.search(r'player\.src\(\[\{.*src:\s*"([^"]+)"', script.text)
            if match:
                video_source = match.group(1)
                if video_source.startswith('//'):
                    video_source = f"https:{video_source}"
                elif not video_source.startswith('https://'):
                    video_source = f"https://video.sibnet.ru{video_source}"
                return video_source
    print_status("Could not extract video source from Sibnet", "warning")
    return None

def extract_oneupload_video_source(html_content):
    if not html_content:
        return None
    soup = BeautifulSoup(html_content, 'html.parser')
    script_tags = soup.find_all('script', type='text/javascript')
    for script in script_tags:
        if script.string and 'jwplayer' in script.string:
            url_match = re.search(r'file:"(https?://.*?)"', script.string)
            if url_match:
                m3u8_url = url_match.group(1)
                return m3u8_url
    print_status("Could not extract video source from OneUpload", "warning")
    return None

def extract_vidmoly_video_source(html_content):
    if not html_content:
        return None
    soup = BeautifulSoup(html_content, 'html.parser')
    script_tags = soup.find_all('script', type='text/javascript')
    for script in script_tags:
        if script.string and 'jwplayer' in script.string:
            url_match = re.search(r'file:"(https?://.*?)"', script.string)
            if url_match:
                m3u8_url = url_match.group(1)
                return m3u8_url
    print_status("Could not extract video source from Vidmoly", "warning")
    return None

def unpack_js_for_ts_file(packed_code, base, count, words):
    def to_base(num, base):
        if num == 0:
            return '0'
        digits = []
        while num:
            digits.append(str(num % base) if num % base < 10 else chr(ord('a') + num % base - 10))
            num //= base
        return ''.join(reversed(digits))
    
    replacements = {}
    for i in range(count):
        key = to_base(i, base)
        if i < len(words) and words[i]:
            replacements[key] = words[i]
    
    unpacked = packed_code
    for key, value in replacements.items():
        pattern = r'\b' + re.escape(key) + r'\b'
        unpacked = re.sub(pattern, value, unpacked)
    
    return unpacked

def extract_packed_code_for_ts(html_content):
    pattern = r"eval\(function\(p,a,c,k,e,d\)\{.*?\}\('(.*?)',(\d+),(\d+),'(.*?)'\.split\('\|'\)\)\)"
    match = re.search(pattern, html_content, re.DOTALL)
    if match:
        packed_code = match.group(1)
        base = int(match.group(2))
        count = int(match.group(3))
        words = match.group(4).split('|')
        return packed_code, base, count, words
    else:
        print("No packed JavaScript code found.")
        return None, None, None, None

def fetch_html_for_ts(url):
    try:
        headers = {
    'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/129.0.0.0 Safari/537.36',
    'Accept': 'text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7',
    'Accept-Language': 'en-US,en;q=0.9',
    'Referer': url.split('/embed/')[0],
    'Connection': 'keep-alive',
}
        response = requests.get(url, headers=headers)
        response.raise_for_status()
        return response.text
    except requests.exceptions.RequestException as e:
        print(f"Error fetching the page: {e}")
        return None
    
def extract_hls_url(unpacked_code):

    pattern = r'var\s+links\s*=\s*\{[^}]*"hls[^"]*"\s*:\s*"([^"]+)"'
    match = re.search(pattern, unpacked_code)
    if match:
        return match.group(1)
    
    pattern2 = r'file\s*:\s*[^"]*"([^"]*\.m3u8[^"]*)"'
    match2 = re.search(pattern2, unpacked_code)
    if match2:
        return match2.group(1)
    
    return None

def extract_movearnpre_video_source(embed_url):
    html_content = fetch_html_for_ts(embed_url)
    if not html_content:
        return
    
    packed_code, base, count, words = extract_packed_code_for_ts(html_content)
    if not packed_code:
        return
    
    unpacked_code = unpack_js_for_ts_file(packed_code, base, count, words)
    
    hls_url = extract_hls_url(unpacked_code)
    if hls_url:
        return hls_url
    else:
        print_status("No HLS URL found in the unpacked code", "error")
        return None
    

