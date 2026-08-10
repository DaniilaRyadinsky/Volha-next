'use client'

import type { ButtonHTMLAttributes, ReactNode } from 'react'
import clsx from 'clsx'
import { Spinner } from '../Spinner/Spinner'
import styles from './IconButton.module.scss'

export type IconButtonProps = ButtonHTMLAttributes<HTMLButtonElement> & {
  label: string
  badge?: ReactNode
  variant?: 'default' | 'dark' | 'accent' | 'ghost' | 'danger'
  size?: 'sm' | 'md' | 'lg'
  loading?: boolean
}

export function IconButton({
  badge,
  children,
  className,
  disabled,
  label,
  loading = false,
  size = 'md',
  title,
  type = 'button',
  variant = 'default',
  ...props
}: IconButtonProps) {
  return (
    <button
      aria-busy={loading || undefined}
      aria-label={label}
      className={clsx(styles.button, styles[size], styles[variant], className)}
      disabled={disabled || loading}
      title={title ?? label}
      type={type}
      {...props}
    >
      {loading ? <Spinner label="Loading" size="sm" /> : children}
      {badge !== undefined && badge !== null && badge !== false && <span className={styles.badge}>{badge}</span>}
    </button>
  )
}
