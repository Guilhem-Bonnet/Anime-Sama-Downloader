"""
Anime Search Engine for Anime-Sama
Scrapes the catalogue dynamically to find animes
"""

import re
import requests
from difflib import SequenceMatcher
from bs4 import BeautifulSoup
from utils.var import print_status, Colors


SITE_BASE_URL = "https://anime-sama.tv"


def _canonicalize_anime_sama_url(url: str) -> str:
    """Force anime-sama.tv comme domaine canonique (org/fr -> tv)."""
    return re.sub(r"^https://anime-sama\.(?:org|fr)/", "https://anime-sama.tv/", url)

# Local cache of popular animes on anime-sama.org
# This avoids scraping issues and provides instant search
ANIME_CACHE = [
    {"title": "Kaiju No. 8", "url": "https://anime-sama.org/catalogue/kaiju-n8/"},
    {"title": "Sword Art Online", "url": "https://anime-sama.org/catalogue/sword-art-online/"},
    {"title": "Shingeki no Kyojin (L'Attaque des Titans)", "url": "https://anime-sama.org/catalogue/shingeki-no-kyojin/"},
    {"title": "One Piece", "url": "https://anime-sama.org/catalogue/one-piece/"},
    {"title": "Naruto", "url": "https://anime-sama.org/catalogue/naruto/"},
    {"title": "Naruto Shippuden", "url": "https://anime-sama.org/catalogue/naruto-shippuden/"},
    {"title": "Demon Slayer (Kimetsu no Yaiba)", "url": "https://anime-sama.org/catalogue/kimetsu-no-yaiba/"},
    {"title": "My Hero Academia (Boku no Hero Academia)", "url": "https://anime-sama.org/catalogue/boku-no-hero-academia/"},
    {"title": "Death Note", "url": "https://anime-sama.org/catalogue/death-note/"},
    {"title": "Tokyo Ghoul", "url": "https://anime-sama.org/catalogue/tokyo-ghoul/"},
    {"title": "Fullmetal Alchemist Brotherhood", "url": "https://anime-sama.org/catalogue/fullmetal-alchemist-brotherhood/"},
    {"title": "Hunter x Hunter", "url": "https://anime-sama.org/catalogue/hunter-x-hunter/"},
    {"title": "Steins;Gate", "url": "https://anime-sama.org/catalogue/steins-gate/"},
    {"title": "Code Geass", "url": "https://anime-sama.org/catalogue/code-geass/"},
    {"title": "Cowboy Bebop", "url": "https://anime-sama.org/catalogue/cowboy-bebop/"},
    {"title": "Dragon Ball Z", "url": "https://anime-sama.org/catalogue/dragon-ball-z/"},
    {"title": "Dragon Ball Super", "url": "https://anime-sama.org/catalogue/dragon-ball-super/"},
    {"title": "One Punch Man", "url": "https://anime-sama.org/catalogue/one-punch-man/"},
    {"title": "Mob Psycho 100", "url": "https://anime-sama.org/catalogue/mob-psycho-100/"},
    {"title": "Jujutsu Kaisen", "url": "https://anime-sama.org/catalogue/jujutsu-kaisen/"},
    {"title": "Chainsaw Man", "url": "https://anime-sama.org/catalogue/chainsaw-man/"},
    {"title": "Spy x Family", "url": "https://anime-sama.org/catalogue/spy-x-family/"},
    {"title": "Vinland Saga", "url": "https://anime-sama.org/catalogue/vinland-saga/"},
    {"title": "Bleach", "url": "https://anime-sama.org/catalogue/bleach/"},
    {"title": "Fairy Tail", "url": "https://anime-sama.org/catalogue/fairy-tail/"},
    {"title": "Black Clover", "url": "https://anime-sama.org/catalogue/black-clover/"},
    {"title": "Dr. Stone", "url": "https://anime-sama.org/catalogue/dr-stone/"},
    {"title": "The Promised Neverland (Yakusoku no Neverland)", "url": "https://anime-sama.org/catalogue/yakusoku-no-neverland/"},
    {"title": "Violet Evergarden", "url": "https://anime-sama.org/catalogue/violet-evergarden/"},
    {"title": "Made in Abyss", "url": "https://anime-sama.org/catalogue/made-in-abyss/"},
    {"title": "Re:Zero", "url": "https://anime-sama.org/catalogue/re-zero/"},
    {"title": "Overlord", "url": "https://anime-sama.org/catalogue/overlord/"},
    {"title": "Konosuba", "url": "https://anime-sama.org/catalogue/kono-subarashii-sekai-ni-shukufuku-wo/"},
    {"title": "No Game No Life", "url": "https://anime-sama.org/catalogue/no-game-no-life/"},
    {"title": "Tokyo Revengers", "url": "https://anime-sama.org/catalogue/tokyo-revengers/"},
    {"title": "Blue Lock", "url": "https://anime-sama.org/catalogue/blue-lock/"},
    {"title": "Haikyuu", "url": "https://anime-sama.org/catalogue/haikyuu/"},
    {"title": "Kuroko no Basket", "url": "https://anime-sama.org/catalogue/kuroko-no-basket/"},
    {"title": "Slam Dunk", "url": "https://anime-sama.org/catalogue/slam-dunk/"},
    {"title": "Assassination Classroom (Ansatsu Kyoushitsu)", "url": "https://anime-sama.org/catalogue/ansatsu-kyoushitsu/"},
    {"title": "Parasyte (Kiseijuu)", "url": "https://anime-sama.org/catalogue/kiseijuu/"},
    {"title": "Erased (Boku dake ga Inai Machi)", "url": "https://anime-sama.org/catalogue/boku-dake-ga-inai-machi/"},
    {"title": "Your Name (Kimi no Na wa)", "url": "https://anime-sama.org/catalogue/kimi-no-na-wa/"},
    {"title": "A Silent Voice (Koe no Katachi)", "url": "https://anime-sama.org/catalogue/koe-no-katachi/"},
    {"title": "Weathering With You (Tenki no Ko)", "url": "https://anime-sama.org/catalogue/tenki-no-ko/"},
    {"title": "Frieren", "url": "https://anime-sama.org/catalogue/frieren/"},
    {"title": "Solo Leveling", "url": "https://anime-sama.org/catalogue/solo-leveling/"},
    {"title": "Horimiya", "url": "https://anime-sama.org/catalogue/horimiya/"},
    {"title": "Kaguya-sama Love is War", "url": "https://anime-sama.org/catalogue/kaguya-sama-wa-kokurasetai/"},
    {"title": "Toradora", "url": "https://anime-sama.org/catalogue/toradora/"},
    {"title": "Clannad", "url": "https://anime-sama.org/catalogue/clannad/"},
    {"title": "Angel Beats", "url": "https://anime-sama.org/catalogue/angel-beats/"},
    {"title": "Anohana", "url": "https://anime-sama.org/catalogue/anohana/"},
    {"title": "Your Lie in April (Shigatsu wa Kimi no Uso)", "url": "https://anime-sama.org/catalogue/shigatsu-wa-kimi-no-uso/"},
    {"title": "Roshidere", "url": "https://anime-sama.org/catalogue/roshidere/"},
    {"title": "Dandadan", "url": "https://anime-sama.org/catalogue/dandadan/"},
    {"title": "Hero Without a Class (Hazurewaku no Joutai Ijou Skill)", "url": "https://anime-sama.org/catalogue/hazurewaku-no-joutai-ijou-skill/"},
    {"title": "Shangri-La Frontier", "url": "https://anime-sama.org/catalogue/shangri-la-frontier/"},
    {"title": "Mashle", "url": "https://anime-sama.org/catalogue/mashle/"},
    {"title": "The Eminence in Shadow (Kage no Jitsuryokusha)", "url": "https://anime-sama.org/catalogue/kage-no-jitsuryokusha/"},
    {"title": "Classroom of the Elite (Youkoso Jitsuryoku)", "url": "https://anime-sama.org/catalogue/youkoso-jitsuryoku/"},
    {"title": "86 Eighty-Six", "url": "https://anime-sama.org/catalogue/86-eighty-six/"},
    {"title": "Goblin Slayer", "url": "https://anime-sama.org/catalogue/goblin-slayer/"},
    {"title": "That Time I Got Reincarnated as a Slime", "url": "https://anime-sama.org/catalogue/tensei-shitara-slime-datta-ken/"},
    {"title": "Rising of the Shield Hero", "url": "https://anime-sama.org/catalogue/tate-no-yuusha-no-nariagari/"},
    {"title": "Mushoku Tensei", "url": "https://anime-sama.org/catalogue/mushoku-tensei/"},
    {"title": "Oregairu", "url": "https://anime-sama.org/catalogue/oregairu/"},
    {"title": "Bunny Girl Senpai", "url": "https://anime-sama.org/catalogue/bunny-girl-senpai/"},
    {"title": "A Wild Last Boss Appeared", "url": "https://anime-sama.org/catalogue/a-wild-last-boss-appeared/"},
    {"title": "The New Gate", "url": "https://anime-sama.org/catalogue/the-new-gate/"},
    {"title": "Tsukimichi", "url": "https://anime-sama.org/catalogue/tsukimichi/"},
    {"title": "Arifureta", "url": "https://anime-sama.org/catalogue/arifureta/"},
]

# Common translations for better matching
TRANSLATIONS = {
    "attack on titan": "shingeki no kyojin",
    "l'attaque des titans": "shingeki no kyojin",
    "attaque des titans": "shingeki no kyojin",
    "demon slayer": "kimetsu no yaiba",
    "sword art online": "sword art online",
    "sao": "sword art online",
    "one piece": "one piece",
    "naruto": "naruto",
    "dragon ball": "dragon ball",
    "my hero academia": "boku no hero academia",
    "death note": "death note",
    "tokyo ghoul": "tokyo ghoul",
    "fullmetal alchemist": "fullmetal alchemist",
    "hunter x hunter": "hunter x hunter",
    "steins gate": "steins gate",
    "code geass": "code geass",
    "cowboy bebop": "cowboy bebop",
    "kaiju": "kaiju n8",
    "promised neverland": "yakusoku no neverland",
}

def get_anime_list():
    """
    Récupère dynamiquement la liste des animes depuis anime-sama.org
    Combine le cache local + scraping de la page d'accueil
    Utilise la page d'accueil qui n'est pas protégée par Cloudflare
    """
    # Commencer avec le cache local (animes populaires)
    all_animes = {}
    for anime in ANIME_CACHE:
        canon = _canonicalize_anime_sama_url(anime["url"])
        all_animes[canon] = {"title": anime["title"], "url": canon}
    
    try:
        # La page d'accueil fonctionne, pas le catalogue (403 Cloudflare)
        homepage_url = f"{SITE_BASE_URL}/"
        
        headers = {
            'User-Agent': 'Mozilla/5.0 (X11; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/119.0',
            'Accept': 'text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8',
            'Accept-Language': 'fr,fr-FR;q=0.8,en-US;q=0.5,en;q=0.3',
            'Connection': 'keep-alive',
        }
        
        print("⏳ Chargement du catalogue...")
        response = requests.get(homepage_url, headers=headers, timeout=15)
        response.raise_for_status()
        
        soup = BeautifulSoup(response.content, 'html.parser')
        
        # Extraire les animes uniques (sans saison/langue)
        for a in soup.find_all('a', href=True):
            href = a['href']
            href_canon = _canonicalize_anime_sama_url(href) if href.startswith("http") else href
            if '/catalogue/' in href_canon and 'anime-sama.tv' in href_canon:
                # Extraire le nom de l'anime (avant /saison ou /scan)
                match = re.search(r'/catalogue/([^/]+)', href_canon)
                if match:
                    anime_slug = match.group(1)
                    anime_url = f'{SITE_BASE_URL}/catalogue/{anime_slug}/'
                    
                    # Éviter les doublons (cache prioritaire)
                    if anime_url not in all_animes:
                        # Récupérer le titre depuis le texte du lien
                        title = a.get_text(strip=True)
                        if not title or len(title) > 100 or title.startswith('Anime') or title.startswith('Manga'):
                            # Utiliser le slug comme titre (nettoyé)
                            title = anime_slug.replace('-', ' ').title()
                        
                        all_animes[anime_url] = {
                            'title': title,
                            'url': anime_url
                        }
        
        animes = list(all_animes.values())
        print(f"✅ {len(animes)} animes trouvés ({len(ANIME_CACHE)} cache + {len(animes) - len(ANIME_CACHE)} en ligne)")
        return animes
            
    except requests.exceptions.RequestException as e:
        print(f"❌ Échec de récupération: {e}")
        print(f"⚠️ Utilisation du cache local uniquement ({len(ANIME_CACHE)} animes)")
        return list(all_animes.values())
    except Exception as e:
        print(f"❌ Erreur d'analyse: {e}")
        print(f"⚠️ Utilisation du cache local uniquement ({len(ANIME_CACHE)} animes)")
        return list(all_animes.values())

def normalize_title(title):
    """Normalize title for better matching"""
    # Remove special characters and convert to lowercase
    title = title.lower().strip()
    title = re.sub(r'[:\-–—]', ' ', title)
    title = re.sub(r'\s+', ' ', title)
    # Remove common words that don't help matching
    remove_words = ['saison', 'season', 'vostfr', 'vf', 'vo', 'the', 'le', 'la', 'les']
    words = title.split()
    words = [w for w in words if w not in remove_words]
    return ' '.join(words)

def similarity_score(a, b):
    """Calculate similarity between two strings (0-1)"""
    return SequenceMatcher(None, normalize_title(a), normalize_title(b)).ratio()

def search_anime(query, limit=10):
    """
    Search for anime by query string
    
    Args:
        query: Search query (e.g., "kaiju", "attack on titan")
        limit: Maximum number of results to return
    
    Returns:
        List of matching animes sorted by relevance
    """
    # Check if query matches a known translation
    query_lower = query.lower().strip()
    if query_lower in TRANSLATIONS:
        query = TRANSLATIONS[query_lower]
        print_status(f"Translation: '{query_lower}' → '{query}'", "info")
    
    # Get anime list from cache
    animes = get_anime_list()
    
    if not animes:
        print_status("Anime cache is empty", "error")
        return []
    
    # Normalize query
    normalized_query = normalize_title(query)
    
    # Calculate similarity scores
    results = []
    for anime in animes:
        anime_normalized = normalize_title(anime['title'])
        score = similarity_score(query, anime['title'])
        
        # Bonus for exact substring match
        if normalized_query in anime_normalized:
            score += 0.3
        
        # Bonus for start of title
        if anime_normalized.startswith(normalized_query):
            score += 0.2
        
        # Bonus for exact word match
        query_words = set(normalized_query.split())
        anime_words = set(anime_normalized.split())
        if query_words & anime_words:  # Intersection
            score += 0.1 * len(query_words & anime_words)
        
        results.append({
            'title': anime['title'],
            'url': anime['url'],
            'score': score
        })
    
    # Sort by score (best first)
    results.sort(key=lambda x: x['score'], reverse=True)
    
    # Return top results
    return results[:limit]

def interactive_search():
    """Interactive search with user selection"""
    print(f"\n{Colors.BOLD}{Colors.HEADER}🔍 ANIME SEARCH{Colors.ENDC}")
    print(f"{Colors.OKCYAN}Search the catalogue by name{Colors.ENDC}")
    print(f"{Colors.OKCYAN}Or press Enter to skip and enter URL manually{Colors.ENDC}\n")
    
    query = input(f"{Colors.BOLD}Enter anime name to search (or Enter to skip): {Colors.ENDC}").strip()
    
    if not query:
        print_status("Skipping search", "info")
        return None
    
    results = search_anime(query, limit=10)
    
    if not results:
        print_status("No animes found matching your search", "warning")
        retry = input(f"{Colors.BOLD}Try another search? (Y/n): {Colors.ENDC}").strip().lower()
        if retry != 'n':
            return interactive_search()
        return None
    
    # Show results
    print(f"\n{Colors.BOLD}{Colors.OKGREEN}Search Results:{Colors.ENDC}\n")
    
    for i, anime in enumerate(results, 1):
        # Color code by confidence
        if anime['score'] > 0.7:
            color = Colors.OKGREEN  # High confidence
        elif anime['score'] > 0.4:
            color = Colors.OKCYAN   # Medium confidence
        else:
            color = Colors.WARNING   # Low confidence
        
        confidence = int(anime['score'] * 100)
        print(f"{color}{i}. {anime['title']} ({confidence}% match){Colors.ENDC}")
    
    print(f"\n{Colors.BOLD}0. Cancel{Colors.ENDC}")
    
    # User selection
    while True:
        try:
            choice = input(f"\n{Colors.BOLD}Select anime (1-{len(results)}, or 0 to cancel): {Colors.ENDC}").strip()
            
            if not choice:
                continue
            
            choice = int(choice)
            
            if choice == 0:
                print_status("Search cancelled", "warning")
                return None
            
            if 1 <= choice <= len(results):
                selected = results[choice - 1]
                print_status(f"Selected: {selected['title']}", "success")
                return selected['url']
            else:
                print_status(f"Please enter a number between 0 and {len(results)}", "error")
        
        except ValueError:
            print_status("Please enter a valid number", "error")

def quick_search(query):
    """
    Quick search - returns best match URL or None
    
    Args:
        query: Search query string
    
    Returns:
        URL of best match, or None if no good match found
    """
    results = search_anime(query, limit=1)
    
    if results and results[0]['score'] > 0.5:
        return results[0]['url']
    
    return None
