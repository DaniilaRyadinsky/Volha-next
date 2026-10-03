import type { HTMLAttributes } from 'react'
import clsx from 'clsx'
import styles from './Spinner.module.scss'

export type SpinnerProps = HTMLAttributes<HTMLSpanElement> & {
  label?: string
  size?: 'sm' | 'md' | 'lg'
}

export function Spinner({ className, label = 'Loading', size = 'md', ...props }: SpinnerProps) {
  return (
    <span
      aria-label={label}
      className={clsx(styles.spinner, styles[size], className)}
      role="status"
      {...props}
    />
  )
}
