export interface ApiResponse<T> {
  data?: T;
  error?: {
    code: string;
    message: string;
  };
}

export interface SearchResult {
  animeId: string;
  title: string;
  episodes: number;
  source: string;
  imageUrl?: string;
}

export interface Download {
  downloadId: string;
  animeId: string;
  episodeNumber: number;
  status: 'pending' | 'running' | 'completed' | 'failed';
  progress: number;
}

export interface Job {
  jobId: string;
  status: 'pending' | 'running' | 'completed' | 'failed';
  progressPercent: number;
  errorMessage?: string;
}

class ApiClient {
  private baseUrl: string;

  constructor(baseUrl = '/api') {
    this.baseUrl = baseUrl;
  }

  private async request<T>(
    method: string,
    path: string,
    body?: unknown
  ): Promise<ApiResponse<T>> {
    const options: RequestInit = { method };
    if (body) {
      options.headers = { 'Content-Type': 'application/json' };
      options.body = JSON.stringify(body);
    }

    try {
      const response = await fetch(`${this.baseUrl}${path}`, options);
      if (!response.ok) {
        const error = await response.json();
        throw new Error(error.error?.message || 'API request failed');
      }
      return response.json();
    } catch (error) {
      return {
        error: {
          code: 'NETWORK_ERROR',
          message: error instanceof Error ? error.message : 'Unknown error',
        },
      };
    }
  }

  async search(query: string): Promise<SearchResult[]> {
    const response = await this.request<SearchResult[]>('GET', `/search?q=${encodeURIComponent(query)}`);
    return response.data || [];
  }

  async listDownloads(): Promise<Download[]> {
    const response = await this.request<Download[]>('GET', '/downloads');
    return response.data || [];
  }

  async createDownload(animeId: string, episodeNumber: number): Promise<Download> {
    const response = await this.request<Download>('POST', '/downloads', {
      anime_id: animeId,
      episode_number: episodeNumber,
    });
    return response.data!;
  }

  async getDownload(downloadId: string): Promise<Download | null> {
    const response = await this.request<Download>('GET', `/downloads/${downloadId}`);
    return response.data || null;
  }

  subscribeToJobProgress(jobId: string): EventSource {
    return new EventSource(`${this.baseUrl}/jobs/${jobId}/progress`);
  }
}

export const apiClient = new ApiClient();
