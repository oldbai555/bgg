declare module 'dplayer' {
  export interface DPlayerOptions {
    container: HTMLElement
    video: {
      url: string
      type?: 'auto' | 'hls' | 'flv' | 'dash' | 'webtorrent' | 'normal'
      pic?: string
      thumbnails?: string
      subtitle?: string
      chapters?: string[]
    }
    autoplay?: boolean
    theme?: string
    loop?: boolean
    lang?: string
    screenshot?: boolean
    hotkey?: boolean
    preload?: 'none' | 'metadata' | 'auto'
    volume?: number
    mutex?: boolean
    playbackSpeed?: number[]
    hlsConfig?: Record<string, unknown>
  }

  export default class DPlayer {
    constructor(options: DPlayerOptions)
    video: HTMLVideoElement | null
    destroy(): void
    play(): void
    pause(): void
    seek(time: number): void
    toggle(): void
    on(event: string, handler: (...args: any[]) => void): void
    off(event: string, handler: (...args: any[]) => void): void
  }
}
