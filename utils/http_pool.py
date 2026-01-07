"""
HTTP Connection Pool and Cache Manager
Optimizes network requests with connection pooling and response caching.
"""

import requests
from requests.adapters import HTTPAdapter
from urllib3.util.retry import Retry
from functools import lru_cache
import hashlib

class HTTPPool:
    """Manages HTTP connection pool with retry logic."""
    
    _instance = None
    
    def __new__(cls):
        if cls._instance is None:
            cls._instance = super(HTTPPool, cls).__new__(cls)
            cls._instance._initialized = False
        return cls._instance
    
    def __init__(self):
        if self._initialized:
            return
            
        self._initialized = True
        
        # Create session with connection pooling
        self.session = requests.Session()
        
        # Configure retry strategy
        retry_strategy = Retry(
            total=3,
            backoff_factor=1,
            status_forcelist=[429, 500, 502, 503, 504],
            allowed_methods=["HEAD", "GET", "OPTIONS"]
        )
        
        # Configure adapter with connection pooling
        adapter = HTTPAdapter(
            pool_connections=10,
            pool_maxsize=20,
            max_retries=retry_strategy,
            pool_block=False
        )
        
        # Mount adapter for both http and https
        self.session.mount("http://", adapter)
        self.session.mount("https://", adapter)
    
    def get(self, url, **kwargs):
        """Make GET request with connection pooling."""
        return self.session.get(url, **kwargs)
    
    def post(self, url, **kwargs):
        """Make POST request with connection pooling."""
        return self.session.post(url, **kwargs)
    
    def close(self):
        """Close session and cleanup resources."""
        self.session.close()


class ResponseCache:
    """Simple response cache for frequently accessed URLs."""
    
    _cache = {}
    _max_cache_size = 100
    
    @classmethod
    def _get_key(cls, url, headers=None):
        """Generate cache key from URL and headers."""
        key_str = url
        if headers:
            key_str += str(sorted(headers.items()))
        return hashlib.md5(key_str.encode()).hexdigest()
    
    @classmethod
    def get(cls, url, headers=None):
        """Get cached response if available."""
        key = cls._get_key(url, headers)
        return cls._cache.get(key)
    
    @classmethod
    def set(cls, url, response_data, headers=None):
        """Cache response data."""
        # Limit cache size
        if len(cls._cache) >= cls._max_cache_size:
            # Remove oldest entry (simple FIFO)
            cls._cache.pop(next(iter(cls._cache)))
        
        key = cls._get_key(url, headers)
        cls._cache[key] = response_data
    
    @classmethod
    def clear(cls):
        """Clear all cached responses."""
        cls._cache.clear()


# Global HTTP pool instance
_http_pool = HTTPPool()


def get_http_pool():
    """Get the global HTTP pool instance."""
    return _http_pool


def cached_get(url, headers=None, use_cache=True, **kwargs):
    """
    Make GET request with caching support.
    
    Args:
        url: URL to fetch
        headers: Optional headers dict
        use_cache: Whether to use cache (default: True)
        **kwargs: Additional arguments for requests.get
    
    Returns:
        Response object
    """
    # Check cache first
    if use_cache:
        cached = ResponseCache.get(url, headers)
        if cached:
            return cached
    
    # Make request with connection pool
    pool = get_http_pool()
    response = pool.get(url, headers=headers, **kwargs)
    
    # Cache successful responses
    if use_cache and response.status_code == 200:
        ResponseCache.set(url, response, headers)
    
    return response
