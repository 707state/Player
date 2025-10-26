import type { Album, Book, Movie } from "./types.ts";
const API_BASE_URL = "http://localhost:8080";
// API service for music
export const musicApi = {
  // Get all albums
  async getAlbums(params?: {
    title?: string;
    artist?: string;
    genre?: string;
    year?: number;
    rating?: number;
  }): Promise<Album[]> {
    const queryParams = new URLSearchParams();
    if (params) {
      Object.entries(params).forEach(([key, value]) => {
        if (value !== undefined && value !== "") {
          queryParams.append(key, value.toString());
        }
      });
    }

    const response = await fetch(`${API_BASE_URL}/music?${queryParams}`);
    if (!response.ok) {
      throw new Error(`Failed to fetch albums: ${response.statusText}`);
    }
    return response.json();
  },

  // Create or update album
  async saveAlbum(album: Album): Promise<{ message: string }> {
    const response = await fetch(`${API_BASE_URL}/music`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(album),
    });
    if (!response.ok) {
      throw new Error(`Failed to save album: ${response.statusText}`);
    }
    return response.json();
  },

  // Delete album
  async deleteAlbum(album: {
    title: string;
    artist: string;
  }): Promise<{ message: string }> {
    const response = await fetch(`${API_BASE_URL}/music`, {
      method: "DELETE",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(album),
    });
    if (!response.ok) {
      throw new Error(`Failed to delete album: ${response.statusText}`);
    }
    return response.json();
  },
};

// API service for books
export const booksApi = {
  // Get all books
  async getBooks(params?: {
    title?: string;
    author?: string;
    genre?: string;
    year?: number;
    rating?: number;
  }): Promise<Book[]> {
    const queryParams = new URLSearchParams();
    if (params) {
      Object.entries(params).forEach(([key, value]) => {
        if (value !== undefined && value !== "") {
          queryParams.append(key, value.toString());
        }
      });
    }

    const response = await fetch(`${API_BASE_URL}/books?${queryParams}`);
    if (!response.ok) {
      throw new Error(`Failed to fetch books: ${response.statusText}`);
    }
    return response.json();
  },

  // Create or update book
  async saveBook(book: Book): Promise<{ message: string }> {
    const response = await fetch(`${API_BASE_URL}/books`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(book),
    });
    if (!response.ok) {
      throw new Error(`Failed to save book: ${response.statusText}`);
    }
    return response.json();
  },

  // Delete book
  async deleteBook(book: {
    title: string;
    author: string;
  }): Promise<{ message: string }> {
    const response = await fetch(`${API_BASE_URL}/books`, {
      method: "DELETE",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(book),
    });
    if (!response.ok) {
      throw new Error(`Failed to delete book: ${response.statusText}`);
    }
    return response.json();
  },
};

// API service for movies
export const moviesApi = {
  // Get all movies
  async getMovies(params?: {
    title?: string;
    director?: string;
    genre?: string;
    year?: number;
    rating?: number;
  }): Promise<Movie[]> {
    const queryParams = new URLSearchParams();
    if (params) {
      Object.entries(params).forEach(([key, value]) => {
        if (value !== undefined && value !== "") {
          queryParams.append(key, value.toString());
        }
      });
    }

    const response = await fetch(`${API_BASE_URL}/movies?${queryParams}`);
    if (!response.ok) {
      throw new Error(`Failed to fetch movies: ${response.statusText}`);
    }
    return response.json();
  },

  // Create or update movie
  async saveMovie(movie: Movie): Promise<{ message: string }> {
    const response = await fetch(`${API_BASE_URL}/movies`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(movie),
    });
    if (!response.ok) {
      throw new Error(`Failed to save movie: ${response.statusText}`);
    }
    return response.json();
  },

  // Delete movie
  async deleteMovie(movie: {
    title: string;
    director: string;
  }): Promise<{ message: string }> {
    const response = await fetch(`${API_BASE_URL}/movies`, {
      method: "DELETE",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(movie),
    });
    if (!response.ok) {
      throw new Error(`Failed to delete movie: ${response.statusText}`);
    }
    return response.json();
  },
};
