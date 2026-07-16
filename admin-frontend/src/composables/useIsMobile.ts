import {onMounted, onUnmounted, ref} from 'vue'
import {MOBILE_BREAKPOINT} from '@/constants/breakpoints'

export function useIsMobile() {
  const isMobile = ref(false)

  const checkMobile = () => {
    if (typeof window !== 'undefined') {
      isMobile.value = window.innerWidth <= MOBILE_BREAKPOINT
    }
  }

  const handleResize = () => {
    checkMobile()
  }

  onMounted(() => {
    checkMobile()
    if (typeof window !== 'undefined') {
      window.addEventListener('resize', handleResize)
    }
  })

  onUnmounted(() => {
    if (typeof window !== 'undefined') {
      window.removeEventListener('resize', handleResize)
    }
  })

  return {isMobile, checkMobile}
}
