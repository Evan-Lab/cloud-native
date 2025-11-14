export enum PixelColor {
  RED = '#FF0000',
  BLUE = '#0000FF',
  YELLOW = '#FFFF00',
  GREEN = '#00FF00',
  BLACK = '#000000',
  WHITE = '#FFFFFF',
}

export interface Pixel {
  x: number
  y: number
  color: PixelColor
}

export interface GridState {
  width: number
  height: number
  pixels: Map<string, PixelColor>
}

export const GRID_WIDTH = 1000
export const GRID_HEIGHT = 1000
export const DEFAULT_COLOR = PixelColor.WHITE
