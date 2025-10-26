export interface Album {
  title: string
  artist: string
  genre: string
  year: number
  cuts: string[]
  url: string
  artwork: string
  comment: string
  rating: number
}

export interface Book {
  title: string
  author: string
  genre: string
  year: number
  url: string
  cover: string
  comment: string
  rating: number
}

export interface Movie {
  title: string
  director: string
  genre: string
  year: number
  url: string
  comment: string
  rating: number
}

export interface DialogFormData {
  title: string
  artist?: string
  author?: string
  director?: string
  genre: string
  year: number
  url: string
  artwork?: string
  cover?: string
  comment: string
  rating: number
  cuts?: string[]
}

export type TabType = 'music' | 'books' | 'movies'
