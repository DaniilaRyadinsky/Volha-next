import type { HTMLAttributes, ReactNode } from 'react'
import clsx from 'clsx'
import styles from './Badge.module.scss'

export type BadgeTone =
  | 'default'
  | 'accent'
  | 'success'
  | 'muted'
  | 'danger'
  | 'warning'
  | 'info'
  | 'easy'
  | 'medium'
  | 'hard'
  | 'hit'

export type BadgeProps = HTMLAttributes<HTMLSpanElement> & {
  tone?: BadgeTone
  size?: 'sm' | 'md'
  icon?: ReactNode
}

export function Badge({ children, className, icon, size = 'md', tone = 'default', ...props }: BadgeProps) {
  return (
    <span className={clsx(styles.badge, styles[size], styles[tone], className)} {...props}>
      {icon && <span className={styles.icon}>{icon}</span>}
      <span>{children}</span>
    </span>
  )
}
