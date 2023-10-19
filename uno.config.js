import extractorSvelte from '@unocss/extractor-svelte'
import {
  presetUno,
  presetWebFonts,
  transformerDirectives,
  transformerVariantGroup,
} from 'unocss'

export default {
  presets: [
    presetUno(),
    presetWebFonts({
      provider: 'google',
      fonts: {
        mono: ['Fira Code'],
      },
    }),
  ],
  transformers: [transformerDirectives(), transformerVariantGroup()],
  safelist: [],
  theme: {},
  rules: [],
  shortcuts: [
    ['full', 'w-full h-full'],
    ['screen', 'w-screen h-screen'],
    ['max-full', 'max-w-full max-h-full'],
    ['max-screen', 'max-w-screen max-h-screen'],
  ],
  extractors: [extractorSvelte],
}
