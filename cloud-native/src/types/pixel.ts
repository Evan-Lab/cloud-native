export enum PixelColor {
  RED = '#EF4444',
  BLUE = '#3B82F6',
  YELLOW = '#FBBF24',
  GREEN = '#10B981',
  PURPLE = '#8B5CF6',
  PINK = '#EC4899',

  ORANGE = '#F97316',
  CYAN = '#06B6D4',
  LIME = '#84CC16',
  INDIGO = '#6366F1',
  ROSE = '#F43F5E',
  TEAL = '#14B8A6',

  WHITE = '#FFFFFF',
  LIGHT_GRAY = '#D1D5DB',
  GRAY = '#6B7280',
  DARK_GRAY = '#374151',
  BLACK = '#000000',

  BROWN = '#92400E',
  BEIGE = '#FDE68A',
  CORAL = '#FB7185',
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

export const GRID_WIDTH = 100
export const GRID_HEIGHT = 100
export const PIXEL_SIZE = 50
export const DEFAULT_COLOR = PixelColor.WHITE

export type Tool = 'brush' | 'eraser' | 'eyedropper'
