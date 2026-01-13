/// <reference types="vite/client" />

declare module 'dplayer' {
  export interface DPlayerOptions {
    container: HTMLElement
    video: {
      url: string
      type?: 'auto' | 'hls' | 'flv' | 'dash' | 'webtorrent' | 'normal'
      pic?: string
      thumbnails?: string
      subtitle?: {
        url: string
        type?: string
        fontSize?: string
        bottom?: string
        color?: string
      }
    }
    autoplay?: boolean
    theme?: string
    loop?: boolean
    lang?: string | 'zh-cn' | 'zh-tw' | 'en' | 'ja'
    screenshot?: boolean
    hotkey?: boolean
    preload?: 'none' | 'metadata' | 'auto'
    volume?: number
    mutex?: boolean
    playbackSpeed?: number[]
    hlsConfig?: {
      xhrSetup?: (xhr: XMLHttpRequest) => void
    }
  }

  export default class DPlayer {
    constructor(options: DPlayerOptions)
    video: HTMLVideoElement
    destroy(): void
    on(event: string, handler: () => void): void
    off(event: string, handler: () => void): void
    play(): void
    pause(): void
    seek(time: number): void
    toggle(): void
    volume(percentage: number, notice?: boolean): void
    speed(rate: number): void
    notice(text: string, time?: number, opacity?: number): void
    switchQuality(index: number): void
    switchVideo(video: DPlayerOptions['video'], danmaku?: unknown): void
  }
}
