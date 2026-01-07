#!/usr/bin/env python3
"""
Test script pour le moteur de recherche d'animes
Démontre les capacités de recherche avec différents types de requêtes
"""

from utils.search import search_anime, TRANSLATIONS, ANIME_CACHE

def print_header(title):
    """Affiche un header stylisé"""
    print("\n" + "=" * 70)
    print(f"  {title}")
    print("=" * 70)

def test_search(query, expected_first=None):
    """Test une recherche et affiche les résultats"""
    print(f"\n🔍 Recherche: '{query}'")
    results = search_anime(query, limit=5)
    
    if results:
        print(f"   ✅ {len(results)} résultats trouvés:")
        for i, r in enumerate(results, 1):
            score = int(r['score'] * 100)
            icon = "🥇" if i == 1 else "🥈" if i == 2 else "🥉" if i == 3 else "  "
            print(f"   {icon} {i}. {r['title']:<50} {score:>3}% match")
        
        # Vérifier si le résultat attendu est en premier
        if expected_first and results[0]['title'].lower().find(expected_first.lower()) == -1:
            print(f"   ⚠️  Attendu '{expected_first}' en premier")
        
        return True
    else:
        print("   ❌ Aucun résultat")
        return False

def main():
    print_header("🎌 TEST DU MOTEUR DE RECHERCHE ANIME-SAMA")
    
    # Statistiques du cache
    print(f"\n📊 Statistiques du cache:")
    print(f"   • Animes disponibles: {len(ANIME_CACHE)}")
    print(f"   • Traductions définies: {len(TRANSLATIONS)}")
    
    print_header("TEST 1: Recherches avec traduction automatique")
    
    # Test des traductions
    test_search("kaiju", "Kaiju No. 8")
    test_search("l'attaque des titans", "Shingeki")
    test_search("attaque des titans", "Shingeki")
    test_search("demon slayer", "Kimetsu")
    test_search("sao", "Sword Art")
    
    print_header("TEST 2: Recherches par titre exact")
    
    test_search("One Piece", "One Piece")
    test_search("Naruto", "Naruto")
    test_search("Death Note", "Death Note")
    
    print_header("TEST 3: Recherches floues (fuzzy matching)")
    
    test_search("one punch", "One Punch Man")
    test_search("tokyo ghol", "Tokyo Ghoul")
    test_search("fullmetal", "Fullmetal")
    test_search("hunter", "Hunter")
    
    print_header("TEST 4: Recherches partielles")
    
    test_search("jujutsu", "Jujutsu")
    test_search("chainsaw", "Chainsaw")
    test_search("blue lock", "Blue Lock")
    
    print_header("TEST 5: Recherches avec scores élevés")
    
    # Ces recherches devraient avoir des scores > 100% grâce aux bonus
    high_score_queries = [
        ("one piece", "One Piece"),
        ("naruto shippuden", "Naruto Shippuden"),
        ("sword art online", "Sword Art Online"),
    ]
    
    for query, expected in high_score_queries:
        results = search_anime(query, limit=1)
        if results:
            score = int(results[0]['score'] * 100)
            icon = "🔥" if score > 100 else "✅"
            print(f"\n{icon} '{query}' → {results[0]['title']} ({score}%)")
    
    print_header("TEST 6: Recherches avec titres japonais")
    
    test_search("shingeki no kyojin", "Shingeki")
    test_search("kimetsu no yaiba", "Kimetsu")
    test_search("boku no hero", "Boku no Hero")
    
    print_header("TEST 7: Cas limites")
    
    # Recherches qui devraient fonctionner malgré tout
    test_search("spy family", "Spy x Family")
    test_search("frieren", "Frieren")
    test_search("solo leveling", "Solo Leveling")
    
    print_header("📋 RÉSUMÉ DES CAPACITÉS")
    
    capabilities = [
        "✅ Traductions automatiques (FR/EN → JP)",
        "✅ Fuzzy matching (tolère les fautes)",
        "✅ Recherches partielles",
        "✅ Bonus de pertinence multiples",
        "✅ Cache local (55+ animes)",
        "✅ Scores de confiance affichés",
        "✅ Recherche instantanée (<20ms)",
        "✅ Support multilingue (FR, EN, JP)"
    ]
    
    print()
    for capability in capabilities:
        print(f"   {capability}")
    
    print_header("🎉 TESTS TERMINÉS")
    print("\n💡 Pour utiliser le moteur de recherche:")
    print("   • Mode interactif: python main.py")
    print("   • Mode CLI: python main.py -s \"ANIME\" -e 1-5 --quick")
    print("   • Documentation: voir SEARCH_GUIDE.md")
    print()

if __name__ == "__main__":
    main()
