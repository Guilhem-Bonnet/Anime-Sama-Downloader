"""
Anime Search Engine for Anime-Sama
Scrapes the catalogue dynamically to find animes
"""

import re
import unicodedata
import requests
from difflib import SequenceMatcher
from bs4 import BeautifulSoup
from utils.var import print_status, Colors
from utils.http_pool import cached_get
from utils.config import get_site_base_url_override, load_config, save_config
from utils.anime_db import anilist_search_titles


SITE_BASE_URL = get_site_base_url_override() or "https://anime-sama.si"


def canonicalize_site_url(url: str) -> str:
    """Force le domaine canonique Anime-Sama.

    Note: anime-sama.tv redirige actuellement vers anime-sama.si sans conserver le chemin,
    donc on remplace explicitement le domaine pour préserver /catalogue/.../.
    """
    return re.sub(r"^https://anime-sama\.(?:tv|org|fr|si)/", SITE_BASE_URL + "/", (url or ""))

# Local cache of popular animes (canonicalized to anime-sama.tv)
# This avoids scraping issues and provides instant search
ANIME_CACHE = [
    {"title": "Kaiju No. 8", "url": "https://anime-sama.tv/catalogue/kaiju-n8/"},
    {"title": "Sword Art Online", "url": "https://anime-sama.tv/catalogue/sword-art-online/"},
    {"title": "Shingeki no Kyojin (L'Attaque des Titans)", "url": "https://anime-sama.tv/catalogue/shingeki-no-kyojin/"},
    {"title": "One Piece", "url": "https://anime-sama.tv/catalogue/one-piece/"},
    {"title": "Naruto", "url": "https://anime-sama.tv/catalogue/naruto/"},
    {"title": "Naruto Shippuden", "url": "https://anime-sama.tv/catalogue/naruto-shippuden/"},
    {"title": "Demon Slayer (Kimetsu no Yaiba)", "url": "https://anime-sama.tv/catalogue/kimetsu-no-yaiba/"},
    {"title": "My Hero Academia (Boku no Hero Academia)", "url": "https://anime-sama.tv/catalogue/boku-no-hero-academia/"},
    {"title": "Death Note", "url": "https://anime-sama.tv/catalogue/death-note/"},
    {"title": "Tokyo Ghoul", "url": "https://anime-sama.tv/catalogue/tokyo-ghoul/"},
    {"title": "Fullmetal Alchemist Brotherhood", "url": "https://anime-sama.tv/catalogue/fullmetal-alchemist-brotherhood/"},
    {"title": "Hunter x Hunter", "url": "https://anime-sama.tv/catalogue/hunter-x-hunter/"},
    {"title": "Steins;Gate", "url": "https://anime-sama.tv/catalogue/steins-gate/"},
    {"title": "Code Geass", "url": "https://anime-sama.tv/catalogue/code-geass/"},
    {"title": "Cowboy Bebop", "url": "https://anime-sama.tv/catalogue/cowboy-bebop/"},
    {"title": "Dragon Ball Z", "url": "https://anime-sama.tv/catalogue/dragon-ball-z/"},
    {"title": "Dragon Ball Super", "url": "https://anime-sama.tv/catalogue/dragon-ball-super/"},
    {"title": "One Punch Man", "url": "https://anime-sama.tv/catalogue/one-punch-man/"},
    {"title": "Mob Psycho 100", "url": "https://anime-sama.tv/catalogue/mob-psycho-100/"},
    {"title": "Jujutsu Kaisen", "url": "https://anime-sama.tv/catalogue/jujutsu-kaisen/"},
    {"title": "Chainsaw Man", "url": "https://anime-sama.tv/catalogue/chainsaw-man/"},
    {"title": "Spy x Family", "url": "https://anime-sama.tv/catalogue/spy-x-family/"},
    {"title": "Vinland Saga", "url": "https://anime-sama.tv/catalogue/vinland-saga/"},
    {"title": "Bleach", "url": "https://anime-sama.tv/catalogue/bleach/"},
    {"title": "Fairy Tail", "url": "https://anime-sama.tv/catalogue/fairy-tail/"},
    {"title": "Black Clover", "url": "https://anime-sama.tv/catalogue/black-clover/"},
    {"title": "Dr. Stone", "url": "https://anime-sama.tv/catalogue/dr-stone/"},
    {"title": "The Promised Neverland (Yakusoku no Neverland)", "url": "https://anime-sama.tv/catalogue/yakusoku-no-neverland/"},
    {"title": "Violet Evergarden", "url": "https://anime-sama.tv/catalogue/violet-evergarden/"},
    {"title": "Made in Abyss", "url": "https://anime-sama.tv/catalogue/made-in-abyss/"},
    {"title": "Re:Zero", "url": "https://anime-sama.tv/catalogue/re-zero/"},
    {"title": "Overlord", "url": "https://anime-sama.tv/catalogue/overlord/"},
    {"title": "Konosuba", "url": "https://anime-sama.tv/catalogue/kono-subarashii-sekai-ni-shukufuku-wo/"},
    {"title": "No Game No Life", "url": "https://anime-sama.tv/catalogue/no-game-no-life/"},
    {"title": "Tokyo Revengers", "url": "https://anime-sama.tv/catalogue/tokyo-revengers/"},
    {"title": "Blue Lock", "url": "https://anime-sama.tv/catalogue/blue-lock/"},
    {"title": "Haikyuu", "url": "https://anime-sama.tv/catalogue/haikyuu/"},
    {"title": "Kuroko no Basket", "url": "https://anime-sama.tv/catalogue/kuroko-no-basket/"},
    {"title": "Slam Dunk", "url": "https://anime-sama.tv/catalogue/slam-dunk/"},
    {"title": "Assassination Classroom (Ansatsu Kyoushitsu)", "url": "https://anime-sama.tv/catalogue/ansatsu-kyoushitsu/"},
    {"title": "Parasyte (Kiseijuu)", "url": "https://anime-sama.tv/catalogue/kiseijuu/"},
    {"title": "Erased (Boku dake ga Inai Machi)", "url": "https://anime-sama.tv/catalogue/boku-dake-ga-inai-machi/"},
    {"title": "Your Name (Kimi no Na wa)", "url": "https://anime-sama.tv/catalogue/kimi-no-na-wa/"},
    {"title": "A Silent Voice (Koe no Katachi)", "url": "https://anime-sama.tv/catalogue/koe-no-katachi/"},
    {"title": "Weathering With You (Tenki no Ko)", "url": "https://anime-sama.tv/catalogue/tenki-no-ko/"},
    {"title": "Frieren", "url": "https://anime-sama.tv/catalogue/frieren/"},
    {"title": "Solo Leveling", "url": "https://anime-sama.tv/catalogue/solo-leveling/"},
    {"title": "Horimiya", "url": "https://anime-sama.tv/catalogue/horimiya/"},
    {"title": "Kaguya-sama Love is War", "url": "https://anime-sama.tv/catalogue/kaguya-sama-wa-kokurasetai/"},
    {"title": "Toradora", "url": "https://anime-sama.tv/catalogue/toradora/"},
    {"title": "Clannad", "url": "https://anime-sama.tv/catalogue/clannad/"},
    {"title": "Angel Beats", "url": "https://anime-sama.tv/catalogue/angel-beats/"},
    {"title": "Anohana", "url": "https://anime-sama.tv/catalogue/anohana/"},
    {"title": "Your Lie in April (Shigatsu wa Kimi no Uso)", "url": "https://anime-sama.tv/catalogue/shigatsu-wa-kimi-no-uso/"},
    {"title": "Roshidere", "url": "https://anime-sama.tv/catalogue/roshidere/"},
    {"title": "Dandadan", "url": "https://anime-sama.tv/catalogue/dandadan/"},
    {"title": "Hero Without a Class (Hazurewaku no Joutai Ijou Skill)", "url": "https://anime-sama.tv/catalogue/hazurewaku-no-joutai-ijou-skill/"},
    {"title": "Shangri-La Frontier", "url": "https://anime-sama.tv/catalogue/shangri-la-frontier/"},
    {"title": "Mashle", "url": "https://anime-sama.tv/catalogue/mashle/"},
    {"title": "The Eminence in Shadow (Kage no Jitsuryokusha)", "url": "https://anime-sama.tv/catalogue/kage-no-jitsuryokusha/"},
    {"title": "Classroom of the Elite (Youkoso Jitsuryoku)", "url": "https://anime-sama.tv/catalogue/youkoso-jitsuryoku/"},
    {"title": "86 Eighty-Six", "url": "https://anime-sama.tv/catalogue/86-eighty-six/"},
    {"title": "Goblin Slayer", "url": "https://anime-sama.tv/catalogue/goblin-slayer/"},
    {"title": "That Time I Got Reincarnated as a Slime", "url": "https://anime-sama.tv/catalogue/tensei-shitara-slime-datta-ken/"},
    {"title": "Rising of the Shield Hero", "url": "https://anime-sama.tv/catalogue/tate-no-yuusha-no-nariagari/"},
    {"title": "Mushoku Tensei", "url": "https://anime-sama.tv/catalogue/mushoku-tensei/"},
    {"title": "Oregairu", "url": "https://anime-sama.tv/catalogue/oregairu/"},
    {"title": "Bunny Girl Senpai", "url": "https://anime-sama.tv/catalogue/bunny-girl-senpai/"},
    {"title": "A Wild Last Boss Appeared", "url": "https://anime-sama.tv/catalogue/a-wild-last-boss-appeared/"},
    {"title": "The New Gate", "url": "https://anime-sama.tv/catalogue/the-new-gate/"},
    {"title": "Tsukimichi", "url": "https://anime-sama.tv/catalogue/tsukimichi/"},
    {"title": "Arifureta", "url": "https://anime-sama.tv/catalogue/arifureta/"},
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

def _slugify_title_for_anime_sama(title: str) -> str:
    """Best-effort slugify to match anime-sama catalogue URLs."""
    s = (title or "").strip().lower()
    if not s:
        return ""

    # Normalize unicode (remove accents)
    s = unicodedata.normalize("NFKD", s)
    s = "".join(ch for ch in s if not unicodedata.combining(ch))

    # Replace common separators/punctuations
    s = s.replace("'", " ")
    s = re.sub(r"[/:+&]", " ", s)
    s = re.sub(r"[^a-z0-9\s-]", " ", s)
    s = re.sub(r"\s+", " ", s).strip()

    s = s.replace(" ", "-")
    s = re.sub(r"-+", "-", s).strip("-")
    return s


def _slug_variants(slug: str) -> list[str]:
    """Génère quelques variantes de slug réalistes (borné).

    Exemple: "hells-paradise" <-> "hell-s-paradise".
    """
    base = (slug or "").strip().lower()
    if not base:
        return []

    variants: list[str] = [base]

    # normaliser les doubles hyphens
    variants.append(re.sub(r"-+", "-", base))

    # apostrophe-s fréquent: hells-... vs hell-s-...
    variants.append(re.sub(r"([a-z])s-([a-z])", r"\1-s-\2", base))
    variants.append(re.sub(r"([a-z])-s-([a-z])", r"\1s-\2", base))

    # dédoublonne en conservant l'ordre
    out: list[str] = []
    seen: set[str] = set()
    for v in variants:
        v = v.strip("-")
        if not v or v in seen:
            continue
        seen.add(v)
        out.append(v)
    return out[:6]


def _config_get(config: dict, *path: str, default=None):
    cur = config
    for key in path:
        if not isinstance(cur, dict) or key not in cur:
            return default
        cur = cur[key]
    return cur


def _config_set(config: dict, value, *path: str) -> None:
    cur = config
    for key in path[:-1]:
        cur = cur.setdefault(key, {})
    cur[path[-1]] = value


def _looks_like_cloudflare(resp: requests.Response) -> bool:
    try:
        if "cdn-cgi" in (resp.url or ""):
            return True
        server = (resp.headers.get("server") or "").lower()
        if "cloudflare" in server:
            text = (resp.text or "").lower()
            if "attention required" in text or "cf-browser-verification" in text:
                return True
    except Exception:
        return False
    return False


def _probe_anime_sama_catalogue_slug(slug: str) -> str | None:
    """Return canonical anime-sama URL if slug exists, else None."""
    slug = (slug or "").strip().strip("/")
    if not slug:
        return None

    url = f"{SITE_BASE_URL}/catalogue/{slug}/"
    headers = {
        "User-Agent": "Mozilla/5.0 (X11; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/119.0",
        "Accept": "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
        "Accept-Language": "fr,fr-FR;q=0.8,en-US;q=0.5,en;q=0.3",
        "Connection": "keep-alive",
    }

    try:
        resp = cached_get(url, headers=headers, use_cache=True, timeout=10, allow_redirects=True)
        if resp.status_code != 200:
            return None
        if _looks_like_cloudflare(resp):
            return None
        return url
    except Exception:
        return None


def resolve_anime_sama_base_url(query: str, provider: str = "anilist") -> str | None:
    """Resolve a base catalogue URL (no season/lang) from a free-text query.

    Strategy:
    1) Use an external DB (default AniList) to get title variants/synonyms.
    2) Slugify and probe anime-sama.tv/catalogue/<slug>/.
    3) Fallback to existing fuzzy search + cache/homepage scraping.
    """

    q = (query or "").strip()
    if not q:
        return None

    config = load_config()
    cache_key = q.lower()
    cached_url = _config_get(config, "search", "resolved_urls", cache_key)
    if isinstance(cached_url, str) and cached_url.startswith(SITE_BASE_URL):
        return cached_url

    titles: list[str] = [q]
    if provider == "anilist":
        try:
            titles.extend(anilist_search_titles(q).as_list())
        except Exception:
            # Silent fallback to local search
            pass

    # Generate slug candidates (bounded to avoid many probes)
    slug_candidates: list[str] = []
    seen: set[str] = set()
    for t in titles:
        base_slug = _slugify_title_for_anime_sama(t)
        for slug in _slug_variants(base_slug):
            if not slug:
                continue
            if slug in seen:
                continue
            seen.add(slug)
            slug_candidates.append(slug)
            if len(slug_candidates) >= 25:
                break
        if len(slug_candidates) >= 25:
            break

    for slug in slug_candidates:
        found = _probe_anime_sama_catalogue_slug(slug)
        if found:
            _config_set(config, found, "search", "resolved_urls", cache_key)
            save_config(config)
            return found

    # Fallback: existing fuzzy search
    results = search_anime(q, limit=1)
    if results and results[0]["score"] > 0.5:
        url = results[0]["url"]
        url = canonicalize_site_url(url)
        _config_set(config, url, "search", "resolved_urls", cache_key)
        save_config(config)
        return url

    return None


def get_anime_list():
    """
    Récupère dynamiquement la liste des animes depuis anime-sama.tv
    Combine le cache local + scraping de la page d'accueil
    Utilise la page d'accueil qui n'est pas protégée par Cloudflare
    """
    # Commencer avec le cache local (animes populaires)
    all_animes = {}
    for anime in ANIME_CACHE:
        canon = canonicalize_site_url(anime["url"])
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
            href_canon = canonicalize_site_url(href) if href.startswith("http") else href
            if '/catalogue/' in href_canon and href_canon.startswith(SITE_BASE_URL):
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

    # Try a direct resolution (DB-backed) first to reduce friction.
    try:
        direct = resolve_anime_sama_base_url(query, provider="anilist")
    except Exception:
        direct = None

    if direct:
        print_status(f"Direct match found: {direct}", "success")
        confirm = input(f"{Colors.BOLD}Use this URL? (Y/n): {Colors.ENDC}").strip().lower()
        if confirm != 'n':
            return direct
    
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

def quick_search(query: str, provider: str = "anilist"):
    """
    Quick search - returns best match URL or None
    
    Args:
        query: Search query string
    
    Returns:
        URL of best match, or None if no good match found
    """
    if provider == "local":
        results = search_anime(query, limit=1)
        if results and results[0]['score'] > 0.5:
            return canonicalize_site_url(results[0]['url'])
        return None

    return resolve_anime_sama_base_url(query, provider="anilist")
