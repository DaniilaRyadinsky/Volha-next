import type { HTMLAttributes } from 'react'
import clsx from 'clsx'
import styles from './VolhaLogo.module.scss'

export type VolhaLogoProps = HTMLAttributes<HTMLSpanElement> & {
  light?: boolean
  compact?: boolean
}

export function VolhaLogo({ className, compact = false, light = false, ...props }: VolhaLogoProps) {
  return (
    <span className={clsx(styles.logo, light && styles.light, compact && styles.compact, className)} {...props}>
      <svg aria-hidden="true" fill="none" viewBox="0 0 26 26">
        <rect fill="currentColor" height="2.5" width="22" x="2" y="4" />
        <rect fill="currentColor" height="2.5" width="22" x="2" y="11.75" />
        <rect fill="currentColor" height="2.5" width="22" x="2" y="19.5" />
        <rect fill="currentColor" height="18" width="2.5" x="2" y="4" />
        <rect fill="currentColor" height="18" width="2.5" x="21.5" y="4" />
      </svg>
      {!compact && <strong>ВОЛХА</strong>}
    </span>
  )
}
