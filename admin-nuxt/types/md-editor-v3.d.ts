declare module 'md-editor-v3' {
  import { Component } from 'vue'
  
  export interface MdPreviewProps {
    editorId?: string
    modelValue?: string
    previewTheme?: string
    onHtmlChanged?: (html: string) => void
  }
  
  export const MdPreview: Component<MdPreviewProps>
}
